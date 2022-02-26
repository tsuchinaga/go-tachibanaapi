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

func Test_LoginRequest_request(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request LoginRequest
		arg1    int64
		arg2    time.Time
		want1   loginRequest
	}{
		{name: "LoginRequestの値と引数の値を持ったloginRequestが生成される",
			request: LoginRequest{
				UserId:   "user-id-001",
				Password: "password-1234",
			},
			arg1: 1,
			arg2: time.Date(2022, 2, 10, 9, 0, 0, 0, time.Local),
			want1: loginRequest{
				commonRequest: commonRequest{
					No:          1,
					SendDate:    RequestTime{Time: time.Date(2022, 2, 10, 9, 0, 0, 0, time.Local)},
					FeatureType: FeatureTypeLoginRequest,
				},
				UserId:   "user-id-001",
				Password: "password-1234",
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

func Test_LoginResponse_Session(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		response LoginResponse
		want1    *Session
		want2    error
	}{
		{name: "LoginResponseから必要な項目を取り出してSessionを作って返す",
			response: LoginResponse{
				CommonResponse: CommonResponse{
					No:           1,
					SendDate:     time.Date(2022, 2, 24, 8, 36, 17, 778000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 24, 8, 36, 17, 725000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeLoginResponse,
				},
				ResultCode:            "0",
				ResultText:            "",
				AccountType:           AccountTypeSpecific,
				LastLoginDateTime:     time.Date(2022, 2, 24, 8, 33, 17, 0, time.Local),
				GeneralAccount:        true,
				SafekeepingAccount:    true,
				TransferAccount:       true,
				ForeignAccount:        true,
				MRFAccount:            false,
				StockSpecificAccount:  SpecificAccountTypeNothing,
				MarginSpecificAccount: SpecificAccountTypeGeneral,
				DividendAccount:       false,
				SpecificAccount:       true,
				MarginAccount:         false,
				FutureOptionAccount:   false,
				MMFAccount:            false,
				ChinaForeignAccount:   false,
				FXAccount:             false,
				NISAAccount:           false,
				UnreadDocument:        false,
				RequestURL:            "https://kabuka.e-shiten.jp/e_api_v4r2/request/MjU3NDEzOTA4MjQwMi0xMjMtNDQ4ODU=/",
				EventURL:              "https://kabuka.e-shiten.jp/e_api_v4r2/event/MjU3NDEzOTA4MjQwMi0xMjMtNDQ4ODU=/",
			},
			want1: &Session{
				lastRequestNo: 1,
				RequestURL:    "https://kabuka.e-shiten.jp/e_api_v4r2/request/MjU3NDEzOTA4MjQwMi0xMjMtNDQ4ODU=/",
				EventURL:      "https://kabuka.e-shiten.jp/e_api_v4r2/event/MjU3NDEzOTA4MjQwMi0xMjMtNDQ4ODU=/",
			},
			want2: nil},
		{name: "LoginResponseがエラーの場合はエラーを返す",
			response: LoginResponse{
				CommonResponse: CommonResponse{
					No:           1,
					SendDate:     time.Date(2022, 2, 24, 8, 36, 17, 778000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 24, 8, 36, 17, 725000000, time.Local),
					ErrorNo:      ErrorBadRequest,
					ErrorMessage: "引数エラー。",
					FeatureType:  FeatureTypeLoginRequest,
				},
				ResultCode:            "",
				ResultText:            "",
				AccountType:           AccountTypeUnspecified,
				LastLoginDateTime:     time.Time{},
				GeneralAccount:        false,
				SafekeepingAccount:    false,
				TransferAccount:       false,
				ForeignAccount:        false,
				MRFAccount:            false,
				StockSpecificAccount:  SpecificAccountTypeUnspecified,
				MarginSpecificAccount: SpecificAccountTypeUnspecified,
				DividendAccount:       false,
				SpecificAccount:       false,
				MarginAccount:         false,
				FutureOptionAccount:   false,
				MMFAccount:            false,
				ChinaForeignAccount:   false,
				FXAccount:             false,
				NISAAccount:           false,
				UnreadDocument:        false,
				RequestURL:            "",
				EventURL:              "",
			},
			want1: nil,
			want2: CanNotCreateSessionErr},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got1, got2 := test.response.Session()
			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_loginResponse_response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		response loginResponse
		want1    LoginResponse
	}{
		{name: "変換できる",
			response: loginResponse{
				commonResponse: commonResponse{
					No:           1,
					SendDate:     RequestTime{Time: time.Date(2022, 2, 24, 8, 36, 17, 778000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 2, 24, 8, 36, 17, 725000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeLoginResponse,
				},
				ResultCode:            "0",
				ResultText:            "",
				AccountType:           AccountTypeSpecific,
				LastLoginDateTime:     YmdHms{Time: time.Date(2022, 2, 24, 8, 33, 17, 0, time.Local)},
				GeneralAccount:        NumberBoolTrue,
				SafekeepingAccount:    NumberBoolTrue,
				TransferAccount:       NumberBoolTrue,
				ForeignAccount:        NumberBoolTrue,
				MRFAccount:            NumberBoolFalse,
				StockSpecificAccount:  SpecificAccountTypeNothing,
				MarginSpecificAccount: SpecificAccountTypeGeneral,
				DividendAccount:       NumberBoolFalse,
				SpecificAccount:       NumberBoolTrue,
				MarginAccount:         NumberBoolFalse,
				FutureOptionAccount:   NumberBoolFalse,
				MMFAccount:            NumberBoolFalse,
				ChinaForeignAccount:   NumberBoolFalse,
				FXAccount:             NumberBoolFalse,
				NISAAccount:           NumberBoolFalse,
				UnreadDocument:        NumberBoolFalse,
				RequestURL:            "https://kabuka.e-shiten.jp/e_api_v4r2/request/NzI1MTczNjA4MjQwMi0xMjMtNjMxMTc=/",
				EventURL:              "https://kabuka.e-shiten.jp/e_api_v4r2/event/NzI1MTczNjA4MjQwMi0xMjMtNjMxMTc=/",
			},
			want1: LoginResponse{
				CommonResponse: CommonResponse{
					No:           1,
					SendDate:     time.Date(2022, 2, 24, 8, 36, 17, 778000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 24, 8, 36, 17, 725000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeLoginResponse,
				},
				ResultCode:            "0",
				ResultText:            "",
				AccountType:           AccountTypeSpecific,
				LastLoginDateTime:     time.Date(2022, 2, 24, 8, 33, 17, 0, time.Local),
				GeneralAccount:        true,
				SafekeepingAccount:    true,
				TransferAccount:       true,
				ForeignAccount:        true,
				MRFAccount:            false,
				StockSpecificAccount:  SpecificAccountTypeNothing,
				MarginSpecificAccount: SpecificAccountTypeGeneral,
				DividendAccount:       false,
				SpecificAccount:       true,
				MarginAccount:         false,
				FutureOptionAccount:   false,
				MMFAccount:            false,
				ChinaForeignAccount:   false,
				FXAccount:             false,
				NISAAccount:           false,
				UnreadDocument:        false,
				RequestURL:            "https://kabuka.e-shiten.jp/e_api_v4r2/request/NzI1MTczNjA4MjQwMi0xMjMtNjMxMTc=/",
				EventURL:              "https://kabuka.e-shiten.jp/e_api_v4r2/event/NzI1MTczNjA4MjQwMi0xMjMtNjMxMTc=/",
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

func Test_client_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		clock  *testClock
		status int
		body   []byte
		arg1   context.Context
		arg2   LoginRequest
		want1  *LoginResponse
		want2  error
	}{
		{name: "正常レスポンスをパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 8, 33, 0, 0, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 48, 56, 58, 51, 54, 58, 49, 55, 46, 55, 55, 56, 34, 44, 34, 49, 55, 53, 34, 58, 34, 49, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 48, 56, 58, 51, 54, 58, 49, 55, 46, 55, 50, 53, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 105, 110, 65, 99, 107, 34, 44, 34, 53, 51, 52, 34, 58, 34, 48, 34, 44, 34, 53, 51, 53, 34, 58, 34, 34, 44, 34, 55, 52, 52, 34, 58, 34, 49, 34, 44, 34, 53, 52, 53, 34, 58, 34, 48, 34, 44, 34, 52, 48, 49, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 52, 48, 56, 51, 51, 49, 55, 34, 44, 34, 53, 56, 48, 34, 58, 34, 49, 34, 44, 34, 50, 56, 55, 34, 58, 34, 49, 34, 44, 34, 50, 51, 50, 34, 58, 34, 49, 34, 44, 34, 50, 51, 52, 34, 58, 34, 49, 34, 44, 34, 52, 48, 52, 34, 58, 34, 48, 34, 44, 34, 54, 52, 53, 34, 58, 34, 49, 34, 44, 34, 54, 52, 54, 34, 58, 34, 48, 34, 44, 34, 54, 52, 55, 34, 58, 34, 49, 34, 44, 34, 54, 52, 50, 34, 58, 34, 48, 34, 44, 34, 54, 52, 52, 34, 58, 34, 49, 34, 44, 34, 53, 54, 53, 34, 58, 34, 48, 34, 44, 34, 53, 52, 50, 34, 58, 34, 48, 34, 44, 34, 52, 48, 51, 34, 58, 34, 48, 34, 44, 34, 54, 54, 48, 34, 58, 34, 48, 34, 44, 34, 51, 52, 56, 34, 58, 34, 48, 34, 44, 34, 50, 56, 49, 34, 58, 34, 48, 34, 44, 34, 51, 55, 52, 34, 58, 34, 48, 34, 44, 34, 54, 56, 57, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 107, 97, 98, 117, 107, 97, 46, 101, 45, 115, 104, 105, 116, 101, 110, 46, 106, 112, 47, 101, 95, 97, 112, 105, 95, 118, 52, 114, 50, 47, 114, 101, 113, 117, 101, 115, 116, 47, 78, 122, 73, 49, 77, 84, 99, 122, 78, 106, 65, 52, 77, 106, 81, 119, 77, 105, 48, 120, 77, 106, 77, 116, 78, 106, 77, 120, 77, 84, 99, 61, 47, 34, 44, 34, 54, 56, 56, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 107, 97, 98, 117, 107, 97, 46, 101, 45, 115, 104, 105, 116, 101, 110, 46, 106, 112, 47, 101, 95, 97, 112, 105, 95, 118, 52, 114, 50, 47, 101, 118, 101, 110, 116, 47, 78, 122, 73, 49, 77, 84, 99, 122, 78, 106, 65, 52, 77, 106, 81, 119, 77, 105, 48, 120, 77, 106, 77, 116, 78, 106, 77, 120, 77, 84, 99, 61, 47, 34, 125},
			arg1:   context.Background(),
			arg2:   LoginRequest{UserId: "user-id", Password: "password"},
			want1: &LoginResponse{
				CommonResponse: CommonResponse{
					No:           1,
					SendDate:     time.Date(2022, 2, 24, 8, 36, 17, 778000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 24, 8, 36, 17, 725000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeLoginResponse,
				},
				ResultCode:            "0",
				ResultText:            "",
				AccountType:           AccountTypeSpecific,
				LastLoginDateTime:     time.Date(2022, 2, 24, 8, 33, 17, 0, time.Local),
				GeneralAccount:        true,
				SafekeepingAccount:    true,
				TransferAccount:       true,
				ForeignAccount:        true,
				MRFAccount:            false,
				StockSpecificAccount:  SpecificAccountTypeNothing,
				MarginSpecificAccount: SpecificAccountTypeGeneral,
				DividendAccount:       false,
				SpecificAccount:       true,
				MarginAccount:         false,
				FutureOptionAccount:   false,
				MMFAccount:            false,
				ChinaForeignAccount:   false,
				FXAccount:             false,
				NISAAccount:           false,
				UnreadDocument:        false,
				RequestURL:            "https://kabuka.e-shiten.jp/e_api_v4r2/request/NzI1MTczNjA4MjQwMi0xMjMtNjMxMTc=/",
				EventURL:              "https://kabuka.e-shiten.jp/e_api_v4r2/event/NzI1MTczNjA4MjQwMi0xMjMtNjMxMTc=/",
			},
			want2: nil,
		},
		{name: "ログイン失敗をパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 8, 33, 0, 0, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 49, 48, 58, 51, 53, 58, 52, 48, 46, 55, 53, 49, 34, 44, 34, 49, 55, 53, 34, 58, 34, 49, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 49, 48, 58, 51, 53, 58, 52, 48, 46, 55, 50, 51, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 105, 110, 65, 99, 107, 34, 44, 34, 53, 51, 52, 34, 58, 34, 49, 48, 48, 51, 49, 34, 44, 34, 53, 51, 53, 34, 58, 34, 131, 134, 129, 91, 131, 85, 73, 68, 130, 169, 136, 195, 143, 216, 148, 212, 141, 134, 130, 240, 130, 168, 138, 212, 136, 225, 130, 166, 130, 197, 130, 183, 129, 66, 130, 178, 138, 109, 148, 70, 130, 204, 143, 227, 129, 65, 141, 196, 147, 120, 130, 178, 147, 252, 151, 205, 137, 186, 130, 179, 130, 162, 129, 66, 130, 200, 130, 168, 129, 65, 130, 168, 138, 212, 136, 225, 130, 166, 130, 204, 137, 241, 144, 148, 130, 170, 149, 190, 142, 208, 139, 75, 146, 232, 137, 241, 144, 148, 130, 240, 146, 180, 130, 166, 130, 233, 130, 198, 129, 65, 131, 141, 131, 79, 131, 67, 131, 147, 146, 226, 142, 126, 130, 198, 130, 200, 130, 232, 130, 220, 130, 183, 130, 204, 130, 197, 130, 178, 146, 141, 136, 211, 137, 186, 130, 179, 130, 162, 129, 66, 40, 131, 141, 131, 79, 131, 67, 131, 147, 146, 226, 142, 126, 130, 204, 137, 240, 143, 156, 130, 205, 129, 65, 131, 82, 129, 91, 131, 139, 131, 90, 131, 147, 131, 94, 129, 91, 130, 220, 130, 197, 130, 168, 147, 100, 152, 98, 137, 186, 130, 179, 130, 162, 129, 66, 41, 34, 44, 34, 55, 52, 52, 34, 58, 34, 34, 44, 34, 53, 52, 53, 34, 58, 34, 34, 44, 34, 52, 48, 49, 34, 58, 34, 34, 44, 34, 53, 56, 48, 34, 58, 34, 34, 44, 34, 50, 56, 55, 34, 58, 34, 34, 44, 34, 50, 51, 50, 34, 58, 34, 34, 44, 34, 50, 51, 52, 34, 58, 34, 34, 44, 34, 52, 48, 52, 34, 58, 34, 34, 44, 34, 54, 52, 53, 34, 58, 34, 34, 44, 34, 54, 52, 54, 34, 58, 34, 34, 44, 34, 54, 52, 55, 34, 58, 34, 34, 44, 34, 54, 52, 50, 34, 58, 34, 34, 44, 34, 54, 52, 52, 34, 58, 34, 34, 44, 34, 53, 54, 53, 34, 58, 34, 34, 44, 34, 53, 52, 50, 34, 58, 34, 34, 44, 34, 52, 48, 51, 34, 58, 34, 34, 44, 34, 54, 54, 48, 34, 58, 34, 34, 44, 34, 51, 52, 56, 34, 58, 34, 34, 44, 34, 50, 56, 49, 34, 58, 34, 34, 44, 34, 51, 55, 52, 34, 58, 34, 34, 44, 34, 54, 56, 57, 34, 58, 34, 34, 44, 34, 54, 56, 56, 34, 58, 34, 34, 125},
			arg1:   context.Background(),
			arg2:   LoginRequest{UserId: "user-id", Password: "password"},
			want1: &LoginResponse{
				CommonResponse: CommonResponse{
					No:           1,
					SendDate:     time.Date(2022, 2, 24, 10, 35, 40, 751000000, time.Local),
					ReceiveDate:  time.Date(2022, 2, 24, 10, 35, 40, 723000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeLoginResponse,
				},
				ResultCode:            "10031",
				ResultText:            "ユーザIDか暗証番号をお間違えです。ご確認の上、再度ご入力下さい。なお、お間違えの回数が弊社規定回数を超えると、ログイン停止となりますのでご注意下さい。(ログイン停止の解除は、コールセンターまでお電話下さい。)",
				AccountType:           "",
				LastLoginDateTime:     time.Time{},
				GeneralAccount:        false,
				SafekeepingAccount:    false,
				TransferAccount:       false,
				ForeignAccount:        false,
				MRFAccount:            false,
				StockSpecificAccount:  SpecificAccountTypeUnspecified,
				MarginSpecificAccount: SpecificAccountTypeUnspecified,
				DividendAccount:       false,
				SpecificAccount:       false,
				MarginAccount:         false,
				FutureOptionAccount:   false,
				MMFAccount:            false,
				ChinaForeignAccount:   false,
				FXAccount:             false,
				NISAAccount:           false,
				UnreadDocument:        false,
				RequestURL:            "",
				EventURL:              "",
			},
			want2: nil,
		},
		{name: "200 OK以外が返ったらエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 8, 33, 0, 0, time.Local)},
			status: http.StatusInternalServerError,
			body:   nil,
			arg1:   context.Background(),
			arg2:   LoginRequest{UserId: "user-id", Password: "password"},
			want1:  nil,
			want2:  StatusNotOkErr,
		},
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

			req := &client{clock: test.clock, auth: ts.URL}
			got1, got2 := req.Login(test.arg1, test.arg2)
			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %v\ngot: %+v, %v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_Login_Execute(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client := NewClient(EnvironmentProduction, ApiVersionLatest)
	got1, got2 := client.Login(context.Background(), LoginRequest{
		UserId:   "user-id",
		Password: "password",
	})
	log.Printf("%+v, %+v\n", got1, got2)
}
