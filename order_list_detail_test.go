package tachibana

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func Test_orderListDetail_request(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request OrderListDetailRequest
		arg1    int64
		arg2    time.Time
		want1   orderListDetailRequest
	}{
		{name: "変換できる",
			request: OrderListDetailRequest{
				OrderNumber:   "28002795",
				ExecutionDate: time.Time{},
			},
			arg1: 123,
			arg2: time.Date(2022, 2, 27, 10, 21, 15, 0, time.Local),
			want1: orderListDetailRequest{
				commonRequest: commonRequest{
					No:             123,
					SendDate:       RequestTime{Time: time.Date(2022, 2, 27, 10, 21, 15, 0, time.Local)},
					FeatureType:    FeatureTypeOrderListDetail,
					ResponseFormat: commonResponseFormat,
				},
				OrderNumber:   "28002795",
				ExecutionDate: Ymd{},
			}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got1 := test.request.request(test.arg1, test.arg2)
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}

func Test_orderListDetailResponse_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		arg1         []byte
		wantResponse orderListDetailResponse
		hasError     bool
	}{
		{name: "正常系レスポンスのパース",
			arg1: []byte(`{
	"p_sd_date":"2022.03.01-14:45:58.421",
	"p_no":"2",
	"p_rv_date":"2022.03.01-14:45:58.284",
	"p_errno":"0",
	"p_err":"",
	"sCLMID":"CLMOrderListDetail",
	"sOrderNumber":"28010833",
	"sEigyouDay":"20220228",
	"sResultCode":"0",
	"sResultText":"",
	"sWarningCode":"0",
	"sWarningText":"",
	"sIssueCode":"1475",
	"sOrderSizyouC":"00",
	"sOrderBaibaiKubun":"1",
	"sGenkinSinyouKubun":"4",
	"sOrderBensaiKubun":"26",
	"sOrderCondition":"6",
	"sOrderOrderPriceKubun":"2",
	"sOrderOrderPrice":"1913.0000",
	"sOrderOrderSuryou":"1",
	"sOrderCurrentSuryou":"0",
	"sOrderStatusCode":"10",
	"sOrderStatus":"全部約定",
	"sOrderOrderDateTime":"20220228111323",
	"sOrderOrderExpireDay":"00000000",
	"sChannel":"1",
	"sGenbutuZyoutoekiKazeiC":"1",
	"sSinyouZyoutoekiKazeiC":"1",
	"sGyakusasiOrderType":"0",
	"sGyakusasiZyouken":"0.0000",
	"sGyakusasiKubun":" ",
	"sGyakusasiPrice":"0.0000",
	"sTriggerType":"0",
	"sTriggerTime":"00000000000000",
	"sUkewatasiDay":"20220302",
	"sYakuzyouPrice":"1912.0000",
	"sYakuzyouSuryou":"1",
	"sBaiBaiDaikin":"1912",
	"sUtidekiKubun":" ",
	"sGaisanDaikin":"3",
	"sBaiBaiTesuryo":"0",
	"sShouhizei":"0",
	"sTatebiType":"1",
	"sSizyouErrorCode":"",
	"sZougen":"",
	"sOrderAcceptTime":"20220228111323",
	"aYakuzyouSikkouList":
	[
	{
		"sYakuzyouWarningCode":"0",
		"sYakuzyouWarningText":"",
		"sYakuzyouSuryou":"1",
		"sYakuzyouPrice":"1912.0000",
		"sYakuzyouDate":"20220228113000"
	}
	],
	"aKessaiOrderTategyokuList":
	[
	{
		"sKessaiWarningCode":"0",
		"sKessaiWarningText":"",
		"sKessaiTatebiZyuni":"1",
		"sKessaiTategyokuDay":"20220228",
		"sKessaiTategyokuPrice":"1909.0000",
		"sKessaiOrderSuryo":"1",
		"sKessaiYakuzyouSuryo":"1",
		"sKessaiYakuzyouPrice":"1912.0000",
		"sKessaiTateTesuryou":"0",
		"sKessaiZyunHibu":"0",
		"sKessaiGyakuhibu":"0",
		"sKessaiKakikaeryou":"0",
		"sKessaiKanrihi":"0",
		"sKessaiKasikaburyou":"0",
		"sKessaiSonota":"0",
		"sKessaiSoneki":"3"
	}
	]
}`),
			wantResponse: orderListDetailResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 3, 1, 14, 45, 58, 421000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 3, 1, 14, 45, 58, 284000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderListDetail,
				},
				OrderNumber:            "28010833",
				ExecutionDate:          Ymd{Time: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local)},
				ResultCode:             "0",
				ResultText:             "",
				WarningCode:            "0",
				WarningText:            "",
				IssueCode:              "1475",
				Exchange:               ExchangeToushou,
				Side:                   SideSell,
				TradeType:              TradeTypeSystemExit,
				ExitTermType:           ExitTermTypeSystemMargin6m,
				ExecutionTiming:        ExecutionTimingFunari,
				ExecutionType:          ExecutionTypeLimit,
				Price:                  1913,
				OrderQuantity:          1,
				CurrentQuantity:        0,
				OrderStatus:            OrderStatusDone,
				OrderStatusText:        "全部約定",
				OrderDateTime:          YmdHms{Time: time.Date(2022, 2, 28, 11, 13, 23, 0, time.Local)},
				ExpireDate:             Ymd{},
				Channel:                ChannelPC,
				StockAccountType:       AccountTypeSpecific,
				MarginAccountType:      AccountTypeSpecific,
				StopOrderType:          StopOrderTypeNormal,
				StopTriggerPrice:       0,
				StopOrderExecutionType: ExecutionTypeUnused,
				StopOrderPrice:         0,
				TriggerType:            TriggerTypeNoFired,
				TriggerDateTime:        YmdHms{},
				DeliveryDate:           Ymd{Time: time.Date(2022, 3, 2, 0, 0, 0, 0, time.Local)},
				ContractPrice:          1912,
				ContractQuantity:       1,
				TradingAmount:          1912,
				PartContractType:       PartContractTypeUnused,
				EstimationAmount:       3,
				Commission:             0,
				CommissionTax:          0,
				ExitOrderType:          ExitOrderTypePositionNumber,
				ExchangeErrorCode:      "",
				ExchangeOrderDateTime:  YmdHms{Time: time.Date(2022, 2, 28, 11, 13, 23, 0, time.Local)},
				Contracts: []contract{{
					WarningCode: "0",
					WarningText: "",
					Quantity:    1,
					Price:       1912,
					DateTime:    YmdHms{Time: time.Date(2022, 2, 28, 11, 30, 0, 0, time.Local)},
				}},
				HoldPositions: []holdPosition{{
					WarningCode:   "0",
					WarningText:   "",
					SortOrder:     1,
					ContractDate:  Ymd{Time: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local)},
					EntryPrice:    1909,
					HoldQuantity:  1,
					ExitQuantity:  1,
					ExitPrice:     1912,
					Commission:    0,
					Interest:      0,
					Premiums:      0,
					RewritingFee:  0,
					ManagementFee: 0,
					LendingFee:    0,
					OtherFee:      0,
					Profit:        3,
				}},
			},
			hasError: false},
		{name: "エラーレスポンスのパース",
			arg1: []byte(`{
	"p_sd_date":"2022.03.01-14:46:11.276",
	"p_no":"2",
	"p_rv_date":"2022.03.01-14:46:11.251",
	"p_errno":"0",
	"p_err":"",
	"sCLMID":"CLMOrderListDetail",
	"sOrderNumber":"",
	"sEigyouDay":"",
	"sResultCode":"991002",
	"sResultText":"只今、一時的にこの業務はご利用できません。",
	"sWarningCode":"",
	"sWarningText":"",
	"sIssueCode":"",
	"sOrderSizyouC":"",
	"sOrderBaibaiKubun":"",
	"sGenkinSinyouKubun":"",
	"sOrderBensaiKubun":"",
	"sOrderCondition":"",
	"sOrderOrderPriceKubun":"",
	"sOrderOrderPrice":"",
	"sOrderOrderSuryou":"",
	"sOrderCurrentSuryou":"",
	"sOrderStatusCode":"",
	"sOrderStatus":"",
	"sOrderOrderDateTime":"",
	"sOrderOrderExpireDay":"",
	"sChannel":"",
	"sGenbutuZyoutoekiKazeiC":"",
	"sSinyouZyoutoekiKazeiC":"",
	"sGyakusasiOrderType":"",
	"sGyakusasiZyouken":"",
	"sGyakusasiKubun":"",
	"sGyakusasiPrice":"",
	"sTriggerType":"",
	"sTriggerTime":"",
	"sUkewatasiDay":"",
	"sYakuzyouPrice":"",
	"sYakuzyouSuryou":"",
	"sBaiBaiDaikin":"",
	"sUtidekiKubun":"",
	"sGaisanDaikin":"",
	"sBaiBaiTesuryo":"",
	"sShouhizei":"",
	"sTatebiType":"",
	"sSizyouErrorCode":"",
	"sZougen":"",
	"sOrderAcceptTime":"",
	"aYakuzyouSikkouList":"",
	"aKessaiOrderTategyokuList":""
}`),
			wantResponse: orderListDetailResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 3, 1, 14, 46, 11, 276000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 3, 1, 14, 46, 11, 251000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderListDetail,
				},
				OrderNumber:            "",
				ExecutionDate:          Ymd{},
				ResultCode:             "991002",
				ResultText:             "只今、一時的にこの業務はご利用できません。",
				WarningCode:            "",
				WarningText:            "",
				IssueCode:              "",
				Exchange:               "",
				Side:                   "",
				TradeType:              "",
				ExitTermType:           "",
				ExecutionTiming:        "",
				ExecutionType:          "",
				Price:                  0,
				OrderQuantity:          0,
				CurrentQuantity:        0,
				OrderStatus:            "",
				OrderStatusText:        "",
				OrderDateTime:          YmdHms{},
				ExpireDate:             Ymd{},
				Channel:                "",
				StockAccountType:       "",
				MarginAccountType:      "",
				StopOrderType:          "",
				StopTriggerPrice:       0,
				StopOrderExecutionType: "",
				StopOrderPrice:         0,
				TriggerType:            "",
				TriggerDateTime:        YmdHms{},
				DeliveryDate:           Ymd{},
				ContractPrice:          0,
				ContractQuantity:       0,
				TradingAmount:          0,
				PartContractType:       "",
				EstimationAmount:       0,
				Commission:             0,
				CommissionTax:          0,
				ExitOrderType:          "",
				ExchangeErrorCode:      "",
				ExchangeOrderDateTime:  YmdHms{},
				Contracts:              []contract{},
				HoldPositions:          []holdPosition{},
			},
			hasError: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var res orderListDetailResponse
			got1 := json.Unmarshal(test.arg1, &res)
			if !reflect.DeepEqual(test.wantResponse, res) || (got1 != nil) != test.hasError {
				t.Errorf("%s error\nwant: %+v\ngot: %+v, %+v\n", t.Name(), test.wantResponse, res, got1)
			}
		})
	}
}

func Test_orderListDetailResponse_response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		response orderListDetailResponse
		want1    OrderListDetailResponse
	}{
		{name: "変換できる",
			response: orderListDetailResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 3, 1, 8, 45, 50, 282000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 3, 1, 8, 45, 49, 557000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderListDetail,
				},
				OrderNumber:            "28010833",
				ExecutionDate:          Ymd{Time: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local)},
				ResultCode:             "0",
				ResultText:             "",
				WarningCode:            "0",
				WarningText:            "",
				IssueCode:              "1475",
				Exchange:               ExchangeToushou,
				Side:                   SideSell,
				TradeType:              TradeTypeSystemExit,
				ExitTermType:           ExitTermTypeSystemMargin6m,
				ExecutionTiming:        ExecutionTimingFunari,
				ExecutionType:          ExecutionTypeLimit,
				Price:                  1913,
				OrderQuantity:          1,
				CurrentQuantity:        0,
				OrderStatus:            OrderStatusDone,
				OrderStatusText:        "全部約定",
				OrderDateTime:          YmdHms{Time: time.Date(2022, 2, 28, 11, 13, 23, 0, time.Local)},
				ExpireDate:             Ymd{},
				Channel:                ChannelPC,
				StockAccountType:       AccountTypeSpecific,
				MarginAccountType:      AccountTypeSpecific,
				StopOrderType:          StopOrderTypeNormal,
				StopTriggerPrice:       0,
				StopOrderExecutionType: ExecutionTypeUnused,
				StopOrderPrice:         0,
				TriggerType:            TriggerTypeNoFired,
				TriggerDateTime:        YmdHms{},
				DeliveryDate:           Ymd{Time: time.Date(2022, 3, 2, 0, 0, 0, 0, time.Local)},
				ContractPrice:          1912,
				ContractQuantity:       1,
				TradingAmount:          1912,
				PartContractType:       PartContractTypeUnused,
				EstimationAmount:       3,
				Commission:             0,
				CommissionTax:          0,
				ExitOrderType:          ExitOrderTypePositionNumber,
				ExchangeErrorCode:      "",
				ExchangeOrderDateTime:  YmdHms{Time: time.Date(2022, 2, 28, 11, 13, 23, 0, time.Local)},
				Contracts: []contract{{
					WarningCode: "0",
					WarningText: "",
					Quantity:    1,
					Price:       1912,
					DateTime:    YmdHms{Time: time.Date(2022, 2, 28, 11, 30, 0, 0, time.Local)},
				}},
				HoldPositions: []holdPosition{{
					WarningCode:   "0",
					WarningText:   "",
					SortOrder:     1,
					ContractDate:  Ymd{Time: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local)},
					EntryPrice:    1909,
					HoldQuantity:  1,
					ExitQuantity:  1,
					ExitPrice:     1912,
					Commission:    0,
					Interest:      0,
					Premiums:      0,
					RewritingFee:  0,
					ManagementFee: 0,
					LendingFee:    0,
					OtherFee:      0,
					Profit:        3,
				}},
			},
			want1: OrderListDetailResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 1, 8, 45, 50, 282000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 1, 8, 45, 49, 557000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderListDetail,
				},
				OrderNumber:            "28010833",
				ExecutionDate:          time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local),
				ResultCode:             "0",
				ResultText:             "",
				WarningCode:            "0",
				WarningText:            "",
				IssueCode:              "1475",
				Exchange:               ExchangeToushou,
				Side:                   SideSell,
				TradeType:              TradeTypeSystemExit,
				ExitTermType:           ExitTermTypeSystemMargin6m,
				ExecutionTiming:        ExecutionTimingFunari,
				ExecutionType:          ExecutionTypeLimit,
				Price:                  1913,
				OrderQuantity:          1,
				CurrentQuantity:        0,
				OrderStatus:            OrderStatusDone,
				OrderStatusText:        "全部約定",
				OrderDateTime:          time.Date(2022, 2, 28, 11, 13, 23, 0, time.Local),
				ExpireDate:             time.Time{},
				Channel:                ChannelPC,
				StockAccountType:       AccountTypeSpecific,
				MarginAccountType:      AccountTypeSpecific,
				StopOrderType:          StopOrderTypeNormal,
				StopTriggerPrice:       0,
				StopOrderExecutionType: ExecutionTypeUnused,
				StopOrderPrice:         0,
				TriggerType:            TriggerTypeNoFired,
				TriggerDateTime:        time.Time{},
				DeliveryDate:           time.Date(2022, 3, 2, 0, 0, 0, 0, time.Local),
				ContractPrice:          1912,
				ContractQuantity:       1,
				TradingAmount:          1912,
				PartContractType:       PartContractTypeUnused,
				EstimationAmount:       3,
				Commission:             0,
				CommissionTax:          0,
				ExitOrderType:          ExitOrderTypePositionNumber,
				ExchangeErrorCode:      "",
				ExchangeOrderDateTime:  time.Date(2022, 2, 28, 11, 13, 23, 0, time.Local),
				Contracts: []Contract{{
					WarningCode: "0",
					WarningText: "",
					Quantity:    1,
					Price:       1912,
					DateTime:    time.Date(2022, 2, 28, 11, 30, 0, 0, time.Local),
				}},
				HoldPositions: []HoldPosition{{
					WarningCode:   "0",
					WarningText:   "",
					SortOrder:     1,
					ContractDate:  time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local),
					EntryPrice:    1909,
					HoldQuantity:  1,
					ExitQuantity:  1,
					ExitPrice:     1912,
					Commission:    0,
					Interest:      0,
					Premiums:      0,
					RewritingFee:  0,
					ManagementFee: 0,
					LendingFee:    0,
					OtherFee:      0,
					Profit:        3,
				}},
			}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got1 := test.response.response()
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}

func Test_contract_response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		response contract
		want1    Contract
	}{
		{name: "変換できる",
			response: contract{
				WarningCode: "0",
				WarningText: "",
				Quantity:    1,
				Price:       1912,
				DateTime:    YmdHms{Time: time.Date(2022, 2, 28, 11, 30, 0, 0, time.Local)},
			},
			want1: Contract{
				WarningCode: "0",
				WarningText: "",
				Quantity:    1,
				Price:       1912,
				DateTime:    time.Date(2022, 2, 28, 11, 30, 0, 0, time.Local),
			}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got1 := test.response.response()
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}

func Test_holdPosition_response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		response holdPosition
		want1    HoldPosition
	}{
		{name: "変換できる",
			response: holdPosition{
				WarningCode:   "0",
				WarningText:   "",
				SortOrder:     1,
				ContractDate:  Ymd{Time: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local)},
				EntryPrice:    1909,
				HoldQuantity:  1,
				ExitQuantity:  1,
				ExitPrice:     1912,
				Commission:    0,
				Interest:      0,
				Premiums:      0,
				RewritingFee:  0,
				ManagementFee: 0,
				LendingFee:    0,
				OtherFee:      0,
				Profit:        3,
			},
			want1: HoldPosition{
				WarningCode:   "0",
				WarningText:   "",
				SortOrder:     1,
				ContractDate:  time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local),
				EntryPrice:    1909,
				HoldQuantity:  1,
				ExitQuantity:  1,
				ExitPrice:     1912,
				Commission:    0,
				Interest:      0,
				Premiums:      0,
				RewritingFee:  0,
				ManagementFee: 0,
				LendingFee:    0,
				OtherFee:      0,
				Profit:        3,
			}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got1 := test.response.response()
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}

func Test_client_OrderListDetail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		clock  *testClock
		status int
		body   []byte
		arg1   context.Context
		arg2   *Session
		arg3   OrderListDetailRequest
		want1  *OrderListDetailResponse
		want2  error
	}{
		{name: "正常レスポンスをパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 10, 9, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 52, 53, 58, 53, 56, 46, 52, 50, 49, 34, 44, 10, 9, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 10, 9, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 52, 53, 58, 53, 56, 46, 50, 56, 52, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 10, 9, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 79, 114, 100, 101, 114, 76, 105, 115, 116, 68, 101, 116, 97, 105, 108, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 78, 117, 109, 98, 101, 114, 34, 58, 34, 50, 56, 48, 49, 48, 56, 51, 51, 34, 44, 10, 9, 34, 115, 69, 105, 103, 121, 111, 117, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 53, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 83, 105, 122, 121, 111, 117, 67, 34, 58, 34, 48, 48, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 66, 97, 105, 98, 97, 105, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 71, 101, 110, 107, 105, 110, 83, 105, 110, 121, 111, 117, 75, 117, 98, 117, 110, 34, 58, 34, 52, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 66, 101, 110, 115, 97, 105, 75, 117, 98, 117, 110, 34, 58, 34, 50, 54, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 67, 111, 110, 100, 105, 116, 105, 111, 110, 34, 58, 34, 54, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 79, 114, 100, 101, 114, 80, 114, 105, 99, 101, 75, 117, 98, 117, 110, 34, 58, 34, 50, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 79, 114, 100, 101, 114, 80, 114, 105, 99, 101, 34, 58, 34, 49, 57, 49, 51, 46, 48, 48, 48, 48, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 79, 114, 100, 101, 114, 83, 117, 114, 121, 111, 117, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 67, 117, 114, 114, 101, 110, 116, 83, 117, 114, 121, 111, 117, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 83, 116, 97, 116, 117, 115, 67, 111, 100, 101, 34, 58, 34, 49, 48, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 83, 116, 97, 116, 117, 115, 34, 58, 34, 145, 83, 149, 148, 150, 241, 146, 232, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 79, 114, 100, 101, 114, 68, 97, 116, 101, 84, 105, 109, 101, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 49, 49, 49, 51, 50, 51, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 79, 114, 100, 101, 114, 69, 120, 112, 105, 114, 101, 68, 97, 121, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 10, 9, 34, 115, 67, 104, 97, 110, 110, 101, 108, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 71, 101, 110, 98, 117, 116, 117, 90, 121, 111, 117, 116, 111, 101, 107, 105, 75, 97, 122, 101, 105, 67, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 83, 105, 110, 121, 111, 117, 90, 121, 111, 117, 116, 111, 101, 107, 105, 75, 97, 122, 101, 105, 67, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 71, 121, 97, 107, 117, 115, 97, 115, 105, 79, 114, 100, 101, 114, 84, 121, 112, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 71, 121, 97, 107, 117, 115, 97, 115, 105, 90, 121, 111, 117, 107, 101, 110, 34, 58, 34, 48, 46, 48, 48, 48, 48, 34, 44, 10, 9, 34, 115, 71, 121, 97, 107, 117, 115, 97, 115, 105, 75, 117, 98, 117, 110, 34, 58, 34, 32, 34, 44, 10, 9, 34, 115, 71, 121, 97, 107, 117, 115, 97, 115, 105, 80, 114, 105, 99, 101, 34, 58, 34, 48, 46, 48, 48, 48, 48, 34, 44, 10, 9, 34, 115, 84, 114, 105, 103, 103, 101, 114, 84, 121, 112, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 84, 114, 105, 103, 103, 101, 114, 84, 105, 109, 101, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 10, 9, 34, 115, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 48, 50, 34, 44, 10, 9, 34, 115, 89, 97, 107, 117, 122, 121, 111, 117, 80, 114, 105, 99, 101, 34, 58, 34, 49, 57, 49, 50, 46, 48, 48, 48, 48, 34, 44, 10, 9, 34, 115, 89, 97, 107, 117, 122, 121, 111, 117, 83, 117, 114, 121, 111, 117, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 66, 97, 105, 66, 97, 105, 68, 97, 105, 107, 105, 110, 34, 58, 34, 49, 57, 49, 50, 34, 44, 10, 9, 34, 115, 85, 116, 105, 100, 101, 107, 105, 75, 117, 98, 117, 110, 34, 58, 34, 32, 34, 44, 10, 9, 34, 115, 71, 97, 105, 115, 97, 110, 68, 97, 105, 107, 105, 110, 34, 58, 34, 51, 34, 44, 10, 9, 34, 115, 66, 97, 105, 66, 97, 105, 84, 101, 115, 117, 114, 121, 111, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 83, 104, 111, 117, 104, 105, 122, 101, 105, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 84, 97, 116, 101, 98, 105, 84, 121, 112, 101, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 83, 105, 122, 121, 111, 117, 69, 114, 114, 111, 114, 67, 111, 100, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 90, 111, 117, 103, 101, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 65, 99, 99, 101, 112, 116, 84, 105, 109, 101, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 49, 49, 49, 51, 50, 51, 34, 44, 10, 9, 34, 97, 89, 97, 107, 117, 122, 121, 111, 117, 83, 105, 107, 107, 111, 117, 76, 105, 115, 116, 34, 58, 10, 9, 91, 10, 9, 123, 10, 9, 9, 34, 115, 89, 97, 107, 117, 122, 121, 111, 117, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 89, 97, 107, 117, 122, 121, 111, 117, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 9, 34, 115, 89, 97, 107, 117, 122, 121, 111, 117, 83, 117, 114, 121, 111, 117, 34, 58, 34, 49, 34, 44, 10, 9, 9, 34, 115, 89, 97, 107, 117, 122, 121, 111, 117, 80, 114, 105, 99, 101, 34, 58, 34, 49, 57, 49, 50, 46, 48, 48, 48, 48, 34, 44, 10, 9, 9, 34, 115, 89, 97, 107, 117, 122, 121, 111, 117, 68, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 49, 49, 51, 48, 48, 48, 34, 10, 9, 125, 10, 9, 93, 44, 10, 9, 34, 97, 75, 101, 115, 115, 97, 105, 79, 114, 100, 101, 114, 84, 97, 116, 101, 103, 121, 111, 107, 117, 76, 105, 115, 116, 34, 58, 10, 9, 91, 10, 9, 123, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 84, 97, 116, 101, 98, 105, 90, 121, 117, 110, 105, 34, 58, 34, 49, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 84, 97, 116, 101, 103, 121, 111, 107, 117, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 84, 97, 116, 101, 103, 121, 111, 107, 117, 80, 114, 105, 99, 101, 34, 58, 34, 49, 57, 48, 57, 46, 48, 48, 48, 48, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 79, 114, 100, 101, 114, 83, 117, 114, 121, 111, 34, 58, 34, 49, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 89, 97, 107, 117, 122, 121, 111, 117, 83, 117, 114, 121, 111, 34, 58, 34, 49, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 89, 97, 107, 117, 122, 121, 111, 117, 80, 114, 105, 99, 101, 34, 58, 34, 49, 57, 49, 50, 46, 48, 48, 48, 48, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 84, 97, 116, 101, 84, 101, 115, 117, 114, 121, 111, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 90, 121, 117, 110, 72, 105, 98, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 71, 121, 97, 107, 117, 104, 105, 98, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 75, 97, 107, 105, 107, 97, 101, 114, 121, 111, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 75, 97, 110, 114, 105, 104, 105, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 75, 97, 115, 105, 107, 97, 98, 117, 114, 121, 111, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 83, 111, 110, 111, 116, 97, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 75, 101, 115, 115, 97, 105, 83, 111, 110, 101, 107, 105, 34, 58, 34, 51, 34, 10, 9, 125, 10, 9, 93, 10, 125, 10, 10},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListDetailRequest{},
			want1: &OrderListDetailResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 1, 14, 45, 58, 421000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 1, 14, 45, 58, 284000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderListDetail,
				},
				OrderNumber:            "28010833",
				ExecutionDate:          time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local),
				ResultCode:             "0",
				ResultText:             "",
				WarningCode:            "0",
				WarningText:            "",
				IssueCode:              "1475",
				Exchange:               ExchangeToushou,
				Side:                   SideSell,
				TradeType:              TradeTypeSystemExit,
				ExitTermType:           ExitTermTypeSystemMargin6m,
				ExecutionTiming:        ExecutionTimingFunari,
				ExecutionType:          ExecutionTypeLimit,
				Price:                  1913,
				OrderQuantity:          1,
				CurrentQuantity:        0,
				OrderStatus:            OrderStatusDone,
				OrderStatusText:        "全部約定",
				OrderDateTime:          time.Date(2022, 2, 28, 11, 13, 23, 0, time.Local),
				ExpireDate:             time.Time{},
				Channel:                ChannelPC,
				StockAccountType:       AccountTypeSpecific,
				MarginAccountType:      AccountTypeSpecific,
				StopOrderType:          StopOrderTypeNormal,
				StopTriggerPrice:       0,
				StopOrderExecutionType: ExecutionTypeUnused,
				StopOrderPrice:         0,
				TriggerType:            TriggerTypeNoFired,
				TriggerDateTime:        time.Time{},
				DeliveryDate:           time.Date(2022, 3, 2, 0, 0, 0, 0, time.Local),
				ContractPrice:          1912,
				ContractQuantity:       1,
				TradingAmount:          1912,
				PartContractType:       PartContractTypeUnused,
				EstimationAmount:       3,
				Commission:             0,
				CommissionTax:          0,
				ExitOrderType:          ExitOrderTypePositionNumber,
				ExchangeErrorCode:      "",
				ExchangeOrderDateTime:  time.Date(2022, 2, 28, 11, 13, 23, 0, time.Local),
				Contracts: []Contract{{
					WarningCode: "0",
					WarningText: "",
					Quantity:    1,
					Price:       1912,
					DateTime:    time.Date(2022, 2, 28, 11, 30, 0, 0, time.Local),
				}},
				HoldPositions: []HoldPosition{{
					WarningCode:   "0",
					WarningText:   "",
					SortOrder:     1,
					ContractDate:  time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local),
					EntryPrice:    1909,
					HoldQuantity:  1,
					ExitQuantity:  1,
					ExitPrice:     1912,
					Commission:    0,
					Interest:      0,
					Premiums:      0,
					RewritingFee:  0,
					ManagementFee: 0,
					LendingFee:    0,
					OtherFee:      0,
					Profit:        3,
				}},
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   nil,
			arg3:   OrderListDetailRequest{},
			want1:  nil,
			want2:  NilArgumentErr},
		{name: "利用できないとう情報をパースできる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 10, 9, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 52, 54, 58, 49, 49, 46, 50, 55, 54, 34, 44, 10, 9, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 10, 9, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 52, 54, 58, 49, 49, 46, 50, 53, 49, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 10, 9, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 79, 114, 100, 101, 114, 76, 105, 115, 116, 68, 101, 116, 97, 105, 108, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 78, 117, 109, 98, 101, 114, 34, 58, 34, 34, 44, 10, 9, 34, 115, 69, 105, 103, 121, 111, 117, 68, 97, 121, 34, 58, 34, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 57, 57, 49, 48, 48, 50, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 145, 252, 141, 161, 129, 65, 136, 234, 142, 158, 147, 73, 130, 201, 130, 177, 130, 204, 139, 198, 150, 177, 130, 205, 130, 178, 151, 152, 151, 112, 130, 197, 130, 171, 130, 220, 130, 185, 130, 241, 129, 66, 34, 44, 10, 9, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 83, 105, 122, 121, 111, 117, 67, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 66, 97, 105, 98, 97, 105, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 71, 101, 110, 107, 105, 110, 83, 105, 110, 121, 111, 117, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 66, 101, 110, 115, 97, 105, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 67, 111, 110, 100, 105, 116, 105, 111, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 79, 114, 100, 101, 114, 80, 114, 105, 99, 101, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 79, 114, 100, 101, 114, 80, 114, 105, 99, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 79, 114, 100, 101, 114, 83, 117, 114, 121, 111, 117, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 67, 117, 114, 114, 101, 110, 116, 83, 117, 114, 121, 111, 117, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 83, 116, 97, 116, 117, 115, 67, 111, 100, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 83, 116, 97, 116, 117, 115, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 79, 114, 100, 101, 114, 68, 97, 116, 101, 84, 105, 109, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 79, 114, 100, 101, 114, 69, 120, 112, 105, 114, 101, 68, 97, 121, 34, 58, 34, 34, 44, 10, 9, 34, 115, 67, 104, 97, 110, 110, 101, 108, 34, 58, 34, 34, 44, 10, 9, 34, 115, 71, 101, 110, 98, 117, 116, 117, 90, 121, 111, 117, 116, 111, 101, 107, 105, 75, 97, 122, 101, 105, 67, 34, 58, 34, 34, 44, 10, 9, 34, 115, 83, 105, 110, 121, 111, 117, 90, 121, 111, 117, 116, 111, 101, 107, 105, 75, 97, 122, 101, 105, 67, 34, 58, 34, 34, 44, 10, 9, 34, 115, 71, 121, 97, 107, 117, 115, 97, 115, 105, 79, 114, 100, 101, 114, 84, 121, 112, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 71, 121, 97, 107, 117, 115, 97, 115, 105, 90, 121, 111, 117, 107, 101, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 71, 121, 97, 107, 117, 115, 97, 115, 105, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 71, 121, 97, 107, 117, 115, 97, 115, 105, 80, 114, 105, 99, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 84, 114, 105, 103, 103, 101, 114, 84, 121, 112, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 84, 114, 105, 103, 103, 101, 114, 84, 105, 109, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 34, 44, 10, 9, 34, 115, 89, 97, 107, 117, 122, 121, 111, 117, 80, 114, 105, 99, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 89, 97, 107, 117, 122, 121, 111, 117, 83, 117, 114, 121, 111, 117, 34, 58, 34, 34, 44, 10, 9, 34, 115, 66, 97, 105, 66, 97, 105, 68, 97, 105, 107, 105, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 85, 116, 105, 100, 101, 107, 105, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 71, 97, 105, 115, 97, 110, 68, 97, 105, 107, 105, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 66, 97, 105, 66, 97, 105, 84, 101, 115, 117, 114, 121, 111, 34, 58, 34, 34, 44, 10, 9, 34, 115, 83, 104, 111, 117, 104, 105, 122, 101, 105, 34, 58, 34, 34, 44, 10, 9, 34, 115, 84, 97, 116, 101, 98, 105, 84, 121, 112, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 83, 105, 122, 121, 111, 117, 69, 114, 114, 111, 114, 67, 111, 100, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 90, 111, 117, 103, 101, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 79, 114, 100, 101, 114, 65, 99, 99, 101, 112, 116, 84, 105, 109, 101, 34, 58, 34, 34, 44, 10, 9, 34, 97, 89, 97, 107, 117, 122, 121, 111, 117, 83, 105, 107, 107, 111, 117, 76, 105, 115, 116, 34, 58, 34, 34, 44, 10, 9, 34, 97, 75, 101, 115, 115, 97, 105, 79, 114, 100, 101, 114, 84, 97, 116, 101, 103, 121, 111, 107, 117, 76, 105, 115, 116, 34, 58, 34, 34, 10, 125, 10, 10},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListDetailRequest{},
			want1: &OrderListDetailResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 1, 14, 46, 11, 276000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 1, 14, 46, 11, 251000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderListDetail,
				},
				ResultCode:    "991002",
				ResultText:    "只今、一時的にこの業務はご利用できません。",
				Contracts:     []Contract{},
				HoldPositions: []HoldPosition{},
			},
			want2: nil},
		{name: "200 OK以外が返ったらエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusInternalServerError,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListDetailRequest{},
			want1:  nil,
			want2:  StatusNotOkErr},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.status)
				_, _ = w.Write(test.body)
			})
			ts := httptest.NewServer(mux)
			defer ts.Close()
			if test.arg2 != nil {
				test.arg2.RequestURL = ts.URL
			}

			client := &client{clock: test.clock}
			got1, got2 := client.OrderListDetail(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_OrderListDetail_Execute(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client := NewClient(EnvironmentProduction, ApiVersionLatest)
	got1, got2 := client.Login(context.Background(), LoginRequest{
		UserId:   "user-id",
		Password: "password",
	})
	log.Printf("%+v, %+v\n", got1, got2)
	if got2 != nil {
		return
	}

	session, err := got1.Session()
	if err != nil {
		t.Errorf("session: %+v\n", err)
	}

	got3, got4 := client.OrderListDetail(context.Background(), session, OrderListDetailRequest{
		OrderNumber:   "28010833",
		ExecutionDate: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local),
	})
	log.Printf("%+v, %+v\n", got3, got4)
}
