package tachibana

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"
	"time"
)

func Test_client_StockWallet(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		clock     *testClock
		requester *testRequester
		arg1      context.Context
		arg2      *Session
		arg3      StockWalletRequest
		want1     *StockWalletResponse
		want2     error
	}{
		{name: "成功レスポンスをパース出来る",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 49, 48, 58, 52, 51, 58, 50, 48, 46, 56, 52, 56, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 49, 48, 58, 52, 51, 58, 50, 48, 46, 56, 48, 50, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 90, 97, 110, 75, 97, 105, 75, 97, 110, 111, 117, 103, 97, 107, 117, 34, 44, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 53, 34, 44, 34, 115, 83, 105, 122, 121, 111, 117, 67, 34, 58, 34, 48, 48, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 85, 112, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 49, 49, 48, 52, 51, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 71, 101, 110, 107, 97, 98, 117, 75, 97, 105, 116, 117, 107, 101, 34, 58, 34, 49, 48, 48, 48, 48, 49, 49, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 78, 105, 115, 97, 75, 97, 105, 116, 117, 107, 101, 75, 97, 110, 111, 117, 103, 97, 107, 117, 34, 58, 34, 48, 34, 44, 34, 115, 72, 117, 115, 111, 107, 117, 107, 105, 110, 72, 97, 115, 115, 101, 105, 70, 108, 103, 34, 58, 34, 48, 34, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockWalletRequest{},
			want1: &StockWalletResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 11, 10, 43, 20, 848000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 11, 10, 43, 20, 802000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeStockWallet,
				},
				IssueCode:      "1475",
				Exchange:       ExchangeToushou,
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				UpdateDateTime: time.Date(2022, 3, 11, 10, 43, 0, 0, time.Local),
				StockWallet:    1000011,
				NisaWallet:     0,
				Shortage:       false,
			},
			want2: nil},
		{name: "銘柄未指定の成功レスポンスをパース出来る",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 49, 48, 58, 52, 51, 58, 48, 55, 46, 57, 55, 51, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 49, 48, 58, 52, 51, 58, 48, 55, 46, 57, 51, 52, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 90, 97, 110, 75, 97, 105, 75, 97, 110, 111, 117, 103, 97, 107, 117, 34, 44, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 34, 44, 34, 115, 83, 105, 122, 121, 111, 117, 67, 34, 58, 34, 48, 48, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 85, 112, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 49, 49, 48, 52, 51, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 71, 101, 110, 107, 97, 98, 117, 75, 97, 105, 116, 117, 107, 101, 34, 58, 34, 49, 48, 48, 48, 48, 49, 49, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 78, 105, 115, 97, 75, 97, 105, 116, 117, 107, 101, 75, 97, 110, 111, 117, 103, 97, 107, 117, 34, 58, 34, 48, 34, 44, 34, 115, 72, 117, 115, 111, 107, 117, 107, 105, 110, 72, 97, 115, 115, 101, 105, 70, 108, 103, 34, 58, 34, 48, 34, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockWalletRequest{},
			want1: &StockWalletResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 11, 10, 43, 7, 973000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 11, 10, 43, 7, 934000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeStockWallet,
				},
				IssueCode:      "",
				Exchange:       ExchangeToushou,
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				UpdateDateTime: time.Date(2022, 3, 11, 10, 43, 0, 0, time.Local),
				StockWallet:    1000011,
				NisaWallet:     0,
				Shortage:       false,
			},
			want2: nil},
		{name: "存在しない銘柄を指定した失敗レスポンスをパース出来る",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 49, 48, 58, 52, 50, 58, 53, 49, 46, 52, 55, 55, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 49, 48, 58, 52, 50, 58, 53, 49, 46, 52, 51, 57, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 90, 97, 110, 75, 97, 105, 75, 97, 110, 111, 117, 103, 97, 107, 117, 34, 44, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 34, 44, 34, 115, 83, 105, 122, 121, 111, 117, 67, 34, 58, 34, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 57, 57, 49, 48, 48, 51, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 233, 138, 152, 230, 159, 132, 227, 130, 179, 227, 131, 188, 227, 131, 137, 227, 129, 171, 232, 170, 164, 227, 130, 138, 227, 129, 140, 227, 129, 130, 227, 130, 138, 227, 129, 190, 227, 129, 153, 227, 128, 130, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 85, 112, 100, 97, 116, 101, 34, 58, 34, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 71, 101, 110, 107, 97, 98, 117, 75, 97, 105, 116, 117, 107, 101, 34, 58, 34, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 78, 105, 115, 97, 75, 97, 105, 116, 117, 107, 101, 75, 97, 110, 111, 117, 103, 97, 107, 117, 34, 58, 34, 34, 44, 34, 115, 72, 117, 115, 111, 107, 117, 107, 105, 110, 72, 97, 115, 115, 101, 105, 70, 108, 103, 34, 58, 34, 34, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockWalletRequest{},
			want1: &StockWalletResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 11, 10, 42, 51, 477000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 11, 10, 42, 51, 439000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeStockWallet,
				},
				IssueCode:      "",
				Exchange:       "",
				ResultCode:     "991003",
				ResultText:     "銘柄コードに誤りがあります。",
				WarningCode:    "",
				WarningText:    "",
				UpdateDateTime: time.Time{},
				StockWallet:    0,
				NisaWallet:     0,
				Shortage:       false,
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			arg1:  context.Background(),
			arg2:  nil,
			arg3:  StockWalletRequest{},
			want1: nil,
			want2: NilArgumentErr},
		{name: "リクエストでエラーが返されたらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get2: StatusNotOkErr},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockWalletRequest{},
			want1:     nil,
			want2:     StatusNotOkErr},
		{name: "レスポンスのUnmarshalに失敗したらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get1: []byte{}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockWalletRequest{},
			want1:     nil,
			want2:     UnmarshalFailedErr},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &client{clock: test.clock, requester: test.requester}
			got1, got2 := client.StockWallet(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_StockWallet_Execute(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

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

	got3, got4 := client.StockWallet(context.Background(), session, StockWalletRequest{
		IssueCode: "1475",
		Exchange:  ExchangeToushou,
	})
	log.Printf("%+v, %+v\n", got3, got4)
	if got3.ResultCode != "0" {
		return
	}
}
