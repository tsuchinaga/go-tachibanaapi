package tachibana

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"
	"time"
)

func Test_client_StockSellable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		clock     *testClock
		requester *testRequester
		arg1      context.Context
		arg2      *Session
		arg3      StockSellableRequest
		want1     *StockSellableResponse
		want2     error
	}{
		{name: "成功レスポンスをパース出来る",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 50, 51, 58, 48, 55, 58, 50, 57, 46, 57, 48, 51, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 50, 51, 58, 48, 55, 58, 50, 57, 46, 56, 55, 57, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 90, 97, 110, 85, 114, 105, 75, 97, 110, 111, 117, 115, 117, 117, 34, 44, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 53, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 85, 112, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 49, 50, 51, 48, 55, 34, 44, 34, 115, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 85, 114, 105, 75, 97, 110, 111, 117, 73, 112, 112, 97, 110, 34, 58, 34, 48, 34, 44, 34, 115, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 85, 114, 105, 75, 97, 110, 111, 117, 84, 111, 107, 117, 116, 101, 105, 34, 58, 34, 48, 34, 44, 34, 115, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 85, 114, 105, 75, 97, 110, 111, 117, 78, 105, 115, 97, 34, 58, 34, 48, 34, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockSellableRequest{},
			want1: &StockSellableResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 11, 23, 7, 29, 903000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 11, 23, 7, 29, 879000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockSellable,
				},
				IssueCode:        "1475",
				ResultCode:       "0",
				ResultText:       "",
				WarningCode:      "0",
				WarningText:      "",
				UpdateDateTime:   time.Date(2022, 3, 11, 23, 7, 0, 0, time.Local),
				GeneralQuantity:  0,
				SpecificQuantity: 0,
				NisaQuantity:     0,
			},
			want2: nil},
		{name: "存在しない銘柄を指定した失敗レスポンスをパース出来る",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 50, 51, 58, 48, 53, 58, 49, 48, 46, 53, 57, 57, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 50, 51, 58, 48, 53, 58, 49, 48, 46, 53, 56, 50, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 90, 97, 110, 85, 114, 105, 75, 97, 110, 111, 117, 115, 117, 117, 34, 44, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 57, 57, 49, 48, 48, 51, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 233, 138, 152, 230, 159, 132, 227, 130, 179, 227, 131, 188, 227, 131, 137, 227, 129, 171, 232, 170, 164, 227, 130, 138, 227, 129, 140, 227, 129, 130, 227, 130, 138, 227, 129, 190, 227, 129, 153, 227, 128, 130, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 85, 112, 100, 97, 116, 101, 34, 58, 34, 34, 44, 34, 115, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 85, 114, 105, 75, 97, 110, 111, 117, 73, 112, 112, 97, 110, 34, 58, 34, 34, 44, 34, 115, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 85, 114, 105, 75, 97, 110, 111, 117, 84, 111, 107, 117, 116, 101, 105, 34, 58, 34, 34, 44, 34, 115, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 85, 114, 105, 75, 97, 110, 111, 117, 78, 105, 115, 97, 34, 58, 34, 34, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockSellableRequest{},
			want1: &StockSellableResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 11, 23, 5, 10, 599000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 11, 23, 5, 10, 582000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockSellable,
				},
				IssueCode:        "",
				ResultCode:       "991003",
				ResultText:       "銘柄コードに誤りがあります。",
				WarningCode:      "",
				WarningText:      "",
				UpdateDateTime:   time.Time{},
				GeneralQuantity:  0,
				SpecificQuantity: 0,
				NisaQuantity:     0,
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			arg1:  context.Background(),
			arg2:  nil,
			arg3:  StockSellableRequest{},
			want1: nil,
			want2: NilArgumentErr},
		{name: "リクエストでエラーが返されたらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get2: StatusNotOkErr},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockSellableRequest{},
			want1:     nil,
			want2:     StatusNotOkErr},
		{name: "レスポンスのUnmarshalに失敗したらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get1: []byte{}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockSellableRequest{},
			want1:     nil,
			want2:     UnmarshalFailedErr},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &client{clock: test.clock, requester: test.requester}
			got1, got2 := client.StockSellable(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_StockSellable_Execute(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	userId := "user-id"
	password := "password"

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

	got3, got4 := client.StockSellable(context.Background(), session, StockSellableRequest{
		IssueCode: "1475",
	})
	log.Printf("%+v, %+v\n", got3, got4)
	if got3.ResultCode != "0" {
		return
	}
}
