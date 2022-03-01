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
					No:          123,
					SendDate:    RequestTime{Time: time.Date(2022, 2, 27, 10, 21, 15, 0, time.Local)},
					FeatureType: FeatureTypeOrderListDetail,
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
			arg1: []byte(`{"177":"2022.03.01-08:45:50.282","175":"2","176":"2022.03.01-08:45:49.557","174":"0","173":"","192":"CLMOrderListDetail","490":"28010833","227":"20220228","534":"0","535":"","692":"0","693":"","328":"1475","501":"00","467":"1","255":"4","468":"26","469":"6","495":"2","494":"1913.0000","496":"1","471":"0","504":"10","503":"全部約定","491":"20220228111323","492":"00000000","193":"1","248":"1","575":"1","259":"0","263":"0.0000","258":" ","260":"0.0000","659":"0","658":"00000000000000","662":"20220302","695":"1912.0000","696":"1","182":"1912","691":" ","235":"3","183":"0","558":"0","620":"1","577":"","741":"","466":"20220228111323","57":[{"697":"0","698":"","696":"1","695":"1912.0000","694":"20220228113000"}],"53":[{"365":"0","366":"","360":"1","361":"20220228","362":"1909.0000","356":"1","368":"1","367":"1912.0000","359":"0","369":"0","352":"0","353":"0","354":"0","355":"0","358":"0","357":"3"}]}`),
			wantResponse: orderListDetailResponse{
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
				SymbolCode:             "1475",
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
				ExitOrderType:          ExitOrderTypeSpecified,
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
			arg1: []byte(`{"177":"2022.02.28-16:08:21.894","175":"2","176":"2022.02.28-16:08:21.878","174":"0","173":"","192":"CLMOrderListDetail","490":"","227":"","534":"991012","535":"只今、一時的にこの業務はご利用できません。","692":"","693":"","328":"","501":"","467":"","255":"","468":"","469":"","495":"","494":"","496":"","471":"","504":"","503":"","491":"","492":"","193":"","248":"","575":"","259":"","263":"","258":"","260":"","659":"","658":"","662":"","695":"","696":"","182":"","691":"","235":"","183":"","558":"","620":"","577":"","741":"","466":"","57":"","53":""}`),
			wantResponse: orderListDetailResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 2, 28, 16, 8, 21, 894000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 2, 28, 16, 8, 21, 878000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderListDetail,
				},
				ResultCode:    "991012",
				ResultText:    "只今、一時的にこの業務はご利用できません。",
				Contracts:     []contract{},
				HoldPositions: []holdPosition{},
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
				SymbolCode:             "1475",
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
				ExitOrderType:          ExitOrderTypeSpecified,
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
				SymbolCode:             "1475",
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
				ExitOrderType:          ExitOrderTypeSpecified,
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
		{name: "正常レスポンスをパースして返せる", // TODO 正常系のレスポンスのテスト
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 48, 56, 58, 52, 53, 58, 53, 48, 46, 50, 56, 50, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 48, 56, 58, 52, 53, 58, 52, 57, 46, 53, 53, 55, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 79, 114, 100, 101, 114, 76, 105, 115, 116, 68, 101, 116, 97, 105, 108, 34, 44, 34, 52, 57, 48, 34, 58, 34, 50, 56, 48, 49, 48, 56, 51, 51, 34, 44, 34, 50, 50, 55, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 34, 44, 34, 53, 51, 52, 34, 58, 34, 48, 34, 44, 34, 53, 51, 53, 34, 58, 34, 34, 44, 34, 54, 57, 50, 34, 58, 34, 48, 34, 44, 34, 54, 57, 51, 34, 58, 34, 34, 44, 34, 51, 50, 56, 34, 58, 34, 49, 52, 55, 53, 34, 44, 34, 53, 48, 49, 34, 58, 34, 48, 48, 34, 44, 34, 52, 54, 55, 34, 58, 34, 49, 34, 44, 34, 50, 53, 53, 34, 58, 34, 52, 34, 44, 34, 52, 54, 56, 34, 58, 34, 50, 54, 34, 44, 34, 52, 54, 57, 34, 58, 34, 54, 34, 44, 34, 52, 57, 53, 34, 58, 34, 50, 34, 44, 34, 52, 57, 52, 34, 58, 34, 49, 57, 49, 51, 46, 48, 48, 48, 48, 34, 44, 34, 52, 57, 54, 34, 58, 34, 49, 34, 44, 34, 52, 55, 49, 34, 58, 34, 48, 34, 44, 34, 53, 48, 52, 34, 58, 34, 49, 48, 34, 44, 34, 53, 48, 51, 34, 58, 34, 145, 83, 149, 148, 150, 241, 146, 232, 34, 44, 34, 52, 57, 49, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 49, 49, 49, 51, 50, 51, 34, 44, 34, 52, 57, 50, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 49, 57, 51, 34, 58, 34, 49, 34, 44, 34, 50, 52, 56, 34, 58, 34, 49, 34, 44, 34, 53, 55, 53, 34, 58, 34, 49, 34, 44, 34, 50, 53, 57, 34, 58, 34, 48, 34, 44, 34, 50, 54, 51, 34, 58, 34, 48, 46, 48, 48, 48, 48, 34, 44, 34, 50, 53, 56, 34, 58, 34, 32, 34, 44, 34, 50, 54, 48, 34, 58, 34, 48, 46, 48, 48, 48, 48, 34, 44, 34, 54, 53, 57, 34, 58, 34, 48, 34, 44, 34, 54, 53, 56, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 54, 54, 50, 34, 58, 34, 50, 48, 50, 50, 48, 51, 48, 50, 34, 44, 34, 54, 57, 53, 34, 58, 34, 49, 57, 49, 50, 46, 48, 48, 48, 48, 34, 44, 34, 54, 57, 54, 34, 58, 34, 49, 34, 44, 34, 49, 56, 50, 34, 58, 34, 49, 57, 49, 50, 34, 44, 34, 54, 57, 49, 34, 58, 34, 32, 34, 44, 34, 50, 51, 53, 34, 58, 34, 51, 34, 44, 34, 49, 56, 51, 34, 58, 34, 48, 34, 44, 34, 53, 53, 56, 34, 58, 34, 48, 34, 44, 34, 54, 50, 48, 34, 58, 34, 49, 34, 44, 34, 53, 55, 55, 34, 58, 34, 34, 44, 34, 55, 52, 49, 34, 58, 34, 34, 44, 34, 52, 54, 54, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 49, 49, 49, 51, 50, 51, 34, 44, 34, 53, 55, 34, 58, 91, 123, 34, 54, 57, 55, 34, 58, 34, 48, 34, 44, 34, 54, 57, 56, 34, 58, 34, 34, 44, 34, 54, 57, 54, 34, 58, 34, 49, 34, 44, 34, 54, 57, 53, 34, 58, 34, 49, 57, 49, 50, 46, 48, 48, 48, 48, 34, 44, 34, 54, 57, 52, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 49, 49, 51, 48, 48, 48, 34, 125, 93, 44, 34, 53, 51, 34, 58, 91, 123, 34, 51, 54, 53, 34, 58, 34, 48, 34, 44, 34, 51, 54, 54, 34, 58, 34, 34, 44, 34, 51, 54, 48, 34, 58, 34, 49, 34, 44, 34, 51, 54, 49, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 34, 44, 34, 51, 54, 50, 34, 58, 34, 49, 57, 48, 57, 46, 48, 48, 48, 48, 34, 44, 34, 51, 53, 54, 34, 58, 34, 49, 34, 44, 34, 51, 54, 56, 34, 58, 34, 49, 34, 44, 34, 51, 54, 55, 34, 58, 34, 49, 57, 49, 50, 46, 48, 48, 48, 48, 34, 44, 34, 51, 53, 57, 34, 58, 34, 48, 34, 44, 34, 51, 54, 57, 34, 58, 34, 48, 34, 44, 34, 51, 53, 50, 34, 58, 34, 48, 34, 44, 34, 51, 53, 51, 34, 58, 34, 48, 34, 44, 34, 51, 53, 52, 34, 58, 34, 48, 34, 44, 34, 51, 53, 53, 34, 58, 34, 48, 34, 44, 34, 51, 53, 56, 34, 58, 34, 48, 34, 44, 34, 51, 53, 55, 34, 58, 34, 51, 34, 125, 93, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListDetailRequest{},
			want1: &OrderListDetailResponse{
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
				SymbolCode:             "1475",
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
				ExitOrderType:          ExitOrderTypeSpecified,
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
		{name: "失敗をパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 48, 57, 58, 51, 57, 58, 49, 53, 46, 56, 51, 55, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 48, 57, 58, 51, 57, 58, 49, 53, 46, 56, 48, 53, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 79, 114, 100, 101, 114, 76, 105, 115, 116, 68, 101, 116, 97, 105, 108, 34, 44, 34, 52, 57, 48, 34, 58, 34, 34, 44, 34, 50, 50, 55, 34, 58, 34, 34, 44, 34, 53, 51, 52, 34, 58, 34, 57, 57, 49, 48, 48, 49, 34, 44, 34, 53, 51, 53, 34, 58, 34, 146, 141, 149, 182, 148, 212, 141, 134, 130, 201, 140, 235, 130, 232, 130, 170, 130, 160, 130, 232, 130, 220, 130, 183, 129, 66, 34, 44, 34, 54, 57, 50, 34, 58, 34, 34, 44, 34, 54, 57, 51, 34, 58, 34, 34, 44, 34, 51, 50, 56, 34, 58, 34, 34, 44, 34, 53, 48, 49, 34, 58, 34, 34, 44, 34, 52, 54, 55, 34, 58, 34, 34, 44, 34, 50, 53, 53, 34, 58, 34, 34, 44, 34, 52, 54, 56, 34, 58, 34, 34, 44, 34, 52, 54, 57, 34, 58, 34, 34, 44, 34, 52, 57, 53, 34, 58, 34, 34, 44, 34, 52, 57, 52, 34, 58, 34, 34, 44, 34, 52, 57, 54, 34, 58, 34, 34, 44, 34, 52, 55, 49, 34, 58, 34, 34, 44, 34, 53, 48, 52, 34, 58, 34, 34, 44, 34, 53, 48, 51, 34, 58, 34, 34, 44, 34, 52, 57, 49, 34, 58, 34, 34, 44, 34, 52, 57, 50, 34, 58, 34, 34, 44, 34, 49, 57, 51, 34, 58, 34, 34, 44, 34, 50, 52, 56, 34, 58, 34, 34, 44, 34, 53, 55, 53, 34, 58, 34, 34, 44, 34, 50, 53, 57, 34, 58, 34, 34, 44, 34, 50, 54, 51, 34, 58, 34, 34, 44, 34, 50, 53, 56, 34, 58, 34, 34, 44, 34, 50, 54, 48, 34, 58, 34, 34, 44, 34, 54, 53, 57, 34, 58, 34, 34, 44, 34, 54, 53, 56, 34, 58, 34, 34, 44, 34, 54, 54, 50, 34, 58, 34, 34, 44, 34, 54, 57, 53, 34, 58, 34, 34, 44, 34, 54, 57, 54, 34, 58, 34, 34, 44, 34, 49, 56, 50, 34, 58, 34, 34, 44, 34, 54, 57, 49, 34, 58, 34, 34, 44, 34, 50, 51, 53, 34, 58, 34, 34, 44, 34, 49, 56, 51, 34, 58, 34, 34, 44, 34, 53, 53, 56, 34, 58, 34, 34, 44, 34, 54, 50, 48, 34, 58, 34, 34, 44, 34, 53, 55, 55, 34, 58, 34, 34, 44, 34, 55, 52, 49, 34, 58, 34, 34, 44, 34, 52, 54, 54, 34, 58, 34, 34, 44, 34, 53, 55, 34, 58, 34, 34, 44, 34, 53, 51, 34, 58, 34, 34, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListDetailRequest{},
			want1: &OrderListDetailResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 28, 9, 39, 15, 837000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 28, 9, 39, 15, 805000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderListDetail,
				},
				ResultCode:    "991001",
				ResultText:    "注文番号に誤りがあります。",
				Contracts:     []Contract{},
				HoldPositions: []HoldPosition{},
			},
			want2: nil},
		{name: "利用できないとう情報をパースできる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 48, 57, 58, 51, 57, 58, 51, 53, 46, 53, 49, 52, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 48, 57, 58, 51, 57, 58, 51, 53, 46, 52, 55, 50, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 79, 114, 100, 101, 114, 76, 105, 115, 116, 68, 101, 116, 97, 105, 108, 34, 44, 34, 52, 57, 48, 34, 58, 34, 34, 44, 34, 50, 50, 55, 34, 58, 34, 34, 44, 34, 53, 51, 52, 34, 58, 34, 57, 57, 49, 48, 48, 50, 34, 44, 34, 53, 51, 53, 34, 58, 34, 145, 252, 141, 161, 129, 65, 136, 234, 142, 158, 147, 73, 130, 201, 130, 177, 130, 204, 139, 198, 150, 177, 130, 205, 130, 178, 151, 152, 151, 112, 130, 197, 130, 171, 130, 220, 130, 185, 130, 241, 129, 66, 34, 44, 34, 54, 57, 50, 34, 58, 34, 34, 44, 34, 54, 57, 51, 34, 58, 34, 34, 44, 34, 51, 50, 56, 34, 58, 34, 34, 44, 34, 53, 48, 49, 34, 58, 34, 34, 44, 34, 52, 54, 55, 34, 58, 34, 34, 44, 34, 50, 53, 53, 34, 58, 34, 34, 44, 34, 52, 54, 56, 34, 58, 34, 34, 44, 34, 52, 54, 57, 34, 58, 34, 34, 44, 34, 52, 57, 53, 34, 58, 34, 34, 44, 34, 52, 57, 52, 34, 58, 34, 34, 44, 34, 52, 57, 54, 34, 58, 34, 34, 44, 34, 52, 55, 49, 34, 58, 34, 34, 44, 34, 53, 48, 52, 34, 58, 34, 34, 44, 34, 53, 48, 51, 34, 58, 34, 34, 44, 34, 52, 57, 49, 34, 58, 34, 34, 44, 34, 52, 57, 50, 34, 58, 34, 34, 44, 34, 49, 57, 51, 34, 58, 34, 34, 44, 34, 50, 52, 56, 34, 58, 34, 34, 44, 34, 53, 55, 53, 34, 58, 34, 34, 44, 34, 50, 53, 57, 34, 58, 34, 34, 44, 34, 50, 54, 51, 34, 58, 34, 34, 44, 34, 50, 53, 56, 34, 58, 34, 34, 44, 34, 50, 54, 48, 34, 58, 34, 34, 44, 34, 54, 53, 57, 34, 58, 34, 34, 44, 34, 54, 53, 56, 34, 58, 34, 34, 44, 34, 54, 54, 50, 34, 58, 34, 34, 44, 34, 54, 57, 53, 34, 58, 34, 34, 44, 34, 54, 57, 54, 34, 58, 34, 34, 44, 34, 49, 56, 50, 34, 58, 34, 34, 44, 34, 54, 57, 49, 34, 58, 34, 34, 44, 34, 50, 51, 53, 34, 58, 34, 34, 44, 34, 49, 56, 51, 34, 58, 34, 34, 44, 34, 53, 53, 56, 34, 58, 34, 34, 44, 34, 54, 50, 48, 34, 58, 34, 34, 44, 34, 53, 55, 55, 34, 58, 34, 34, 44, 34, 55, 52, 49, 34, 58, 34, 34, 44, 34, 52, 54, 54, 34, 58, 34, 34, 44, 34, 53, 55, 34, 58, 34, 34, 44, 34, 53, 51, 34, 58, 34, 34, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListDetailRequest{},
			want1: &OrderListDetailResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 28, 9, 39, 35, 514000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 28, 9, 39, 35, 472000000, time.Local),
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
