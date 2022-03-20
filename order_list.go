package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// OrderListRequest - 注文一覧リクエスト
type OrderListRequest struct {
	IssueCode          string             // 銘柄コード
	ExecutionDate      time.Time          // 注文執行予定日
	OrderInquiryStatus OrderInquiryStatus // 注文照会状態
}

func (r *OrderListRequest) request(no int64, now time.Time) orderListRequest {
	return orderListRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			MessageType:    MessageTypeOrderList,
			ResponseFormat: commonResponseFormat,
		},
		IssueCode:     r.IssueCode,
		ExecutionDate: Ymd{Time: r.ExecutionDate},
		OrderStatus:   r.OrderInquiryStatus,
	}
}

type orderListRequest struct {
	commonRequest
	IssueCode     string             `json:"sIssueCode,omitempty"`          // 銘柄コード
	ExecutionDate Ymd                `json:"sSikkouDay,omitempty"`          // 注文執行予定日
	OrderStatus   OrderInquiryStatus `json:"sOrderSyoukaiStatus,omitempty"` // 注文照会状態
}

type orderListResponse struct {
	commonResponse
	IssueCode          string             `json:"sIssueCode"`          // 銘柄コード
	ExecutionDate      Ymd                `json:"sSikkouDay"`          // 注文執行予定日
	OrderInquiryStatus OrderInquiryStatus `json:"sOrderSyoukaiStatus"` // 注文照会状態
	ResultCode         string             `json:"sResultCode"`         // 結果コード
	ResultText         string             `json:"sResultText"`         // 結果テキスト
	WarningCode        string             `json:"sWarningCode"`        // 警告コード
	WarningText        string             `json:"sWarningText"`        // 警告テキスト
	Orders             []order            `json:"aOrderList"`          // 注文リスト
}

func (r *orderListResponse) UnmarshalJSON(b []byte) error {
	// 注文一覧が返されない場合は空文字が返されるので、空文字なら空配列に置き換えてからパースする
	replaced := bytes.Replace(b, []byte(`"aOrderList":""`), []byte(`"aOrderList":[]`), -1)

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
		IssueCode:          r.IssueCode,
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
	WarningCode            string            `json:"sOrderWarningCode"`             // 警告コード
	WarningText            string            `json:"sOrderWarningText"`             // 警告テキスト
	OrderNumber            string            `json:"sOrderOrderNumber"`             // 注文番号
	IssueCode              string            `json:"sOrderIssueCode"`               // 銘柄コード
	Exchange               Exchange          `json:"sOrderSizyouC"`                 // 市場
	AccountType            AccountType       `json:"sOrderZyoutoekiKazeiC"`         // 譲渡益課税区分
	TradeType              TradeType         `json:"sGenkinSinyouKubun"`            // 現金信用区分
	ExitTermType           ExitTermType      `json:"sOrderBensaiKubun"`             // 弁済区分
	Side                   Side              `json:"sOrderBaibaiKubun"`             // 売買区分
	OrderQuantity          float64           `json:"sOrderOrderSuryou,string"`      // 注文株数
	CurrentQuantity        float64           `json:"sOrderCurrentSuryou,string"`    // 有効株数
	Price                  float64           `json:"sOrderOrderPrice,string"`       // 注文単価
	ExecutionTiming        ExecutionTiming   `json:"sOrderCondition"`               // 執行条件
	ExecutionType          ExecutionType     `json:"sOrderOrderPriceKubun"`         // 注文値段区分
	StopOrderType          StopOrderType     `json:"sOrderGyakusasiOrderType"`      // 逆指値注文種別
	StopTriggerPrice       float64           `json:"sOrderGyakusasiZyouken,string"` // 逆指値条件
	StopOrderExecutionType ExecutionType     `json:"sOrderGyakusasiKubun"`          // 逆指値値段区分
	StopOrderPrice         float64           `json:"sOrderGyakusasiPrice,string"`   // 逆指値値段
	TriggerType            TriggerType       `json:"sOrderTriggerType"`             // トリガータイプ
	ExitOrderType          ExitOrderType     `json:"sOrderTatebiType"`              // 建日種類
	ContractQuantity       float64           `json:"sOrderYakuzyouSuryo,string"`    // 成立株数
	ContractPrice          float64           `json:"sOrderYakuzyouPrice,string"`    // 成立単価
	PartContractType       PartContractType  `json:"sOrderUtidekiKbn"`              // 内出来区分
	ExecutionDate          Ymd               `json:"sOrderSikkouDay"`               // 執行日
	OrderStatus            OrderStatus       `json:"sOrderStatusCode"`              // 状態コード
	OrderStatusText        string            `json:"sOrderStatus"`                  // 状態
	ContractStatus         ContractStatus    `json:"sOrderYakuzyouStatus"`          // 約定ステータス
	OrderDateTime          YmdHms            `json:"sOrderOrderDateTime"`           // 注文日付
	ExpireDate             Ymd               `json:"sOrderOrderExpireDay"`          // 有効期限
	CarryOverType          CarryOverType     `json:"sOrderKurikosiOrderFlg"`        // 繰越注文フラグ
	CorrectCancelType      CorrectCancelType `json:"sOrderCorrectCancelKahiFlg"`    // 訂正取消可否フラグ
	EstimationAmount       float64           `json:"sGaisanDaikin,string"`          // 概算代金
}

func (r *order) response() Order {
	return Order{
		WarningCode:            r.WarningCode,
		WarningText:            r.WarningText,
		OrderNumber:            r.OrderNumber,
		IssueCode:              r.IssueCode,
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
	IssueCode          string             // 銘柄コード
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
	IssueCode              string            // 銘柄コード
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

	b, err := c.requester.get(ctx, session.RequestURL, r)
	if err != nil {
		return nil, err
	}
	var res orderListResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, fmt.Errorf("%s: %w", err, UnmarshalFailedErr)
	}

	Res := res.response()
	return &Res, nil
}
