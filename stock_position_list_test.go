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
			request: StockPositionListRequest{SymbolCode: "1475"},
			arg1:    123,
			arg2:    time.Date(2022, 3, 1, 9, 0, 0, 0, time.Local),
			want1: stockPositionListRequest{
				commonRequest: commonRequest{
					No:          123,
					SendDate:    RequestTime{Time: time.Date(2022, 3, 1, 9, 0, 0, 0, time.Local)},
					FeatureType: FeatureTypeStockPositionList,
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

func Test_stockPositionListResponse_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		arg1         []byte
		wantResponse stockPositionListResponse
		hasError     bool
	}{
		{name: "正常系のパース",
			arg1: []byte(`{"177":"2022.02.28-11:11:32.066","175":"2","176":"2022.02.28-11:11:32.031","174":"0","173":"","192":"CLMGenbutuKabuList","328":"","534":"0","535":"","692":"0","693":"","641":"1911","313":"0","429":"0","651":"1911","640":"-78","312":"0","428":"0","650":"-78","49":[{"681":"0","682":"","679":"1475","684":"1","683":"1","680":"1","674":"1989.0000","678":"1911.0000","677":"1911","675":"-78","676":"-3.92","600":"1914","735":"-3","736":"-0.15","668":"07","431":"2226"}]}`),
			wantResponse: stockPositionListResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 2, 28, 11, 11, 32, 66000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 2, 28, 11, 11, 32, 31000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockPositionList,
				},
				SymbolCode:     "",
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
					SymbolCode:         "1475",
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
		{name: "空配列でもパースできる",
			arg1: []byte(`{"177":"2022.02.28-10:09:27.637","175":"2","176":"2022.02.28-10:09:27.605","174":"0","173":"","192":"CLMGenbutuKabuList","328":"","534":"0","535":"","692":"0","693":"","641":"0","313":"0","429":"0","651":"0","640":"0","312":"0","428":"0","650":"0","49":""}`),
			wantResponse: stockPositionListResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 2, 28, 10, 9, 27, 637000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 2, 28, 10, 9, 27, 605000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockPositionList,
				},
				SymbolCode:     "",
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				SpecificAmount: 0,
				GeneralAmount:  0,
				NisaAmount:     0,
				TotalAmount:    0,
				SpecificProfit: 0,
				GeneralProfit:  0,
				NisaProfit:     0,
				TotalProfit:    0,
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
				SymbolCode:     "",
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
					SymbolCode:         "1475",
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
				SymbolCode:     "",
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
					SymbolCode:         "1475",
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
				SymbolCode:         "1475",
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
				SymbolCode:         "1475",
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
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 49, 49, 58, 49, 50, 58, 49, 55, 46, 51, 52, 53, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 49, 49, 58, 49, 50, 58, 49, 55, 46, 51, 48, 54, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 71, 101, 110, 98, 117, 116, 117, 75, 97, 98, 117, 76, 105, 115, 116, 34, 44, 34, 51, 50, 56, 34, 58, 34, 34, 44, 34, 53, 51, 52, 34, 58, 34, 48, 34, 44, 34, 53, 51, 53, 34, 58, 34, 34, 44, 34, 54, 57, 50, 34, 58, 34, 48, 34, 44, 34, 54, 57, 51, 34, 58, 34, 34, 44, 34, 54, 52, 49, 34, 58, 34, 49, 57, 48, 57, 34, 44, 34, 51, 49, 51, 34, 58, 34, 48, 34, 44, 34, 52, 50, 57, 34, 58, 34, 48, 34, 44, 34, 54, 53, 49, 34, 58, 34, 49, 57, 48, 57, 34, 44, 34, 54, 52, 48, 34, 58, 34, 45, 56, 48, 34, 44, 34, 51, 49, 50, 34, 58, 34, 48, 34, 44, 34, 52, 50, 56, 34, 58, 34, 48, 34, 44, 34, 54, 53, 48, 34, 58, 34, 45, 56, 48, 34, 44, 34, 52, 57, 34, 58, 91, 123, 34, 54, 56, 49, 34, 58, 34, 48, 34, 44, 34, 54, 56, 50, 34, 58, 34, 34, 44, 34, 54, 55, 57, 34, 58, 34, 49, 52, 55, 53, 34, 44, 34, 54, 56, 52, 34, 58, 34, 49, 34, 44, 34, 54, 56, 51, 34, 58, 34, 49, 34, 44, 34, 54, 56, 48, 34, 58, 34, 48, 34, 44, 34, 54, 55, 52, 34, 58, 34, 49, 57, 56, 57, 46, 48, 48, 48, 48, 34, 44, 34, 54, 55, 56, 34, 58, 34, 49, 57, 48, 57, 46, 48, 48, 48, 48, 34, 44, 34, 54, 55, 55, 34, 58, 34, 49, 57, 48, 57, 34, 44, 34, 54, 55, 53, 34, 58, 34, 45, 56, 48, 34, 44, 34, 54, 55, 54, 34, 58, 34, 45, 52, 46, 48, 50, 34, 44, 34, 54, 48, 48, 34, 58, 34, 49, 57, 49, 52, 34, 44, 34, 55, 51, 53, 34, 58, 34, 45, 53, 34, 44, 34, 55, 51, 54, 34, 58, 34, 45, 48, 46, 50, 54, 34, 44, 34, 54, 54, 56, 34, 58, 34, 48, 55, 34, 44, 34, 52, 51, 49, 34, 58, 34, 50, 50, 50, 54, 34, 125, 93, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   StockPositionListRequest{},
			want1: &StockPositionListResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 28, 11, 12, 17, 345000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 28, 11, 12, 17, 306000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockPositionList,
				},
				SymbolCode:     "",
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				SpecificAmount: 1909,
				GeneralAmount:  0,
				NisaAmount:     0,
				TotalAmount:    1909,
				SpecificProfit: -80,
				GeneralProfit:  0,
				NisaProfit:     0,
				TotalProfit:    -80,
				Positions: []StockPosition{
					{
						WarningCode:        "0",
						WarningText:        "",
						SymbolCode:         "1475",
						AccountType:        AccountTypeSpecific,
						OwnedQuantity:      1,
						UnHoldQuantity:     0,
						UnitValuation:      1909,
						TotalValuation:     1909,
						Profit:             -80,
						ProfitRatio:        -4.02,
						PrevClosePrice:     1914,
						PrevCloseRatio:     -5,
						PrevClosePercent:   -0.26,
						PrevCloseRatioType: PrevCloseRatioTypeUnder0,
						MarginBalance:      2226,
					},
				},
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
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 49, 48, 58, 48, 57, 58, 53, 55, 46, 48, 50, 49, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 49, 48, 58, 48, 57, 58, 53, 54, 46, 57, 56, 54, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 71, 101, 110, 98, 117, 116, 117, 75, 97, 98, 117, 76, 105, 115, 116, 34, 44, 34, 51, 50, 56, 34, 58, 34, 42, 34, 44, 34, 53, 51, 52, 34, 58, 34, 48, 34, 44, 34, 53, 51, 53, 34, 58, 34, 34, 44, 34, 54, 57, 50, 34, 58, 34, 48, 34, 44, 34, 54, 57, 51, 34, 58, 34, 34, 44, 34, 54, 52, 49, 34, 58, 34, 48, 34, 44, 34, 51, 49, 51, 34, 58, 34, 48, 34, 44, 34, 52, 50, 57, 34, 58, 34, 48, 34, 44, 34, 54, 53, 49, 34, 58, 34, 48, 34, 44, 34, 54, 52, 48, 34, 58, 34, 48, 34, 44, 34, 51, 49, 50, 34, 58, 34, 48, 34, 44, 34, 52, 50, 56, 34, 58, 34, 48, 34, 44, 34, 54, 53, 48, 34, 58, 34, 48, 34, 44, 34, 52, 57, 34, 58, 34, 34, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   StockPositionListRequest{},
			want1: &StockPositionListResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 28, 10, 9, 57, 21000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 28, 10, 9, 56, 986000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeStockPositionList,
				},
				SymbolCode:     "*",
				ResultCode:     "0",
				ResultText:     "",
				WarningCode:    "0",
				WarningText:    "",
				SpecificAmount: 0,
				GeneralAmount:  0,
				NisaAmount:     0,
				TotalAmount:    0,
				SpecificProfit: 0,
				GeneralProfit:  0,
				NisaProfit:     0,
				TotalProfit:    0,
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
		SymbolCode: "",
	})
	log.Printf("%+v, %+v\n", got3, got4)
}
