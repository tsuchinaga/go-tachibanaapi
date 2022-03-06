package tachibana

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func Test_StockMasterRequest_request(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request StockMasterRequest
		arg1    int64
		arg2    time.Time
		want1   stockMasterRequest
	}{
		{name: "カラム指定がnilなら全指定になる",
			request: StockMasterRequest{
				Columns: nil,
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 4, 14, 0, 0, 0, time.Local),
			want1: stockMasterRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 4, 14, 0, 0, 0, time.Local)},
					FeatureType:    FeatureTypeMasterData,
					ResponseFormat: commonResponseFormat,
				},
				TargetFeature: fmt.Sprintf(`"%s"`, FeatureTypeStockMaster),
				Columns:       "sIssueCode,sIssueName,sIssueNameRyaku,sIssueNameKana,sIssueNameEizi,sTokuteiF,sHikazeiC,sZyouzyouHakkouKabusu,sKenriotiFlag,sKenritukiSaisyuDay,sZyouzyouNyusatuC,sNyusatuKaizyoDay,sNyusatuDay,sBaibaiTani,sBaibaiTaniYoku,sBaibaiTeisiC,sHakkouKaisiDay,sHakkouSaisyuDay,sKessanC,sKessanDay,sZyouzyouOutouDay,sNiruiKizituC,sOogutiKabusu,sOogutiKingmaker,sBadenpyouOutputYNC,sHosyoukinDaiyouKakeme,sDaiyouHyoukaTanka,sKikoSankaC,sKarikessaiC,sYusenSizyou,sMukigenC,sGyousyuCode,sGyousyuName,sSorC,sCreateDate,sUpdateDate,sUpdateNumber",
			}},
		{name: "カラム指定がなければ全指定になる",
			request: StockMasterRequest{
				Columns: []StockMasterColumn{},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 4, 14, 0, 0, 0, time.Local),
			want1: stockMasterRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 4, 14, 0, 0, 0, time.Local)},
					FeatureType:    FeatureTypeMasterData,
					ResponseFormat: commonResponseFormat,
				},
				TargetFeature: fmt.Sprintf(`"%s"`, FeatureTypeStockMaster),
				Columns:       "sIssueCode,sIssueName,sIssueNameRyaku,sIssueNameKana,sIssueNameEizi,sTokuteiF,sHikazeiC,sZyouzyouHakkouKabusu,sKenriotiFlag,sKenritukiSaisyuDay,sZyouzyouNyusatuC,sNyusatuKaizyoDay,sNyusatuDay,sBaibaiTani,sBaibaiTaniYoku,sBaibaiTeisiC,sHakkouKaisiDay,sHakkouSaisyuDay,sKessanC,sKessanDay,sZyouzyouOutouDay,sNiruiKizituC,sOogutiKabusu,sOogutiKingmaker,sBadenpyouOutputYNC,sHosyoukinDaiyouKakeme,sDaiyouHyoukaTanka,sKikoSankaC,sKarikessaiC,sYusenSizyou,sMukigenC,sGyousyuCode,sGyousyuName,sSorC,sCreateDate,sUpdateDate,sUpdateNumber",
			}},
		{name: "変換できる",
			request: StockMasterRequest{
				Columns: []StockMasterColumn{StockMasterColumnCode, StockMasterColumnName},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 4, 14, 0, 0, 0, time.Local),
			want1: stockMasterRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 4, 14, 0, 0, 0, time.Local)},
					FeatureType:    FeatureTypeMasterData,
					ResponseFormat: commonResponseFormat,
				},
				TargetFeature: fmt.Sprintf(`"%s"`, FeatureTypeStockMaster),
				Columns:       "sIssueCode,sIssueName",
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

func Test_stockMasterResponse_response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		response stockMasterResponse
		want1    StockMasterResponse
	}{
		{name: "変換できる",
			response: stockMasterResponse{
				commonResponse: commonResponse{
					No:           2,
					SendDate:     RequestTime{Time: time.Date(2022, 3, 6, 10, 42, 21, 909000000, time.Local)},
					ReceiveDate:  RequestTime{Time: time.Date(2022, 3, 6, 10, 42, 21, 759000000, time.Local)},
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMasterData,
				},
				StockMasters: []stockMaster{
					{Code: "1475", Name: "ｉシェアーズＴＯＰＩＸ"},
					{Code: "1476", Name: "ｉシェアーズＪリート"},
				},
			},
			want1: StockMasterResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 6, 10, 42, 21, 909000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 6, 10, 42, 21, 759000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMasterData,
				},
				StockMasters: []StockMaster{
					{Code: "1475", Name: "ｉシェアーズＴＯＰＩＸ"},
					{Code: "1476", Name: "ｉシェアーズＪリート"},
				},
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

func Test_client_StockMaster(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		clock  *testClock
		status int
		body   []byte
		arg1   context.Context
		arg2   *Session
		arg3   StockMasterRequest
		want1  *StockMasterResponse
		want2  error
	}{
		{name: "正常レスポンスをパースして返せる",
			clock:  &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			status: http.StatusOK,
			body:   []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 54, 45, 49, 48, 58, 52, 50, 58, 50, 49, 46, 57, 48, 57, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 54, 45, 49, 48, 58, 52, 50, 58, 50, 49, 46, 55, 53, 57, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 77, 102, 100, 115, 71, 101, 116, 77, 97, 115, 116, 101, 114, 68, 97, 116, 97, 34, 44, 34, 67, 76, 77, 73, 115, 115, 117, 101, 77, 115, 116, 75, 97, 98, 117, 34, 58, 91, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 53, 34, 44, 34, 115, 73, 115, 115, 117, 101, 78, 97, 109, 101, 34, 58, 34, 130, 137, 131, 86, 131, 70, 131, 65, 129, 91, 131, 89, 130, 115, 130, 110, 130, 111, 130, 104, 130, 119, 34, 125, 44, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 54, 34, 44, 34, 115, 73, 115, 115, 117, 101, 78, 97, 109, 101, 34, 58, 34, 130, 137, 131, 86, 131, 70, 131, 65, 129, 91, 131, 89, 130, 105, 131, 138, 129, 91, 131, 103, 34, 125, 93, 125},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   StockMasterRequest{},
			want1: &StockMasterResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 6, 10, 42, 21, 909000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 6, 10, 42, 21, 759000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					FeatureType:  FeatureTypeMasterData,
				},
				StockMasters: []StockMaster{
					{Code: "1475", Name: "ｉシェアーズＴＯＰＩＸ"},
					{Code: "1476", Name: "ｉシェアーズＪリート"},
				},
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusOK,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   nil,
			arg3:   StockMasterRequest{},
			want1:  nil,
			want2:  NilArgumentErr},
		{name: "200 OK以外が返ったらエラー",
			clock:  &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			status: http.StatusInternalServerError,
			body:   []byte{},
			arg1:   context.Background(),
			arg2:   &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:   StockMasterRequest{},
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
			got1, got2 := client.StockMaster(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_StockMaster_Execute(t *testing.T) {
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

	got3, got4 := client.StockMaster(context.Background(), session, StockMasterRequest{
		Columns: []StockMasterColumn{StockMasterColumnCode, StockMasterColumnName},
	})

	log.Printf("%+v, %+v\n", got3, got4)
}
