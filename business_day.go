package tachibana

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// BusinessDayRequest - 日付情報リクエスト
type BusinessDayRequest struct{}

func (r *BusinessDayRequest) request(no int64, now time.Time) businessDayRequest {
	return businessDayRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			FeatureType:    FeatureTypeEventDownload,
			ResponseFormat: commonResponseFormat,
		},
		TargetFeatures: string(FeatureTypeBusinessDay),
	}
}

type businessDayRequest struct {
	commonRequest
	TargetFeatures string `json:"sTargetCLMID"` // 取得対象マスタリスト
}

type businessDayResponse struct {
	commonResponse
	DayKey                 DayKey `json:"sDayKey"`               // 日付KEY
	PrevDay1               Ymd    `json:"sMaeEigyouDay_1"`       // 1営業日前
	PrevDay2               Ymd    `json:"sMaeEigyouDay_2"`       // 2営業日前
	PrevDay3               Ymd    `json:"sMaeEigyouDay_3"`       // 3営業日前
	Today                  Ymd    `json:"sTheDay"`               // 当日日付
	NextDay1               Ymd    `json:"sYokuEigyouDay_1"`      // 翌1営業日
	NextDay2               Ymd    `json:"sYokuEigyouDay_2"`      // 翌2営業日
	NextDay3               Ymd    `json:"sYokuEigyouDay_3"`      // 翌3営業日
	NextDay4               Ymd    `json:"sYokuEigyouDay_4"`      // 翌4営業日
	NextDay5               Ymd    `json:"sYokuEigyouDay_5"`      // 翌5営業日
	NextDay6               Ymd    `json:"sYokuEigyouDay_6"`      // 翌6営業日
	NextDay7               Ymd    `json:"sYokuEigyouDay_7"`      // 翌7営業日
	NextDay8               Ymd    `json:"sYokuEigyouDay_8"`      // 翌8営業日
	NextDay9               Ymd    `json:"sYokuEigyouDay_9"`      // 翌9営業日
	NextDay10              Ymd    `json:"sYokuEigyouDay_10"`     // 翌10営業日
	DeliveryDay            Ymd    `json:"sKabuUkewatasiDay"`     // 株式受渡日
	ProvisionalDeliveryDay Ymd    `json:"sKabuKariUkewatasiDay"` // 株式仮決受渡日
	BondDeliveryDay        Ymd    `json:"sBondUkewatasiDay"`     // 債券受渡日
}

func (r *businessDayResponse) response() BusinessDayResponse {
	return BusinessDayResponse{
		CommonResponse:         r.commonResponse.response(),
		DayKey:                 r.DayKey,
		PrevDay1:               r.PrevDay1.Time,
		PrevDay2:               r.PrevDay2.Time,
		PrevDay3:               r.PrevDay3.Time,
		Today:                  r.Today.Time,
		NextDay1:               r.NextDay1.Time,
		NextDay2:               r.NextDay2.Time,
		NextDay3:               r.NextDay3.Time,
		NextDay4:               r.NextDay4.Time,
		NextDay5:               r.NextDay5.Time,
		NextDay6:               r.NextDay6.Time,
		NextDay7:               r.NextDay7.Time,
		NextDay8:               r.NextDay8.Time,
		NextDay9:               r.NextDay9.Time,
		NextDay10:              r.NextDay10.Time,
		DeliveryDay:            r.DeliveryDay.Time,
		ProvisionalDeliveryDay: r.ProvisionalDeliveryDay.Time,
		BondDeliveryDay:        r.BondDeliveryDay.Time,
	}
}

// BusinessDayResponse - 日付情報レスポンス
type BusinessDayResponse struct {
	CommonResponse
	DayKey                 DayKey    // 日付KEY
	PrevDay1               time.Time // 1営業日前
	PrevDay2               time.Time // 2営業日前
	PrevDay3               time.Time // 3営業日前
	Today                  time.Time // 当日日付
	NextDay1               time.Time // 翌1営業日
	NextDay2               time.Time // 翌2営業日
	NextDay3               time.Time // 翌3営業日
	NextDay4               time.Time // 翌4営業日
	NextDay5               time.Time // 翌5営業日
	NextDay6               time.Time // 翌6営業日
	NextDay7               time.Time // 翌7営業日
	NextDay8               time.Time // 翌8営業日
	NextDay9               time.Time // 翌9営業日
	NextDay10              time.Time // 翌10営業日
	DeliveryDay            time.Time // 株式受渡日
	ProvisionalDeliveryDay time.Time // 株式仮決受渡日
	BondDeliveryDay        time.Time // 債券受渡日
}

// BusinessDay - 日付情報
func (c *client) BusinessDay(ctx context.Context, session *Session, req BusinessDayRequest) ([]*BusinessDayResponse, error) {
	if session == nil {
		return nil, NilArgumentErr
	}
	session.mtx.Lock()
	defer session.mtx.Unlock()

	session.lastRequestNo++
	r := req.request(session.lastRequestNo, c.clock.Now())

	// 終了通知を受け取って終了するために、停止可能にしておく
	cCtx, cf := context.WithCancel(ctx)
	defer cf()

	ch, errCh := c.requester.stream(cCtx, session.RequestURL, r)
	var responses []*BusinessDayResponse
	for {
		select {
		case err, ok := <-errCh:
			if ok {
				return nil, err
			}
		case b, ok := <-ch:
			// chanがcloseされたら抜ける
			if !ok {
				return responses, nil
			}

			var res businessDayResponse
			if err := json.Unmarshal(b, &res); err != nil {
				return nil, fmt.Errorf("%s: %w", err, UnmarshalFailedErr)
			}

			// データ終了の合図が届いたらループを抜ける
			if res.FeatureType == FeatureTypeEventDownloadComplete {
				return responses, nil
			}

			Res := res.response()
			responses = append(responses, &Res)
		}
	}
}
