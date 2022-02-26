package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

// OrderListRequest - 注文一覧リクエスト
type OrderListRequest struct {
	SymbolCode         string             // 銘柄コード
	ExecutionDate      time.Time          // 注文執行予定日
	OrderInquiryStatus OrderInquiryStatus // 注文照会状態
}

func (r *OrderListRequest) request(no int64, now time.Time) orderListRequest {
	return orderListRequest{
		commonRequest: commonRequest{
			No:          no,
			SendDate:    RequestTime{Time: now},
			FeatureType: FeatureTypeOrderList,
		},
		SymbolCode:    r.SymbolCode,
		ExecutionDate: Ymd{Time: r.ExecutionDate},
		OrderStatus:   r.OrderInquiryStatus,
	}
}

type orderListRequest struct {
	commonRequest
	SymbolCode    string             `json:"328,omitempty"` // 銘柄コード
	ExecutionDate Ymd                `json:"559,omitempty"` // 注文執行予定日
	OrderStatus   OrderInquiryStatus `json:"508,omitempty"` // 注文照会状態
}

type orderListResponse struct {
	commonResponse
	SymbolCode         string             `json:"328"` // 銘柄コード
	ExecutionDate      Ymd                `json:"559"` // 注文執行予定日
	OrderInquiryStatus OrderInquiryStatus `json:"508"` // 注文照会状態
	ResultCode         string             `json:"534"` // 結果コード
	ResultText         string             `json:"535"` // 結果テキスト
	WarningCode        string             `json:"692"` // 警告コード
	WarningText        string             `json:"693"` // 警告テキスト
	Orders             []order            `json:"55"`  // 注文リスト
}

func (r *orderListResponse) UnmarshalJSON(b []byte) error {
	// 注文一覧が返されない場合は空文字が返されるので、空文字なら空配列に置き換えてからパースする
	replaced := bytes.Replace(b, []byte(`"55":""`), []byte(`"55":[]`), -1)

	type alias orderListResponse
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *orderListResponse) response() OrderListResponse {
	orders := make([]Order, len(r.Orders))
	for i, o := range r.Orders {
		orders[i] = o.response()
	}

	return OrderListResponse{
		CommonResponse:     r.commonResponse.response(),
		SymbolCode:         r.SymbolCode,
		ExecutionDate:      r.ExecutionDate.Time,
		OrderInquiryStatus: r.OrderInquiryStatus,
		ResultCode:         r.ResultCode,
		ResultText:         r.ResultText,
		WarningCode:        r.WarningCode,
		WarningText:        r.WarningText,
		Orders:             orders,
	}
}

type order struct {
	WarningCode            string            `json:"521"`        // 警告コード
	WarningText            string            `json:"522"`        // 警告テキスト
	OrderNumber            string            `json:"493"`        // 注文番号
	SymbolCode             string            `json:"485"`        // 銘柄コード
	Exchange               Exchange          `json:"501"`        // 市場
	AccountType            AccountType       `json:"528"`        // 譲渡益課税区分
	TradeType              TradeType         `json:"255"`        // 現金信用区分
	ExitTermType           ExitTermType      `json:"468"`        // 弁済区分
	Side                   Side              `json:"467"`        // 売買区分
	OrderQuantity          float64           `json:"496,string"` // 注文株数
	CurrentQuantity        float64           `json:"471,string"` // 有効株数
	Price                  float64           `json:"494,string"` // 注文単価
	ExecutionTiming        ExecutionTiming   `json:"469"`        // 執行条件
	ExecutionType          ExecutionType     `json:"495"`        // 注文値段区分
	StopOrderType          StopOrderType     `json:"480"`        // 逆指値注文種別
	StopTriggerPrice       float64           `json:"482,string"` // 逆指値条件
	StopOrderExecutionType ExecutionType     `json:"479"`        // 逆指値値段区分
	StopOrderPrice         float64           `json:"481,string"` // 逆指値値段
	TriggerType            TriggerType       `json:"517"`        // トリガータイプ
	ExitOrderType          ExitOrderType     `json:"510"`        // 建日種類
	ContractQuantity       float64           `json:"526,string"` // 成立株数
	ContractPrice          float64           `json:"524,string"` // 成立単価
	PartContractType       PartContractType  `json:"520"`        // 内出来区分
	ExecutionDate          Ymd               `json:"500"`        // 執行日
	OrderStatus            OrderStatus       `json:"504"`        // 状態コード
	OrderStatusText        string            `json:"503"`        // 状態
	ContractStatus         ContractStatus    `json:"525"`        // 約定ステータス
	OrderDateTime          YmdHms            `json:"491"`        // 注文日付
	ExpireDate             Ymd               `json:"492"`        // 有効期限
	CarryOverType          CarryOverType     `json:"489"`        // 繰越注文フラグ
	CorrectCancelType      CorrectCancelType `json:"470"`        // 訂正取消可否フラグ
	EstimationAmount       float64           `json:"235,string"` // 概算代金
}

func (r *order) response() Order {
	return Order{
		WarningCode:            r.WarningCode,
		WarningText:            r.WarningText,
		OrderNumber:            r.OrderNumber,
		SymbolCode:             r.SymbolCode,
		Exchange:               r.Exchange,
		AccountType:            r.AccountType,
		TradeType:              r.TradeType,
		ExitTermType:           r.ExitTermType,
		Side:                   r.Side,
		OrderQuantity:          r.OrderQuantity,
		CurrentQuantity:        r.CurrentQuantity,
		Price:                  r.Price,
		ExecutionTiming:        r.ExecutionTiming,
		ExecutionType:          r.ExecutionType,
		StopOrderType:          r.StopOrderType,
		StopTriggerPrice:       r.StopTriggerPrice,
		StopOrderExecutionType: r.StopOrderExecutionType,
		StopOrderPrice:         r.StopOrderPrice,
		TriggerType:            r.TriggerType,
		ExitOrderType:          r.ExitOrderType,
		ContractQuantity:       r.ContractQuantity,
		ContractPrice:          r.ContractPrice,
		PartContractType:       r.PartContractType,
		ExecutionDate:          r.ExecutionDate.Time,
		OrderStatus:            r.OrderStatus,
		OrderStatusText:        r.OrderStatusText,
		ContractStatus:         r.ContractStatus,
		OrderDateTime:          r.OrderDateTime.Time,
		ExpireDate:             r.ExpireDate.Time,
		CarryOverType:          r.CarryOverType,
		CorrectCancelType:      r.CorrectCancelType,
		EstimationAmount:       r.EstimationAmount,
	}
}

// OrderListResponse - 注文一覧レスポンス
type OrderListResponse struct {
	CommonResponse
	SymbolCode         string             // 銘柄コード
	ExecutionDate      time.Time          // 注文執行予定日
	OrderInquiryStatus OrderInquiryStatus // 注文照会状態
	ResultCode         string             // 結果コード
	ResultText         string             // 結果テキスト
	WarningCode        string             // 警告コード
	WarningText        string             // 警告テキスト
	Orders             []Order            // 注文リスト
}

// Order - 注文
type Order struct {
	WarningCode            string            // 警告コード
	WarningText            string            // 警告テキスト
	OrderNumber            string            // 注文番号
	SymbolCode             string            // 銘柄コード
	Exchange               Exchange          // 市場
	AccountType            AccountType       // 譲渡益課税区分
	TradeType              TradeType         // 現金信用区分
	ExitTermType           ExitTermType      // 弁済区分
	Side                   Side              // 弁済区分
	OrderQuantity          float64           // 注文株数
	CurrentQuantity        float64           // 有効株数
	Price                  float64           // 注文単価
	ExecutionTiming        ExecutionTiming   // 執行条件
	ExecutionType          ExecutionType     // 注文値段区分
	StopOrderType          StopOrderType     // 逆指値注文種別
	StopTriggerPrice       float64           // 逆指値条件
	StopOrderExecutionType ExecutionType     // 逆指値値段区分
	StopOrderPrice         float64           // 逆指値値段
	TriggerType            TriggerType       // トリガータイプ
	ExitOrderType          ExitOrderType     // 建日種類
	ContractQuantity       float64           // 成立株数
	ContractPrice          float64           // 成立単価
	PartContractType       PartContractType  // 内出来区分
	ExecutionDate          time.Time         // 執行日
	OrderStatus            OrderStatus       // 状態コード
	OrderStatusText        string            // 状態
	ContractStatus         ContractStatus    // 約定ステータス
	OrderDateTime          time.Time         // 注文日付
	ExpireDate             time.Time         // 有効期限
	CarryOverType          CarryOverType     // 繰越注文フラグ
	CorrectCancelType      CorrectCancelType // 訂正取消可否フラグ
	EstimationAmount       float64           // 概算代金
}

// OrderList - 注文一覧
func (c *client) OrderList(ctx context.Context, session *Session, req OrderListRequest) (*OrderListResponse, error) {
	if session == nil {
		return nil, NilArgumentErr
	}
	session.mtx.Lock()
	defer session.mtx.Unlock()

	session.lastRequestNo++
	r := req.request(session.lastRequestNo, c.clock.Now())

	var res orderListResponse
	if err := c.get(ctx, session.RequestURL, r, &res); err != nil {
		return nil, err
	}

	Res := res.response()
	return &Res, nil
}
