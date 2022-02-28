package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

// OrderListDetailRequest - 注文約定一覧(詳細)リクエスト
type OrderListDetailRequest struct {
	OrderNumber   string    // 注文番号
	ExecutionDate time.Time // 営業日
}

func (r *OrderListDetailRequest) request(no int64, now time.Time) orderListDetailRequest {
	return orderListDetailRequest{
		commonRequest: commonRequest{
			No:          no,
			SendDate:    RequestTime{Time: now},
			FeatureType: FeatureTypeOrderListDetail,
		},
		OrderNumber:   r.OrderNumber,
		ExecutionDate: Ymd{Time: r.ExecutionDate},
	}
}

type orderListDetailRequest struct {
	commonRequest
	OrderNumber   string `json:"490,omitempty"` // 注文番号
	ExecutionDate Ymd    `json:"227,omitempty"` // 営業日
}

type orderListDetailResponse struct {
	commonResponse
	ResultCode             string           `json:"534"`        // 結果コード
	ResultText             string           `json:"535"`        // 結果テキスト
	WarningCode            string           `json:"692"`        // 警告コード
	WarningText            string           `json:"693"`        // 警告テキスト
	SymbolCode             string           `json:"328"`        // 銘柄CODE
	Exchange               Exchange         `json:"501"`        // 市場
	Side                   Side             `json:"467"`        // 売買区分
	TradeType              TradeType        `json:"255"`        // 現金信用区分
	ExitTermType           ExitTermType     `json:"468"`        // 弁済区分
	ExecutionTiming        ExecutionTiming  `json:"469"`        // 執行条件
	ExecutionType          ExecutionType    `json:"495"`        // 注文値段区分
	Price                  float64          `json:"494,string"` // 注文単価
	OrderQuantity          float64          `json:"496,string"` // 注文株数
	CurrentQuantity        float64          `json:"471,string"` // 有効株数
	OrderStatus            OrderStatus      `json:"504"`        // 状態コード
	OrderStatusText        string           `json:"503"`        // 状態
	OrderDateTime          YmdHms           `json:"491"`        // 注文日付
	ExpireDate             Ymd              `json:"492"`        // 有効期限
	Channel                Channel          `json:"193"`        // チャネル
	StockAccountType       AccountType      `json:"248"`        // 現物口座区分
	MarginAccountType      AccountType      `json:"575"`        // 建玉口座区分
	StopOrderType          StopOrderType    `json:"259"`        // 逆指値注文種別
	StopTriggerPrice       float64          `json:"263,string"` // 逆指値条件
	StopOrderExecutionType ExecutionType    `json:"258"`        // 逆指値値段区分
	StopOrderPrice         float64          `json:"260,string"` // 逆指値値段
	TriggerType            TriggerType      `json:"659"`        // トリガータイプ
	TriggerDateTime        YmdHms           `json:"658"`        // トリガー日時
	DeliveryDate           Ymd              `json:"662"`        // 受渡日
	ContractPrice          float64          `json:"695,string"` // 約定単価
	ContractQuantity       float64          `json:"696,string"` // 約定株数
	TradingAmount          float64          `json:"182,string"` // 売買代金
	PartContractType       PartContractType `json:"691"`        // 内出来区分
	EstimationAmount       float64          `json:"235,string"` // 概算代金
	Commission             float64          `json:"183,string"` // 手数料
	CommissionTax          float64          `json:"558,string"` // 消費税
	ExitOrderType          ExitOrderType    `json:"620"`        // 建日種類
	ExchangeErrorCode      string           `json:"577"`        // 市場/取次ErrorCode
	ExchangeOrderDateTime  YmdHms           `json:"466"`        // 市場注文受付時刻
	Contracts              []contract       `json:"57"`         // 約定失効リスト
	HoldPositions          []holdPosition   `json:"53"`         // 決済注文建株指定リスト
}

func (r *orderListDetailResponse) UnmarshalJSON(b []byte) error {
	// 文字列でないところに空文字を返されることがあるので、置換しておく
	replaced := b
	replaces := map[string]string{
		`"494":""`: `"494":"0"`,
		`"496":""`: `"496":"0"`,
		`"471":""`: `"471":"0"`,
		`"263":""`: `"263":"0"`,
		`"260":""`: `"260":"0"`,
		`"695":""`: `"695":"0"`,
		`"696":""`: `"696":"0"`,
		`"182":""`: `"182":"0"`,
		`"235":""`: `"235":"0"`,
		`"183":""`: `"183":"0"`,
		`"558":""`: `"558":"0"`,
		`"57":""`:  `"57":[]`,
		`"53":""`:  `"53":[]`,
	}
	for o, n := range replaces {
		replaced = bytes.Replace(replaced, []byte(o), []byte(n), -1)
	}

	type alias orderListDetailResponse
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *orderListDetailResponse) response() OrderListDetailResponse {
	contracts := make([]Contract, len(r.Contracts))
	for i, c := range r.Contracts {
		contracts[i] = c.response()
	}

	holdPositions := make([]HoldPosition, len(r.HoldPositions))
	for i, hp := range r.HoldPositions {
		holdPositions[i] = hp.response()
	}

	return OrderListDetailResponse{
		CommonResponse:         r.commonResponse.response(),
		ResultCode:             r.ResultCode,
		ResultText:             r.ResultText,
		WarningCode:            r.WarningCode,
		WarningText:            r.WarningText,
		SymbolCode:             r.SymbolCode,
		Exchange:               r.Exchange,
		Side:                   r.Side,
		TradeType:              r.TradeType,
		ExitTermType:           r.ExitTermType,
		ExecutionTiming:        r.ExecutionTiming,
		ExecutionType:          r.ExecutionType,
		Price:                  r.Price,
		OrderQuantity:          r.OrderQuantity,
		CurrentQuantity:        r.CurrentQuantity,
		OrderStatus:            r.OrderStatus,
		OrderStatusText:        r.OrderStatusText,
		OrderDateTime:          r.OrderDateTime.Time,
		ExpireDate:             r.ExpireDate.Time,
		Channel:                r.Channel,
		StockAccountType:       r.StockAccountType,
		MarginAccountType:      r.MarginAccountType,
		StopOrderType:          r.StopOrderType,
		StopTriggerPrice:       r.StopTriggerPrice,
		StopOrderExecutionType: r.StopOrderExecutionType,
		StopOrderPrice:         r.StopOrderPrice,
		TriggerType:            r.TriggerType,
		TriggerDateTime:        r.TriggerDateTime.Time,
		DeliveryDate:           r.DeliveryDate.Time,
		ContractPrice:          r.ContractPrice,
		ContractQuantity:       r.ContractQuantity,
		TradingAmount:          r.TradingAmount,
		PartContractType:       r.PartContractType,
		EstimationAmount:       r.EstimationAmount,
		Commission:             r.Commission,
		CommissionTax:          r.CommissionTax,
		ExitOrderType:          r.ExitOrderType,
		ExchangeErrorCode:      r.ExchangeErrorCode,
		ExchangeOrderDateTime:  r.ExchangeOrderDateTime.Time,
		Contracts:              contracts,
		HoldPositions:          holdPositions,
	}
}

type contract struct {
	WarningCode string  `json:"697"`        // 警告コード
	WarningText string  `json:"698"`        // 警告テキスト
	Quantity    float64 `json:"696,string"` // 約定数量
	Price       float64 `json:"695,string"` // 約定価格
	DateTime    YmdHms  `json:"694"`        // 約定日時
}

func (r *contract) response() Contract {
	return Contract{
		WarningCode: r.WarningCode,
		WarningText: r.WarningText,
		Quantity:    r.Quantity,
		Price:       r.Price,
		DateTime:    r.DateTime.Time,
	}
}

type holdPosition struct {
	WarningCode   string  `json:"365"`        // 警告コード
	WarningText   string  `json:"366"`        // 警告テキスト
	SortOrder     int     `json:"360,string"` // 順位
	ContractDate  Ymd     `json:"361"`        // 建日
	EntryPrice    float64 `json:"362"`        // 建単価
	HoldQuantity  float64 `json:"356,string"` // 返済注文株数
	ExitQuantity  float64 `json:"368,string"` // 約定株数
	ExitPrice     float64 `json:"367,string"` // 約定単価
	Commission    float64 `json:"359,string"` // 建手数料
	Interest      float64 `json:"369,string"` // 順日歩
	Premiums      float64 `json:"352,string"` // 逆日歩
	RewritingFee  float64 `json:"353,string"` // 書換料
	ManagementFee float64 `json:"354,string"` // 管理費
	LendingFee    float64 `json:"355,string"` // 貸株料
	OtherFee      float64 `json:"358,string"` // その他
	Profit        float64 `json:"357,string"` // 決済損益/受渡代金
}

func (r *holdPosition) response() HoldPosition {
	return HoldPosition{
		WarningCode:   r.WarningCode,
		WarningText:   r.WarningText,
		SortOrder:     r.SortOrder,
		ContractDate:  r.ContractDate.Time,
		EntryPrice:    r.EntryPrice,
		HoldQuantity:  r.HoldQuantity,
		ExitQuantity:  r.ExitQuantity,
		ExitPrice:     r.ExitPrice,
		Commission:    r.Commission,
		Interest:      r.Interest,
		Premiums:      r.Premiums,
		RewritingFee:  r.RewritingFee,
		ManagementFee: r.ManagementFee,
		LendingFee:    r.LendingFee,
		OtherFee:      r.OtherFee,
		Profit:        r.Profit,
	}
}

// OrderListDetailResponse - 注文約定一覧(詳細)レスポンス
type OrderListDetailResponse struct {
	CommonResponse
	ResultCode             string           // 結果コード
	ResultText             string           // 結果テキスト
	WarningCode            string           // 警告コード
	WarningText            string           // 警告テキスト
	SymbolCode             string           // 銘柄CODE
	Exchange               Exchange         // 市場
	Side                   Side             // 売買区分
	TradeType              TradeType        // 現金信用区分
	ExitTermType           ExitTermType     // 弁済区分
	ExecutionTiming        ExecutionTiming  // 執行条件
	ExecutionType          ExecutionType    // 注文値段区分
	Price                  float64          // 注文単価
	OrderQuantity          float64          // 注文株数
	CurrentQuantity        float64          // 有効株数
	OrderStatus            OrderStatus      // 状態コード
	OrderStatusText        string           // 状態
	OrderDateTime          time.Time        // 注文日付
	ExpireDate             time.Time        // 有効期限
	Channel                Channel          // チャネル
	StockAccountType       AccountType      // 現物口座区分
	MarginAccountType      AccountType      // 建玉口座区分
	StopOrderType          StopOrderType    // 逆指値注文種別
	StopTriggerPrice       float64          // 逆指値条件
	StopOrderExecutionType ExecutionType    // 逆指値値段区分
	StopOrderPrice         float64          // 逆指値値段
	TriggerType            TriggerType      // トリガータイプ
	TriggerDateTime        time.Time        // トリガー日時
	DeliveryDate           time.Time        // 受渡日
	ContractPrice          float64          // 約定単価
	ContractQuantity       float64          // 約定株数
	TradingAmount          float64          // 売買代金
	PartContractType       PartContractType // 内出来区分
	EstimationAmount       float64          // 概算代金
	Commission             float64          // 手数料
	CommissionTax          float64          // 消費税
	ExitOrderType          ExitOrderType    // 建日種類
	ExchangeErrorCode      string           // 市場/取次ErrorCode
	ExchangeOrderDateTime  time.Time        // 市場注文受付時刻
	Contracts              []Contract       // 約定失効リスト
	HoldPositions          []HoldPosition   // 決済注文建株指定リスト
}

// Contract - 約定失効
type Contract struct {
	WarningCode string    // 警告コード
	WarningText string    // 警告テキスト
	Quantity    float64   // 約定数量
	Price       float64   // 約定価格
	DateTime    time.Time // 約定日時
}

// HoldPosition - 決済注文建株
type HoldPosition struct {
	WarningCode   string    // 警告コード
	WarningText   string    // 警告テキスト
	SortOrder     int       // 順位
	ContractDate  time.Time // 建日
	EntryPrice    float64   // 建単価
	HoldQuantity  float64   // 返済注文株数
	ExitQuantity  float64   // 約定株数
	ExitPrice     float64   // 約定単価
	Commission    float64   // 建手数料
	Interest      float64   // 順日歩
	Premiums      float64   // 逆日歩
	RewritingFee  float64   // 書換料
	ManagementFee float64   // 管理費
	LendingFee    float64   // 貸株料
	OtherFee      float64   // その他
	Profit        float64   // 決済損益/受渡代金
}

// OrderListDetail - 注文約定一覧(詳細)
func (c *client) OrderListDetail(ctx context.Context, session *Session, req OrderListDetailRequest) (*OrderListDetailResponse, error) {
	if session == nil {
		return nil, NilArgumentErr
	}
	session.mtx.Lock()
	defer session.mtx.Unlock()

	session.lastRequestNo++
	r := req.request(session.lastRequestNo, c.clock.Now())

	var res orderListDetailResponse
	if err := c.get(ctx, session.RequestURL, r, &res); err != nil {
		return nil, err
	}

	Res := res.response()
	return &Res, nil
}
