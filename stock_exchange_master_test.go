package tachibana

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"
	"time"
)

func Test_StockExchangeMasterRequest_request(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request StockExchangeMasterRequest
		arg1    int64
		arg2    time.Time
		want1   stockExchangeMasterRequest
	}{
		{name: "カラム指定がnilなら全指定になる",
			request: StockExchangeMasterRequest{
				Columns: nil,
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 4, 14, 0, 0, 0, time.Local),
			want1: stockExchangeMasterRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 4, 14, 0, 0, 0, time.Local)},
					MessageType:    MessageTypeMasterData,
					ResponseFormat: commonResponseFormat,
				},
				TargetFeature: string(MessageTypeStockExchangeMaster),
				Columns:       "",
			}},
		{name: "カラム指定がなければ全指定になる",
			request: StockExchangeMasterRequest{
				Columns: []StockExchangeMasterColumn{},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 4, 14, 0, 0, 0, time.Local),
			want1: stockExchangeMasterRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 4, 14, 0, 0, 0, time.Local)},
					MessageType:    MessageTypeMasterData,
					ResponseFormat: commonResponseFormat,
				},
				TargetFeature: string(MessageTypeStockExchangeMaster),
				Columns:       "",
			}},
		{name: "変換できる",
			request: StockExchangeMasterRequest{
				Columns: []StockExchangeMasterColumn{StockExchangeMasterColumnIssueCode, StockExchangeMasterColumnExchange},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 4, 14, 0, 0, 0, time.Local),
			want1: stockExchangeMasterRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 4, 14, 0, 0, 0, time.Local)},
					MessageType:    MessageTypeMasterData,
					ResponseFormat: commonResponseFormat,
				},
				TargetFeature: string(MessageTypeStockExchangeMaster),
				Columns:       "sIssueCode,sZyouzyouSizyou",
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

func Test_client_StockExchangeMaster(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		clock     *testClock
		requester *testRequester
		arg1      context.Context
		arg2      *Session
		arg3      StockExchangeMasterRequest
		want1     *StockExchangeMasterResponse
		want2     error
	}{
		{name: "正常レスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 50, 51, 45, 48, 57, 58, 52, 50, 58, 49, 57, 46, 53, 53, 55, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 50, 51, 45, 48, 57, 58, 52, 50, 58, 49, 57, 46, 52, 50, 51, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 77, 102, 100, 115, 71, 101, 116, 77, 97, 115, 116, 101, 114, 68, 97, 116, 97, 34, 44, 34, 67, 76, 77, 73, 115, 115, 117, 101, 83, 105, 122, 121, 111, 117, 77, 115, 116, 75, 97, 98, 117, 34, 58, 91, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 51, 48, 49, 34, 44, 34, 115, 90, 121, 111, 117, 122, 121, 111, 117, 83, 105, 122, 121, 111, 117, 34, 58, 34, 48, 48, 34, 125, 44, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 51, 48, 53, 34, 44, 34, 115, 90, 121, 111, 117, 122, 121, 111, 117, 83, 105, 122, 121, 111, 117, 34, 58, 34, 48, 48, 34, 125, 93, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockExchangeMasterRequest{},
			want1: &StockExchangeMasterResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 23, 9, 42, 19, 557000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 23, 9, 42, 19, 423000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeMasterData,
				},
				StockExchangeMasters: []StockExchangeMaster{
					{IssueCode: "1301", Exchange: ExchangeToushou},
					{IssueCode: "1305", Exchange: ExchangeToushou},
				},
			},
			want2: nil},
		{name: "全項目含む正常レスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 50, 51, 45, 48, 57, 58, 52, 50, 58, 49, 57, 46, 53, 53, 55, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 50, 51, 45, 48, 57, 58, 52, 50, 58, 49, 57, 46, 52, 50, 51, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 77, 102, 100, 115, 71, 101, 116, 77, 97, 115, 116, 101, 114, 68, 97, 116, 97, 34, 44, 34, 67, 76, 77, 73, 115, 115, 117, 101, 83, 105, 122, 121, 111, 117, 77, 115, 116, 75, 97, 98, 117, 34, 58, 91, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 54, 57, 54, 50, 34, 44, 34, 115, 90, 121, 111, 117, 122, 121, 111, 117, 83, 105, 122, 121, 111, 117, 34, 58, 34, 48, 48, 34, 44, 34, 115, 83, 121, 115, 116, 101, 109, 67, 34, 58, 34, 49, 34, 44, 34, 115, 78, 101, 104, 97, 98, 97, 77, 105, 110, 34, 58, 34, 55, 54, 50, 46, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 78, 101, 104, 97, 98, 97, 77, 97, 120, 34, 58, 34, 49, 51, 54, 50, 46, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 73, 115, 115, 117, 101, 75, 117, 98, 117, 110, 67, 34, 58, 34, 32, 34, 44, 34, 115, 78, 101, 104, 97, 98, 97, 83, 105, 122, 121, 111, 117, 67, 34, 58, 34, 48, 48, 34, 44, 34, 115, 83, 105, 110, 121, 111, 117, 67, 34, 58, 34, 49, 34, 44, 34, 115, 83, 105, 110, 107, 105, 90, 121, 111, 117, 122, 121, 111, 117, 68, 97, 121, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 78, 101, 104, 97, 98, 97, 75, 105, 103, 101, 110, 68, 97, 121, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 78, 101, 104, 97, 98, 97, 75, 105, 115, 101, 105, 67, 34, 58, 34, 32, 34, 44, 34, 115, 78, 101, 104, 97, 98, 97, 75, 105, 115, 101, 105, 84, 105, 34, 58, 34, 48, 46, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 78, 101, 104, 97, 98, 97, 67, 104, 101, 99, 107, 75, 97, 104, 105, 67, 34, 58, 34, 49, 34, 44, 34, 115, 73, 115, 115, 117, 101, 66, 117, 98, 101, 116, 117, 67, 34, 58, 34, 49, 34, 44, 34, 115, 90, 101, 110, 122, 105, 116, 117, 79, 119, 97, 114, 105, 110, 101, 34, 58, 34, 49, 48, 54, 50, 46, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 78, 101, 104, 97, 98, 97, 83, 97, 110, 115, 121, 117, 116, 117, 83, 105, 122, 121, 111, 117, 67, 34, 58, 34, 32, 34, 44, 34, 115, 73, 115, 115, 117, 101, 75, 105, 115, 101, 105, 49, 67, 34, 58, 34, 32, 34, 44, 34, 115, 73, 115, 115, 117, 101, 75, 105, 115, 101, 105, 50, 67, 34, 58, 34, 32, 34, 44, 34, 115, 90, 121, 111, 117, 122, 121, 111, 117, 75, 117, 98, 117, 110, 34, 58, 34, 48, 49, 34, 44, 34, 115, 90, 121, 111, 117, 122, 121, 111, 117, 72, 97, 105, 115, 105, 68, 97, 121, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 83, 105, 122, 121, 111, 117, 98, 101, 116, 117, 66, 97, 105, 98, 97, 105, 84, 97, 110, 105, 34, 58, 34, 48, 34, 44, 34, 115, 83, 105, 122, 121, 111, 117, 98, 101, 116, 117, 66, 97, 105, 98, 97, 105, 84, 97, 110, 105, 89, 111, 107, 117, 34, 58, 34, 48, 34, 44, 34, 115, 89, 111, 98, 105, 110, 101, 84, 97, 110, 105, 78, 117, 109, 98, 101, 114, 34, 58, 34, 49, 48, 49, 34, 44, 34, 115, 89, 111, 98, 105, 110, 101, 84, 97, 110, 105, 78, 117, 109, 98, 101, 114, 89, 111, 107, 117, 34, 58, 34, 49, 48, 49, 34, 44, 34, 115, 90, 121, 111, 117, 104, 111, 117, 83, 111, 117, 114, 99, 101, 34, 58, 34, 51, 34, 44, 34, 115, 90, 121, 111, 117, 104, 111, 117, 67, 111, 100, 101, 34, 58, 34, 69, 54, 57, 54, 50, 35, 48, 47, 84, 34, 44, 34, 115, 75, 111, 117, 98, 111, 80, 114, 105, 99, 101, 34, 58, 34, 48, 46, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 67, 114, 101, 97, 116, 101, 68, 97, 116, 101, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 115, 85, 112, 100, 97, 116, 101, 68, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 50, 49, 57, 48, 56, 49, 57, 34, 44, 34, 115, 85, 112, 100, 97, 116, 101, 78, 117, 109, 98, 101, 114, 34, 58, 34, 52, 51, 49, 50, 34, 125, 93, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockExchangeMasterRequest{},
			want1: &StockExchangeMasterResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 23, 9, 42, 19, 557000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 23, 9, 42, 19, 423000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeMasterData,
				},
				StockExchangeMasters: []StockExchangeMaster{
					{
						IssueCode:                   "6962",
						Exchange:                    ExchangeToushou,
						StockSystemType:             "1",
						UnderLimitPrice:             762,
						UpperLimitPrice:             1362,
						SymbolCategory:              " ",
						LimitPriceExchange:          ExchangeToushou,
						MarginType:                  MarginTypeMarginTrading,
						ListingDate:                 time.Time{},
						LimitPriceDate:              time.Time{},
						LimitPriceCategory:          " ",
						LimitPriceValue:             0,
						ConfirmLimitPrice:           true,
						Section:                     "1",
						PrevClosePrice:              1062,
						CalculateLimitPriceExchange: " ",
						Regulation1:                 " ",
						Regulation2:                 " ",
						SectionType:                 "01",
						DelistingDate:               time.Time{},
						TradingUnit:                 0,
						NextTradingUnit:             0,
						TickGroupType:               TickGroupTypeStock1,
						NextTickGroupType:           TickGroupTypeStock1,
						InformationSource:           "3",
						InformationCode:             "E6962#0/T",
						OfferPrice:                  0,
						CreateDateTime:              time.Time{},
						UpdateDateTime:              time.Date(2022, 3, 22, 19, 8, 19, 0, time.Local),
						UpdateNumber:                "4312",
					},
				},
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			arg1:  context.Background(),
			arg2:  nil,
			arg3:  StockExchangeMasterRequest{},
			want1: nil,
			want2: NilArgumentErr},
		{name: "リクエストでエラーが返されたらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get2: StatusNotOkErr},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockExchangeMasterRequest{},
			want1:     nil,
			want2:     StatusNotOkErr},
		{name: "レスポンスのUnmarshalに失敗したらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get1: []byte{}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      StockExchangeMasterRequest{},
			want1:     nil,
			want2:     UnmarshalFailedErr},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &client{clock: test.clock, requester: test.requester}
			got1, got2 := client.StockExchangeMaster(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_StockExchangeMaster_Execute(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	client := NewClient(EnvironmentProduction, ApiVersionLatest)
	got1, got2 := client.Login(context.Background(), LoginRequest{
		UserId:   userId,
		Password: password,
	})
	log.Printf("%+v, %+v\n", got1, got2)
	if got1.ResultCode != "0" {
		return
	}

	session, err := got1.Session()
	if err != nil {
		t.Errorf("session: %+v\n", err)
	}

	ctx := context.Background()
	ctx, cf := context.WithTimeout(ctx, 30*time.Second)
	defer cf()

	got3, got4 := client.StockExchangeMaster(ctx, session, StockExchangeMasterRequest{
		Columns: []StockExchangeMasterColumn{
			StockExchangeMasterColumnIssueCode,
			StockExchangeMasterColumnExchange,
			StockExchangeMasterColumnTickGroupType,
		},
	})
	log.Printf("%+v, %+v\n", got3, got4)
}
