package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"time"
)

// CorrectOrderRequest - 訂正注文リクエスト
type CorrectOrderRequest struct {
	OrderNumber        string          `json:"sOrderNumber"` // 注文番号
	ExecutionDay       time.Time       `json:"sEigyouDay"`   // 営業日
	ExecutionTiming    ExecutionTiming `json:"sCondition"`   // 執行条件
	OrderPrice         float64         `json:"sOrderPrice"`  // 注文値段
	OrderQuantity      float64         `json:"sOrderSuryou"` // 注文数量
	ExpireDate         time.Time       // 注文期日
	ExpireDateIsToday  bool            // 注文期日を当日
	ExpireDateNoChange bool            // 注文期日を変更しない
	TriggerPrice       float64         `json:"sGyakusasiZyouken,string"` // 逆指値条件
	StopOrderPrice     float64         `json:"sGyakusasiPrice"`          // 逆指値値段
	SecondPassword     string          `json:"sSecondPassword"`          // 第二パスワード
}

func (r *CorrectOrderRequest) request(no int64, now time.Time) correctOrderRequest {
	orderPrice := "*"
	if r.OrderPrice != NoChangeFloat {
		orderPrice = strconv.FormatFloat(r.OrderPrice, 'f', -1, 64)
	}

	orderQuantity := "*"
	if r.OrderQuantity != NoChangeFloat {
		orderQuantity = strconv.FormatFloat(r.OrderQuantity, 'f', -1, 64)
	}

	triggerPrice := "*"
	if r.TriggerPrice != NoChangeFloat {
		triggerPrice = strconv.FormatFloat(r.TriggerPrice, 'f', -1, 64)
	}

	stopOrderPrice := "*"
	if r.StopOrderPrice != NoChangeFloat {
		stopOrderPrice = strconv.FormatFloat(r.StopOrderPrice, 'f', -1, 64)
	}

	return correctOrderRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			FeatureType:    FeatureTypeCorrectOrder,
			ResponseFormat: commonResponseFormat,
		},
		OrderNumber:     r.OrderNumber,
		ExecutionDay:    Ymd{Time: r.ExecutionDay},
		ExecutionTiming: r.ExecutionTiming,
		OrderPrice:      orderPrice,
		OrderQuantity:   orderQuantity,
		ExpireDate: Ymd{
			Time:       r.ExpireDate,
			isNoChange: r.ExpireDateNoChange,
			isToday:    r.ExpireDateIsToday,
		},
		TriggerPrice:   triggerPrice,
		StopOrderPrice: stopOrderPrice,
		SecondPassword: r.SecondPassword,
	}
}

type correctOrderRequest struct {
	commonRequest
	OrderNumber     string          `json:"sOrderNumber"`      // 注文番号
	ExecutionDay    Ymd             `json:"sEigyouDay"`        // 営業日
	ExecutionTiming ExecutionTiming `json:"sCondition"`        // 執行条件
	OrderPrice      string          `json:"sOrderPrice"`       // 注文値段
	OrderQuantity   string          `json:"sOrderSuryou"`      // 注文数量
	ExpireDate      Ymd             `json:"sOrderExpireDay"`   // 注文期日
	TriggerPrice    string          `json:"sGyakusasiZyouken"` // 逆指値条件
	StopOrderPrice  string          `json:"sGyakusasiPrice"`   // 逆指値値段
	SecondPassword  string          `json:"sSecondPassword"`   // 第二パスワード
}

type correctOrderResponse struct {
	commonResponse
	ResultCode     string  `json:"sResultCode"`                   // 結果コード
	ResultText     string  `json:"sResultText"`                   // 結果テキスト
	OrderNumber    string  `json:"sOrderNumber"`                  // 注文番号
	ExecutionDay   Ymd     `json:"sEigyouDay"`                    // 営業日
	DeliveryAmount float64 `json:"sOrderUkewatasiKingaku,string"` // 注文受渡金額
	Commission     float64 `json:"sOrderTesuryou,string"`         // 注文手数料
	CommissionTax  float64 `json:"sOrderSyouhizei,string"`        // 注文消費税
	OrderDateTime  YmdHms  `json:"sOrderDate"`                    // 注文日時
}

func (r *correctOrderResponse) UnmarshalJSON(b []byte) error {
	// 文字列でないところに空文字を返されることがあるので、置換しておく
	replaced := b
	replaces := map[string]string{
		`"sOrderUkewatasiKingaku":""`: `"sOrderUkewatasiKingaku":"0"`,
		`"sOrderTesuryou":""`:         `"sOrderTesuryou":"0"`,
		`"sOrderSyouhizei":""`:        `"sOrderSyouhizei":"0"`,
	}
	for o, n := range replaces {
		replaced = bytes.Replace(replaced, []byte(o), []byte(n), -1)
	}

	type alias correctOrderResponse
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *correctOrderResponse) response() CorrectOrderResponse {
	return CorrectOrderResponse{
		CommonResponse: r.commonResponse.response(),
		ResultCode:     r.ResultCode,
		ResultText:     r.ResultText,
		OrderNumber:    r.OrderNumber,
		ExecutionDay:   r.ExecutionDay.Time,
		DeliveryAmount: r.DeliveryAmount,
		Commission:     r.Commission,
		CommissionTax:  r.CommissionTax,
		OrderDateTime:  r.OrderDateTime.Time,
	}
}

// CorrectOrderResponse - 訂正注文レスポンス
type CorrectOrderResponse struct {
	CommonResponse
	ResultCode     string    // 結果コード
	ResultText     string    // 結果テキスト
	OrderNumber    string    // 注文番号
	ExecutionDay   time.Time // 営業日
	DeliveryAmount float64   // 注文受渡金額
	Commission     float64   // 注文手数料
	CommissionTax  float64   // 注文消費税
	OrderDateTime  time.Time // 注文日時
}

// CorrectOrder - 訂正注文
func (c *client) CorrectOrder(ctx context.Context, session *Session, req CorrectOrderRequest) (*CorrectOrderResponse, error) {
	if session == nil {
		return nil, NilArgumentErr
	}
	session.mtx.Lock()
	defer session.mtx.Unlock()

	session.lastRequestNo++
	r := req.request(session.lastRequestNo, c.clock.Now())

	var res correctOrderResponse
	if err := c.get(ctx, session.RequestURL, r, &res); err != nil {
		return nil, err
	}

	Res := res.response()
	return &Res, nil
}
