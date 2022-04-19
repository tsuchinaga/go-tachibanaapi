package tachibana

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// TickGroupRequest - 呼値リクエスト
type TickGroupRequest struct{}

func (r *TickGroupRequest) request(no int64, now time.Time) tickGroupRequest {
	return tickGroupRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			MessageType:    MessageTypeEventDownload,
			ResponseFormat: commonResponseFormat,
		},
		TargetFeatures: string(MessageTypeTickGroup),
	}
}

type tickGroupRequest struct {
	commonRequest
	TargetFeatures string `json:"sTargetCLMID"` // 取得対象マスタリスト
}

type tickGroupResponse struct {
	commonResponse
	TickGroupType TickGroupType `json:"sYobineTaniNumber"`      // 呼値の単位番号
	StartDate     Ymd           `json:"sTekiyouDay"`            // 適用日
	BasePrice1    float64       `json:"sKizunPrice_1,string"`   // 基準値段1
	BasePrice2    float64       `json:"sKizunPrice_2,string"`   // 基準値段2
	BasePrice3    float64       `json:"sKizunPrice_3,string"`   // 基準値段3
	BasePrice4    float64       `json:"sKizunPrice_4,string"`   // 基準値段4
	BasePrice5    float64       `json:"sKizunPrice_5,string"`   // 基準値段5
	BasePrice6    float64       `json:"sKizunPrice_6,string"`   // 基準値段6
	BasePrice7    float64       `json:"sKizunPrice_7,string"`   // 基準値段7
	BasePrice8    float64       `json:"sKizunPrice_8,string"`   // 基準値段8
	BasePrice9    float64       `json:"sKizunPrice_9,string"`   // 基準値段9
	BasePrice10   float64       `json:"sKizunPrice_10,string"`  // 基準値段10
	BasePrice11   float64       `json:"sKizunPrice_11,string"`  // 基準値段11
	BasePrice12   float64       `json:"sKizunPrice_12,string"`  // 基準値段12
	BasePrice13   float64       `json:"sKizunPrice_13,string"`  // 基準値段13
	BasePrice14   float64       `json:"sKizunPrice_14,string"`  // 基準値段14
	BasePrice15   float64       `json:"sKizunPrice_15,string"`  // 基準値段15
	BasePrice16   float64       `json:"sKizunPrice_16,string"`  // 基準値段16
	BasePrice17   float64       `json:"sKizunPrice_17,string"`  // 基準値段17
	BasePrice18   float64       `json:"sKizunPrice_18,string"`  // 基準値段18
	BasePrice19   float64       `json:"sKizunPrice_19,string"`  // 基準値段19
	BasePrice20   float64       `json:"sKizunPrice_20,string"`  // 基準値段20
	UnitPrice1    float64       `json:"sYobineTanka_1,string"`  // 呼値値段1
	UnitPrice2    float64       `json:"sYobineTanka_2,string"`  // 呼値値段2
	UnitPrice3    float64       `json:"sYobineTanka_3,string"`  // 呼値値段3
	UnitPrice4    float64       `json:"sYobineTanka_4,string"`  // 呼値値段4
	UnitPrice5    float64       `json:"sYobineTanka_5,string"`  // 呼値値段5
	UnitPrice6    float64       `json:"sYobineTanka_6,string"`  // 呼値値段6
	UnitPrice7    float64       `json:"sYobineTanka_7,string"`  // 呼値値段7
	UnitPrice8    float64       `json:"sYobineTanka_8,string"`  // 呼値値段8
	UnitPrice9    float64       `json:"sYobineTanka_9,string"`  // 呼値値段9
	UnitPrice10   float64       `json:"sYobineTanka_10,string"` // 呼値値段10
	UnitPrice11   float64       `json:"sYobineTanka_11,string"` // 呼値値段11
	UnitPrice12   float64       `json:"sYobineTanka_12,string"` // 呼値値段12
	UnitPrice13   float64       `json:"sYobineTanka_13,string"` // 呼値値段13
	UnitPrice14   float64       `json:"sYobineTanka_14,string"` // 呼値値段14
	UnitPrice15   float64       `json:"sYobineTanka_15,string"` // 呼値値段15
	UnitPrice16   float64       `json:"sYobineTanka_16,string"` // 呼値値段16
	UnitPrice17   float64       `json:"sYobineTanka_17,string"` // 呼値値段17
	UnitPrice18   float64       `json:"sYobineTanka_18,string"` // 呼値値段18
	UnitPrice19   float64       `json:"sYobineTanka_19,string"` // 呼値値段19
	UnitPrice20   float64       `json:"sYobineTanka_20,string"` // 呼値値段20
	Digits1       float64       `json:"sDecimal_1,string"`      // 小数点桁数1
	Digits2       float64       `json:"sDecimal_2,string"`      // 小数点桁数2
	Digits3       float64       `json:"sDecimal_3,string"`      // 小数点桁数3
	Digits4       float64       `json:"sDecimal_4,string"`      // 小数点桁数4
	Digits5       float64       `json:"sDecimal_5,string"`      // 小数点桁数5
	Digits6       float64       `json:"sDecimal_6,string"`      // 小数点桁数6
	Digits7       float64       `json:"sDecimal_7,string"`      // 小数点桁数7
	Digits8       float64       `json:"sDecimal_8,string"`      // 小数点桁数8
	Digits9       float64       `json:"sDecimal_9,string"`      // 小数点桁数9
	Digits10      float64       `json:"sDecimal_10,string"`     // 小数点桁数10
	Digits11      float64       `json:"sDecimal_11,string"`     // 小数点桁数11
	Digits12      float64       `json:"sDecimal_12,string"`     // 小数点桁数12
	Digits13      float64       `json:"sDecimal_13,string"`     // 小数点桁数13
	Digits14      float64       `json:"sDecimal_14,string"`     // 小数点桁数14
	Digits15      float64       `json:"sDecimal_15,string"`     // 小数点桁数15
	Digits16      float64       `json:"sDecimal_16,string"`     // 小数点桁数16
	Digits17      float64       `json:"sDecimal_17,string"`     // 小数点桁数17
	Digits18      float64       `json:"sDecimal_18,string"`     // 小数点桁数18
	Digits19      float64       `json:"sDecimal_19,string"`     // 小数点桁数19
	Digits20      float64       `json:"sDecimal_20,string"`     // 小数点桁数20
	CreateDate    Hms           `json:"sCreateDate"`            // 作成日時
	UpdateDate    Hms           `json:"sUpdateDate"`            // 更新日時
}

func (r *tickGroupResponse) response() TickGroupResponse {
	return TickGroupResponse{
		CommonResponse: r.commonResponse.response(),
		TickGroupType:  r.TickGroupType,
		StartDate:      r.StartDate.Time,
		TickGroups: [20]TickGroup{
			{Number: 1, BasePrice: r.BasePrice1, UnitPrice: r.UnitPrice1, Digits: r.Digits1},
			{Number: 2, BasePrice: r.BasePrice2, UnitPrice: r.UnitPrice2, Digits: r.Digits2},
			{Number: 3, BasePrice: r.BasePrice3, UnitPrice: r.UnitPrice3, Digits: r.Digits3},
			{Number: 4, BasePrice: r.BasePrice4, UnitPrice: r.UnitPrice4, Digits: r.Digits4},
			{Number: 5, BasePrice: r.BasePrice5, UnitPrice: r.UnitPrice5, Digits: r.Digits5},
			{Number: 6, BasePrice: r.BasePrice6, UnitPrice: r.UnitPrice6, Digits: r.Digits6},
			{Number: 7, BasePrice: r.BasePrice7, UnitPrice: r.UnitPrice7, Digits: r.Digits7},
			{Number: 8, BasePrice: r.BasePrice8, UnitPrice: r.UnitPrice8, Digits: r.Digits8},
			{Number: 9, BasePrice: r.BasePrice9, UnitPrice: r.UnitPrice9, Digits: r.Digits9},
			{Number: 10, BasePrice: r.BasePrice10, UnitPrice: r.UnitPrice10, Digits: r.Digits10},
			{Number: 11, BasePrice: r.BasePrice11, UnitPrice: r.UnitPrice11, Digits: r.Digits11},
			{Number: 12, BasePrice: r.BasePrice12, UnitPrice: r.UnitPrice12, Digits: r.Digits12},
			{Number: 13, BasePrice: r.BasePrice13, UnitPrice: r.UnitPrice13, Digits: r.Digits13},
			{Number: 14, BasePrice: r.BasePrice14, UnitPrice: r.UnitPrice14, Digits: r.Digits14},
			{Number: 15, BasePrice: r.BasePrice15, UnitPrice: r.UnitPrice15, Digits: r.Digits15},
			{Number: 16, BasePrice: r.BasePrice16, UnitPrice: r.UnitPrice16, Digits: r.Digits16},
			{Number: 17, BasePrice: r.BasePrice17, UnitPrice: r.UnitPrice17, Digits: r.Digits17},
			{Number: 18, BasePrice: r.BasePrice18, UnitPrice: r.UnitPrice18, Digits: r.Digits18},
			{Number: 19, BasePrice: r.BasePrice19, UnitPrice: r.UnitPrice19, Digits: r.Digits19},
			{Number: 20, BasePrice: r.BasePrice20, UnitPrice: r.UnitPrice20, Digits: r.Digits20},
		},
		CreateDate: r.CreateDate.Time,
		UpdateDate: r.UpdateDate.Time,
	}
}

// TickGroupResponse - 呼値レスポンス
type TickGroupResponse struct {
	CommonResponse
	TickGroupType TickGroupType // 呼値の単位番号
	StartDate     time.Time     // 適用日
	TickGroups    [20]TickGroup // 呼値グループリスト
	CreateDate    time.Time     // 作成日時
	UpdateDate    time.Time     // 更新日時
}

// TickGroup - 呼値グループ
type TickGroup struct {
	Number    int     // 番号
	BasePrice float64 // 基準値段
	UnitPrice float64 // 呼値値段
	Digits    float64 // 小数点桁数
}

// TickGroup - 呼値
func (c *client) TickGroup(ctx context.Context, session *Session, req TickGroupRequest) ([]*TickGroupResponse, error) {
	if session == nil {
		return nil, NilArgumentErr
	}
	// マスタ情報ダウンロード系は他のリクエストと並行して実行できるため、ロックしない
	//session.mtx.Lock()
	//defer session.mtx.Unlock()

	session.lastRequestNo++
	r := req.request(session.lastRequestNo, c.clock.Now())

	// 終了通知を受け取って終了するために、停止可能にしておく
	cCtx, cf := context.WithCancel(ctx)
	defer cf()

	ch, errCh := c.requester.stream(cCtx, session.RequestURL, r)
	var responses []*TickGroupResponse
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

			var res tickGroupResponse
			if err := json.Unmarshal(b, &res); err != nil {
				return nil, fmt.Errorf("%s: %w", err, UnmarshalFailedErr)
			}

			// データ終了の合図が届いたらループを抜ける
			if res.MessageType == MessageTypeEventDownloadComplete {
				return responses, nil
			}

			Res := res.response()
			responses = append(responses, &Res)
		}
	}
}
