package tachibana

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"
	"time"
)

func Test_NewOrderRequest_request(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request NewOrderRequest
		arg1    int64
		arg2    time.Time
		want1   newOrderRequest
	}{
		{name: "現物成行を変換できる",
			request: NewOrderRequest{
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideBuy,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          0,
				OrderQuantity:       1,
				TradeType:           TradeTypeStock,
				ExpireDate:          time.Time{},
				ExpireDateIsToday:   true,
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      0,
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local),
			want1: newOrderRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local)},
					MessageType:    MessageTypeNewOrder,
					ResponseFormat: commonResponseFormat,
				},
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideBuy,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          "0",
				OrderQuantity:       1,
				TradeType:           TradeTypeStock,
				ExpireDate:          Ymd{isToday: true},
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      "*",
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			}},
		{name: "現物指値を変換できる",
			request: NewOrderRequest{
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideBuy,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          1700,
				OrderQuantity:       1,
				TradeType:           TradeTypeStock,
				ExpireDate:          time.Time{},
				ExpireDateIsToday:   true,
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      0,
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local),
			want1: newOrderRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local)},
					MessageType:    MessageTypeNewOrder,
					ResponseFormat: commonResponseFormat,
				},
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideBuy,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          "1700",
				OrderQuantity:       1,
				TradeType:           TradeTypeStock,
				ExpireDate:          Ymd{isToday: true},
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      "*",
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			}},
		{name: "現物逆指値を変換できる",
			request: NewOrderRequest{
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideBuy,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          0,
				OrderQuantity:       1,
				TradeType:           TradeTypeStock,
				ExpireDate:          time.Time{},
				ExpireDateIsToday:   true,
				StopOrderType:       StopOrderTypeStop,
				TriggerPrice:        1800,
				StopOrderPrice:      0,
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local),
			want1: newOrderRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local)},
					MessageType:    MessageTypeNewOrder,
					ResponseFormat: commonResponseFormat,
				},
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideBuy,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          "*",
				OrderQuantity:       1,
				TradeType:           TradeTypeStock,
				ExpireDate:          Ymd{isToday: true},
				StopOrderType:       StopOrderTypeStop,
				TriggerPrice:        1800,
				StopOrderPrice:      "0",
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			}},
		{name: "現物OCOを変換できる",
			request: NewOrderRequest{
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideBuy,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          1700,
				OrderQuantity:       1,
				TradeType:           TradeTypeStock,
				ExpireDate:          time.Time{},
				ExpireDateIsToday:   true,
				StopOrderType:       StopOrderTypeOCO,
				TriggerPrice:        1800,
				StopOrderPrice:      0,
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local),
			want1: newOrderRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local)},
					MessageType:    MessageTypeNewOrder,
					ResponseFormat: commonResponseFormat,
				},
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideBuy,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          "1700",
				OrderQuantity:       1,
				TradeType:           TradeTypeStock,
				ExpireDate:          Ymd{isToday: true},
				StopOrderType:       StopOrderTypeOCO,
				TriggerPrice:        1800,
				StopOrderPrice:      "0",
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			}},
		{name: "現物エグジットを変換できる",
			request: NewOrderRequest{
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideSell,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          0,
				OrderQuantity:       1,
				TradeType:           TradeTypeStock,
				ExpireDate:          time.Time{},
				ExpireDateIsToday:   true,
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      0,
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local),
			want1: newOrderRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local)},
					MessageType:    MessageTypeNewOrder,
					ResponseFormat: commonResponseFormat,
				},
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideSell,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          "0",
				OrderQuantity:       1,
				TradeType:           TradeTypeStock,
				ExpireDate:          Ymd{isToday: true},
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      "*",
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			}},
		{name: "信用買いエントリーを変換できる",
			request: NewOrderRequest{
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideBuy,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          0,
				OrderQuantity:       1,
				TradeType:           TradeTypeStandardEntry,
				ExpireDate:          time.Time{},
				ExpireDateIsToday:   true,
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      0,
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local),
			want1: newOrderRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local)},
					MessageType:    MessageTypeNewOrder,
					ResponseFormat: commonResponseFormat,
				},
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideBuy,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          "0",
				OrderQuantity:       1,
				TradeType:           TradeTypeStandardEntry,
				ExpireDate:          Ymd{isToday: true},
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      "*",
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			}},
		{name: "信用売りエグジットを変換できる",
			request: NewOrderRequest{
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideSell,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          0,
				OrderQuantity:       1,
				TradeType:           TradeTypeStandardExit,
				ExpireDate:          time.Time{},
				ExpireDateIsToday:   true,
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      0,
				ExitOrderType:       ExitOrderTypeDayAsc,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local),
			want1: newOrderRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local)},
					MessageType:    MessageTypeNewOrder,
					ResponseFormat: commonResponseFormat,
				},
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideSell,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          "0",
				OrderQuantity:       1,
				TradeType:           TradeTypeStandardExit,
				ExpireDate:          Ymd{isToday: true},
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      "*",
				ExitOrderType:       ExitOrderTypeDayAsc,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			}},
		{name: "信用売りエントリーを変換できる",
			request: NewOrderRequest{
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideSell,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          0,
				OrderQuantity:       1,
				TradeType:           TradeTypeStandardEntry,
				ExpireDate:          time.Time{},
				ExpireDateIsToday:   true,
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      0,
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local),
			want1: newOrderRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local)},
					MessageType:    MessageTypeNewOrder,
					ResponseFormat: commonResponseFormat,
				},
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideSell,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          "0",
				OrderQuantity:       1,
				TradeType:           TradeTypeStandardEntry,
				ExpireDate:          Ymd{isToday: true},
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      "*",
				ExitOrderType:       ExitOrderTypeUnused,
				SecondPassword:      "second-password",
				ExitPositions:       []ExitPosition{},
			}},
		{name: "信用買いエグジットを変換できる",
			request: NewOrderRequest{
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideBuy,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          0,
				OrderQuantity:       1,
				TradeType:           TradeTypeStandardExit,
				ExpireDate:          time.Time{},
				ExpireDateIsToday:   true,
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      0,
				ExitOrderType:       ExitOrderTypePositionNumber,
				SecondPassword:      "second-password",
				ExitPositions: []ExitPosition{
					{PositionNumber: "202203090000557", SequenceNumber: "1", OrderQuantity: "1"},
				},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local),
			want1: newOrderRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 9, 8, 17, 0, 0, time.Local)},
					MessageType:    MessageTypeNewOrder,
					ResponseFormat: commonResponseFormat,
				},
				AccountType:         AccountTypeSpecific,
				DeliveryAccountType: DeliveryAccountTypeUnused,
				IssueCode:           "1475",
				Exchange:            ExchangeToushou,
				Side:                SideBuy,
				ExecutionTiming:     ExecutionTimingNormal,
				OrderPrice:          "0",
				OrderQuantity:       1,
				TradeType:           TradeTypeStandardExit,
				ExpireDate:          Ymd{isToday: true},
				StopOrderType:       StopOrderTypeNormal,
				TriggerPrice:        0,
				StopOrderPrice:      "*",
				ExitOrderType:       ExitOrderTypePositionNumber,
				SecondPassword:      "second-password",
				ExitPositions: []ExitPosition{
					{PositionNumber: "202203090000557", SequenceNumber: "1", OrderQuantity: "1"},
				},
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

func Test_client_NewOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		clock     *testClock
		requester *testRequester
		arg1      context.Context
		arg2      *Session
		arg3      NewOrderRequest
		want1     *NewOrderResponse
		want2     error
	}{
		{name: "現物の成行注文レスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 56, 45, 48, 56, 58, 52, 57, 58, 48, 54, 46, 52, 57, 52, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 56, 45, 48, 56, 58, 52, 57, 58, 48, 54, 46, 52, 48, 51, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 75, 97, 98, 117, 78, 101, 119, 79, 114, 100, 101, 114, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 78, 117, 109, 98, 101, 114, 34, 58, 34, 56, 48, 48, 52, 51, 53, 56, 34, 44, 34, 115, 69, 105, 103, 121, 111, 117, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 48, 56, 34, 44, 34, 115, 79, 114, 100, 101, 114, 85, 107, 101, 119, 97, 116, 97, 115, 105, 75, 105, 110, 103, 97, 107, 117, 34, 58, 34, 50, 51, 48, 51, 34, 44, 34, 115, 79, 114, 100, 101, 114, 84, 101, 115, 117, 114, 121, 111, 117, 34, 58, 34, 55, 48, 34, 44, 34, 115, 79, 114, 100, 101, 114, 83, 121, 111, 117, 104, 105, 122, 101, 105, 34, 58, 34, 55, 34, 44, 34, 115, 75, 105, 110, 114, 105, 34, 58, 34, 45, 34, 44, 34, 115, 79, 114, 100, 101, 114, 68, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 51, 48, 56, 48, 56, 52, 57, 48, 54, 34, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      NewOrderRequest{},
			want1: &NewOrderResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 8, 8, 49, 6, 494000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 8, 8, 49, 6, 403000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeNewOrder,
				},
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				OrderNumber:    "8004358",
				ExecutionDate:  time.Date(2022, 3, 8, 0, 0, 0, 0, time.Local),
				DeliveryAmount: 2303,
				Commission:     70,
				CommissionTax:  7,
				Interest:       0,
				OrderDateTime:  time.Date(2022, 3, 8, 8, 49, 6, 0, time.Local),
			},
			want2: nil},
		{name: "現物の注文失敗レスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 56, 45, 48, 56, 58, 52, 57, 58, 53, 50, 46, 51, 54, 52, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 56, 45, 48, 56, 58, 52, 57, 58, 53, 50, 46, 51, 49, 55, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 75, 97, 98, 117, 78, 101, 119, 79, 114, 100, 101, 114, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 49, 49, 48, 48, 55, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 233, 138, 152, 230, 159, 132, 227, 130, 179, 227, 131, 188, 227, 131, 137, 227, 129, 171, 232, 170, 164, 227, 130, 138, 227, 129, 140, 227, 129, 130, 227, 130, 138, 227, 129, 190, 227, 129, 153, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 78, 117, 109, 98, 101, 114, 34, 58, 34, 34, 44, 34, 115, 69, 105, 103, 121, 111, 117, 68, 97, 121, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 85, 107, 101, 119, 97, 116, 97, 115, 105, 75, 105, 110, 103, 97, 107, 117, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 84, 101, 115, 117, 114, 121, 111, 117, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 83, 121, 111, 117, 104, 105, 122, 101, 105, 34, 58, 34, 34, 44, 34, 115, 75, 105, 110, 114, 105, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 68, 97, 116, 101, 34, 58, 34, 34, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      NewOrderRequest{},
			want1: &NewOrderResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 8, 8, 49, 52, 364000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 8, 8, 49, 52, 317000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeNewOrder,
				},
				ResultCode:     "11007",
				ResultText:     "銘柄コードに誤りがあります",
				WarningCode:    "",
				WarningText:    "",
				OrderNumber:    "",
				ExecutionDate:  time.Time{},
				DeliveryAmount: 0,
				Commission:     0,
				CommissionTax:  0,
				Interest:       0,
				OrderDateTime:  time.Time{},
			},
			want2: nil},
		{name: "信用の指値注文レスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 56, 45, 48, 56, 58, 53, 48, 58, 51, 56, 46, 57, 55, 50, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 56, 45, 48, 56, 58, 53, 48, 58, 51, 56, 46, 56, 56, 51, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 75, 97, 98, 117, 78, 101, 119, 79, 114, 100, 101, 114, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 78, 117, 109, 98, 101, 114, 34, 58, 34, 56, 48, 48, 52, 52, 48, 52, 34, 44, 34, 115, 69, 105, 103, 121, 111, 117, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 48, 56, 34, 44, 34, 115, 79, 114, 100, 101, 114, 85, 107, 101, 119, 97, 116, 97, 115, 105, 75, 105, 110, 103, 97, 107, 117, 34, 58, 34, 49, 55, 53, 48, 34, 44, 34, 115, 79, 114, 100, 101, 114, 84, 101, 115, 117, 114, 121, 111, 117, 34, 58, 34, 48, 34, 44, 34, 115, 79, 114, 100, 101, 114, 83, 121, 111, 117, 104, 105, 122, 101, 105, 34, 58, 34, 48, 34, 44, 34, 115, 75, 105, 110, 114, 105, 34, 58, 34, 49, 46, 54, 34, 44, 34, 115, 79, 114, 100, 101, 114, 68, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 51, 48, 56, 48, 56, 53, 48, 51, 56, 34, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      NewOrderRequest{},
			want1: &NewOrderResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 8, 8, 50, 38, 972000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 8, 8, 50, 38, 883000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeNewOrder,
				},
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				OrderNumber:    "8004404",
				ExecutionDate:  time.Date(2022, 3, 8, 0, 0, 0, 0, time.Local),
				DeliveryAmount: 1750,
				Commission:     0,
				CommissionTax:  0,
				Interest:       1.6,
				OrderDateTime:  time.Date(2022, 3, 8, 8, 50, 38, 0, time.Local),
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			arg1:  context.Background(),
			arg2:  nil,
			arg3:  NewOrderRequest{},
			want1: nil,
			want2: NilArgumentErr},
		{name: "リクエストでエラーが返されたらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get2: StatusNotOkErr},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      NewOrderRequest{},
			want1:     nil,
			want2:     StatusNotOkErr},
		{name: "レスポンスのUnmarshalに失敗したらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get1: []byte{}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      NewOrderRequest{},
			want1:     nil,
			want2:     UnmarshalFailedErr},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &client{clock: test.clock, requester: test.requester}
			got1, got2 := client.NewOrder(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_NewOrder_Execute_Stock_Entry(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	userId := "user-id"
	password := "password"
	secondPassword := "second-password"

	client := NewClient(EnvironmentProduction, ApiVersionLatest)
	got1, got2 := client.Login(context.Background(), LoginRequest{
		UserId:   userId,
		Password: password,
	})
	log.Printf("%+v, %+v\n", got1, got2)

	session, err := got1.Session()
	if err != nil {
		t.Errorf("session: %+v\n", err)
	}

	got3, got4 := client.NewOrder(context.Background(), session, NewOrderRequest{
		AccountType:         AccountTypeSpecific,
		DeliveryAccountType: DeliveryAccountTypeUnused,
		IssueCode:           "1475",
		Exchange:            ExchangeToushou,
		Side:                SideBuy,
		ExecutionTiming:     ExecutionTimingNormal,
		OrderPrice:          0,
		OrderQuantity:       1,
		TradeType:           TradeTypeStock,
		ExpireDate:          time.Time{},
		ExpireDateIsToday:   true,
		StopOrderType:       StopOrderTypeNormal,
		TriggerPrice:        0,
		StopOrderPrice:      0,
		ExitOrderType:       ExitOrderTypeUnused,
		SecondPassword:      secondPassword,
		ExitPositions:       []ExitPosition{},
	})
	log.Printf("%+v, %+v\n", got3, got4)
}

func Test_client_NewOrder_Execute_Stock_Exit(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	userId := "user-id"
	password := "password"
	secondPassword := "second-password"

	client := NewClient(EnvironmentProduction, ApiVersionLatest)
	got1, got2 := client.Login(context.Background(), LoginRequest{
		UserId:   userId,
		Password: password,
	})
	log.Printf("%+v, %+v\n", got1, got2)

	session, err := got1.Session()
	if err != nil {
		t.Errorf("session: %+v\n", err)
	}

	got3, got4 := client.NewOrder(context.Background(), session, NewOrderRequest{
		AccountType:         AccountTypeSpecific,
		DeliveryAccountType: DeliveryAccountTypeUnused,
		IssueCode:           "1475",
		Exchange:            ExchangeToushou,
		Side:                SideSell,
		ExecutionTiming:     ExecutionTimingNormal,
		OrderPrice:          0,
		OrderQuantity:       1,
		TradeType:           TradeTypeStock,
		ExpireDate:          time.Time{},
		ExpireDateIsToday:   true,
		StopOrderType:       StopOrderTypeNormal,
		TriggerPrice:        0,
		StopOrderPrice:      0,
		ExitOrderType:       ExitOrderTypeUnused,
		SecondPassword:      secondPassword,
		ExitPositions:       []ExitPosition{},
	})
	log.Printf("%+v, %+v\n", got3, got4)
}

func Test_client_NewOrder_Execute_Stop(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	userId := "user-id"
	password := "password"
	secondPassword := "second-password"

	client := NewClient(EnvironmentProduction, ApiVersionLatest)
	got1, got2 := client.Login(context.Background(), LoginRequest{
		UserId:   userId,
		Password: password,
	})
	log.Printf("%+v, %+v\n", got1, got2)

	session, err := got1.Session()
	if err != nil {
		t.Errorf("session: %+v\n", err)
	}

	got3, got4 := client.NewOrder(context.Background(), session, NewOrderRequest{
		AccountType:         AccountTypeSpecific,
		DeliveryAccountType: DeliveryAccountTypeUnused,
		IssueCode:           "1475",
		Exchange:            ExchangeToushou,
		Side:                SideBuy,
		ExecutionTiming:     ExecutionTimingNormal,
		OrderPrice:          0,
		OrderQuantity:       1,
		TradeType:           TradeTypeStock,
		ExpireDate:          time.Time{},
		ExpireDateIsToday:   true,
		StopOrderType:       StopOrderTypeStop,
		TriggerPrice:        1800,
		StopOrderPrice:      0,
		ExitOrderType:       ExitOrderTypeUnused,
		SecondPassword:      secondPassword,
		ExitPositions:       []ExitPosition{},
	})
	log.Printf("%+v, %+v\n", got3, got4)
}

func Test_client_NewOrder_Execute_Stock_OCO(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	userId := "user-id"
	password := "password"
	secondPassword := "second-password"

	client := NewClient(EnvironmentProduction, ApiVersionLatest)
	got1, got2 := client.Login(context.Background(), LoginRequest{
		UserId:   userId,
		Password: password,
	})
	log.Printf("%+v, %+v\n", got1, got2)

	session, err := got1.Session()
	if err != nil {
		t.Errorf("session: %+v\n", err)
	}

	got3, got4 := client.NewOrder(context.Background(), session, NewOrderRequest{
		AccountType:         AccountTypeSpecific,
		DeliveryAccountType: DeliveryAccountTypeUnused,
		IssueCode:           "1475",
		Exchange:            ExchangeToushou,
		Side:                SideBuy,
		ExecutionTiming:     ExecutionTimingNormal,
		OrderPrice:          1700,
		OrderQuantity:       1,
		TradeType:           TradeTypeStock,
		ExpireDate:          time.Time{},
		ExpireDateIsToday:   true,
		StopOrderType:       StopOrderTypeOCO,
		TriggerPrice:        1800,
		StopOrderPrice:      0,
		ExitOrderType:       ExitOrderTypeUnused,
		SecondPassword:      secondPassword,
		ExitPositions:       []ExitPosition{},
	})
	log.Printf("%+v, %+v\n", got3, got4)
}

func Test_client_NewOrder_Execute_Margin(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	userId := "user-id"
	password := "password"
	secondPassword := "second-password"

	client := NewClient(EnvironmentProduction, ApiVersionLatest)
	got1, got2 := client.Login(context.Background(), LoginRequest{
		UserId:   userId,
		Password: password,
	})
	log.Printf("%+v, %+v\n", got1, got2)

	session, err := got1.Session()
	if err != nil {
		t.Errorf("session: %+v\n", err)
	}

	got3, got4 := client.NewOrder(context.Background(), session, NewOrderRequest{
		AccountType:         AccountTypeSpecific,
		DeliveryAccountType: DeliveryAccountTypeUnused,
		IssueCode:           "1475",
		Exchange:            ExchangeToushou,
		Side:                SideBuy,
		ExecutionTiming:     ExecutionTimingNormal,
		OrderPrice:          0,
		OrderQuantity:       1,
		TradeType:           TradeTypeStandardEntry,
		ExpireDate:          time.Time{},
		ExpireDateIsToday:   true,
		StopOrderType:       StopOrderTypeNormal,
		TriggerPrice:        0,
		StopOrderPrice:      0,
		ExitOrderType:       ExitOrderTypeUnused,
		SecondPassword:      secondPassword,
		ExitPositions:       []ExitPosition{},
	})
	log.Printf("%+v, %+v\n", got3, got4)
}

func Test_client_NewOrder_Execute_Margin_Exit(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	userId := "user-id"
	password := "password"
	secondPassword := "second-password"

	client := NewClient(EnvironmentProduction, ApiVersionLatest)
	got1, got2 := client.Login(context.Background(), LoginRequest{
		UserId:   userId,
		Password: password,
	})
	log.Printf("%+v, %+v\n", got1, got2)

	session, err := got1.Session()
	if err != nil {
		t.Errorf("session: %+v\n", err)
	}

	got3, got4 := client.NewOrder(context.Background(), session, NewOrderRequest{
		AccountType:         AccountTypeSpecific,
		DeliveryAccountType: DeliveryAccountTypeUnused,
		IssueCode:           "1475",
		Exchange:            ExchangeToushou,
		Side:                SideBuy,
		ExecutionTiming:     ExecutionTimingClosing,
		OrderPrice:          0,
		OrderQuantity:       1,
		TradeType:           TradeTypeStandardExit,
		ExpireDate:          time.Time{},
		ExpireDateIsToday:   true,
		StopOrderType:       StopOrderTypeNormal,
		TriggerPrice:        0,
		StopOrderPrice:      0,
		ExitOrderType:       ExitOrderTypeDayAsc,
		SecondPassword:      secondPassword,
		ExitPositions:       []ExitPosition{},
	})
	log.Printf("%+v, %+v\n", got3, got4)
}
