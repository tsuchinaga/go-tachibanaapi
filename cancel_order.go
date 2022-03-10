package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

// CancelOrderRequest - 取消注文リクエスト
type CancelOrderRequest struct {
	OrderNumber    string    // 注文番号
	BusinessDay    time.Time // 営業日
	SecondPassword string    // 第二パスワード
}

func (r *CancelOrderRequest) request(no int64, now time.Time) cancelOrderRequest {
	return cancelOrderRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			FeatureType:    FeatureTypeCancelOrder,
			ResponseFormat: commonResponseFormat,
		},
		OrderNumber:    r.OrderNumber,
		BusinessDay:    Ymd{Time: r.BusinessDay},
		SecondPassword: r.SecondPassword,
	}
}

type cancelOrderRequest struct {
	commonRequest
	OrderNumber    string `json:"sOrderNumber"`    // 注文番号
	BusinessDay    Ymd    `json:"sEigyouDay"`      // 営業日
	SecondPassword string `json:"sSecondPassword"` // 第二パスワード
}

type cancelOrderResponse struct {
	commonResponse
	ResultCode     string  `json:"sResultCode"`                   // 結果コード
	ResultText     string  `json:"sResultText"`                   // 結果テキスト
	OrderNumber    string  `json:"sOrderNumber"`                  // 注文番号
	BusinessDay    Ymd     `json:"sEigyouDay"`                    // 営業日
	DeliveryAmount float64 `json:"sOrderUkewatasiKingaku,string"` // 注文受渡金額
	OrderDateTime  YmdHms  `json:"sOrderDate"`                    // 注文日時
}

func (r *cancelOrderResponse) UnmarshalJSON(b []byte) error {
	// 文字列でないところに空文字を返されることがあるので、置換しておく
	replaced := b
	replaces := map[string]string{
		`"sOrderUkewatasiKingaku":""`: `"sOrderUkewatasiKingaku":"0"`,
	}
	for o, n := range replaces {
		replaced = bytes.Replace(replaced, []byte(o), []byte(n), -1)
	}

	type alias cancelOrderResponse
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *cancelOrderResponse) response() CancelOrderResponse {
	return CancelOrderResponse{
		CommonResponse: r.commonResponse.response(),
		ResultCode:     r.ResultCode,
		ResultText:     r.ResultText,
		OrderNumber:    r.OrderNumber,
		BusinessDay:    r.BusinessDay.Time,
		DeliveryAmount: r.DeliveryAmount,
		OrderDateTime:  r.OrderDateTime.Time,
	}
}

// CancelOrderResponse - 取消注文レスポンス
type CancelOrderResponse struct {
	CommonResponse
	ResultCode     string    // 結果コード
	ResultText     string    // 結果テキスト
	OrderNumber    string    // 注文番号
	BusinessDay    time.Time // 営業日
	DeliveryAmount float64   // 注文受渡金額
	OrderDateTime  time.Time // 注文日時
}

// CancelOrder - 取消注文
func (c *client) CancelOrder(ctx context.Context, session *Session, req CancelOrderRequest) (*CancelOrderResponse, error) {
	if session == nil {
		return nil, NilArgumentErr
	}
	session.mtx.Lock()
	defer session.mtx.Unlock()

	session.lastRequestNo++
	r := req.request(session.lastRequestNo, c.clock.Now())

	var res cancelOrderResponse
	if err := c.get(ctx, session.RequestURL, r, &res); err != nil {
		return nil, err
	}

	Res := res.response()
	return &Res, nil
}
