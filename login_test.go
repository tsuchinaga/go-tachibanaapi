package tachibana

import (
	"context"
	"errors"
	"log"
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
					No:             1,
					SendDate:       RequestTime{Time: time.Date(2022, 2, 10, 9, 0, 0, 0, time.Local)},
					FeatureType:    FeatureTypeLoginRequest,
					ResponseFormat: commonResponseFormat,
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
				ResultCode:                "0",
				ResultText:                "",
				AccountType:               AccountTypeSpecific,
				SecondPasswordOmit:        NumberBoolFalse,
				LastLoginDateTime:         YmdHms{Time: time.Date(2022, 2, 24, 8, 33, 17, 0, time.Local)},
				GeneralAccount:            NumberBoolTrue,
				SafekeepingAccount:        NumberBoolTrue,
				TransferAccount:           NumberBoolTrue,
				ForeignAccount:            NumberBoolTrue,
				MRFAccount:                NumberBoolFalse,
				StockSpecificAccount:      SpecificAccountTypeNothing,
				MarginSpecificAccount:     SpecificAccountTypeGeneral,
				InvestmentSpecificAccount: SpecificAccountTypeNothing,
				DividendAccount:           NumberBoolFalse,
				SpecificAccount:           NumberBoolTrue,
				MarginAccount:             NumberBoolFalse,
				FutureOptionAccount:       NumberBoolFalse,
				MMFAccount:                NumberBoolFalse,
				ChinaForeignAccount:       NumberBoolFalse,
				FXAccount:                 NumberBoolFalse,
				NISAAccount:               NumberBoolFalse,
				UnreadDocument:            NumberBoolFalse,
				RequestURL:                "https://kabuka.e-shiten.jp/e_api_v4r2/request/NzI1MTczNjA4MjQwMi0xMjMtNjMxMTc=/",
				EventURL:                  "https://kabuka.e-shiten.jp/e_api_v4r2/event/NzI1MTczNjA4MjQwMi0xMjMtNjMxMTc=/",
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
				ResultCode:                "0",
				ResultText:                "",
				AccountType:               AccountTypeSpecific,
				SecondPasswordOmit:        false,
				LastLoginDateTime:         time.Date(2022, 2, 24, 8, 33, 17, 0, time.Local),
				GeneralAccount:            true,
				SafekeepingAccount:        true,
				TransferAccount:           true,
				ForeignAccount:            true,
				MRFAccount:                false,
				StockSpecificAccount:      SpecificAccountTypeNothing,
				MarginSpecificAccount:     SpecificAccountTypeGeneral,
				InvestmentSpecificAccount: SpecificAccountTypeNothing,
				DividendAccount:           false,
				SpecificAccount:           true,
				MarginAccount:             false,
				FutureOptionAccount:       false,
				MMFAccount:                false,
				ChinaForeignAccount:       false,
				FXAccount:                 false,
				NISAAccount:               false,
				UnreadDocument:            false,
				RequestURL:                "https://kabuka.e-shiten.jp/e_api_v4r2/request/NzI1MTczNjA4MjQwMi0xMjMtNjMxMTc=/",
				EventURL:                  "https://kabuka.e-shiten.jp/e_api_v4r2/event/NzI1MTczNjA4MjQwMi0xMjMtNjMxMTc=/",
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
		name      string
		clock     *testClock
		requester *testRequester
		arg1      context.Context
		arg2      LoginRequest
		want1     *LoginResponse
		want2     error
	}{
		{name: "正常レスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 8, 33, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 10, 9, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 48, 58, 53, 55, 58, 49, 54, 46, 55, 54, 56, 34, 44, 10, 9, 34, 112, 95, 110, 111, 34, 58, 34, 49, 34, 44, 10, 9, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 48, 58, 53, 55, 58, 49, 54, 46, 54, 57, 51, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 10, 9, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 105, 110, 65, 99, 107, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 90, 121, 111, 117, 116, 111, 101, 107, 105, 75, 97, 122, 101, 105, 67, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 83, 101, 99, 111, 110, 100, 80, 97, 115, 115, 119, 111, 114, 100, 79, 109, 105, 116, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 76, 97, 115, 116, 76, 111, 103, 105, 110, 68, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 51, 48, 49, 49, 48, 53, 54, 52, 49, 34, 44, 10, 9, 34, 115, 83, 111, 103, 111, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 72, 111, 103, 111, 65, 100, 117, 107, 97, 114, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 70, 117, 114, 105, 107, 97, 101, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 71, 97, 105, 107, 111, 107, 117, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 77, 82, 70, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 71, 101, 110, 98, 117, 116, 117, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 83, 105, 110, 121, 111, 117, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 84, 111, 117, 115, 105, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 72, 97, 105, 116, 111, 117, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 75, 97, 110, 114, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 83, 105, 110, 121, 111, 117, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 83, 97, 107, 111, 112, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 77, 77, 70, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 84, 121, 117, 107, 111, 107, 117, 102, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 75, 97, 119, 97, 115, 101, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 72, 105, 107, 97, 122, 101, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 75, 105, 110, 115, 121, 111, 117, 104, 111, 117, 77, 105, 100, 111, 107, 117, 70, 108, 103, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 85, 114, 108, 82, 101, 113, 117, 101, 115, 116, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 107, 97, 98, 117, 107, 97, 46, 101, 45, 115, 104, 105, 116, 101, 110, 46, 106, 112, 47, 101, 95, 97, 112, 105, 95, 118, 52, 114, 50, 47, 114, 101, 113, 117, 101, 115, 116, 47, 78, 106, 107, 48, 77, 84, 89, 49, 78, 122, 69, 119, 77, 68, 69, 119, 77, 121, 48, 120, 77, 106, 77, 116, 78, 84, 77, 48, 78, 122, 89, 61, 47, 34, 44, 10, 9, 34, 115, 85, 114, 108, 69, 118, 101, 110, 116, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 107, 97, 98, 117, 107, 97, 46, 101, 45, 115, 104, 105, 116, 101, 110, 46, 106, 112, 47, 101, 95, 97, 112, 105, 95, 118, 52, 114, 50, 47, 101, 118, 101, 110, 116, 47, 78, 106, 107, 48, 77, 84, 89, 49, 78, 122, 69, 119, 77, 68, 69, 119, 77, 121, 48, 120, 77, 106, 77, 116, 78, 84, 77, 48, 78, 122, 89, 61, 47, 34, 10, 125, 10, 10}},
			arg1:      context.Background(),
			arg2:      LoginRequest{UserId: "user-id", Password: "password"},
			want1: &LoginResponse{
				CommonResponse: CommonResponse{
					No:           1,
					SendDate:     time.Date(2022, 3, 1, 10, 57, 16, 768000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 1, 10, 57, 16, 693000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeLoginResponse,
				},
				ResultCode:                "0",
				ResultText:                "",
				AccountType:               AccountTypeSpecific,
				SecondPasswordOmit:        false,
				LastLoginDateTime:         time.Date(2022, 3, 1, 10, 56, 41, 0, time.Local),
				GeneralAccount:            true,
				SafekeepingAccount:        true,
				TransferAccount:           true,
				ForeignAccount:            true,
				MRFAccount:                false,
				StockSpecificAccount:      SpecificAccountTypeNothing,
				MarginSpecificAccount:     SpecificAccountTypeNothing,
				InvestmentSpecificAccount: SpecificAccountTypeNothing,
				DividendAccount:           false,
				SpecificAccount:           true,
				MarginAccount:             true,
				FutureOptionAccount:       false,
				MMFAccount:                false,
				ChinaForeignAccount:       false,
				FXAccount:                 false,
				NISAAccount:               false,
				UnreadDocument:            false,
				RequestURL:                "https://kabuka.e-shiten.jp/e_api_v4r2/request/Njk0MTY1NzEwMDEwMy0xMjMtNTM0NzY=/",
				EventURL:                  "https://kabuka.e-shiten.jp/e_api_v4r2/event/Njk0MTY1NzEwMDEwMy0xMjMtNTM0NzY=/",
			},
			want2: nil,
		},
		{name: "ログイン失敗をパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 8, 33, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 10, 9, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 49, 58, 48, 55, 58, 49, 55, 46, 54, 53, 48, 34, 44, 10, 9, 34, 112, 95, 110, 111, 34, 58, 34, 49, 34, 44, 10, 9, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 49, 58, 48, 55, 58, 49, 55, 46, 54, 50, 51, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 10, 9, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 105, 110, 65, 99, 107, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 49, 48, 48, 51, 49, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 227, 131, 166, 227, 131, 188, 227, 130, 182, 73, 68, 227, 129, 139, 230, 154, 151, 232, 168, 188, 231, 149, 170, 229, 143, 183, 227, 130, 146, 227, 129, 138, 233, 150, 147, 233, 129, 149, 227, 129, 136, 227, 129, 167, 227, 129, 153, 227, 128, 130, 227, 129, 148, 231, 162, 186, 232, 170, 141, 227, 129, 174, 228, 184, 138, 227, 128, 129, 229, 134, 141, 229, 186, 166, 227, 129, 148, 229, 133, 165, 229, 138, 155, 228, 184, 139, 227, 129, 149, 227, 129, 132, 227, 128, 130, 227, 129, 170, 227, 129, 138, 227, 128, 129, 227, 129, 138, 233, 150, 147, 233, 129, 149, 227, 129, 136, 227, 129, 174, 229, 155, 158, 230, 149, 176, 227, 129, 140, 229, 188, 138, 231, 164, 190, 232, 166, 143, 229, 174, 154, 229, 155, 158, 230, 149, 176, 227, 130, 146, 232, 182, 133, 227, 129, 136, 227, 130, 139, 227, 129, 168, 227, 128, 129, 227, 131, 173, 227, 130, 176, 227, 130, 164, 227, 131, 179, 229, 129, 156, 230, 173, 162, 227, 129, 168, 227, 129, 170, 227, 130, 138, 227, 129, 190, 227, 129, 153, 227, 129, 174, 227, 129, 167, 227, 129, 148, 230, 179, 168, 230, 132, 143, 228, 184, 139, 227, 129, 149, 227, 129, 132, 227, 128, 130, 40, 227, 131, 173, 227, 130, 176, 227, 130, 164, 227, 131, 179, 229, 129, 156, 230, 173, 162, 227, 129, 174, 232, 167, 163, 233, 153, 164, 227, 129, 175, 227, 128, 129, 227, 130, 179, 227, 131, 188, 227, 131, 171, 227, 130, 187, 227, 131, 179, 227, 130, 191, 227, 131, 188, 227, 129, 190, 227, 129, 167, 227, 129, 138, 233, 155, 187, 232, 169, 177, 228, 184, 139, 227, 129, 149, 227, 129, 132, 227, 128, 130, 41, 34, 44, 10, 9, 34, 115, 90, 121, 111, 117, 116, 111, 101, 107, 105, 75, 97, 122, 101, 105, 67, 34, 58, 34, 34, 44, 10, 9, 34, 115, 83, 101, 99, 111, 110, 100, 80, 97, 115, 115, 119, 111, 114, 100, 79, 109, 105, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 76, 97, 115, 116, 76, 111, 103, 105, 110, 68, 97, 116, 101, 34, 58, 34, 34, 44, 10, 9, 34, 115, 83, 111, 103, 111, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 72, 111, 103, 111, 65, 100, 117, 107, 97, 114, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 70, 117, 114, 105, 107, 97, 101, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 71, 97, 105, 107, 111, 107, 117, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 77, 82, 70, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 71, 101, 110, 98, 117, 116, 117, 34, 58, 34, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 83, 105, 110, 121, 111, 117, 34, 58, 34, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 84, 111, 117, 115, 105, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 72, 97, 105, 116, 111, 117, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 75, 97, 110, 114, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 83, 105, 110, 121, 111, 117, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 83, 97, 107, 111, 112, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 77, 77, 70, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 84, 121, 117, 107, 111, 107, 117, 102, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 75, 97, 119, 97, 115, 101, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 72, 105, 107, 97, 122, 101, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 34, 44, 10, 9, 34, 115, 75, 105, 110, 115, 121, 111, 117, 104, 111, 117, 77, 105, 100, 111, 107, 117, 70, 108, 103, 34, 58, 34, 34, 44, 10, 9, 34, 115, 85, 114, 108, 82, 101, 113, 117, 101, 115, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 85, 114, 108, 69, 118, 101, 110, 116, 34, 58, 34, 34, 10, 125, 10, 10}},
			arg1:      context.Background(),
			arg2:      LoginRequest{UserId: "user-id", Password: "password"},
			want1: &LoginResponse{
				CommonResponse: CommonResponse{
					No:           1,
					SendDate:     time.Date(2022, 3, 1, 11, 7, 17, 650000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 1, 11, 7, 17, 623000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeLoginResponse,
				},
				ResultCode:                "10031",
				ResultText:                "ユーザIDか暗証番号をお間違えです。ご確認の上、再度ご入力下さい。なお、お間違えの回数が弊社規定回数を超えると、ログイン停止となりますのでご注意下さい。(ログイン停止の解除は、コールセンターまでお電話下さい。)",
				AccountType:               "",
				SecondPasswordOmit:        false,
				LastLoginDateTime:         time.Time{},
				GeneralAccount:            false,
				SafekeepingAccount:        false,
				TransferAccount:           false,
				ForeignAccount:            false,
				MRFAccount:                false,
				StockSpecificAccount:      SpecificAccountTypeUnspecified,
				MarginSpecificAccount:     SpecificAccountTypeUnspecified,
				InvestmentSpecificAccount: "",
				DividendAccount:           false,
				SpecificAccount:           false,
				MarginAccount:             false,
				FutureOptionAccount:       false,
				MMFAccount:                false,
				ChinaForeignAccount:       false,
				FXAccount:                 false,
				NISAAccount:               false,
				UnreadDocument:            false,
				RequestURL:                "",
				EventURL:                  "",
			},
			want2: nil,
		},
		{name: "リクエストでエラーが返されたらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get2: StatusNotOkErr},
			arg1:      context.Background(),
			arg2:      LoginRequest{},
			want1:     nil,
			want2:     StatusNotOkErr},
		{name: "レスポンスのUnmarshalに失敗したらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get1: []byte{}},
			arg1:      context.Background(),
			arg2:      LoginRequest{},
			want1:     nil,
			want2:     UnmarshalFailedErr},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			req := &client{clock: test.clock, requester: test.requester}
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
