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

func Test_OrderListRequest_request(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request OrderListRequest
		arg1    int64
		arg2    time.Time
		want1   orderListRequest
	}{
		{name: "変換できる",
			request: OrderListRequest{
				SymbolCode:         "1475",
				ExecutionDate:      time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local),
				OrderInquiryStatus: OrderInquiryStatusInOrder,
			},
			arg1: 123,
			arg2: time.Date(2022, 2, 28, 8, 5, 30, 123000000, time.Local),
			want1: orderListRequest{
				commonRequest: commonRequest{
					No:          123,
					SendDate:    RequestTime{Time: time.Date(2022, 2, 28, 8, 5, 30, 123000000, time.Local)},
					FeatureType: FeatureTypeOrderList,
				},
				SymbolCode:    "1475",
				ExecutionDate: Ymd{Time: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local)},
				OrderStatus:   OrderInquiryStatusInOrder,
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

func Test_orderListResponse_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name         string
		arg1         []byte
		wantResponse orderListResponse
		hasError     bool
	}{
		{name: "注文一覧がある場合のパースができる",
			arg1: []byte(`{"177":"2022.02.26-08:01:09.476","175":"2","176":"2022.02.26-08:01:09.458","174":"0","173":"","192":"CLMOrderList","328":"1475","559":"20220228","508":"","534":"0","535":"","692":"0","693":"","55":[{"521":"0","522":"","493":"28002177","485":"1475","501":"00","528":"1","255":"0","468":"00","467":"3","496":"1","471":"1","494":"0.0000","469":"0","495":"1","480":"0","482":"0.0000","479":" ","481":"0.0000","517":"0","510":" ","527":"","526":"0","524":"0.0000","520":" ","500":"20220228","504":"0","503":"受付未済","525":"0","491":"20220226074344","492":"00000000","489":"0","470":"0","235":"2391"}]}`),
			wantResponse: orderListResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 2, 26, 8, 1, 9, 476000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 2, 26, 8, 1, 9, 458000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderList,
				},
				SymbolCode:         "1475",
				ExecutionDate:      Ymd{Time: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local)},
				OrderInquiryStatus: OrderInquiryStatusUnspecified,
				ResultCode:         "0",
				ResultText:         "",
				WarningCode:        "0",
				WarningText:        "",
				Orders: []order{{
					WarningCode:            "0",
					WarningText:            "",
					OrderNumber:            "28002177",
					SymbolCode:             "1475",
					Exchange:               ExchangeToushou,
					AccountType:            AccountTypeSpecific,
					TradeType:              TradeTypeStock,
					ExitTermType:           ExitTermTypeNoLimit,
					Side:                   SideBuy,
					OrderQuantity:          1,
					CurrentQuantity:        1,
					Price:                  0,
					ExecutionTiming:        ExecutionTimingNormal,
					ExecutionType:          ExecutionTypeMarket,
					StopOrderType:          StopOrderTypeNormal,
					StopTriggerPrice:       0,
					StopOrderExecutionType: ExecutionTypeUnused,
					StopOrderPrice:         0,
					TriggerType:            TriggerTypeNoFired,
					ExitOrderType:          ExitOrderTypeUnused,
					ContractQuantity:       0,
					ContractPrice:          0,
					PartContractType:       PartContractTypeUnused,
					ExecutionDate:          Ymd{Time: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local)},
					OrderStatus:            OrderStatusReceived,
					OrderStatusText:        "受付未済",
					ContractStatus:         ContractStatusInOrder,
					OrderDateTime:          YmdHms{Time: time.Date(2022, 2, 26, 7, 43, 44, 0, time.Local)},
					ExpireDate:             Ymd{},
					CarryOverType:          CarryOverTypeToday,
					CorrectCancelType:      CorrectCancelTypeCorrectable,
					EstimationAmount:       2391,
				}},
			}},
		{name: "注文一覧が空文字の場合でもパース出来る",
			arg1: []byte(`{"177":"2022.02.26-08:01:09.476","175":"2","176":"2022.02.26-08:01:09.458","174":"0","173":"","192":"CLMOrderList","328":"1475","559":"20220228","508":"","534":"0","535":"","692":"0","693":"","55":""}`),
			wantResponse: orderListResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 2, 26, 8, 1, 9, 476000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 2, 26, 8, 1, 9, 458000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderList,
				},
				SymbolCode:         "1475",
				ExecutionDate:      Ymd{Time: time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local)},
				OrderInquiryStatus: OrderInquiryStatusUnspecified,
				ResultCode:         "0",
				ResultText:         "",
				WarningCode:        "0",
				WarningText:        "",
				Orders:             []order{},
			}},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var res orderListResponse
			got1 := json.Unmarshal(test.arg1, &res)
			if !reflect.DeepEqual(test.wantResponse, res) || (got1 != nil) != test.hasError {
				t.Errorf("%s error\nwant: %+v\ngot: %+v, %+v\n", t.Name(), test.wantResponse, res, got1)
			}
		})
	}
}

func Test_client_OrderList(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		clock  *testClock
		status int
		body   []byte
		arg1   context.Context
		arg2   *Session
		arg3   OrderListRequest
		want1  *OrderListResponse
		want2  error
	}{
		{name: "正常レスポンスをパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 55, 45, 48, 55, 58, 49, 56, 58, 52, 49, 46, 49, 52, 54, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 55, 45, 48, 55, 58, 49, 56, 58, 52, 49, 46, 49, 49, 49, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 79, 114, 100, 101, 114, 76, 105, 115, 116, 34, 44, 34, 51, 50, 56, 34, 58, 34, 34, 44, 34, 53, 53, 57, 34, 58, 34, 34, 44, 34, 53, 48, 56, 34, 58, 34, 34, 44, 34, 53, 51, 52, 34, 58, 34, 48, 34, 44, 34, 53, 51, 53, 34, 58, 34, 34, 44, 34, 54, 57, 50, 34, 58, 34, 48, 34, 44, 34, 54, 57, 51, 34, 58, 34, 34, 44, 34, 53, 53, 34, 58, 91, 123, 34, 53, 50, 49, 34, 58, 34, 48, 34, 44, 34, 53, 50, 50, 34, 58, 34, 34, 44, 34, 52, 57, 51, 34, 58, 34, 50, 56, 48, 48, 50, 49, 55, 55, 34, 44, 34, 52, 56, 53, 34, 58, 34, 49, 52, 55, 53, 34, 44, 34, 53, 48, 49, 34, 58, 34, 48, 48, 34, 44, 34, 53, 50, 56, 34, 58, 34, 49, 34, 44, 34, 50, 53, 53, 34, 58, 34, 48, 34, 44, 34, 52, 54, 56, 34, 58, 34, 48, 48, 34, 44, 34, 52, 54, 55, 34, 58, 34, 51, 34, 44, 34, 52, 57, 54, 34, 58, 34, 49, 34, 44, 34, 52, 55, 49, 34, 58, 34, 48, 34, 44, 34, 52, 57, 52, 34, 58, 34, 48, 46, 48, 48, 48, 48, 34, 44, 34, 52, 54, 57, 34, 58, 34, 48, 34, 44, 34, 52, 57, 53, 34, 58, 34, 49, 34, 44, 34, 52, 56, 48, 34, 58, 34, 48, 34, 44, 34, 52, 56, 50, 34, 58, 34, 48, 46, 48, 48, 48, 48, 34, 44, 34, 52, 55, 57, 34, 58, 34, 32, 34, 44, 34, 52, 56, 49, 34, 58, 34, 48, 46, 48, 48, 48, 48, 34, 44, 34, 53, 49, 55, 34, 58, 34, 48, 34, 44, 34, 53, 49, 48, 34, 58, 34, 32, 34, 44, 34, 53, 50, 55, 34, 58, 34, 34, 44, 34, 53, 50, 54, 34, 58, 34, 48, 34, 44, 34, 53, 50, 52, 34, 58, 34, 48, 46, 48, 48, 48, 48, 34, 44, 34, 53, 50, 48, 34, 58, 34, 32, 34, 44, 34, 53, 48, 48, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 56, 34, 44, 34, 53, 48, 52, 34, 58, 34, 55, 34, 44, 34, 53, 48, 51, 34, 58, 34, 142, 230, 143, 193, 138, 174, 151, 185, 34, 44, 34, 53, 50, 53, 34, 58, 34, 48, 34, 44, 34, 52, 57, 49, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 54, 48, 55, 52, 51, 52, 52, 34, 44, 34, 52, 57, 50, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 52, 56, 57, 34, 58, 34, 48, 34, 44, 34, 52, 55, 48, 34, 58, 34, 49, 34, 44, 34, 50, 51, 53, 34, 58, 34, 48, 34, 125, 93, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListRequest{},
			want1: &OrderListResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 27, 7, 18, 41, 146000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 27, 7, 18, 41, 111000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderList,
				},
				SymbolCode:         "",
				ExecutionDate:      time.Time{},
				OrderInquiryStatus: OrderInquiryStatusUnspecified,
				ResultCode:         "0",
				ResultText:         "",
				WarningCode:        "0",
				WarningText:        "",
				Orders: []Order{
					{
						WarningCode:            "0",
						WarningText:            "",
						OrderNumber:            "28002177",
						SymbolCode:             "1475",
						Exchange:               ExchangeToushou,
						AccountType:            AccountTypeSpecific,
						TradeType:              TradeTypeStock,
						ExitTermType:           ExitTermTypeNoLimit,
						Side:                   SideBuy,
						OrderQuantity:          1,
						CurrentQuantity:        0,
						Price:                  0,
						ExecutionTiming:        ExecutionTimingNormal,
						ExecutionType:          ExecutionTypeMarket,
						StopOrderType:          StopOrderTypeNormal,
						StopTriggerPrice:       0,
						StopOrderExecutionType: ExecutionTypeUnused,
						StopOrderPrice:         0,
						TriggerType:            TriggerTypeNoFired,
						ExitOrderType:          ExitOrderTypeUnused,
						ContractQuantity:       0,
						ContractPrice:          0,
						PartContractType:       PartContractTypeUnused,
						ExecutionDate:          time.Date(2022, 2, 28, 0, 0, 0, 0, time.Local),
						OrderStatus:            OrderStatusCanceled,
						OrderStatusText:        "取消完了",
						ContractStatus:         ContractStatusInOrder,
						OrderDateTime:          time.Date(2022, 2, 26, 7, 43, 44, 0, time.Local),
						ExpireDate:             time.Time{},
						CarryOverType:          CarryOverTypeToday,
						CorrectCancelType:      CorrectCancelTypeCancelable,
						EstimationAmount:       0,
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
			arg3:   OrderListRequest{},
			want1:  nil,
			want2:  NilArgumentErr},
		{name: "失敗をパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 55, 45, 48, 55, 58, 50, 55, 58, 49, 48, 46, 49, 51, 51, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 55, 45, 48, 55, 58, 50, 55, 58, 49, 48, 46, 48, 57, 53, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 79, 114, 100, 101, 114, 76, 105, 115, 116, 34, 44, 34, 51, 50, 56, 34, 58, 34, 34, 44, 34, 53, 53, 57, 34, 58, 34, 34, 44, 34, 53, 48, 56, 34, 58, 34, 34, 44, 34, 53, 51, 52, 34, 58, 34, 57, 57, 49, 48, 48, 51, 34, 44, 34, 53, 51, 53, 34, 58, 34, 150, 193, 149, 191, 131, 82, 129, 91, 131, 104, 130, 201, 140, 235, 130, 232, 130, 170, 130, 160, 130, 232, 130, 220, 130, 183, 129, 66, 34, 44, 34, 54, 57, 50, 34, 58, 34, 34, 44, 34, 54, 57, 51, 34, 58, 34, 34, 44, 34, 53, 53, 34, 58, 34, 34, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListRequest{},
			want1: &OrderListResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 27, 7, 27, 10, 133000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 27, 7, 27, 10, 95000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderList,
				},
				SymbolCode:         "",
				ExecutionDate:      time.Time{},
				OrderInquiryStatus: OrderInquiryStatusUnspecified,
				ResultCode:         "991003",
				ResultText:         "銘柄コードに誤りがあります。",
				WarningCode:        "",
				WarningText:        "",
				Orders:             []Order{},
			},
			want2: nil},
		{name: "200 OK以外が返ったらエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusInternalServerError,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListRequest{},
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
			got1, got2 := client.OrderList(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_OrderList_Execute(t *testing.T) {
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

	got3, got4 := client.OrderList(context.Background(), session, OrderListRequest{
		SymbolCode:         "",
		ExecutionDate:      time.Time{},
		OrderInquiryStatus: OrderInquiryStatusUnspecified,
	})
	log.Printf("%+v, %+v\n", got3, got4)
}
