package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"time"
)

// NewOrderRequest - 新規注文リクエスト
type NewOrderRequest struct {
	StockAccountType  AccountType     // 譲渡益課税区分
	MarginAccountType AccountType     // 建玉譲渡益課税区分
	SymbolCode        string          // 銘柄コード
	Exchange          Exchange        // 市場
	Side              Side            // 売買区分
	ExecutionTiming   ExecutionTiming // 執行条件
	OrderPrice        float64         // 注文値段
	OrderQuantity     float64         // 注文数量
	TradeType         TradeType       // 現金信用区分
	ExpireDate        time.Time       // 注文期日
	ExpireDateIsToday bool            // 注文期日を当日
	StopOrderType     StopOrderType   // 逆指値注文種別
	TriggerPrice      float64         // 逆指値条件
	StopOrderPrice    float64         // 逆指値値段
	ExitOrderType     ExitOrderType   // 建日種類
	SecondPassword    string          // 第二パスワード
	ExitPositions     []ExitPosition  // 返済リスト
}

func (r NewOrderRequest) request(no int64, now time.Time) newOrderRequest {
	exitPositions := make([]ExitPosition, len(r.ExitPositions))
	for i, p := range r.ExitPositions {
		exitPositions[i] = p
	}

	orderPrice := "*" // 指定なし
	if r.StopOrderType != StopOrderTypeStop {
		orderPrice = strconv.FormatFloat(r.OrderPrice, 'f', -1, 64)
	}

	stopOrderPrice := "*" // 指定なし
	if r.StopOrderType != StopOrderTypeNormal {
		stopOrderPrice = strconv.FormatFloat(r.StopOrderPrice, 'f', -1, 64)
	}

	return newOrderRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			FeatureType:    FeatureTypeNewOrder,
			ResponseFormat: commonResponseFormat,
		},
		StockAccountType:  r.StockAccountType,
		MarginAccountType: r.MarginAccountType,
		SymbolCode:        r.SymbolCode,
		Exchange:          r.Exchange,
		Side:              r.Side,
		ExecutionTiming:   r.ExecutionTiming,
		OrderPrice:        orderPrice,
		OrderQuantity:     r.OrderQuantity,
		TradeType:         r.TradeType,
		ExpireDate:        Ymd{Time: r.ExpireDate, isToday: r.ExpireDateIsToday},
		StopOrderType:     r.StopOrderType,
		TriggerPrice:      r.TriggerPrice,
		StopOrderPrice:    stopOrderPrice,
		ExitOrderType:     r.ExitOrderType,
		SecondPassword:    r.SecondPassword,
		ExitPositions:     exitPositions,
	}
}

type newOrderRequest struct {
	commonRequest
	StockAccountType  AccountType     `json:"sZyoutoekiKazeiC"`          // 譲渡益課税区分
	MarginAccountType AccountType     `json:"sTategyokuZyoutoekiKazeiC"` // 建玉譲渡益課税区分
	SymbolCode        string          `json:"sIssueCode"`                // 銘柄コード
	Exchange          Exchange        `json:"sSizyouC"`                  // 市場
	Side              Side            `json:"sBaibaiKubun"`              // 売買区分
	ExecutionTiming   ExecutionTiming `json:"sCondition"`                // 執行条件
	OrderPrice        string          `json:"sOrderPrice"`               // 注文値段
	OrderQuantity     float64         `json:"sOrderSuryou,string"`       // 注文数量
	TradeType         TradeType       `json:"sGenkinShinyouKubun"`       // 現金信用区分
	ExpireDate        Ymd             `json:"sOrderExpireDay"`           // 注文期日
	StopOrderType     StopOrderType   `json:"sGyakusasiOrderType"`       // 逆指値注文種別
	TriggerPrice      float64         `json:"sGyakusasiZyouken,string"`  // 逆指値条件
	StopOrderPrice    string          `json:"sGyakusasiPrice"`           // 逆指値値段
	ExitOrderType     ExitOrderType   `json:"sTatebiType"`               // 建日種類
	SecondPassword    string          `json:"sSecondPassword"`           // 第二パスワード
	ExitPositions     []ExitPosition  `json:"aCLMKabuHensaiData"`        // 返済リスト
}

type ExitPosition struct {
	PositionNumber string `json:"sTategyokuNumber"` // 新規建玉番号
	SequenceNumber string `json:"sTatebiZyuni"`     // 建日順位
	OrderQuantity  string `json:"sOrderSuryou"`     // 注文数量
}

type newOrderResponse struct {
	commonResponse
	ResultCode     string  `json:"sResultCode"`                   // 結果コード
	ResultText     string  `json:"sResultText"`                   // 結果テキスト
	WarningCode    string  `json:"sWarningCode"`                  // 警告コード
	WarningText    string  `json:"sWarningText"`                  // 警告テキスト
	OrderNumber    string  `json:"sOrderNumber"`                  // 注文番号
	BusinessDay    Ymd     `json:"sEigyouDay"`                    // 営業日
	DeliveryAmount float64 `json:"sOrderUkewatasiKingaku,string"` // 注文受渡金額
	Commission     float64 `json:"sOrderTesuryou,string"`         // 注文手数料
	CommissionTax  float64 `json:"sOrderSyouhizei,string"`        // 注文消費税
	Interest       float64 `json:"sKinri,string"`                 // 金利
	OrderDateTime  YmdHms  `json:"sOrderDate"`                    // 注文日時
}

func (r *newOrderResponse) UnmarshalJSON(b []byte) error {
	// 文字列でないところに空文字を返されることがあるので、置換しておく
	replaced := b
	replaces := map[string]string{
		`"sOrderUkewatasiKingaku":""`: `"sOrderUkewatasiKingaku":"0"`,
		`"sOrderTesuryou":""`:         `"sOrderTesuryou":"0"`,
		`"sOrderSyouhizei":""`:        `"sOrderSyouhizei":"0"`,
		`"sKinri":""`:                 `"sKinri":"0"`,
		`"sKinri":"-"`:                `"sKinri":"0"`,
	}
	for o, n := range replaces {
		replaced = bytes.Replace(replaced, []byte(o), []byte(n), -1)
	}

	type alias newOrderResponse
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *newOrderResponse) response() NewOrderResponse {
	return NewOrderResponse{
		CommonResponse: r.commonResponse.response(),
		ResultCode:     r.ResultCode,
		ResultText:     r.ResultText,
		WarningCode:    r.WarningCode,
		WarningText:    r.WarningText,
		OrderNumber:    r.OrderNumber,
		BusinessDay:    r.BusinessDay.Time,
		DeliveryAmount: r.DeliveryAmount,
		Commission:     r.Commission,
		CommissionTax:  r.CommissionTax,
		Interest:       r.Interest,
		OrderDateTime:  r.OrderDateTime.Time,
	}
}

// NewOrderResponse - 新規注文レスポンス
type NewOrderResponse struct {
	CommonResponse
	ResultCode     string    // 結果コード
	ResultText     string    // 結果テキスト
	WarningCode    string    // 警告コード
	WarningText    string    // 警告テキスト
	OrderNumber    string    // 注文番号
	BusinessDay    time.Time // 営業日
	DeliveryAmount float64   // 注文受渡金額
	Commission     float64   // 注文手数料
	CommissionTax  float64   // 注文消費税
	Interest       float64   // 金利
	OrderDateTime  time.Time // 注文日時
}

// NewOrder - 新規注文
func (c *client) NewOrder(ctx context.Context, session *Session, req NewOrderRequest) (*NewOrderResponse, error) {
	if session == nil {
		return nil, NilArgumentErr
	}
	session.mtx.Lock()
	defer session.mtx.Unlock()

	session.lastRequestNo++
	r := req.request(session.lastRequestNo, c.clock.Now())

	var res newOrderResponse
	if err := c.get(ctx, session.RequestURL, r, &res); err != nil {
		return nil, err
	}

	Res := res.response()
	return &Res, nil
}
