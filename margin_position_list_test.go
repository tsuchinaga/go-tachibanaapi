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

func Test_MarginPositionListRequest_request(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request MarginPositionListRequest
		arg1    int64
		arg2    time.Time
		want1   marginPositionListRequest
	}{
		{name: "変換できる",
			request: MarginPositionListRequest{
				SymbolCode: "1475",
			},
			arg1: 123,
			arg2: time.Date(2022, 3, 1, 9, 40, 0, 0, time.Local),
			want1: marginPositionListRequest{
				commonRequest: commonRequest{
					No:          123,
					SendDate:    RequestTime{Time: time.Date(2022, 3, 1, 9, 40, 0, 0, time.Local)},
					FeatureType: FeatureTypeMarginPositionList,
				},
				SymbolCode: "1475",
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

func Test_marginPositionListResponse_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		arg1         []byte
		wantResponse marginPositionListResponse
		hasError     bool
	}{
		{name: "正常系のレスポンスをパースできる",
			arg1: []byte(`{"177":"2022.02.28-11:13:31.300","175":"2","176":"2022.02.28-11:13:31.268","174":"0","173":"","192":"CLMShinyouTategyokuList","328":"","534":"0","535":"","692":"0","693":"","686":"0","344":"1909","649":"1909","309":"0","308":"1","652":"1","643":"1","314":"0","56":[{"521":"0","522":"","513":"202202280006156","485":"1475","501":"00","467":"3","468":"26","528":"1","514":"1","515":"1909.0000","484":"1910.0000","475":"1","476":"0.05","622":"1909","509":"0","529":"0","478":"0","486":"0","487":"0","488":"0","502":"0","511":"20220228","512":"20220825","624":"1","523":"0","477":"0","496":"1","483":"0","600":"1914","735":"-4","736":"-0.20","668":"07"}]}`),
			wantResponse: marginPositionListResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 2, 28, 11, 13, 31, 300000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 2, 28, 11, 13, 31, 268000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMarginPositionList,
				},
				SymbolCode:            "",
				ResultCode:            "0",
				ResultText:            "",
				WarningCode:           "0",
				WarningText:           "",
				TotalSellAmount:       0,
				TotalBuyAmount:        1909,
				TotalAmount:           1909,
				TotalSellProfit:       0,
				TotalBuyProfit:        1,
				TotalProfit:           1,
				SpecificAccountProfit: 1,
				GeneralAccountProfit:  0,
				Positions: []marginPosition{{
					WarningCode:        "0",
					WarningText:        "",
					PositionCode:       "202202280006156",
					SymbolCode:         "1475",
					Exchange:           ExchangeToushou,
					Side:               SideBuy,
					ExitTermType:       ExitTermTypeSystemMargin6m,
					AccountType:        AccountTypeSpecific,
					OrderQuantity:      1,
					UnitPrice:          1909,
					CurrentPrice:       1910,
					Profit:             1,
					ProfitRatio:        0.05,
					TotalPrice:         1909,
					Commission:         0,
					Interest:           0,
					Premiums:           0,
					RewritingFee:       0,
					ManagementFee:      0,
					LendingFee:         0,
					OtherFee:           0,
					ContractDate:       Ymd{Time: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local)},
					ExitTerm:           Ymd{Time: time.Date(2022, 8, 25, 0, 0, 0, 0, time.Local)},
					OwnedQuantity:      1,
					ExitQuantity:       0,
					DeliveryQuantity:   0,
					HoldQuantity:       1,
					ReturnableQuantity: 0,
					PrevClosePrice:     1914,
					PrevCloseRatio:     -4,
					PrevClosePercent:   -0.20,
					PrevCloseRatioType: PrevCloseRatioTypeUnder0,
				}},
			}},
		{name: "異常系のレスポンスをパースできる",
			arg1: []byte(`{"177":"2022.02.28-10:12:34.242","175":"2","176":"2022.02.28-10:12:34.206","174":"0","173":"","192":"CLMShinyouTategyokuList","328":"*","534":"0","535":"","692":"0","693":"","686":"0","344":"0","649":"0","309":"0","308":"0","652":"0","643":"0","314":"0","56":""}`),
			wantResponse: marginPositionListResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 2, 28, 10, 12, 34, 242000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 2, 28, 10, 12, 34, 206000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMarginPositionList,
				},
				SymbolCode:            "*",
				ResultCode:            "0",
				ResultText:            "",
				WarningCode:           "0",
				WarningText:           "",
				TotalSellAmount:       0,
				TotalBuyAmount:        0,
				TotalAmount:           0,
				TotalSellProfit:       0,
				TotalBuyProfit:        0,
				TotalProfit:           0,
				SpecificAccountProfit: 0,
				GeneralAccountProfit:  0,
				Positions:             []marginPosition{},
			}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var res marginPositionListResponse
			got1 := json.Unmarshal(test.arg1, &res)
			if !reflect.DeepEqual(test.wantResponse, res) || (got1 != nil) != test.hasError {
				t.Errorf("%s error\nwant: %+v\ngot: %+v, %+v\n", t.Name(), test.wantResponse, res, got1)
			}
		})
	}
}

func Test_marginPositionListResponse_response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		response marginPositionListResponse
		want1    MarginPositionListResponse
	}{
		{name: "変換できる",
			response: marginPositionListResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 2, 28, 11, 13, 31, 300000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 2, 28, 11, 13, 31, 268000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMarginPositionList,
				},
				SymbolCode:            "",
				ResultCode:            "0",
				ResultText:            "",
				WarningCode:           "0",
				WarningText:           "",
				TotalSellAmount:       0,
				TotalBuyAmount:        1909,
				TotalAmount:           1909,
				TotalSellProfit:       0,
				TotalBuyProfit:        1,
				TotalProfit:           1,
				SpecificAccountProfit: 1,
				GeneralAccountProfit:  0,
				Positions: []marginPosition{{
					WarningCode:        "0",
					WarningText:        "",
					PositionCode:       "202202280006156",
					SymbolCode:         "1475",
					Exchange:           ExchangeToushou,
					Side:               SideBuy,
					ExitTermType:       ExitTermTypeSystemMargin6m,
					AccountType:        AccountTypeSpecific,
					OrderQuantity:      1,
					UnitPrice:          1909,
					CurrentPrice:       1910,
					Profit:             1,
					ProfitRatio:        0.05,
					TotalPrice:         1909,
					Commission:         0,
					Interest:           0,
					Premiums:           0,
					RewritingFee:       0,
					ManagementFee:      0,
					LendingFee:         0,
					OtherFee:           0,
					ContractDate:       Ymd{Time: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local)},
					ExitTerm:           Ymd{Time: time.Date(2022, 8, 25, 0, 0, 0, 0, time.Local)},
					OwnedQuantity:      1,
					ExitQuantity:       0,
					DeliveryQuantity:   0,
					HoldQuantity:       1,
					ReturnableQuantity: 0,
					PrevClosePrice:     1914,
					PrevCloseRatio:     -4,
					PrevClosePercent:   -0.20,
					PrevCloseRatioType: PrevCloseRatioTypeUnder0,
				}},
			},
			want1: MarginPositionListResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 28, 11, 13, 31, 300000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 28, 11, 13, 31, 268000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMarginPositionList,
				},
				SymbolCode:            "",
				ResultCode:            "0",
				ResultText:            "",
				WarningCode:           "0",
				WarningText:           "",
				TotalSellAmount:       0,
				TotalBuyAmount:        1909,
				TotalAmount:           1909,
				TotalSellProfit:       0,
				TotalBuyProfit:        1,
				TotalProfit:           1,
				SpecificAccountProfit: 1,
				GeneralAccountProfit:  0,
				Positions: []MarginPosition{{
					WarningCode:        "0",
					WarningText:        "",
					PositionCode:       "202202280006156",
					SymbolCode:         "1475",
					Exchange:           ExchangeToushou,
					Side:               SideBuy,
					ExitTermType:       ExitTermTypeSystemMargin6m,
					AccountType:        AccountTypeSpecific,
					OrderQuantity:      1,
					UnitPrice:          1909,
					CurrentPrice:       1910,
					Profit:             1,
					ProfitRatio:        0.05,
					TotalPrice:         1909,
					Commission:         0,
					Interest:           0,
					Premiums:           0,
					RewritingFee:       0,
					ManagementFee:      0,
					LendingFee:         0,
					OtherFee:           0,
					ContractDate:       time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local),
					ExitTerm:           time.Date(2022, 8, 25, 0, 0, 0, 0, time.Local),
					OwnedQuantity:      1,
					ExitQuantity:       0,
					DeliveryQuantity:   0,
					HoldQuantity:       1,
					ReturnableQuantity: 0,
					PrevClosePrice:     1914,
					PrevCloseRatio:     -4,
					PrevClosePercent:   -0.20,
					PrevCloseRatioType: PrevCloseRatioTypeUnder0,
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

func Test_marginPosition_response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		response marginPosition
		want1    MarginPosition
	}{
		{name: "変換できる",
			response: marginPosition{
				WarningCode:        "0",
				WarningText:        "",
				PositionCode:       "202202280006156",
				SymbolCode:         "1475",
				Exchange:           ExchangeToushou,
				Side:               SideBuy,
				ExitTermType:       ExitTermTypeSystemMargin6m,
				AccountType:        AccountTypeSpecific,
				OrderQuantity:      1,
				UnitPrice:          1909,
				CurrentPrice:       1910,
				Profit:             1,
				ProfitRatio:        0.05,
				TotalPrice:         1909,
				Commission:         0,
				Interest:           0,
				Premiums:           0,
				RewritingFee:       0,
				ManagementFee:      0,
				LendingFee:         0,
				OtherFee:           0,
				ContractDate:       Ymd{Time: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local)},
				ExitTerm:           Ymd{Time: time.Date(2022, 8, 25, 0, 0, 0, 0, time.Local)},
				OwnedQuantity:      1,
				ExitQuantity:       0,
				DeliveryQuantity:   0,
				HoldQuantity:       1,
				ReturnableQuantity: 0,
				PrevClosePrice:     1914,
				PrevCloseRatio:     -4,
				PrevClosePercent:   -0.20,
				PrevCloseRatioType: PrevCloseRatioTypeUnder0,
			},
			want1: MarginPosition{
				WarningCode:        "0",
				WarningText:        "",
				PositionCode:       "202202280006156",
				SymbolCode:         "1475",
				Exchange:           ExchangeToushou,
				Side:               SideBuy,
				ExitTermType:       ExitTermTypeSystemMargin6m,
				AccountType:        AccountTypeSpecific,
				OrderQuantity:      1,
				UnitPrice:          1909,
				CurrentPrice:       1910,
				Profit:             1,
				ProfitRatio:        0.05,
				TotalPrice:         1909,
				Commission:         0,
				Interest:           0,
				Premiums:           0,
				RewritingFee:       0,
				ManagementFee:      0,
				LendingFee:         0,
				OtherFee:           0,
				ContractDate:       time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local),
				ExitTerm:           time.Date(2022, 8, 25, 0, 0, 0, 0, time.Local),
				OwnedQuantity:      1,
				ExitQuantity:       0,
				DeliveryQuantity:   0,
				HoldQuantity:       1,
				ReturnableQuantity: 0,
				PrevClosePrice:     1914,
				PrevCloseRatio:     -4,
				PrevClosePercent:   -0.20,
				PrevCloseRatioType: PrevCloseRatioTypeUnder0,
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

func Test_client_MarginPositionList(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name   string
		clock  *testClock
		status int
		body   []byte
		arg1   context.Context
		arg2   *Session
		arg3   MarginPositionListRequest
		want1  *MarginPositionListResponse
		want2  error
	}{
		{name: "正常レスポンスをパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 49, 49, 58, 49, 51, 58, 51, 49, 46, 51, 48, 48, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 49, 49, 58, 49, 51, 58, 51, 49, 46, 50, 54, 56, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 83, 104, 105, 110, 121, 111, 117, 84, 97, 116, 101, 103, 121, 111, 107, 117, 76, 105, 115, 116, 34, 44, 34, 51, 50, 56, 34, 58, 34, 34, 44, 34, 53, 51, 52, 34, 58, 34, 48, 34, 44, 34, 53, 51, 53, 34, 58, 34, 34, 44, 34, 54, 57, 50, 34, 58, 34, 48, 34, 44, 34, 54, 57, 51, 34, 58, 34, 34, 44, 34, 54, 56, 54, 34, 58, 34, 48, 34, 44, 34, 51, 52, 52, 34, 58, 34, 49, 57, 48, 57, 34, 44, 34, 54, 52, 57, 34, 58, 34, 49, 57, 48, 57, 34, 44, 34, 51, 48, 57, 34, 58, 34, 48, 34, 44, 34, 51, 48, 56, 34, 58, 34, 49, 34, 44, 34, 54, 53, 50, 34, 58, 34, 49, 34, 44, 34, 54, 52, 51, 34, 58, 34, 49, 34, 44, 34, 51, 49, 52, 34, 58, 34, 48, 34, 44, 34, 53, 54, 34, 58, 91, 123, 34, 53, 50, 49, 34, 58, 34, 48, 34, 44, 34, 53, 50, 50, 34, 58, 34, 34, 44, 34, 53, 49, 51, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 48, 48, 48, 54, 49, 53, 54, 34, 44, 34, 52, 56, 53, 34, 58, 34, 49, 52, 55, 53, 34, 44, 34, 53, 48, 49, 34, 58, 34, 48, 48, 34, 44, 34, 52, 54, 55, 34, 58, 34, 51, 34, 44, 34, 52, 54, 56, 34, 58, 34, 50, 54, 34, 44, 34, 53, 50, 56, 34, 58, 34, 49, 34, 44, 34, 53, 49, 52, 34, 58, 34, 49, 34, 44, 34, 53, 49, 53, 34, 58, 34, 49, 57, 48, 57, 46, 48, 48, 48, 48, 34, 44, 34, 52, 56, 52, 34, 58, 34, 49, 57, 49, 48, 46, 48, 48, 48, 48, 34, 44, 34, 52, 55, 53, 34, 58, 34, 49, 34, 44, 34, 52, 55, 54, 34, 58, 34, 48, 46, 48, 53, 34, 44, 34, 54, 50, 50, 34, 58, 34, 49, 57, 48, 57, 34, 44, 34, 53, 48, 57, 34, 58, 34, 48, 34, 44, 34, 53, 50, 57, 34, 58, 34, 48, 34, 44, 34, 52, 55, 56, 34, 58, 34, 48, 34, 44, 34, 52, 56, 54, 34, 58, 34, 48, 34, 44, 34, 52, 56, 55, 34, 58, 34, 48, 34, 44, 34, 52, 56, 56, 34, 58, 34, 48, 34, 44, 34, 53, 48, 50, 34, 58, 34, 48, 34, 44, 34, 53, 49, 49, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 34, 44, 34, 53, 49, 50, 34, 58, 34, 50, 48, 50, 50, 48, 56, 50, 53, 34, 44, 34, 54, 50, 52, 34, 58, 34, 49, 34, 44, 34, 53, 50, 51, 34, 58, 34, 48, 34, 44, 34, 52, 55, 55, 34, 58, 34, 48, 34, 44, 34, 52, 57, 54, 34, 58, 34, 49, 34, 44, 34, 52, 56, 51, 34, 58, 34, 48, 34, 44, 34, 54, 48, 48, 34, 58, 34, 49, 57, 49, 52, 34, 44, 34, 55, 51, 53, 34, 58, 34, 45, 52, 34, 44, 34, 55, 51, 54, 34, 58, 34, 45, 48, 46, 50, 48, 34, 44, 34, 54, 54, 56, 34, 58, 34, 48, 55, 34, 125, 93, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   MarginPositionListRequest{},
			want1: &MarginPositionListResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 28, 11, 13, 31, 300000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 28, 11, 13, 31, 268000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMarginPositionList,
				},
				SymbolCode:            "",
				ResultCode:            "0",
				ResultText:            "",
				WarningCode:           "0",
				WarningText:           "",
				TotalSellAmount:       0,
				TotalBuyAmount:        1909,
				TotalAmount:           1909,
				TotalSellProfit:       0,
				TotalBuyProfit:        1,
				TotalProfit:           1,
				SpecificAccountProfit: 1,
				GeneralAccountProfit:  0,
				Positions: []MarginPosition{{
					WarningCode:        "0",
					WarningText:        "",
					PositionCode:       "202202280006156",
					SymbolCode:         "1475",
					Exchange:           ExchangeToushou,
					Side:               SideBuy,
					ExitTermType:       ExitTermTypeSystemMargin6m,
					AccountType:        AccountTypeSpecific,
					OrderQuantity:      1,
					UnitPrice:          1909,
					CurrentPrice:       1910,
					Profit:             1,
					ProfitRatio:        0.05,
					TotalPrice:         1909,
					Commission:         0,
					Interest:           0,
					Premiums:           0,
					RewritingFee:       0,
					ManagementFee:      0,
					LendingFee:         0,
					OtherFee:           0,
					ContractDate:       time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local),
					ExitTerm:           time.Date(2022, 8, 25, 0, 0, 0, 0, time.Local),
					OwnedQuantity:      1,
					ExitQuantity:       0,
					DeliveryQuantity:   0,
					HoldQuantity:       1,
					ReturnableQuantity: 0,
					PrevClosePrice:     1914,
					PrevCloseRatio:     -4,
					PrevClosePercent:   -0.20,
					PrevCloseRatioType: PrevCloseRatioTypeUnder0,
				}},
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   nil,
			arg3:   MarginPositionListRequest{},
			want1:  nil,
			want2:  NilArgumentErr},
		{name: "失敗をパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 49, 48, 58, 49, 50, 58, 51, 52, 46, 50, 52, 50, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 49, 48, 58, 49, 50, 58, 51, 52, 46, 50, 48, 54, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 83, 104, 105, 110, 121, 111, 117, 84, 97, 116, 101, 103, 121, 111, 107, 117, 76, 105, 115, 116, 34, 44, 34, 51, 50, 56, 34, 58, 34, 42, 34, 44, 34, 53, 51, 52, 34, 58, 34, 48, 34, 44, 34, 53, 51, 53, 34, 58, 34, 34, 44, 34, 54, 57, 50, 34, 58, 34, 48, 34, 44, 34, 54, 57, 51, 34, 58, 34, 34, 44, 34, 54, 56, 54, 34, 58, 34, 48, 34, 44, 34, 51, 52, 52, 34, 58, 34, 48, 34, 44, 34, 54, 52, 57, 34, 58, 34, 48, 34, 44, 34, 51, 48, 57, 34, 58, 34, 48, 34, 44, 34, 51, 48, 56, 34, 58, 34, 48, 34, 44, 34, 54, 53, 50, 34, 58, 34, 48, 34, 44, 34, 54, 52, 51, 34, 58, 34, 48, 34, 44, 34, 51, 49, 52, 34, 58, 34, 48, 34, 44, 34, 53, 54, 34, 58, 34, 34, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   MarginPositionListRequest{},
			want1: &MarginPositionListResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 28, 10, 12, 34, 242000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 28, 10, 12, 34, 206000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMarginPositionList,
				},
				SymbolCode:            "*",
				ResultCode:            "0",
				ResultText:            "",
				WarningCode:           "0",
				WarningText:           "",
				TotalSellAmount:       0,
				TotalBuyAmount:        0,
				TotalAmount:           0,
				TotalSellProfit:       0,
				TotalBuyProfit:        0,
				TotalProfit:           0,
				SpecificAccountProfit: 0,
				GeneralAccountProfit:  0,
				Positions:             []MarginPosition{},
			},
			want2: nil},
		{name: "200 OK以外が返ったらエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusInternalServerError,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   MarginPositionListRequest{},
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
			got1, got2 := client.MarginPositionList(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_MarginPositionList_Execute(t *testing.T) {
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

	got3, got4 := client.MarginPositionList(context.Background(), session, MarginPositionListRequest{
		SymbolCode: "",
	})
	log.Printf("%+v, %+v\n", got3, got4)
}
