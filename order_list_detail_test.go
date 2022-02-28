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

func Test_orderListDetail_request(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request OrderListDetailRequest
		arg1    int64
		arg2    time.Time
		want1   orderListDetailRequest
	}{
		{name: "変換できる",
			request: OrderListDetailRequest{
				OrderNumber:   "28002795",
				ExecutionDate: time.Time{},
			},
			arg1: 123,
			arg2: time.Date(2022, 2, 27, 10, 21, 15, 0, time.Local),
			want1: orderListDetailRequest{
				commonRequest: commonRequest{
					No:          123,
					SendDate:    RequestTime{Time: time.Date(2022, 2, 27, 10, 21, 15, 0, time.Local)},
					FeatureType: FeatureTypeOrderListDetail,
				},
				OrderNumber:   "28002795",
				ExecutionDate: Ymd{},
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

// TODO 各種レスポンスの変換のテスト

func Test_client_OrderListDetail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		clock  *testClock
		status int
		body   []byte
		arg1   context.Context
		arg2   *Session
		arg3   OrderListDetailRequest
		want1  *OrderListDetailResponse
		want2  error
	}{
		{name: "正常レスポンスをパースして返せる", // TODO 正常系のレスポンスのテスト
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 48, 57, 58, 51, 57, 58, 49, 53, 46, 56, 51, 55, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 48, 57, 58, 51, 57, 58, 49, 53, 46, 56, 48, 53, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 79, 114, 100, 101, 114, 76, 105, 115, 116, 68, 101, 116, 97, 105, 108, 34, 44, 34, 52, 57, 48, 34, 58, 34, 34, 44, 34, 50, 50, 55, 34, 58, 34, 34, 44, 34, 53, 51, 52, 34, 58, 34, 57, 57, 49, 48, 48, 49, 34, 44, 34, 53, 51, 53, 34, 58, 34, 146, 141, 149, 182, 148, 212, 141, 134, 130, 201, 140, 235, 130, 232, 130, 170, 130, 160, 130, 232, 130, 220, 130, 183, 129, 66, 34, 44, 34, 54, 57, 50, 34, 58, 34, 34, 44, 34, 54, 57, 51, 34, 58, 34, 34, 44, 34, 51, 50, 56, 34, 58, 34, 34, 44, 34, 53, 48, 49, 34, 58, 34, 34, 44, 34, 52, 54, 55, 34, 58, 34, 34, 44, 34, 50, 53, 53, 34, 58, 34, 34, 44, 34, 52, 54, 56, 34, 58, 34, 34, 44, 34, 52, 54, 57, 34, 58, 34, 34, 44, 34, 52, 57, 53, 34, 58, 34, 34, 44, 34, 52, 57, 52, 34, 58, 34, 34, 44, 34, 52, 57, 54, 34, 58, 34, 34, 44, 34, 52, 55, 49, 34, 58, 34, 34, 44, 34, 53, 48, 52, 34, 58, 34, 34, 44, 34, 53, 48, 51, 34, 58, 34, 34, 44, 34, 52, 57, 49, 34, 58, 34, 34, 44, 34, 52, 57, 50, 34, 58, 34, 34, 44, 34, 49, 57, 51, 34, 58, 34, 34, 44, 34, 50, 52, 56, 34, 58, 34, 34, 44, 34, 53, 55, 53, 34, 58, 34, 34, 44, 34, 50, 53, 57, 34, 58, 34, 34, 44, 34, 50, 54, 51, 34, 58, 34, 34, 44, 34, 50, 53, 56, 34, 58, 34, 34, 44, 34, 50, 54, 48, 34, 58, 34, 34, 44, 34, 54, 53, 57, 34, 58, 34, 34, 44, 34, 54, 53, 56, 34, 58, 34, 34, 44, 34, 54, 54, 50, 34, 58, 34, 34, 44, 34, 54, 57, 53, 34, 58, 34, 34, 44, 34, 54, 57, 54, 34, 58, 34, 34, 44, 34, 49, 56, 50, 34, 58, 34, 34, 44, 34, 54, 57, 49, 34, 58, 34, 34, 44, 34, 50, 51, 53, 34, 58, 34, 34, 44, 34, 49, 56, 51, 34, 58, 34, 34, 44, 34, 53, 53, 56, 34, 58, 34, 34, 44, 34, 54, 50, 48, 34, 58, 34, 34, 44, 34, 53, 55, 55, 34, 58, 34, 34, 44, 34, 55, 52, 49, 34, 58, 34, 34, 44, 34, 52, 54, 54, 34, 58, 34, 34, 44, 34, 53, 55, 34, 58, 34, 34, 44, 34, 53, 51, 34, 58, 34, 34, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListDetailRequest{},
			want1: &OrderListDetailResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 28, 9, 39, 15, 837000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 28, 9, 39, 15, 805000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderListDetail,
				},
				ResultCode:    "991001",
				ResultText:    "注文番号に誤りがあります。",
				Contracts:     []Contract{},
				HoldPositions: []HoldPosition{},
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   nil,
			arg3:   OrderListDetailRequest{},
			want1:  nil,
			want2:  NilArgumentErr},
		{name: "失敗をパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 48, 57, 58, 51, 57, 58, 49, 53, 46, 56, 51, 55, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 48, 57, 58, 51, 57, 58, 49, 53, 46, 56, 48, 53, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 79, 114, 100, 101, 114, 76, 105, 115, 116, 68, 101, 116, 97, 105, 108, 34, 44, 34, 52, 57, 48, 34, 58, 34, 34, 44, 34, 50, 50, 55, 34, 58, 34, 34, 44, 34, 53, 51, 52, 34, 58, 34, 57, 57, 49, 48, 48, 49, 34, 44, 34, 53, 51, 53, 34, 58, 34, 146, 141, 149, 182, 148, 212, 141, 134, 130, 201, 140, 235, 130, 232, 130, 170, 130, 160, 130, 232, 130, 220, 130, 183, 129, 66, 34, 44, 34, 54, 57, 50, 34, 58, 34, 34, 44, 34, 54, 57, 51, 34, 58, 34, 34, 44, 34, 51, 50, 56, 34, 58, 34, 34, 44, 34, 53, 48, 49, 34, 58, 34, 34, 44, 34, 52, 54, 55, 34, 58, 34, 34, 44, 34, 50, 53, 53, 34, 58, 34, 34, 44, 34, 52, 54, 56, 34, 58, 34, 34, 44, 34, 52, 54, 57, 34, 58, 34, 34, 44, 34, 52, 57, 53, 34, 58, 34, 34, 44, 34, 52, 57, 52, 34, 58, 34, 34, 44, 34, 52, 57, 54, 34, 58, 34, 34, 44, 34, 52, 55, 49, 34, 58, 34, 34, 44, 34, 53, 48, 52, 34, 58, 34, 34, 44, 34, 53, 48, 51, 34, 58, 34, 34, 44, 34, 52, 57, 49, 34, 58, 34, 34, 44, 34, 52, 57, 50, 34, 58, 34, 34, 44, 34, 49, 57, 51, 34, 58, 34, 34, 44, 34, 50, 52, 56, 34, 58, 34, 34, 44, 34, 53, 55, 53, 34, 58, 34, 34, 44, 34, 50, 53, 57, 34, 58, 34, 34, 44, 34, 50, 54, 51, 34, 58, 34, 34, 44, 34, 50, 53, 56, 34, 58, 34, 34, 44, 34, 50, 54, 48, 34, 58, 34, 34, 44, 34, 54, 53, 57, 34, 58, 34, 34, 44, 34, 54, 53, 56, 34, 58, 34, 34, 44, 34, 54, 54, 50, 34, 58, 34, 34, 44, 34, 54, 57, 53, 34, 58, 34, 34, 44, 34, 54, 57, 54, 34, 58, 34, 34, 44, 34, 49, 56, 50, 34, 58, 34, 34, 44, 34, 54, 57, 49, 34, 58, 34, 34, 44, 34, 50, 51, 53, 34, 58, 34, 34, 44, 34, 49, 56, 51, 34, 58, 34, 34, 44, 34, 53, 53, 56, 34, 58, 34, 34, 44, 34, 54, 50, 48, 34, 58, 34, 34, 44, 34, 53, 55, 55, 34, 58, 34, 34, 44, 34, 55, 52, 49, 34, 58, 34, 34, 44, 34, 52, 54, 54, 34, 58, 34, 34, 44, 34, 53, 55, 34, 58, 34, 34, 44, 34, 53, 51, 34, 58, 34, 34, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListDetailRequest{},
			want1: &OrderListDetailResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 28, 9, 39, 15, 837000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 28, 9, 39, 15, 805000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderListDetail,
				},
				ResultCode:    "991001",
				ResultText:    "注文番号に誤りがあります。",
				Contracts:     []Contract{},
				HoldPositions: []HoldPosition{},
			},
			want2: nil},
		{name: "利用できないとう情報をパースできる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 48, 57, 58, 51, 57, 58, 51, 53, 46, 53, 49, 52, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 56, 45, 48, 57, 58, 51, 57, 58, 51, 53, 46, 52, 55, 50, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 79, 114, 100, 101, 114, 76, 105, 115, 116, 68, 101, 116, 97, 105, 108, 34, 44, 34, 52, 57, 48, 34, 58, 34, 34, 44, 34, 50, 50, 55, 34, 58, 34, 34, 44, 34, 53, 51, 52, 34, 58, 34, 57, 57, 49, 48, 48, 50, 34, 44, 34, 53, 51, 53, 34, 58, 34, 145, 252, 141, 161, 129, 65, 136, 234, 142, 158, 147, 73, 130, 201, 130, 177, 130, 204, 139, 198, 150, 177, 130, 205, 130, 178, 151, 152, 151, 112, 130, 197, 130, 171, 130, 220, 130, 185, 130, 241, 129, 66, 34, 44, 34, 54, 57, 50, 34, 58, 34, 34, 44, 34, 54, 57, 51, 34, 58, 34, 34, 44, 34, 51, 50, 56, 34, 58, 34, 34, 44, 34, 53, 48, 49, 34, 58, 34, 34, 44, 34, 52, 54, 55, 34, 58, 34, 34, 44, 34, 50, 53, 53, 34, 58, 34, 34, 44, 34, 52, 54, 56, 34, 58, 34, 34, 44, 34, 52, 54, 57, 34, 58, 34, 34, 44, 34, 52, 57, 53, 34, 58, 34, 34, 44, 34, 52, 57, 52, 34, 58, 34, 34, 44, 34, 52, 57, 54, 34, 58, 34, 34, 44, 34, 52, 55, 49, 34, 58, 34, 34, 44, 34, 53, 48, 52, 34, 58, 34, 34, 44, 34, 53, 48, 51, 34, 58, 34, 34, 44, 34, 52, 57, 49, 34, 58, 34, 34, 44, 34, 52, 57, 50, 34, 58, 34, 34, 44, 34, 49, 57, 51, 34, 58, 34, 34, 44, 34, 50, 52, 56, 34, 58, 34, 34, 44, 34, 53, 55, 53, 34, 58, 34, 34, 44, 34, 50, 53, 57, 34, 58, 34, 34, 44, 34, 50, 54, 51, 34, 58, 34, 34, 44, 34, 50, 53, 56, 34, 58, 34, 34, 44, 34, 50, 54, 48, 34, 58, 34, 34, 44, 34, 54, 53, 57, 34, 58, 34, 34, 44, 34, 54, 53, 56, 34, 58, 34, 34, 44, 34, 54, 54, 50, 34, 58, 34, 34, 44, 34, 54, 57, 53, 34, 58, 34, 34, 44, 34, 54, 57, 54, 34, 58, 34, 34, 44, 34, 49, 56, 50, 34, 58, 34, 34, 44, 34, 54, 57, 49, 34, 58, 34, 34, 44, 34, 50, 51, 53, 34, 58, 34, 34, 44, 34, 49, 56, 51, 34, 58, 34, 34, 44, 34, 53, 53, 56, 34, 58, 34, 34, 44, 34, 54, 50, 48, 34, 58, 34, 34, 44, 34, 53, 55, 55, 34, 58, 34, 34, 44, 34, 55, 52, 49, 34, 58, 34, 34, 44, 34, 52, 54, 54, 34, 58, 34, 34, 44, 34, 53, 55, 34, 58, 34, 34, 44, 34, 53, 51, 34, 58, 34, 34, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListDetailRequest{},
			want1: &OrderListDetailResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 28, 9, 39, 35, 514000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 28, 9, 39, 35, 472000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeOrderListDetail,
				},
				ResultCode:    "991002",
				ResultText:    "只今、一時的にこの業務はご利用できません。",
				Contracts:     []Contract{},
				HoldPositions: []HoldPosition{},
			},
			want2: nil},
		{name: "200 OK以外が返ったらエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusInternalServerError,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   OrderListDetailRequest{},
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
			got1, got2 := client.OrderListDetail(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_OrderListDetail_Execute(t *testing.T) {
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

	got3, got4 := client.OrderListDetail(context.Background(), session, OrderListDetailRequest{
		OrderNumber:   "28008243",
		ExecutionDate: time.Time{},
	})
	log.Printf("%+v, %+v\n", got3, got4)
}
