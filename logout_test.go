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

func Test_LogoutRequest_request(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request LogoutRequest
		arg1    int64
		arg2    time.Time
		want1   logoutRequest
	}{
		{name: "引数の値を持ったログアウトのリクエストが生成される",
			arg1:    2,
			request: LogoutRequest{},
			arg2:    time.Date(2022, 2, 24, 10, 0, 0, 0, time.Local),
			want1: logoutRequest{commonRequest: commonRequest{
				No:          2,
				SendDate:    RequestTime{Time: time.Date(2022, 2, 24, 10, 0, 0, 0, time.Local)},
				FeatureType: FeatureTypeLogoutRequest,
			}}},
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

func Test_logoutResponse_response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		response logoutResponse
		want1    LogoutResponse
	}{
		{name: "変換できる",
			response: logoutResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 2, 24, 21, 2, 23, 335000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeLogoutResponse,
				},
				ResultCode: "0",
				ResultText: "",
			},
			want1: LogoutResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 24, 21, 2, 23, 335000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeLogoutResponse,
				},
				ResultCode: "0",
				ResultText: "",
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

func Test_client_Logout(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		clock  *testClock
		status int
		body   []byte
		arg1   context.Context
		arg2   *Session
		arg3   LogoutRequest
		want1  *LogoutResponse
		want2  error
	}{
		{name: "正常レスポンスをパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 50, 49, 58, 48, 50, 58, 50, 51, 46, 51, 54, 53, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 50, 49, 58, 48, 50, 58, 50, 51, 46, 51, 51, 53, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 111, 117, 116, 65, 99, 107, 34, 44, 34, 53, 51, 52, 34, 58, 34, 48, 34, 44, 34, 53, 51, 53, 34, 58, 34, 34, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   LogoutRequest{},
			want1: &LogoutResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 24, 21, 2, 23, 335000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeLogoutResponse,
				},
				ResultCode: "0",
				ResultText: "",
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 50, 49, 58, 48, 50, 58, 50, 51, 46, 51, 54, 53, 34, 44, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 50, 49, 58, 48, 50, 58, 50, 51, 46, 51, 51, 53, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 111, 117, 116, 65, 99, 107, 34, 44, 34, 53, 51, 52, 34, 58, 34, 48, 34, 44, 34, 53, 51, 53, 34, 58, 34, 34, 125},
			arg1:   context.Background(),
			arg2:   nil,
			arg3:   LogoutRequest{},
			want1:  nil,
			want2:  NilArgumentErr},
		{name: "失敗をパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 50, 49, 58, 48, 49, 58, 53, 56, 46, 52, 53, 51, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 50, 49, 58, 48, 49, 58, 53, 56, 46, 52, 52, 48, 34, 44, 34, 49, 55, 52, 34, 58, 34, 50, 34, 44, 34, 49, 55, 51, 34, 58, 34, 131, 90, 131, 98, 131, 86, 131, 135, 131, 147, 130, 170, 144, 216, 146, 102, 130, 181, 130, 220, 130, 181, 130, 189, 129, 66, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 111, 117, 116, 82, 101, 113, 117, 101, 115, 116, 34, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   LogoutRequest{},
			want1: &LogoutResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 24, 21, 1, 58, 440000000, time.Local),
					ErrorNo:      ErrorSessionInactive,
					ErrorMessage: "セッションが切断しました。",
					FeatureType:  FeatureTypeLogoutRequest,
				},
				ResultCode: "",
				ResultText: "",
			},
			want2: nil},
		{name: "200 OK以外が返ったらエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusInternalServerError,
			body:   []byte{123, 34, 49, 55, 53, 34, 58, 34, 50, 34, 44, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 50, 49, 58, 48, 49, 58, 53, 56, 46, 52, 53, 51, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 50, 49, 58, 48, 49, 58, 53, 56, 46, 52, 52, 48, 34, 44, 34, 49, 55, 52, 34, 58, 34, 50, 34, 44, 34, 49, 55, 51, 34, 58, 34, 131, 90, 131, 98, 131, 86, 131, 135, 131, 147, 130, 170, 144, 216, 146, 102, 130, 181, 130, 220, 130, 181, 130, 189, 129, 66, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 111, 117, 116, 82, 101, 113, 117, 101, 115, 116, 34, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   LogoutRequest{},
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
			got1, got2 := client.Logout(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_Logout_Execute(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client := NewClient(EnvironmentProduction, ApiVersionLatest)
	got1, got2 := client.Login(context.Background(), LoginRequest{
		UserId:   "user-id",
		Password: "password",
	})
	log.Printf("%+v, %+v\n", got1, got2)

	session, err := got1.Session()
	if err != nil {
		t.Errorf("session: %+v\n", err)
	}

	got3, got4 := client.Logout(context.Background(), session, LogoutRequest{})
	log.Printf("%+v, %+v\n", got3, got4)
}