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

func Test_client_StockSellable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		clock  *testClock
		status int
		body   []byte
		arg1   context.Context
		arg2   *Session
		arg3   StockSellableRequest
		want1  *StockSellableResponse
		want2  error
	}{
		{name: "成功レスポンスをパース出来る",
			clock:  &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 50, 51, 58, 48, 55, 58, 50, 57, 46, 57, 48, 51, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 50, 51, 58, 48, 55, 58, 50, 57, 46, 56, 55, 57, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 90, 97, 110, 85, 114, 105, 75, 97, 110, 111, 117, 115, 117, 117, 34, 44, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 53, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 85, 112, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 49, 50, 51, 48, 55, 34, 44, 34, 115, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 85, 114, 105, 75, 97, 110, 111, 117, 73, 112, 112, 97, 110, 34, 58, 34, 48, 34, 44, 34, 115, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 85, 114, 105, 75, 97, 110, 111, 117, 84, 111, 107, 117, 116, 101, 105, 34, 58, 34, 48, 34, 44, 34, 115, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 85, 114, 105, 75, 97, 110, 111, 117, 78, 105, 115, 97, 34, 58, 34, 48, 34, 125, 10},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   StockSellableRequest{},
			want1: &StockSellableResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 11, 23, 7, 29, 903000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 11, 23, 7, 29, 879000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockSellable,
				},
				SymbolCode:       "1475",
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
			clock:  &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 50, 51, 58, 48, 53, 58, 49, 48, 46, 53, 57, 57, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 49, 49, 45, 50, 51, 58, 48, 53, 58, 49, 48, 46, 53, 56, 50, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 90, 97, 110, 85, 114, 105, 75, 97, 110, 111, 117, 115, 117, 117, 34, 44, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 57, 57, 49, 48, 48, 51, 34, 44, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 150, 193, 149, 191, 131, 82, 129, 91, 131, 104, 130, 201, 140, 235, 130, 232, 130, 170, 130, 160, 130, 232, 130, 220, 130, 183, 129, 66, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 34, 44, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 34, 115, 83, 117, 109, 109, 97, 114, 121, 85, 112, 100, 97, 116, 101, 34, 58, 34, 34, 44, 34, 115, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 85, 114, 105, 75, 97, 110, 111, 117, 73, 112, 112, 97, 110, 34, 58, 34, 34, 44, 34, 115, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 85, 114, 105, 75, 97, 110, 111, 117, 84, 111, 107, 117, 116, 101, 105, 34, 58, 34, 34, 44, 34, 115, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 85, 114, 105, 75, 97, 110, 111, 117, 78, 105, 115, 97, 34, 58, 34, 34, 125, 10},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   StockSellableRequest{},
			want1: &StockSellableResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 11, 23, 5, 10, 599000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 11, 23, 5, 10, 582000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockSellable,
				},
				SymbolCode:       "",
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
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   nil,
			arg3:   StockSellableRequest{},
			want1:  nil,
			want2:  NilArgumentErr},
		{name: "200 OK以外が返ったらエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusInternalServerError,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   StockSellableRequest{},
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
		SymbolCode: "1475",
	})
	log.Printf("%+v, %+v\n", got3, got4)
	if got3.ResultCode != "0" {
		return
	}
}
