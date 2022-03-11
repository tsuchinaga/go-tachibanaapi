package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

// OrderListDetailRequest - 注文約定一覧(詳細)リクエスト
type OrderListDetailRequest struct {
	OrderNumber  string    // 注文番号
	ExecutionDay time.Time // 営業日
}

func (r *OrderListDetailRequest) request(no int64, now time.Time) orderListDetailRequest {
	return orderListDetailRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			FeatureType:    FeatureTypeOrderListDetail,
			ResponseFormat: commonResponseFormat,
		},
		OrderNumber:  r.OrderNumber,
		ExecutionDay: Ymd{Time: r.ExecutionDay},
	}
}

type orderListDetailRequest struct {
	commonRequest
	OrderNumber  string `json:"sOrderNumber,omitempty"` // 注文番号
	ExecutionDay Ymd    `json:"sEigyouDay,omitempty"`   // 営業日
}

type orderListDetailResponse struct {
	commonResponse
	OrderNumber            string           `json:"sOrderNumber"`               // 注文番号
	ExecutionDay           Ymd              `json:"sEigyouDay"`                 // 営業日
	ResultCode             string           `json:"sResultCode"`                // 結果コード
	ResultText             string           `json:"sResultText"`                // 結果テキスト
	WarningCode            string           `json:"sWarningCode"`               // 警告コード
	WarningText            string           `json:"sWarningText"`               // 警告テキスト
	IssueCode              string           `json:"sIssueCode"`                 // 銘柄CODE
	Exchange               Exchange         `json:"sOrderSizyouC"`              // 市場
	Side                   Side             `json:"sOrderBaibaiKubun"`          // 売買区分
	TradeType              TradeType        `json:"sGenkinSinyouKubun"`         // 現金信用区分
	ExitTermType           ExitTermType     `json:"sOrderBensaiKubun"`          // 弁済区分
	ExecutionTiming        ExecutionTiming  `json:"sOrderCondition"`            // 執行条件
	ExecutionType          ExecutionType    `json:"sOrderOrderPriceKubun"`      // 注文値段区分
	Price                  float64          `json:"sOrderOrderPrice,string"`    // 注文単価
	OrderQuantity          float64          `json:"sOrderOrderSuryou,string"`   // 注文株数
	CurrentQuantity        float64          `json:"sOrderCurrentSuryou,string"` // 有効株数
	OrderStatus            OrderStatus      `json:"sOrderStatusCode"`           // 状態コード
	OrderStatusText        string           `json:"sOrderStatus"`               // 状態
	OrderDateTime          YmdHms           `json:"sOrderOrderDateTime"`        // 注文日付
	ExpireDate             Ymd              `json:"sOrderOrderExpireDay"`       // 有効期限
	Channel                Channel          `json:"sChannel"`                   // チャネル
	StockAccountType       AccountType      `json:"sGenbutuZyoutoekiKazeiC"`    // 現物口座区分
	MarginAccountType      AccountType      `json:"sSinyouZyoutoekiKazeiC"`     // 建玉口座区分
	StopOrderType          StopOrderType    `json:"sGyakusasiOrderType"`        // 逆指値注文種別
	StopTriggerPrice       float64          `json:"sGyakusasiZyouken,string"`   // 逆指値条件
	StopOrderExecutionType ExecutionType    `json:"sGyakusasiKubun"`            // 逆指値値段区分
	StopOrderPrice         float64          `json:"sGyakusasiPrice,string"`     // 逆指値値段
	TriggerType            TriggerType      `json:"sTriggerType"`               // トリガータイプ
	TriggerDateTime        YmdHms           `json:"sTriggerTime"`               // トリガー日時
	DeliveryDate           Ymd              `json:"sUkewatasiDay"`              // 受渡日
	ContractPrice          float64          `json:"sYakuzyouPrice,string"`      // 約定単価
	ContractQuantity       float64          `json:"sYakuzyouSuryou,string"`     // 約定株数
	TradingAmount          float64          `json:"sBaiBaiDaikin,string"`       // 売買代金
	PartContractType       PartContractType `json:"sUtidekiKubun"`              // 内出来区分
	EstimationAmount       float64          `json:"sGaisanDaikin,string"`       // 概算代金
	Commission             float64          `json:"sBaiBaiTesuryo,string"`      // 手数料
	CommissionTax          float64          `json:"sShouhizei,string"`          // 消費税
	ExitOrderType          ExitOrderType    `json:"sTatebiType"`                // 建日種類
	ExchangeErrorCode      string           `json:"sSizyouErrorCode"`           // 市場/取次ErrorCode
	ExchangeOrderDateTime  YmdHms           `json:"sOrderAcceptTime"`           // 市場注文受付時刻
	Contracts              []contract       `json:"aYakuzyouSikkouList"`        // 約定失効リスト
	HoldPositions          []holdPosition   `json:"aKessaiOrderTategyokuList"`  // 決済注文建株指定リスト
}

func (r *orderListDetailResponse) UnmarshalJSON(b []byte) error {
	// 文字列でないところに空文字を返されることがあるので、置換しておく
	replaced := b
	replaces := map[string]string{
		`"sOrderOrderPrice":""`:          `"sOrderOrderPrice":"0"`,
		`"sOrderOrderSuryou":""`:         `"sOrderOrderSuryou":"0"`,
		`"sOrderCurrentSuryou":""`:       `"sOrderCurrentSuryou":"0"`,
		`"sGyakusasiZyouken":""`:         `"sGyakusasiZyouken":"0"`,
		`"sGyakusasiPrice":""`:           `"sGyakusasiPrice":"0"`,
		`"sYakuzyouPrice":""`:            `"sYakuzyouPrice":"0"`,
		`"sYakuzyouSuryou":""`:           `"sYakuzyouSuryou":"0"`,
		`"sBaiBaiDaikin":""`:             `"sBaiBaiDaikin":"0"`,
		`"sGaisanDaikin":""`:             `"sGaisanDaikin":"0"`,
		`"sBaiBaiTesuryo":""`:            `"sBaiBaiTesuryo":"0"`,
		`"sShouhizei":""`:                `"sShouhizei":"0"`,
		`"aYakuzyouSikkouList":""`:       `"aYakuzyouSikkouList":[]`,
		`"aKessaiOrderTategyokuList":""`: `"aKessaiOrderTategyokuList":[]`,
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
		OrderNumber:            r.OrderNumber,
		ExecutionDate:          r.ExecutionDay.Time,
		ResultCode:             r.ResultCode,
		ResultText:             r.ResultText,
		WarningCode:            r.WarningCode,
		WarningText:            r.WarningText,
		IssueCode:              r.IssueCode,
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
	WarningCode string  `json:"sYakuzyouWarningCode"`   // 警告コード
	WarningText string  `json:"sYakuzyouWarningText"`   // 警告テキスト
	Quantity    float64 `json:"sYakuzyouSuryou,string"` // 約定数量
	Price       float64 `json:"sYakuzyouPrice,string"`  // 約定価格
	DateTime    YmdHms  `json:"sYakuzyouDate"`          // 約定日時
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
	WarningCode   string  `json:"sKessaiWarningCode"`           // 警告コード
	WarningText   string  `json:"sKessaiWarningText"`           // 警告テキスト
	SortOrder     int     `json:"sKessaiTatebiZyuni,string"`    // 順位
	ContractDate  Ymd     `json:"sKessaiTategyokuDay"`          // 建日
	EntryPrice    float64 `json:"sKessaiTategyokuPrice,string"` // 建単価
	HoldQuantity  float64 `json:"sKessaiOrderSuryo,string"`     // 返済注文株数
	ExitQuantity  float64 `json:"sKessaiYakuzyouSuryo,string"`  // 約定株数
	ExitPrice     float64 `json:"sKessaiYakuzyouPrice,string"`  // 約定単価
	Commission    float64 `json:"sKessaiTateTesuryou,string"`   // 建手数料
	Interest      float64 `json:"sKessaiZyunHibu,string"`       // 順日歩
	Premiums      float64 `json:"sKessaiGyakuhibu,string"`      // 逆日歩
	RewritingFee  float64 `json:"sKessaiKakikaeryou,string"`    // 書換料
	ManagementFee float64 `json:"sKessaiKanrihi,string"`        // 管理費
	LendingFee    float64 `json:"sKessaiKasikaburyou,string"`   // 貸株料
	OtherFee      float64 `json:"sKessaiSonota,string"`         // その他
	Profit        float64 `json:"sKessaiSoneki,string"`         // 決済損益/受渡代金
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
	OrderNumber            string           // 注文番号
	ExecutionDate          time.Time        // 営業日
	ResultCode             string           // 結果コード
	ResultText             string           // 結果テキスト
	WarningCode            string           // 警告コード
	WarningText            string           // 警告テキスト
	IssueCode              string           // 銘柄CODE
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
