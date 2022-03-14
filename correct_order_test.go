package tachibana

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func Test_CorrectOrderRequest_request(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request CorrectOrderRequest
		arg1    int64
		arg2    time.Time
		want1   correctOrderRequest
	}{
		{name: "変更なしを指定した項目が変換できる",
			request: CorrectOrderRequest{
				OrderNumber:        "11002847",
				ExecutionDate:      time.Date(2022, 3, 11, 0, 0, 0, 0, time.Local),
				ExecutionTiming:    ExecutionTimingNoChange,
				OrderPrice:         NoChangeFloat,
				OrderQuantity:      NoChangeFloat,
				ExpireDate:         time.Time{},
				ExpireDateIsToday:  false,
				ExpireDateNoChange: true,
				TriggerPrice:       NoChangeFloat,
				StopOrderPrice:     NoChangeFloat,
				SecondPassword:     "second-password",
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 11, 7, 34, 0, 0, time.Local),
			want1: correctOrderRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 11, 7, 34, 0, 0, time.Local)},
					FeatureType:    FeatureTypeCorrectOrder,
					ResponseFormat: commonResponseFormat,
				},
				OrderNumber:     "11002847",
				ExecutionDate:   Ymd{Time: time.Date(2022, 3, 11, 0, 0, 0, 0, time.Local)},
				ExecutionTiming: "*",
				OrderPrice:      "*",
				OrderQuantity:   "*",
				ExpireDate:      Ymd{isNoChange: true},
				TriggerPrice:    "*",
				StopOrderPrice:  "*",
				SecondPassword:  "second-password",
			}},
		{name: "変更値を指定した項目が変換できる",
			request: CorrectOrderRequest{
				OrderNumber:        "11002847",
				ExecutionDate:      time.Date(2022, 3, 11, 0, 0, 0, 0, time.Local),
				ExecutionTiming:    ExecutionTimingFunari,
				OrderPrice:         2001,
				OrderQuantity:      2,
				ExpireDate:         time.Date(2022, 3, 12, 0, 0, 0, 0, time.Local),
				ExpireDateIsToday:  false,
				ExpireDateNoChange: false,
				TriggerPrice:       1999,
				StopOrderPrice:     1999,
				SecondPassword:     "second-password",
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 11, 7, 34, 0, 0, time.Local),
			want1: correctOrderRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 11, 7, 34, 0, 0, time.Local)},
					FeatureType:    FeatureTypeCorrectOrder,
					ResponseFormat: commonResponseFormat,
				},
				OrderNumber:     "11002847",
				ExecutionDate:   Ymd{Time: time.Date(2022, 3, 11, 0, 0, 0, 0, time.Local)},
				ExecutionTiming: ExecutionTimingFunari,
				OrderPrice:      "2001",
				OrderQuantity:   "2",
				ExpireDate:      Ymd{Time: time.Date(2022, 3, 12, 0, 0, 0, 0, time.Local)},
				TriggerPrice:    "1999",
				StopOrderPrice:  "1999",
				SecondPassword:  "second-password",
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

func Test_client_CorrectOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		clock  *testClock
		status int
		body   []byte
		arg1   context.Context
		arg2   *Session
		arg3   CorrectOrderRequest
		want1  *CorrectOrderResponse
		want2  error
	}{
		{name: "注文変更のレスポンスをパース出来る",
			clock:  &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 48, 55, 58, 50, 55, 58, 53, 50, 46, 55, 55, 49, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 51, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 48, 55, 58, 50, 55, 58, 53, 50, 46, 55, 48, 51, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 75, 97, 98, 117, 67, 111, 114, 114, 101, 99, 116, 79, 114, 100, 101, 114, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 78, 117, 109, 98, 101, 114, 34, 58, 34, 49, 49, 48, 48, 50, 56, 52, 55, 34, 44, 34, 115, 69, 105, 103, 121, 111, 117, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 49, 34, 44, 34, 115, 79, 114, 100, 101, 114, 85, 107, 101, 119, 97, 116, 97, 115, 105, 75, 105, 110, 103, 97, 107, 117, 34, 58, 34, 50, 48, 55, 56, 34, 44, 34, 115, 79, 114, 100, 101, 114, 84, 101, 115, 117, 114, 121, 111, 117, 34, 58, 34, 55, 48, 34, 44, 34, 115, 79, 114, 100, 101, 114, 83, 121, 111, 117, 104, 105, 122, 101, 105, 34, 58, 34, 55, 34, 44, 34, 115, 79, 114, 100, 101, 114, 68, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 49, 48, 55, 50, 55, 53, 50, 34, 125, 10},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   CorrectOrderRequest{},
			want1: &CorrectOrderResponse{
				CommonResponse: CommonResponse{
					No:           3,
					SendDate:     time.Date(2022, 3, 11, 7, 27, 52, 771000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 11, 7, 27, 52, 703000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeCorrectOrder,
				},
				ResultCode:     "0",
				ResultText:     "",
				OrderNumber:    "11002847",
				ExecutionDate:  time.Date(2022, 3, 11, 0, 0, 0, 0, time.Local),
				DeliveryAmount: 2078,
				Commission:     70,
				CommissionTax:  7,
				OrderDateTime:  time.Date(2022, 3, 11, 7, 27, 52, 0, time.Local),
			},
			want2: nil},
		{name: "注文変更の変更項目なしレスポンスをパース出来る",
			clock:  &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 48, 55, 58, 50, 49, 58, 49, 54, 46, 57, 50, 49, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 51, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 48, 55, 58, 50, 49, 58, 49, 54, 46, 56, 56, 49, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 75, 97, 98, 117, 67, 111, 114, 114, 101, 99, 116, 79, 114, 100, 101, 114, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 49, 50, 49, 49, 54, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 149, 207, 141, 88, 141, 128, 150, 218, 130, 170, 130, 160, 130, 232, 130, 220, 130, 185, 130, 241, 34, 44, 34, 115, 79, 114, 100, 101, 114, 78, 117, 109, 98, 101, 114, 34, 58, 34, 34, 44, 34, 115, 69, 105, 103, 121, 111, 117, 68, 97, 121, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 85, 107, 101, 119, 97, 116, 97, 115, 105, 75, 105, 110, 103, 97, 107, 117, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 84, 101, 115, 117, 114, 121, 111, 117, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 83, 121, 111, 117, 104, 105, 122, 101, 105, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 68, 97, 116, 101, 34, 58, 34, 34, 125, 10},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   CorrectOrderRequest{},
			want1: &CorrectOrderResponse{
				CommonResponse: CommonResponse{
					No:           3,
					SendDate:     time.Date(2022, 3, 11, 7, 21, 16, 921000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 11, 7, 21, 16, 881000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeCorrectOrder,
				},
				ResultCode:     "12116",
				ResultText:     "変更項目がありません",
				OrderNumber:    "",
				ExecutionDate:  time.Time{},
				DeliveryAmount: 0,
				Commission:     0,
				CommissionTax:  0,
				OrderDateTime:  time.Time{},
			},
			want2: nil},
		{name: "注文変更の失敗レスポンスをパース出来る",
			clock:  &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 48, 55, 58, 49, 56, 58, 51, 57, 46, 48, 52, 55, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 51, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 48, 55, 58, 49, 56, 58, 51, 57, 46, 48, 49, 50, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 75, 97, 98, 117, 67, 111, 114, 114, 101, 99, 116, 79, 114, 100, 101, 114, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 49, 50, 48, 49, 52, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 139, 116, 142, 119, 146, 108, 143, 240, 140, 143, 130, 201, 140, 235, 130, 232, 130, 170, 130, 160, 130, 232, 130, 220, 130, 183, 34, 44, 34, 115, 79, 114, 100, 101, 114, 78, 117, 109, 98, 101, 114, 34, 58, 34, 34, 44, 34, 115, 69, 105, 103, 121, 111, 117, 68, 97, 121, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 85, 107, 101, 119, 97, 116, 97, 115, 105, 75, 105, 110, 103, 97, 107, 117, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 84, 101, 115, 117, 114, 121, 111, 117, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 83, 121, 111, 117, 104, 105, 122, 101, 105, 34, 58, 34, 34, 44, 34, 115, 79, 114, 100, 101, 114, 68, 97, 116, 101, 34, 58, 34, 34, 125, 10},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   CorrectOrderRequest{},
			want1: &CorrectOrderResponse{
				CommonResponse: CommonResponse{
					No:           3,
					SendDate:     time.Date(2022, 3, 11, 7, 18, 39, 47000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 11, 7, 18, 39, 12000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeCorrectOrder,
				},
				ResultCode:     "12014",
				ResultText:     "逆指値条件に誤りがあります",
				OrderNumber:    "",
				ExecutionDate:  time.Time{},
				DeliveryAmount: 0,
				Commission:     0,
				CommissionTax:  0,
				OrderDateTime:  time.Time{},
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   nil,
			arg3:   CorrectOrderRequest{},
			want1:  nil,
			want2:  NilArgumentErr},
		{name: "200 OK以外が返ったらエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusInternalServerError,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   CorrectOrderRequest{},
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
			got1, got2 := client.CorrectOrder(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_CorrectOrder_Execute_NoChange(t *testing.T) {
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
		StockAccountType:  AccountTypeSpecific,
		MarginAccountType: AccountTypeUnused,
		IssueCode:         "1475",
		Exchange:          ExchangeToushou,
		Side:              SideBuy,
		ExecutionTiming:   ExecutionTimingNormal,
		OrderPrice:        2000,
		OrderQuantity:     1,
		TradeType:         TradeTypeStock,
		ExpireDate:        time.Time{},
		ExpireDateIsToday: true,
		StopOrderType:     StopOrderTypeNormal,
		TriggerPrice:      0,
		StopOrderPrice:    0,
		ExitOrderType:     ExitOrderTypeUnused,
		SecondPassword:    secondPassword,
		ExitPositions:     []ExitPosition{},
	})
	log.Printf("%+v, %+v\n", got3, got4)
	if got3.ResultCode != "0" {
		return
	}

	got5, got6 := client.CorrectOrder(context.Background(), session, CorrectOrderRequest{
		OrderNumber:        got3.OrderNumber,
		ExecutionDate:      got3.ExecutionDate,
		ExecutionTiming:    ExecutionTimingNoChange,
		OrderPrice:         NoChangeFloat,
		OrderQuantity:      NoChangeFloat,
		ExpireDate:         time.Time{},
		ExpireDateIsToday:  false,
		ExpireDateNoChange: true,
		TriggerPrice:       NoChangeFloat,
		StopOrderPrice:     NoChangeFloat,
		SecondPassword:     secondPassword,
	})
	log.Printf("%+v, %+v\n", got5, got6)
	if got5.ResultCode != "0" {
		return
	}
}

func Test_client_CorrectOrder_Execute_Change(t *testing.T) {
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
		StockAccountType:  AccountTypeSpecific,
		MarginAccountType: AccountTypeUnused,
		IssueCode:         "1475",
		Exchange:          ExchangeToushou,
		Side:              SideBuy,
		ExecutionTiming:   ExecutionTimingNormal,
		OrderPrice:        2000,
		OrderQuantity:     1,
		TradeType:         TradeTypeStock,
		ExpireDate:        time.Time{},
		ExpireDateIsToday: true,
		StopOrderType:     StopOrderTypeNormal,
		TriggerPrice:      0,
		StopOrderPrice:    0,
		ExitOrderType:     ExitOrderTypeUnused,
		SecondPassword:    secondPassword,
		ExitPositions:     []ExitPosition{},
	})
	log.Printf("%+v, %+v\n", got3, got4)
	if got3.ResultCode != "0" {
		return
	}

	got5, got6 := client.CorrectOrder(context.Background(), session, CorrectOrderRequest{
		OrderNumber:        got3.OrderNumber,
		ExecutionDate:      got3.ExecutionDate,
		ExecutionTiming:    ExecutionTimingNoChange,
		OrderPrice:         2001,
		OrderQuantity:      NoChangeFloat,
		ExpireDate:         time.Time{},
		ExpireDateIsToday:  false,
		ExpireDateNoChange: true,
		TriggerPrice:       NoChangeFloat,
		StopOrderPrice:     NoChangeFloat,
		SecondPassword:     secondPassword,
	})
	log.Printf("%+v, %+v\n", got5, got6)
	if got5.ResultCode != "0" {
		return
	}
}
