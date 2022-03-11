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

func Test_StockPositionListRequest_request(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request StockPositionListRequest
		arg1    int64
		arg2    time.Time
		want1   stockPositionListRequest
	}{
		{name: "変換できる",
			request: StockPositionListRequest{IssueCode: "1475"},
			arg1:    123,
			arg2:    time.Date(2022, 3, 1, 9, 0, 0, 0, time.Local),
			want1: stockPositionListRequest{
				commonRequest: commonRequest{
					No:             123,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 1, 9, 0, 0, 0, time.Local)},
					FeatureType:    FeatureTypeStockPositionList,
					ResponseFormat: commonResponseFormat,
				},
				IssueCode: "1475",
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

func Test_stockPositionListResponse_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		arg1         []byte
		wantResponse stockPositionListResponse
		hasError     bool
	}{
		{name: "正常系のパース",
			arg1: []byte(`{
	"p_sd_date":"2022.03.01-14:56:35.717",
	"p_no":"2",
	"p_rv_date":"2022.03.01-14:56:35.693",
	"p_errno":"0",
	"p_err":"",
	"sCLMID":"CLMGenbutuKabuList",
	"sIssueCode":"",
	"sResultCode":"0",
	"sResultText":"",
	"sWarningCode":"0",
	"sWarningText":"",
	"sTokuteiGaisanHyoukagakuGoukei":"1932",
	"sIppanGaisanHyoukagakuGoukei":"0",
	"sNisaGaisanHyoukagakuGoukei":"0",
	"sTotalGaisanHyoukagakuGoukei":"1932",
	"sTokuteiGaisanHyoukaSonekiGoukei":"-78",
	"sIppanGaisanHyoukaSonekiGoukei":"0",
	"sNisaGaisanHyoukaSonekiGoukei":"0",
	"sTotalGaisanHyoukaSonekiGoukei":"-78",
	"aGenbutuKabuList":
	[
	{
		"sUriOrderWarningCode":"0",
		"sUriOrderWarningText":"",
		"sUriOrderIssueCode":"1475",
		"sUriOrderZyoutoekiKazeiC":"1",
		"sUriOrderZanKabuSuryou":"1",
		"sUriOrderUritukeKanouSuryou":"0",
		"sUriOrderGaisanBokaTanka":"2010.0000",
		"sUriOrderHyoukaTanka":"1932.0000",
		"sUriOrderGaisanHyoukagaku":"1932",
		"sUriOrderGaisanHyoukaSoneki":"-78",
		"sUriOrderGaisanHyoukaSonekiRitu":"-3.88",
		"sSyuzituOwarine":"1923",
		"sZenzituHi":"9",
		"sZenzituHiPer":"0.46",
		"sUpDownFlag":"05",
		"sNissyoukinKasikabuZan":"2409"
	}
	]
}`),
			wantResponse: stockPositionListResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 3, 1, 14, 56, 35, 717000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 3, 1, 14, 56, 35, 693000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockPositionList,
				},
				IssueCode:      "",
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				SpecificAmount: 1932,
				GeneralAmount:  0,
				NisaAmount:     0,
				TotalAmount:    1932,
				SpecificProfit: -78,
				GeneralProfit:  0,
				NisaProfit:     0,
				TotalProfit:    -78,
				Positions: []stockPosition{{
					WarningCode:        "0",
					WarningText:        "",
					IssueCode:          "1475",
					AccountType:        AccountTypeSpecific,
					OwnedQuantity:      1,
					UnHoldQuantity:     0,
					UnitValuation:      2010,
					BookValuation:      1932,
					TotalValuation:     1932,
					Profit:             -78,
					ProfitRatio:        -3.88,
					PrevClosePrice:     1923,
					PrevCloseRatio:     9,
					PrevClosePercent:   0.46,
					PrevCloseRatioType: PrevCloseRatioTypeOver0,
					MarginBalance:      2409,
				}},
			}},
		{name: "空配列でもパースできる",
			arg1: []byte(`{
	"p_sd_date":"2022.03.01-14:57:42.983",
	"p_no":"2",
	"p_rv_date":"2022.03.01-14:57:42.966",
	"p_errno":"0",
	"p_err":"",
	"sCLMID":"CLMGenbutuKabuList",
	"sIssueCode":"*",
	"sResultCode":"0",
	"sResultText":"",
	"sWarningCode":"0",
	"sWarningText":"",
	"sTokuteiGaisanHyoukagakuGoukei":"1933",
	"sIppanGaisanHyoukagakuGoukei":"0",
	"sNisaGaisanHyoukagakuGoukei":"0",
	"sTotalGaisanHyoukagakuGoukei":"1933",
	"sTokuteiGaisanHyoukaSonekiGoukei":"-77",
	"sIppanGaisanHyoukaSonekiGoukei":"0",
	"sNisaGaisanHyoukaSonekiGoukei":"0",
	"sTotalGaisanHyoukaSonekiGoukei":"-77",
	"aGenbutuKabuList":""
}`),
			wantResponse: stockPositionListResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 3, 1, 14, 57, 42, 983000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 3, 1, 14, 57, 42, 966000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockPositionList,
				},
				IssueCode:      "*",
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				SpecificAmount: 1933,
				GeneralAmount:  0,
				NisaAmount:     0,
				TotalAmount:    1933,
				SpecificProfit: -77,
				GeneralProfit:  0,
				NisaProfit:     0,
				TotalProfit:    -77,
				Positions:      []stockPosition{},
			}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var res stockPositionListResponse
			got1 := json.Unmarshal(test.arg1, &res)
			if !reflect.DeepEqual(test.wantResponse, res) || (got1 != nil) != test.hasError {
				t.Errorf("%s error\nwant: %+v\ngot: %+v, %+v\n", t.Name(), test.wantResponse, res, got1)
			}
		})
	}
}

func Test_stockPositionListResponse_response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		response stockPositionListResponse
		want1    StockPositionListResponse
	}{
		{name: "変換できる",
			response: stockPositionListResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 2, 28, 11, 11, 32, 66000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 2, 28, 11, 11, 32, 31000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockPositionList,
				},
				IssueCode:      "",
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				SpecificAmount: 1911,
				GeneralAmount:  0,
				NisaAmount:     0,
				TotalAmount:    1911,
				SpecificProfit: -78,
				GeneralProfit:  0,
				NisaProfit:     0,
				TotalProfit:    -78,
				Positions: []stockPosition{{
					WarningCode:        "0",
					WarningText:        "",
					IssueCode:          "1475",
					AccountType:        AccountTypeSpecific,
					OwnedQuantity:      1,
					UnHoldQuantity:     1,
					UnitValuation:      1911,
					TotalValuation:     1911,
					Profit:             -78,
					ProfitRatio:        -3.92,
					PrevClosePrice:     1914,
					PrevCloseRatio:     -3,
					PrevClosePercent:   -0.15,
					PrevCloseRatioType: PrevCloseRatioTypeUnder0,
					MarginBalance:      2226,
				}},
			},
			want1: StockPositionListResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 28, 11, 11, 32, 66000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 28, 11, 11, 32, 31000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockPositionList,
				},
				IssueCode:      "",
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				SpecificAmount: 1911,
				GeneralAmount:  0,
				NisaAmount:     0,
				TotalAmount:    1911,
				SpecificProfit: -78,
				GeneralProfit:  0,
				NisaProfit:     0,
				TotalProfit:    -78,
				Positions: []StockPosition{{
					WarningCode:        "0",
					WarningText:        "",
					IssueCode:          "1475",
					AccountType:        AccountTypeSpecific,
					OwnedQuantity:      1,
					UnHoldQuantity:     1,
					UnitValuation:      1911,
					TotalValuation:     1911,
					Profit:             -78,
					ProfitRatio:        -3.92,
					PrevClosePrice:     1914,
					PrevCloseRatio:     -3,
					PrevClosePercent:   -0.15,
					PrevCloseRatioType: PrevCloseRatioTypeUnder0,
					MarginBalance:      2226,
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

func Test_stockPosition_response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		response stockPosition
		want1    StockPosition
	}{
		{name: "変換できる",
			response: stockPosition{
				WarningCode:        "0",
				WarningText:        "",
				IssueCode:          "1475",
				AccountType:        AccountTypeSpecific,
				OwnedQuantity:      1,
				UnHoldQuantity:     1,
				UnitValuation:      1911,
				TotalValuation:     1911,
				Profit:             -78,
				ProfitRatio:        -3.92,
				PrevClosePrice:     1914,
				PrevCloseRatio:     -3,
				PrevClosePercent:   -0.15,
				PrevCloseRatioType: PrevCloseRatioTypeUnder0,
				MarginBalance:      2226,
			},
			want1: StockPosition{
				WarningCode:        "0",
				WarningText:        "",
				IssueCode:          "1475",
				AccountType:        AccountTypeSpecific,
				OwnedQuantity:      1,
				UnHoldQuantity:     1,
				UnitValuation:      1911,
				TotalValuation:     1911,
				Profit:             -78,
				ProfitRatio:        -3.92,
				PrevClosePrice:     1914,
				PrevCloseRatio:     -3,
				PrevClosePercent:   -0.15,
				PrevCloseRatioType: PrevCloseRatioTypeUnder0,
				MarginBalance:      2226,
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

func Test_client_StockPositionList(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name   string
		clock  *testClock
		status int
		body   []byte
		arg1   context.Context
		arg2   *Session
		arg3   StockPositionListRequest
		want1  *StockPositionListResponse
		want2  error
	}{
		{name: "正常レスポンスをパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 10, 9, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 53, 54, 58, 51, 53, 46, 55, 49, 55, 34, 44, 10, 9, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 10, 9, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 53, 54, 58, 51, 53, 46, 54, 57, 51, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 10, 9, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 71, 101, 110, 98, 117, 116, 117, 75, 97, 98, 117, 76, 105, 115, 116, 34, 44, 10, 9, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 103, 97, 107, 117, 71, 111, 117, 107, 101, 105, 34, 58, 34, 49, 57, 51, 50, 34, 44, 10, 9, 34, 115, 73, 112, 112, 97, 110, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 103, 97, 107, 117, 71, 111, 117, 107, 101, 105, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 78, 105, 115, 97, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 103, 97, 107, 117, 71, 111, 117, 107, 101, 105, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 84, 111, 116, 97, 108, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 103, 97, 107, 117, 71, 111, 117, 107, 101, 105, 34, 58, 34, 49, 57, 51, 50, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 45, 55, 56, 34, 44, 10, 9, 34, 115, 73, 112, 112, 97, 110, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 78, 105, 115, 97, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 84, 111, 116, 97, 108, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 45, 55, 56, 34, 44, 10, 9, 34, 97, 71, 101, 110, 98, 117, 116, 117, 75, 97, 98, 117, 76, 105, 115, 116, 34, 58, 10, 9, 91, 10, 9, 123, 10, 9, 9, 34, 115, 85, 114, 105, 79, 114, 100, 101, 114, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 85, 114, 105, 79, 114, 100, 101, 114, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 9, 34, 115, 85, 114, 105, 79, 114, 100, 101, 114, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 53, 34, 44, 10, 9, 9, 34, 115, 85, 114, 105, 79, 114, 100, 101, 114, 90, 121, 111, 117, 116, 111, 101, 107, 105, 75, 97, 122, 101, 105, 67, 34, 58, 34, 49, 34, 44, 10, 9, 9, 34, 115, 85, 114, 105, 79, 114, 100, 101, 114, 90, 97, 110, 75, 97, 98, 117, 83, 117, 114, 121, 111, 117, 34, 58, 34, 49, 34, 44, 10, 9, 9, 34, 115, 85, 114, 105, 79, 114, 100, 101, 114, 85, 114, 105, 116, 117, 107, 101, 75, 97, 110, 111, 117, 83, 117, 114, 121, 111, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 85, 114, 105, 79, 114, 100, 101, 114, 71, 97, 105, 115, 97, 110, 66, 111, 107, 97, 84, 97, 110, 107, 97, 34, 58, 34, 50, 48, 49, 48, 46, 48, 48, 48, 48, 34, 44, 10, 9, 9, 34, 115, 85, 114, 105, 79, 114, 100, 101, 114, 72, 121, 111, 117, 107, 97, 84, 97, 110, 107, 97, 34, 58, 34, 49, 57, 51, 50, 46, 48, 48, 48, 48, 34, 44, 10, 9, 9, 34, 115, 85, 114, 105, 79, 114, 100, 101, 114, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 103, 97, 107, 117, 34, 58, 34, 49, 57, 51, 50, 34, 44, 10, 9, 9, 34, 115, 85, 114, 105, 79, 114, 100, 101, 114, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 34, 58, 34, 45, 55, 56, 34, 44, 10, 9, 9, 34, 115, 85, 114, 105, 79, 114, 100, 101, 114, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 82, 105, 116, 117, 34, 58, 34, 45, 51, 46, 56, 56, 34, 44, 10, 9, 9, 34, 115, 83, 121, 117, 122, 105, 116, 117, 79, 119, 97, 114, 105, 110, 101, 34, 58, 34, 49, 57, 50, 51, 34, 44, 10, 9, 9, 34, 115, 90, 101, 110, 122, 105, 116, 117, 72, 105, 34, 58, 34, 57, 34, 44, 10, 9, 9, 34, 115, 90, 101, 110, 122, 105, 116, 117, 72, 105, 80, 101, 114, 34, 58, 34, 48, 46, 52, 54, 34, 44, 10, 9, 9, 34, 115, 85, 112, 68, 111, 119, 110, 70, 108, 97, 103, 34, 58, 34, 48, 53, 34, 44, 10, 9, 9, 34, 115, 78, 105, 115, 115, 121, 111, 117, 107, 105, 110, 75, 97, 115, 105, 107, 97, 98, 117, 90, 97, 110, 34, 58, 34, 50, 52, 48, 57, 34, 10, 9, 125, 10, 9, 93, 10, 125, 10, 10},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   StockPositionListRequest{},
			want1: &StockPositionListResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 1, 14, 56, 35, 717000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 1, 14, 56, 35, 693000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockPositionList,
				},
				IssueCode:      "",
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				SpecificAmount: 1932,
				GeneralAmount:  0,
				NisaAmount:     0,
				TotalAmount:    1932,
				SpecificProfit: -78,
				GeneralProfit:  0,
				NisaProfit:     0,
				TotalProfit:    -78,
				Positions: []StockPosition{{
					WarningCode:        "0",
					WarningText:        "",
					IssueCode:          "1475",
					AccountType:        AccountTypeSpecific,
					OwnedQuantity:      1,
					UnHoldQuantity:     0,
					UnitValuation:      2010,
					BookValuation:      1932,
					TotalValuation:     1932,
					Profit:             -78,
					ProfitRatio:        -3.88,
					PrevClosePrice:     1923,
					PrevCloseRatio:     9,
					PrevClosePercent:   0.46,
					PrevCloseRatioType: PrevCloseRatioTypeOver0,
					MarginBalance:      2409,
				}},
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   nil,
			arg3:   StockPositionListRequest{},
			want1:  nil,
			want2:  NilArgumentErr},
		{name: "失敗をパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 10, 9, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 53, 55, 58, 52, 50, 46, 57, 56, 51, 34, 44, 10, 9, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 10, 9, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 53, 55, 58, 52, 50, 46, 57, 54, 54, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 10, 9, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 71, 101, 110, 98, 117, 116, 117, 75, 97, 98, 117, 76, 105, 115, 116, 34, 44, 10, 9, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 42, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 103, 97, 107, 117, 71, 111, 117, 107, 101, 105, 34, 58, 34, 49, 57, 51, 51, 34, 44, 10, 9, 34, 115, 73, 112, 112, 97, 110, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 103, 97, 107, 117, 71, 111, 117, 107, 101, 105, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 78, 105, 115, 97, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 103, 97, 107, 117, 71, 111, 117, 107, 101, 105, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 84, 111, 116, 97, 108, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 103, 97, 107, 117, 71, 111, 117, 107, 101, 105, 34, 58, 34, 49, 57, 51, 51, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 45, 55, 55, 34, 44, 10, 9, 34, 115, 73, 112, 112, 97, 110, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 78, 105, 115, 97, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 84, 111, 116, 97, 108, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 45, 55, 55, 34, 44, 10, 9, 34, 97, 71, 101, 110, 98, 117, 116, 117, 75, 97, 98, 117, 76, 105, 115, 116, 34, 58, 34, 34, 10, 125, 10, 10},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   StockPositionListRequest{},
			want1: &StockPositionListResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 1, 14, 57, 42, 983000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 1, 14, 57, 42, 966000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockPositionList,
				},
				IssueCode:      "*",
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				SpecificAmount: 1933,
				GeneralAmount:  0,
				NisaAmount:     0,
				TotalAmount:    1933,
				SpecificProfit: -77,
				GeneralProfit:  0,
				NisaProfit:     0,
				TotalProfit:    -77,
				Positions:      []StockPosition{},
			},
			want2: nil},
		{name: "200 OK以外が返ったらエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusInternalServerError,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   StockPositionListRequest{},
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
			got1, got2 := client.StockPositionList(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_StockPositionList_Execute(t *testing.T) {
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

	got3, got4 := client.StockPositionList(context.Background(), session, StockPositionListRequest{
		IssueCode: "",
	})
	log.Printf("%+v, %+v\n", got3, got4)
}
