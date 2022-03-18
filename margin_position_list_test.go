package tachibana

import (
	"context"
	"encoding/json"
	"errors"
	"log"
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
				IssueCode: "1475",
			},
			arg1: 123,
			arg2: time.Date(2022, 3, 1, 9, 40, 0, 0, time.Local),
			want1: marginPositionListRequest{
				commonRequest: commonRequest{
					No:             123,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 1, 9, 40, 0, 0, time.Local)},
					FeatureType:    FeatureTypeMarginPositionList,
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

func Test_marginPositionListResponse_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		arg1         []byte
		wantResponse marginPositionListResponse
		hasError     bool
	}{
		{name: "正常系のレスポンスをパースできる",
			arg1: []byte(`{
	"p_sd_date":"2022.03.01-14:56:53.904",
	"p_no":"2",
	"p_rv_date":"2022.03.01-14:56:53.886",
	"p_errno":"0",
	"p_err":"",
	"sCLMID":"CLMShinyouTategyokuList",
	"sIssueCode":"",
	"sResultCode":"0",
	"sResultText":"",
	"sWarningCode":"0",
	"sWarningText":"",
	"sUritateDaikin":"0",
	"sKaitateDaikin":"1933",
	"sTotalDaikin":"1933",
	"sHyoukaSonekiGoukeiUridate":"0",
	"sHyoukaSonekiGoukeiKaidate":"-1",
	"sTotalHyoukaSonekiGoukei":"-1",
	"sTokuteiHyoukaSonekiGoukei":"-1",
	"sIppanHyoukaSonekiGoukei":"0",
	"aShinyouTategyokuList":
	[
	{
		"sOrderWarningCode":"0",
		"sOrderWarningText":"",
		"sOrderTategyokuNumber":"202203010010437",
		"sOrderIssueCode":"1475",
		"sOrderSizyouC":"00",
		"sOrderBaibaiKubun":"3",
		"sOrderBensaiKubun":"26",
		"sOrderZyoutoekiKazeiC":"1",
		"sOrderTategyokuSuryou":"1",
		"sOrderTategyokuTanka":"1933.0000",
		"sOrderHyoukaTanka":"1932.0000",
		"sOrderGaisanHyoukaSoneki":"-1",
		"sOrderGaisanHyoukaSonekiRitu":"-0.05",
		"sTategyokuDaikin":"1933",
		"sOrderTateTesuryou":"0",
		"sOrderZyunHibu":"0",
		"sOrderGyakuhibu":"0",
		"sOrderKakikaeryou":"0",
		"sOrderKanrihi":"0",
		"sOrderKasikaburyou":"0",
		"sOrderSonota":"0",
		"sOrderTategyokuDay":"20220301",
		"sOrderTategyokuKizituDay":"20220831",
		"sTategyokuSuryou":"1",
		"sOrderYakuzyouHensaiKabusu":"0",
		"sOrderGenbikiGenwatasiKabusu":"0",
		"sOrderOrderSuryou":"1",
		"sOrderHensaiKanouSuryou":"0",
		"sSyuzituOwarine":"1923",
		"sZenzituHi":"9",
		"sZenzituHiPer":"0.46",
		"sUpDownFlag":"05"
	}
	]
}`),
			wantResponse: marginPositionListResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 3, 1, 14, 56, 53, 904000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 3, 1, 14, 56, 53, 886000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMarginPositionList,
				},
				IssueCode:             "",
				ResultCode:            "0",
				ResultText:            "",
				WarningCode:           "0",
				WarningText:           "",
				TotalSellAmount:       0,
				TotalBuyAmount:        1933,
				TotalAmount:           1933,
				TotalSellProfit:       0,
				TotalBuyProfit:        -1,
				TotalProfit:           -1,
				SpecificAccountProfit: -1,
				GeneralAccountProfit:  0,
				Positions: []marginPosition{{
					WarningCode:        "0",
					WarningText:        "",
					PositionNumber:     "202203010010437",
					IssueCode:          "1475",
					Exchange:           ExchangeToushou,
					Side:               SideBuy,
					ExitTermType:       ExitTermTypeSystemMargin6m,
					AccountType:        AccountTypeSpecific,
					OrderQuantity:      1,
					UnitPrice:          1933,
					CurrentPrice:       1932,
					Profit:             -1,
					ProfitRatio:        -0.05,
					TotalPrice:         1933,
					Commission:         0,
					Interest:           0,
					Premiums:           0,
					RewritingFee:       0,
					ManagementFee:      0,
					LendingFee:         0,
					OtherFee:           0,
					ContractDate:       Ymd{Time: time.Date(2022, 3, 1, 0, 0, 0, 0, time.Local)},
					ExitTerm:           Ymd{Time: time.Date(2022, 8, 31, 0, 0, 0, 0, time.Local)},
					OwnedQuantity:      1,
					ExitQuantity:       0,
					DeliveryQuantity:   0,
					HoldQuantity:       1,
					ReturnableQuantity: 0,
					PrevClosePrice:     1923,
					PrevCloseRatio:     9,
					PrevClosePercent:   0.46,
					PrevCloseRatioType: PrevCloseRatioTypeOver0,
				}},
			}},
		{name: "異常系のレスポンスをパースできる",
			arg1: []byte(`{
	"p_sd_date":"2022.03.01-14:58:00.105",
	"p_no":"2",
	"p_rv_date":"2022.03.01-14:58:00.083",
	"p_errno":"0",
	"p_err":"",
	"sCLMID":"CLMShinyouTategyokuList",
	"sIssueCode":"*",
	"sResultCode":"0",
	"sResultText":"",
	"sWarningCode":"0",
	"sWarningText":"",
	"sUritateDaikin":"0",
	"sKaitateDaikin":"1933",
	"sTotalDaikin":"1933",
	"sHyoukaSonekiGoukeiUridate":"0",
	"sHyoukaSonekiGoukeiKaidate":"1",
	"sTotalHyoukaSonekiGoukei":"1",
	"sTokuteiHyoukaSonekiGoukei":"1",
	"sIppanHyoukaSonekiGoukei":"0",
	"aShinyouTategyokuList":""
}`),
			wantResponse: marginPositionListResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 3, 1, 14, 58, 0, 105000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 3, 1, 14, 58, 0, 83000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMarginPositionList,
				},
				IssueCode:             "*",
				ResultCode:            "0",
				ResultText:            "",
				WarningCode:           "0",
				WarningText:           "",
				TotalSellAmount:       0,
				TotalBuyAmount:        1933,
				TotalAmount:           1933,
				TotalSellProfit:       0,
				TotalBuyProfit:        1,
				TotalProfit:           1,
				SpecificAccountProfit: 1,
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
				IssueCode:             "",
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
					PositionNumber:     "202202280006156",
					IssueCode:          "1475",
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
				IssueCode:             "",
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
					PositionNumber:     "202202280006156",
					IssueCode:          "1475",
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
				PositionNumber:     "202202280006156",
				IssueCode:          "1475",
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
				PositionNumber:     "202202280006156",
				IssueCode:          "1475",
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
		name      string
		clock     *testClock
		requester *testRequester
		arg1      context.Context
		arg2      *Session
		arg3      MarginPositionListRequest
		want1     *MarginPositionListResponse
		want2     error
	}{
		{name: "正常レスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get1: []byte{123, 10, 9, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 53, 54, 58, 53, 51, 46, 57, 48, 52, 34, 44, 10, 9, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 10, 9, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 53, 54, 58, 53, 51, 46, 56, 56, 54, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 10, 9, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 83, 104, 105, 110, 121, 111, 117, 84, 97, 116, 101, 103, 121, 111, 107, 117, 76, 105, 115, 116, 34, 44, 10, 9, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 85, 114, 105, 116, 97, 116, 101, 68, 97, 105, 107, 105, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 75, 97, 105, 116, 97, 116, 101, 68, 97, 105, 107, 105, 110, 34, 58, 34, 49, 57, 51, 51, 34, 44, 10, 9, 34, 115, 84, 111, 116, 97, 108, 68, 97, 105, 107, 105, 110, 34, 58, 34, 49, 57, 51, 51, 34, 44, 10, 9, 34, 115, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 85, 114, 105, 100, 97, 116, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 75, 97, 105, 100, 97, 116, 101, 34, 58, 34, 45, 49, 34, 44, 10, 9, 34, 115, 84, 111, 116, 97, 108, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 45, 49, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 45, 49, 34, 44, 10, 9, 34, 115, 73, 112, 112, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 48, 34, 44, 10, 9, 34, 97, 83, 104, 105, 110, 121, 111, 117, 84, 97, 116, 101, 103, 121, 111, 107, 117, 76, 105, 115, 116, 34, 58, 10, 9, 91, 10, 9, 123, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 84, 97, 116, 101, 103, 121, 111, 107, 117, 78, 117, 109, 98, 101, 114, 34, 58, 34, 50, 48, 50, 50, 48, 51, 48, 49, 48, 48, 49, 48, 52, 51, 55, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 53, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 83, 105, 122, 121, 111, 117, 67, 34, 58, 34, 48, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 66, 97, 105, 98, 97, 105, 75, 117, 98, 117, 110, 34, 58, 34, 51, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 66, 101, 110, 115, 97, 105, 75, 117, 98, 117, 110, 34, 58, 34, 50, 54, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 90, 121, 111, 117, 116, 111, 101, 107, 105, 75, 97, 122, 101, 105, 67, 34, 58, 34, 49, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 84, 97, 116, 101, 103, 121, 111, 107, 117, 83, 117, 114, 121, 111, 117, 34, 58, 34, 49, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 84, 97, 116, 101, 103, 121, 111, 107, 117, 84, 97, 110, 107, 97, 34, 58, 34, 49, 57, 51, 51, 46, 48, 48, 48, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 72, 121, 111, 117, 107, 97, 84, 97, 110, 107, 97, 34, 58, 34, 49, 57, 51, 50, 46, 48, 48, 48, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 34, 58, 34, 45, 49, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 71, 97, 105, 115, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 82, 105, 116, 117, 34, 58, 34, 45, 48, 46, 48, 53, 34, 44, 10, 9, 9, 34, 115, 84, 97, 116, 101, 103, 121, 111, 107, 117, 68, 97, 105, 107, 105, 110, 34, 58, 34, 49, 57, 51, 51, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 84, 97, 116, 101, 84, 101, 115, 117, 114, 121, 111, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 90, 121, 117, 110, 72, 105, 98, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 71, 121, 97, 107, 117, 104, 105, 98, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 75, 97, 107, 105, 107, 97, 101, 114, 121, 111, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 75, 97, 110, 114, 105, 104, 105, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 75, 97, 115, 105, 107, 97, 98, 117, 114, 121, 111, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 83, 111, 110, 111, 116, 97, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 84, 97, 116, 101, 103, 121, 111, 107, 117, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 48, 49, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 84, 97, 116, 101, 103, 121, 111, 107, 117, 75, 105, 122, 105, 116, 117, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 56, 51, 49, 34, 44, 10, 9, 9, 34, 115, 84, 97, 116, 101, 103, 121, 111, 107, 117, 83, 117, 114, 121, 111, 117, 34, 58, 34, 49, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 89, 97, 107, 117, 122, 121, 111, 117, 72, 101, 110, 115, 97, 105, 75, 97, 98, 117, 115, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 71, 101, 110, 98, 105, 107, 105, 71, 101, 110, 119, 97, 116, 97, 115, 105, 75, 97, 98, 117, 115, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 79, 114, 100, 101, 114, 83, 117, 114, 121, 111, 117, 34, 58, 34, 49, 34, 44, 10, 9, 9, 34, 115, 79, 114, 100, 101, 114, 72, 101, 110, 115, 97, 105, 75, 97, 110, 111, 117, 83, 117, 114, 121, 111, 117, 34, 58, 34, 48, 34, 44, 10, 9, 9, 34, 115, 83, 121, 117, 122, 105, 116, 117, 79, 119, 97, 114, 105, 110, 101, 34, 58, 34, 49, 57, 50, 51, 34, 44, 10, 9, 9, 34, 115, 90, 101, 110, 122, 105, 116, 117, 72, 105, 34, 58, 34, 57, 34, 44, 10, 9, 9, 34, 115, 90, 101, 110, 122, 105, 116, 117, 72, 105, 80, 101, 114, 34, 58, 34, 48, 46, 52, 54, 34, 44, 10, 9, 9, 34, 115, 85, 112, 68, 111, 119, 110, 70, 108, 97, 103, 34, 58, 34, 48, 53, 34, 10, 9, 125, 10, 9, 93, 10, 125, 10, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      MarginPositionListRequest{},
			want1: &MarginPositionListResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 1, 14, 56, 53, 904000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 1, 14, 56, 53, 886000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMarginPositionList,
				},
				IssueCode:             "",
				ResultCode:            "0",
				ResultText:            "",
				WarningCode:           "0",
				WarningText:           "",
				TotalSellAmount:       0,
				TotalBuyAmount:        1933,
				TotalAmount:           1933,
				TotalSellProfit:       0,
				TotalBuyProfit:        -1,
				TotalProfit:           -1,
				SpecificAccountProfit: -1,
				GeneralAccountProfit:  0,
				Positions: []MarginPosition{{
					WarningCode:        "0",
					WarningText:        "",
					PositionNumber:     "202203010010437",
					IssueCode:          "1475",
					Exchange:           ExchangeToushou,
					Side:               SideBuy,
					ExitTermType:       ExitTermTypeSystemMargin6m,
					AccountType:        AccountTypeSpecific,
					OrderQuantity:      1,
					UnitPrice:          1933,
					CurrentPrice:       1932,
					Profit:             -1,
					ProfitRatio:        -0.05,
					TotalPrice:         1933,
					Commission:         0,
					Interest:           0,
					Premiums:           0,
					RewritingFee:       0,
					ManagementFee:      0,
					LendingFee:         0,
					OtherFee:           0,
					ContractDate:       time.Date(2022, 3, 1, 0, 0, 0, 0, time.Local),
					ExitTerm:           time.Date(2022, 8, 31, 0, 0, 0, 0, time.Local),
					OwnedQuantity:      1,
					ExitQuantity:       0,
					DeliveryQuantity:   0,
					HoldQuantity:       1,
					ReturnableQuantity: 0,
					PrevClosePrice:     1923,
					PrevCloseRatio:     9,
					PrevClosePercent:   0.46,
					PrevCloseRatioType: PrevCloseRatioTypeOver0,
				}},
			},
			want2: nil},
		{name: "失敗をパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			requester: &testRequester{get1: []byte{123, 10, 9, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 53, 56, 58, 48, 48, 46, 49, 48, 53, 34, 44, 10, 9, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 10, 9, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 53, 56, 58, 48, 48, 46, 48, 56, 51, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 10, 9, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 83, 104, 105, 110, 121, 111, 117, 84, 97, 116, 101, 103, 121, 111, 107, 117, 76, 105, 115, 116, 34, 44, 10, 9, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 42, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 87, 97, 114, 110, 105, 110, 103, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 87, 97, 114, 110, 105, 110, 103, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 85, 114, 105, 116, 97, 116, 101, 68, 97, 105, 107, 105, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 75, 97, 105, 116, 97, 116, 101, 68, 97, 105, 107, 105, 110, 34, 58, 34, 49, 57, 51, 51, 34, 44, 10, 9, 34, 115, 84, 111, 116, 97, 108, 68, 97, 105, 107, 105, 110, 34, 58, 34, 49, 57, 51, 51, 34, 44, 10, 9, 34, 115, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 85, 114, 105, 100, 97, 116, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 75, 97, 105, 100, 97, 116, 101, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 84, 111, 116, 97, 108, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 73, 112, 112, 97, 110, 72, 121, 111, 117, 107, 97, 83, 111, 110, 101, 107, 105, 71, 111, 117, 107, 101, 105, 34, 58, 34, 48, 34, 44, 10, 9, 34, 97, 83, 104, 105, 110, 121, 111, 117, 84, 97, 116, 101, 103, 121, 111, 107, 117, 76, 105, 115, 116, 34, 58, 34, 34, 10, 125, 10, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      MarginPositionListRequest{},
			want1: &MarginPositionListResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 1, 14, 58, 0, 105000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 1, 14, 58, 0, 83000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMarginPositionList,
				},
				IssueCode:             "*",
				ResultCode:            "0",
				ResultText:            "",
				WarningCode:           "0",
				WarningText:           "",
				TotalSellAmount:       0,
				TotalBuyAmount:        1933,
				TotalAmount:           1933,
				TotalSellProfit:       0,
				TotalBuyProfit:        1,
				TotalProfit:           1,
				SpecificAccountProfit: 1,
				GeneralAccountProfit:  0,
				Positions:             []MarginPosition{},
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			arg1:  context.Background(),
			arg2:  nil,
			arg3:  MarginPositionListRequest{},
			want1: nil,
			want2: NilArgumentErr},
		{name: "リクエストでエラーが返されたらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get2: StatusNotOkErr},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      MarginPositionListRequest{},
			want1:     nil,
			want2:     StatusNotOkErr},
		{name: "レスポンスのUnmarshalに失敗したらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get1: []byte{}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      MarginPositionListRequest{},
			want1:     nil,
			want2:     UnmarshalFailedErr},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &client{clock: test.clock, requester: test.requester}
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
		IssueCode: "",
	})
	log.Printf("%+v, %+v\n", got3, got4)
}
