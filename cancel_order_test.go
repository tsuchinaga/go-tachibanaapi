package tachibana

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"
	"time"
)

func Test_client_CancelOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		clock     *testClock
		requester *testRequester
		arg1      context.Context
		arg2      *Session
		arg3      CancelOrderRequest
		want1     *CancelOrderResponse
		want2     error
	}{
		{name: "注文取消のレスポンスをパース出来る",
			clock: &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{
				get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 48, 45, 48, 56, 58, 52, 51, 58, 53, 56, 46, 49, 55, 57, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 51, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 48, 45, 48, 56, 58, 52, 51, 58, 53, 56, 46, 49, 50, 52, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 75, 97, 98, 117, 67, 97, 110, 99, 101, 108, 79, 114, 100, 101, 114, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 78, 117, 109, 98, 101, 114, 34, 58, 34, 49, 48, 48, 48, 52, 50, 51, 55, 34, 44, 34, 115, 69, 105, 103, 121, 111, 117, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 48, 34, 44, 34, 115, 79, 114, 100, 101, 114, 85, 107, 101, 119, 97, 116, 97, 115, 105, 75, 105, 110, 103, 97, 107, 117, 34, 58, 34, 48, 34, 44, 34, 115, 79, 114, 100, 101, 114, 68, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 48, 48, 56, 52, 51, 53, 56, 34, 125, 10},
			},
			arg1: context.Background(),
			arg2: &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3: CancelOrderRequest{},
			want1: &CancelOrderResponse{
				CommonResponse: CommonResponse{
					No:           3,
					SendDate:     time.Date(2022, 3, 10, 8, 43, 58, 179000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 10, 8, 43, 58, 124000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeCancelOrder,
				},
				ResultCode:     "0",
				ResultText:     "",
				OrderNumber:    "10004237",
				ExecutionDate:  time.Date(2022, 3, 10, 0, 0, 0, 0, time.Local),
				DeliveryAmount: 0,
				OrderDateTime:  time.Date(2022, 3, 10, 8, 43, 58, 0, time.Local),
			},
			want2: nil},
		{name: "注文取消失敗のレスポンスをパース出来る",
			clock: &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{
				get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 48, 45, 49, 50, 58, 51, 51, 58, 48, 53, 46, 49, 51, 56, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 48, 45, 49, 50, 58, 51, 51, 58, 48, 53, 46, 49, 48, 56, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 75, 97, 98, 117, 67, 97, 110, 99, 101, 108, 79, 114, 100, 101, 114, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 49, 51, 48, 48, 49, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 230, 179, 168, 230, 150, 135, 231, 149, 170, 229, 143, 183, 227, 129, 171, 232, 170, 164, 227, 130, 138, 227, 129, 140, 227, 129, 130, 227, 130, 138, 227, 129, 190, 227, 129, 153, 34, 44, 34, 115, 79, 114, 100, 101, 114, 78, 117, 109, 98, 101, 114, 34, 58, 34, 34, 44, 34, 115, 69, 105, 103, 121, 111, 117, 68, 97, 121, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 85, 107, 101, 119, 97, 116, 97, 115, 105, 75, 105, 110, 103, 97, 107, 117, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 68, 97, 116, 101, 34, 58, 34, 34, 125, 10},
			},
			arg1: context.Background(),
			arg2: &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3: CancelOrderRequest{},
			want1: &CancelOrderResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 10, 12, 33, 5, 138000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 10, 12, 33, 5, 108000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeCancelOrder,
				},
				ResultCode:     "13001",
				ResultText:     "注文番号に誤りがあります",
				OrderNumber:    "",
				ExecutionDate:  time.Time{},
				DeliveryAmount: 0,
				OrderDateTime:  time.Time{},
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			arg1:  context.Background(),
			arg2:  nil,
			arg3:  CancelOrderRequest{},
			want1: nil,
			want2: NilArgumentErr},
		{name: "リクエストでエラーが返されたらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get2: StatusNotOkErr},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      CancelOrderRequest{},
			want1:     nil,
			want2:     StatusNotOkErr},
		{name: "レスポンスのUnmarshalに失敗したらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get1: []byte{}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      CancelOrderRequest{},
			want1:     nil,
			want2:     UnmarshalFailedErr},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &client{clock: test.clock, requester: test.requester}
			got1, got2 := client.CancelOrder(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_CancelOrder_Execute(t *testing.T) {
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
	if got1.ResultCode != "0" {
		return
	}

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
		OrderPrice:          2000,
		OrderQuantity:       1,
		TradeType:           TradeTypeStock,
		ExpireDate:          time.Time{},
		ExpireDateIsToday:   true,
		StopOrderType:       StopOrderTypeNormal,
		TriggerPrice:        0,
		StopOrderPrice:      0,
		ExitPositionType:    ExitPositionTypeUnused,
		SecondPassword:      secondPassword,
		ExitPositions:       []ExitPosition{},
	})
	log.Printf("%+v, %+v\n", got3, got4)
	if got3.ResultCode != "0" {
		return
	}

	got5, got6 := client.CancelOrder(context.Background(), session, CancelOrderRequest{
		OrderNumber:    got3.OrderNumber,
		ExecutionDate:  got3.ExecutionDate,
		SecondPassword: secondPassword,
	})
	log.Printf("%+v, %+v\n", got5, got6)
	if got5.ResultCode != "0" {
		return
	}
}
