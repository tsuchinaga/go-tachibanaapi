package tachibana

import (
	"context"
	"errors"
	"log"
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
				No:             2,
				SendDate:       RequestTime{Time: time.Date(2022, 2, 24, 10, 0, 0, 0, time.Local)},
				FeatureType:    FeatureTypeLogoutRequest,
				ResponseFormat: commonResponseFormat,
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
		name      string
		clock     *testClock
		requester *testRequester
		arg1      context.Context
		arg2      *Session
		arg3      LogoutRequest
		want1     *LogoutResponse
		want2     error
	}{
		{name: "正常レスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get1: []byte{123, 10, 9, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 48, 50, 58, 53, 55, 46, 56, 51, 51, 34, 44, 10, 9, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 10, 9, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 48, 50, 58, 53, 55, 46, 55, 57, 53, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 10, 9, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 111, 117, 116, 65, 99, 107, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 10, 125, 10, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      LogoutRequest{},
			want1: &LogoutResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 1, 14, 2, 57, 833000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 1, 14, 2, 57, 795000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeLogoutResponse,
				},
				ResultCode: "0",
				ResultText: "",
			},
			want2: nil},
		{name: "失敗をパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 1, 58, 453000000, time.Local)},
			requester: &testRequester{get1: []byte{123, 10, 9, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 10, 9, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 48, 51, 58, 51, 51, 46, 49, 48, 56, 34, 44, 10, 9, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 52, 58, 48, 51, 58, 51, 51, 46, 48, 57, 54, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 50, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 34, 58, 34, 227, 130, 187, 227, 131, 131, 227, 130, 183, 227, 131, 167, 227, 131, 179, 227, 129, 140, 229, 136, 135, 230, 150, 173, 227, 129, 151, 227, 129, 190, 227, 129, 151, 227, 129, 159, 227, 128, 130, 34, 44, 10, 9, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 111, 117, 116, 82, 101, 113, 117, 101, 115, 116, 34, 10, 125, 10, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      LogoutRequest{},
			want1: &LogoutResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 1, 14, 3, 33, 108000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 1, 14, 3, 33, 96000000, time.Local),
					ErrorNo:      ErrorSessionInactive,
					ErrorMessage: "セッションが切断しました。",
					FeatureType:  FeatureTypeLogoutRequest,
				},
				ResultCode: "",
				ResultText: "",
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			arg1:  context.Background(),
			arg2:  nil,
			arg3:  LogoutRequest{},
			want1: nil,
			want2: NilArgumentErr},
		{name: "リクエストでエラーが返されたらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get2: StatusNotOkErr},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      LogoutRequest{},
			want1:     nil,
			want2:     StatusNotOkErr},
		{name: "レスポンスのUnmarshalに失敗したらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get1: []byte{}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      LogoutRequest{},
			want1:     nil,
			want2:     UnmarshalFailedErr},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &client{clock: test.clock, requester: test.requester}
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
