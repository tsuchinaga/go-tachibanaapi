package tachibana

import (
	"context"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func Test_tachibana_authURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  Environment
		arg2  ApiVersion
		want1 string
	}{
		{name: "環境を指定しなければ本番環境",
			arg1:  EnvironmentUnspecified,
			arg2:  ApiVersionV4R2,
			want1: "https://kabuka.e-shiten.jp/e_api_v4r2/auth/"},
		{name: "本番環境を指定すれば本番環境",
			arg1:  EnvironmentProduction,
			arg2:  ApiVersionV4R2,
			want1: "https://kabuka.e-shiten.jp/e_api_v4r2/auth/"},
		{name: "デモ環境を指定すればデモ環境",
			arg1:  EnvironmentDemo,
			arg2:  ApiVersionV4R2,
			want1: "https://demo-kabuka.e-shiten.jp/e_api_v4r2/auth/"},
		{name: "APIバージョンを指定しなければ最新バージョン",
			arg1:  EnvironmentProduction,
			arg2:  ApiVersionUnspecified,
			want1: "https://kabuka.e-shiten.jp/e_api_v4r2/auth/"},
		{name: "最新のAPIバージョンを指定すれば最新バージョン",
			arg1:  EnvironmentProduction,
			arg2:  ApiVersionLatest,
			want1: "https://kabuka.e-shiten.jp/e_api_v4r2/auth/"},
		{name: "バージョンV4R1を指定すればV4R1",
			arg1:  EnvironmentProduction,
			arg2:  ApiVersionV4R1,
			want1: "https://kabuka.e-shiten.jp/e_api_v4r1/auth/"},
		{name: "バージョンV4R2を指定すればV4R2",
			arg1:  EnvironmentProduction,
			arg2:  ApiVersionV4R2,
			want1: "https://kabuka.e-shiten.jp/e_api_v4r2/auth/"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			client := &client{}
			got1 := client.authURL(test.arg1, test.arg2)
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}

func Test_client_encode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  string
		want1 string
		want2 error
	}{
		{name: "エンコードできる",
			arg1:  `{"p_no":"1","p_sd_date":"2020.07.01-10:00:00.000","sCLMID":"CLMAuthLoginRequest","sUserId":"login-id","sPassword":"pswd","japanese":"ひらがなカタカナ漢字"}`,
			want1: "%7B%22p_no%22%3A%221%22%2C%22p_sd_date%22%3A%222020.07.01-10%3A00%3A00.000%22%2C%22sCLMID%22%3A%22CLMAuthLoginRequest%22%2C%22sUserId%22%3A%22login-id%22%2C%22sPassword%22%3A%22pswd%22%2C%22japanese%22%3A%22%82%D0%82%E7%82%AA%82%C8%83J%83%5E%83J%83i%8A%BF%8E%9A%22%7D",
			want2: nil},
		{name: "UTF-8からShift-JISにエンコードできない文字列を含むとエラー",
			arg1:  "\u1234",
			want1: "",
			want2: EncodeErr},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			client := &client{}
			got1, got2 := client.encode(test.arg1)
			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1,
					test.want2,
					got1, got2)
			}
		})
	}
}

func Test_client_decode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  string
		want1 string
		want2 error
	}{
		{name: "デコードできる",
			arg1:  string([]byte{123, 34, 49, 55, 53, 34, 58, 34, 34, 44, 34, 49, 55, 55, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 48, 56, 58, 51, 52, 58, 52, 54, 46, 53, 55, 51, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 48, 56, 58, 51, 52, 58, 52, 54, 46, 53, 54, 53, 34, 44, 34, 49, 55, 52, 34, 58, 34, 45, 49, 34, 44, 34, 49, 55, 51, 34, 58, 34, 136, 248, 144, 148, 131, 71, 131, 137, 129, 91, 129, 66, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 105, 110, 82, 101, 113, 117, 101, 115, 116, 34, 125}),
			want1: `{"175":"","177":"2022.02.24-08:34:46.573","176":"2022.02.24-08:34:46.565","174":"-1","173":"引数エラー。","192":"CLMAuthLoginRequest"}`,
			want2: nil},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			client := &client{}
			got1, got2 := client.decode(test.arg1)
			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1,
					test.want2,
					got1, got2)
			}
		})
	}
}

func Test_NewClient(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  Environment
		arg2  ApiVersion
		want1 Client
	}{
		{name: "本番へのクライアントの生成",
			arg1: EnvironmentProduction,
			arg2: ApiVersionLatest,
			want1: &client{
				clock: newClock(),
				env:   EnvironmentProduction,
				ver:   ApiVersionLatest,
				auth:  "https://kabuka.e-shiten.jp/e_api_v4r2/auth/",
			}},
		{name: "デモへのクライアントの生成",
			arg1: EnvironmentDemo,
			arg2: ApiVersionV4R2,
			want1: &client{
				clock: newClock(),
				env:   EnvironmentDemo,
				ver:   ApiVersionLatest,
				auth:  "https://demo-kabuka.e-shiten.jp/e_api_v4r2/auth/",
			}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got1 := NewClient(test.arg1, test.arg2)
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}

func Test_client_parseResponse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		arg1         []byte
		arg2         loginResponse
		wantResponse loginResponse
		hasError     bool
	}{
		{name: "ログインレスポンスをパース出来る",
			arg1: []byte{123, 10, 9, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 48, 58, 53, 55, 58, 49, 54, 46, 55, 54, 56, 34, 44, 10, 9, 34, 112, 95, 110, 111, 34, 58, 34, 49, 34, 44, 10, 9, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 49, 45, 49, 48, 58, 53, 55, 58, 49, 54, 46, 54, 57, 51, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 10, 9, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 10, 9, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 105, 110, 65, 99, 107, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 67, 111, 100, 101, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 82, 101, 115, 117, 108, 116, 84, 101, 120, 116, 34, 58, 34, 34, 44, 10, 9, 34, 115, 90, 121, 111, 117, 116, 111, 101, 107, 105, 75, 97, 122, 101, 105, 67, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 83, 101, 99, 111, 110, 100, 80, 97, 115, 115, 119, 111, 114, 100, 79, 109, 105, 116, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 76, 97, 115, 116, 76, 111, 103, 105, 110, 68, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 51, 48, 49, 49, 48, 53, 54, 52, 49, 34, 44, 10, 9, 34, 115, 83, 111, 103, 111, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 72, 111, 103, 111, 65, 100, 117, 107, 97, 114, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 70, 117, 114, 105, 107, 97, 101, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 71, 97, 105, 107, 111, 107, 117, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 77, 82, 70, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 71, 101, 110, 98, 117, 116, 117, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 83, 105, 110, 121, 111, 117, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 84, 111, 117, 115, 105, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 72, 97, 105, 116, 111, 117, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 84, 111, 107, 117, 116, 101, 105, 75, 97, 110, 114, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 83, 105, 110, 121, 111, 117, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 49, 34, 44, 10, 9, 34, 115, 83, 97, 107, 111, 112, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 77, 77, 70, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 84, 121, 117, 107, 111, 107, 117, 102, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 75, 97, 119, 97, 115, 101, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 72, 105, 107, 97, 122, 101, 105, 75, 111, 117, 122, 97, 75, 117, 98, 117, 110, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 75, 105, 110, 115, 121, 111, 117, 104, 111, 117, 77, 105, 100, 111, 107, 117, 70, 108, 103, 34, 58, 34, 48, 34, 44, 10, 9, 34, 115, 85, 114, 108, 82, 101, 113, 117, 101, 115, 116, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 107, 97, 98, 117, 107, 97, 46, 101, 45, 115, 104, 105, 116, 101, 110, 46, 106, 112, 47, 101, 95, 97, 112, 105, 95, 118, 52, 114, 50, 47, 114, 101, 113, 117, 101, 115, 116, 47, 78, 106, 107, 48, 77, 84, 89, 49, 78, 122, 69, 119, 77, 68, 69, 119, 77, 121, 48, 120, 77, 106, 77, 116, 78, 84, 77, 48, 78, 122, 89, 61, 47, 34, 44, 10, 9, 34, 115, 85, 114, 108, 69, 118, 101, 110, 116, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 107, 97, 98, 117, 107, 97, 46, 101, 45, 115, 104, 105, 116, 101, 110, 46, 106, 112, 47, 101, 95, 97, 112, 105, 95, 118, 52, 114, 50, 47, 101, 118, 101, 110, 116, 47, 78, 106, 107, 48, 77, 84, 89, 49, 78, 122, 69, 119, 77, 68, 69, 119, 77, 121, 48, 120, 77, 106, 77, 116, 78, 84, 77, 48, 78, 122, 89, 61, 47, 34, 10, 125, 10, 10},
			arg2: loginResponse{},
			wantResponse: loginResponse{
				commonResponse: commonResponse{
					No:           1,
					SendDate:     RequestTime{Time: time.Date(2022, 3, 1, 10, 57, 16, 768000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 3, 1, 10, 57, 16, 693000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeLoginResponse,
				},
				ResultCode:                "0",
				ResultText:                "",
				AccountType:               AccountTypeSpecific,
				SecondPasswordOmit:        NumberBoolFalse,
				LastLoginDateTime:         YmdHms{Time: time.Date(2022, 3, 1, 10, 56, 41, 0, time.Local)},
				GeneralAccount:            NumberBoolTrue,
				SafekeepingAccount:        NumberBoolTrue,
				TransferAccount:           NumberBoolTrue,
				ForeignAccount:            NumberBoolTrue,
				MRFAccount:                NumberBoolFalse,
				StockSpecificAccount:      SpecificAccountTypeNothing,
				MarginSpecificAccount:     SpecificAccountTypeNothing,
				InvestmentSpecificAccount: SpecificAccountTypeNothing,
				DividendAccount:           NumberBoolFalse,
				SpecificAccount:           NumberBoolTrue,
				MarginAccount:             NumberBoolTrue,
				FutureOptionAccount:       NumberBoolFalse,
				MMFAccount:                NumberBoolFalse,
				ChinaForeignAccount:       NumberBoolFalse,
				FXAccount:                 NumberBoolFalse,
				NISAAccount:               NumberBoolFalse,
				UnreadDocument:            NumberBoolFalse,
				RequestURL:                "https://kabuka.e-shiten.jp/e_api_v4r2/request/Njk0MTY1NzEwMDEwMy0xMjMtNTM0NzY=/",
				EventURL:                  "https://kabuka.e-shiten.jp/e_api_v4r2/event/Njk0MTY1NzEwMDEwMy0xMjMtNTM0NzY=/",
			},
			hasError: false},
		{name: "パース出来なければエラー",
			arg1:         []byte{123, 34, 49, 55, 55, 34, 0, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 48, 56, 58, 51, 54, 58, 49, 55, 46, 55, 55, 56, 34, 44, 34, 49, 55, 53, 34, 58, 34, 49, 34, 44, 34, 49, 55, 54, 34, 58, 34, 50, 48, 50, 50, 46, 48, 50, 46, 50, 52, 45, 48, 56, 58, 51, 54, 58, 49, 55, 46, 55, 50, 53, 34, 44, 34, 49, 55, 52, 34, 58, 34, 48, 34, 44, 34, 49, 55, 51, 34, 58, 34, 34, 44, 34, 49, 57, 50, 34, 58, 34, 67, 76, 77, 65, 117, 116, 104, 76, 111, 103, 105, 110, 65, 99, 107, 34, 44, 34, 53, 51, 52, 34, 58, 34, 48, 34, 44, 34, 53, 51, 53, 34, 58, 34, 34, 44, 34, 55, 52, 52, 34, 58, 34, 49, 34, 44, 34, 53, 52, 53, 34, 58, 34, 48, 34, 44, 34, 52, 48, 49, 34, 58, 34, 50, 48, 50, 50, 48, 50, 50, 52, 48, 56, 51, 51, 49, 55, 34, 44, 34, 53, 56, 48, 34, 58, 34, 49, 34, 44, 34, 50, 56, 55, 34, 58, 34, 49, 34, 44, 34, 50, 51, 50, 34, 58, 34, 49, 34, 44, 34, 50, 51, 52, 34, 58, 34, 49, 34, 44, 34, 52, 48, 52, 34, 58, 34, 48, 34, 44, 34, 54, 52, 53, 34, 58, 34, 49, 34, 44, 34, 54, 52, 54, 34, 58, 34, 48, 34, 44, 34, 54, 52, 55, 34, 58, 34, 49, 34, 44, 34, 54, 52, 50, 34, 58, 34, 48, 34, 44, 34, 54, 52, 52, 34, 58, 34, 49, 34, 44, 34, 53, 54, 53, 34, 58, 34, 48, 34, 44, 34, 53, 52, 50, 34, 58, 34, 48, 34, 44, 34, 52, 48, 51, 34, 58, 34, 48, 34, 44, 34, 54, 54, 48, 34, 58, 34, 48, 34, 44, 34, 51, 52, 56, 34, 58, 34, 48, 34, 44, 34, 50, 56, 49, 34, 58, 34, 48, 34, 44, 34, 51, 55, 52, 34, 58, 34, 48, 34, 44, 34, 54, 56, 57, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 107, 97, 98, 117, 107, 97, 46, 101, 45, 115, 104, 105, 116, 101, 110, 46, 106, 112, 47, 101, 95, 97, 112, 105, 95, 118, 52, 114, 50, 47, 114, 101, 113, 117, 101, 115, 116, 47, 78, 122, 73, 49, 77, 84, 99, 122, 78, 106, 65, 52, 77, 106, 81, 119, 77, 105, 48, 120, 77, 106, 77, 116, 78, 106, 77, 120, 77, 84, 99, 61, 47, 34, 44, 34, 54, 56, 56, 34, 58, 34, 104, 116, 116, 112, 115, 58, 47, 47, 107, 97, 98, 117, 107, 97, 46, 101, 45, 115, 104, 105, 116, 101, 110, 46, 106, 112, 47, 101, 95, 97, 112, 105, 95, 118, 52, 114, 50, 47, 101, 118, 101, 110, 116, 47, 78, 122, 73, 49, 77, 84, 99, 122, 78, 106, 65, 52, 77, 106, 81, 119, 77, 105, 48, 120, 77, 106, 77, 116, 78, 106, 77, 120, 77, 84, 99, 61, 47, 34, 125},
			arg2:         loginResponse{},
			wantResponse: loginResponse{},
			hasError:     true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			client := &client{}
			got1 := client.parseResponse(test.arg1, &test.arg2)
			if !reflect.DeepEqual(test.wantResponse, test.arg2) || (got1 != nil) != test.hasError {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.wantResponse, test.hasError, test.arg2, got1)
			}
		})
	}
}

func Test_client_get(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		status       int
		headers      map[string]string
		body         []byte
		overwriteUrl bool
		arg1         context.Context
		arg2         string
		arg3         interface{}
		arg4         interface{}
		hasError     bool
	}{
		{name: "json.Marshalでエラーがでたらエラーを返す",
			arg1:     context.Background(),
			arg2:     "",
			arg3:     []float64{math.Inf(1)},
			arg4:     nil,
			hasError: true},
		{name: "encodeでエラーがでたらエラーを返す",
			arg1:     context.Background(),
			arg2:     "",
			arg3:     loginRequest{UserId: "\u1234"},
			arg4:     nil,
			hasError: true},
		{name: "contextが設定されていなかったらエラー",
			arg1:     nil,
			arg2:     "http://example",
			arg3:     loginRequest{},
			arg4:     nil,
			hasError: true},
		{name: "urlが設定されておらずリクエストできなかったらエラー",
			overwriteUrl: false,
			arg1:         context.Background(),
			arg2:         "",
			arg3:         loginRequest{},
			arg4:         nil,
			hasError:     true},
		{name: "bodyの読み込みに失敗したらエラー",
			status:       http.StatusOK,
			headers:      map[string]string{"Content-Length": "1"},
			body:         nil,
			overwriteUrl: true,
			arg1:         context.Background(),
			arg2:         "",
			arg3:         loginRequest{},
			arg4:         nil,
			hasError:     true},
		{name: "レスポンスをパースする先がnilならエラー",
			status:       http.StatusOK,
			body:         nil,
			overwriteUrl: true,
			arg1:         context.Background(),
			arg2:         "",
			arg3:         loginRequest{},
			arg4:         nil,
			hasError:     true},
		{name: "statusがOKでなければエラー",
			status:       http.StatusInternalServerError,
			body:         nil,
			overwriteUrl: true,
			arg1:         context.Background(),
			arg2:         "",
			arg3:         loginRequest{},
			arg4:         nil,
			hasError:     true},
		{name: "エラーなく処理が終わればnilを返す",
			status:       http.StatusOK,
			body:         []byte(`{}`),
			overwriteUrl: true,
			arg1:         context.Background(),
			arg2:         "",
			arg3:         loginRequest{},
			arg4:         &loginResponse{},
			hasError:     false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				for k, v := range test.headers {
					w.Header().Set(k, v)
				}
				w.WriteHeader(test.status)
				_, _ = w.Write(test.body)
			})
			ts := httptest.NewServer(mux)
			defer ts.Close()
			if test.overwriteUrl {
				test.arg2 = ts.URL
			}

			client := &client{}
			got1 := client.get(test.arg1, test.arg2, test.arg3, test.arg4)
			if (got1 != nil) != test.hasError {
				t.Errorf("%s error\nerror: %+v\n", t.Name(), got1)
			}
		})
	}
}

func Test_commonResponse_response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		response commonResponse
		want1    CommonResponse
	}{
		{name: "変換できる",
			response: commonResponse{
				No:           2,
				SendDate:     RequestTime{Time: time.Date(2022, 2, 25, 10, 5, 30, 123456789, time.Local)},
				ReceiveDate:  RequestTime{Time: time.Date(2022, 2, 25, 10, 5, 31, 123456789, time.Local)},
				ErrorNo:      ErrorNoProblem,
				ErrorMessage: "",
				FeatureType:  FeatureTypeLoginResponse,
			},
			want1: CommonResponse{
				No:           2,
				SendDate:     time.Date(2022, 2, 25, 10, 5, 30, 123456789, time.Local),
				ReceiveDate:  time.Date(2022, 2, 25, 10, 5, 31, 123456789, time.Local),
				ErrorNo:      ErrorNoProblem,
				ErrorMessage: "",
				FeatureType:  FeatureTypeLoginResponse,
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

func Test_client_host(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  Environment
		want1 string
	}{
		{name: "envが本番なら本番のホストが返される",
			arg1:  EnvironmentProduction,
			want1: "kabuka.e-shiten.jp"},
		{name: "envがデモならデモのホストが返される",
			arg1:  EnvironmentDemo,
			want1: "demo-kabuka.e-shiten.jp"},
		{name: "envが未指定なら本番のホストが返される",
			arg1:  EnvironmentUnspecified,
			want1: "kabuka.e-shiten.jp"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			client := &client{}
			got1 := client.host(test.arg1)
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}
