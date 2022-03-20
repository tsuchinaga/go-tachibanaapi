package tachibana

import (
	"context"
	"errors"
	"log"
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
					MessageType:    MessageTypeMasterData,
					ResponseFormat: commonResponseFormat,
				},
				TargetFeature: string(MessageTypeStockMaster),
				Columns:       "",
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
					MessageType:    MessageTypeMasterData,
					ResponseFormat: commonResponseFormat,
				},
				TargetFeature: string(MessageTypeStockMaster),
				Columns:       "",
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
					MessageType:    MessageTypeMasterData,
					ResponseFormat: commonResponseFormat,
				},
				TargetFeature: string(MessageTypeStockMaster),
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
					MessageType:  MessageTypeMasterData,
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
					MessageType:  MessageTypeMasterData,
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
		name      string
		clock     *testClock
		requester *testRequester
		arg1      context.Context
		arg2      *Session
		arg3      StockMasterRequest
		want1     *StockMasterResponse
		want2     error
	}{
		{name: "正常レスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 54, 45, 49, 48, 58, 52, 50, 58, 50, 49, 46, 57, 48, 57, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 54, 45, 49, 48, 58, 52, 50, 58, 50, 49, 46, 55, 53, 57, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 77, 102, 100, 115, 71, 101, 116, 77, 97, 115, 116, 101, 114, 68, 97, 116, 97, 34, 44, 34, 67, 76, 77, 73, 115, 115, 117, 101, 77, 115, 116, 75, 97, 98, 117, 34, 58, 91, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 53, 34, 44, 34, 115, 73, 115, 115, 117, 101, 78, 97, 109, 101, 34, 58, 34, 239, 189, 137, 227, 130, 183, 227, 130, 167, 227, 130, 162, 227, 131, 188, 227, 130, 186, 239, 188, 180, 239, 188, 175, 239, 188, 176, 239, 188, 169, 239, 188, 184, 34, 125, 44, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 54, 34, 44, 34, 115, 73, 115, 115, 117, 101, 78, 97, 109, 101, 34, 58, 34, 239, 189, 137, 227, 130, 183, 227, 130, 167, 227, 130, 162, 227, 131, 188, 227, 130, 186, 239, 188, 170, 227, 131, 170, 227, 131, 188, 227, 131, 136, 34, 125, 93, 125}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockMasterRequest{},
			want1: &StockMasterResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 6, 10, 42, 21, 909000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 6, 10, 42, 21, 759000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeMasterData,
				},
				StockMasters: []StockMaster{
					{Code: "1475", Name: "ｉシェアーズＴＯＰＩＸ"},
					{Code: "1476", Name: "ｉシェアーズＪリート"},
				},
			},
			want2: nil},
		{name: "全項目含む正常レスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 54, 45, 50, 48, 58, 52, 50, 58, 49, 50, 46, 48, 54, 53, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 54, 45, 50, 48, 58, 52, 50, 58, 49, 49, 46, 57, 48, 51, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 77, 102, 100, 115, 71, 101, 116, 77, 97, 115, 116, 101, 114, 68, 97, 116, 97, 34, 44, 34, 67, 76, 77, 73, 115, 115, 117, 101, 77, 115, 116, 75, 97, 98, 117, 34, 58, 91, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 56, 48, 53, 56, 34, 44, 34, 115, 73, 115, 115, 117, 101, 78, 97, 109, 101, 34, 58, 34, 228, 184, 137, 232, 143, 177, 229, 149, 134, 228, 186, 139, 34, 44, 34, 115, 73, 115, 115, 117, 101, 78, 97, 109, 101, 82, 121, 97, 107, 117, 34, 58, 34, 228, 184, 137, 232, 143, 177, 229, 149, 134, 34, 44, 34, 115, 73, 115, 115, 117, 101, 78, 97, 109, 101, 75, 97, 110, 97, 34, 58, 34, 227, 131, 159, 227, 131, 132, 227, 131, 147, 227, 130, 183, 32, 32, 227, 130, 183, 227, 131, 168, 227, 130, 166, 227, 130, 184, 34, 44, 34, 115, 73, 115, 115, 117, 101, 78, 97, 109, 101, 69, 105, 122, 105, 34, 58, 34, 77, 73, 84, 66, 73, 83, 73, 34, 44, 34, 115, 84, 111, 107, 117, 116, 101, 105, 70, 34, 58, 34, 49, 34, 44, 34, 115, 72, 105, 107, 97, 122, 101, 105, 67, 34, 58, 34, 49, 34, 44, 34, 115, 90, 121, 111, 117, 122, 121, 111, 117, 72, 97, 107, 107, 111, 117, 75, 97, 98, 117, 115, 117, 34, 58, 34, 49, 52, 56, 53, 55, 50, 51, 51, 53, 49, 34, 44, 34, 115, 75, 101, 110, 114, 105, 111, 116, 105, 70, 108, 97, 103, 34, 58, 34, 48, 34, 44, 34, 115, 75, 101, 110, 114, 105, 116, 117, 107, 105, 83, 97, 105, 115, 121, 117, 68, 97, 121, 34, 58, 34, 50, 48, 49, 54, 48, 57, 50, 55, 34, 44, 34, 115, 90, 121, 111, 117, 122, 121, 111, 117, 78, 121, 117, 115, 97, 116, 117, 67, 34, 58, 34, 32, 34, 44, 34, 115, 78, 121, 117, 115, 97, 116, 117, 75, 97, 105, 122, 121, 111, 68, 97, 121, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 78, 121, 117, 115, 97, 116, 117, 68, 97, 121, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 66, 97, 105, 98, 97, 105, 84, 97, 110, 105, 34, 58, 34, 49, 48, 48, 34, 44, 34, 115, 66, 97, 105, 98, 97, 105, 84, 97, 110, 105, 89, 111, 107, 117, 34, 58, 34, 49, 48, 48, 34, 44, 34, 115, 66, 97, 105, 98, 97, 105, 84, 101, 105, 115, 105, 67, 34, 58, 34, 32, 34, 44, 34, 115, 72, 97, 107, 107, 111, 117, 75, 97, 105, 115, 105, 68, 97, 121, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 72, 97, 107, 107, 111, 117, 83, 97, 105, 115, 121, 117, 68, 97, 121, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 75, 101, 115, 115, 97, 110, 67, 34, 58, 34, 48, 49, 34, 44, 34, 115, 75, 101, 115, 115, 97, 110, 68, 97, 121, 34, 58, 34, 50, 48, 49, 54, 48, 57, 50, 55, 34, 44, 34, 115, 90, 121, 111, 117, 122, 121, 111, 117, 79, 117, 116, 111, 117, 68, 97, 121, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 78, 105, 114, 117, 105, 75, 105, 122, 105, 116, 117, 67, 34, 58, 34, 48, 34, 44, 34, 115, 79, 111, 103, 117, 116, 105, 75, 97, 98, 117, 115, 117, 34, 58, 34, 48, 34, 44, 34, 115, 66, 97, 100, 101, 110, 112, 121, 111, 117, 79, 117, 116, 112, 117, 116, 89, 78, 67, 34, 58, 34, 50, 34, 44, 34, 115, 72, 111, 115, 121, 111, 117, 107, 105, 110, 68, 97, 105, 121, 111, 117, 75, 97, 107, 101, 109, 101, 34, 58, 34, 56, 48, 46, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 68, 97, 105, 121, 111, 117, 72, 121, 111, 117, 107, 97, 84, 97, 110, 107, 97, 34, 58, 34, 52, 48, 51, 56, 46, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 75, 105, 107, 111, 83, 97, 110, 107, 97, 67, 34, 58, 34, 49, 34, 44, 34, 115, 75, 97, 114, 105, 107, 101, 115, 115, 97, 105, 67, 34, 58, 34, 32, 34, 44, 34, 115, 89, 117, 115, 101, 110, 83, 105, 122, 121, 111, 117, 34, 58, 34, 48, 48, 34, 44, 34, 115, 77, 117, 107, 105, 103, 101, 110, 67, 34, 58, 34, 32, 34, 44, 34, 115, 71, 121, 111, 117, 115, 121, 117, 67, 111, 100, 101, 34, 58, 34, 54, 48, 53, 48, 34, 44, 34, 115, 71, 121, 111, 117, 115, 121, 117, 78, 97, 109, 101, 34, 58, 34, 229, 141, 184, 229, 163, 178, 230, 165, 173, 34, 44, 34, 115, 83, 111, 114, 67, 34, 58, 34, 32, 34, 44, 34, 115, 67, 114, 101, 97, 116, 101, 68, 97, 116, 101, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 85, 112, 100, 97, 116, 101, 68, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 51, 48, 52, 49, 57, 48, 56, 53, 56, 34, 44, 34, 115, 85, 112, 100, 97, 116, 101, 78, 117, 109, 98, 101, 114, 34, 58, 34, 49, 52, 50, 55, 34, 125, 93, 125}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockMasterRequest{},
			want1: &StockMasterResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 6, 20, 42, 12, 65000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 6, 20, 42, 11, 903000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeMasterData,
				},
				StockMasters: []StockMaster{
					{
						Code:                 "8058",
						Name:                 "三菱商事",
						ShortName:            "三菱商",
						Kana:                 "ミツビシ  シヨウジ",
						Alphabet:             "MITBISI",
						SpecificTarget:       true,
						TaxFree:              TaxFreeValid,
						SharedStocks:         1_485_723_351,
						ExRight:              ExRightTypeNothing,
						LastRightDay:         time.Date(2016, 9, 27, 0, 0, 0, 0, time.Local),
						ListingType:          ListingTypeUnUsed,
						ReleaseTradingDate:   time.Time{},
						TradingDate:          time.Time{},
						TradingUnit:          100,
						NextTradingUnit:      100,
						StopTradingType:      StopTradingTypeUnUsed,
						StartPublicationDate: time.Time{},
						LastPublicationDate:  time.Time{},
						SettlementType:       SettlementTypeCapitalIncrease,
						SettlementDate:       time.Date(2016, 9, 27, 0, 0, 0, 0, time.Local),
						ListingDate:          time.Time{},
						ExpireDate2Type:      "0",
						LargeUnit:            0,
						LargeAmount:          0,
						OutputTicketType:     "2",
						DepositAmount:        80,
						DepositValuation:     4038,
						OrganizationType:     "1",
						ProvisionalType:      " ",
						PrimaryExchange:      ExchangeToushou,
						IndefinitePeriodType: " ",
						IndustryCode:         "6050",
						IndustryName:         "卸売業",
						SORTargetType:        " ",
						CreateDateTime:       time.Time{},
						UpdateDateTime:       time.Date(2022, 3, 4, 19, 8, 58, 0, time.Local),
						UpdateNumber:         "1427",
					},
				},
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			arg1:  context.Background(),
			arg2:  nil,
			arg3:  StockMasterRequest{},
			want1: nil,
			want2: NilArgumentErr},
		{name: "リクエストでエラーが返されたらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get2: StatusNotOkErr},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockMasterRequest{},
			want1:     nil,
			want2:     StatusNotOkErr},
		{name: "レスポンスのUnmarshalに失敗したらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get1: []byte{}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockMasterRequest{},
			want1:     nil,
			want2:     UnmarshalFailedErr},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &client{clock: test.clock, requester: test.requester}
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
