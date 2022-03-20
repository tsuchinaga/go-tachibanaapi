package tachibana

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"
	"time"
)

func Test_MarketPriceRequest_request(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request MarketPriceRequest
		arg1    int64
		arg2    time.Time
		want1   marketPriceRequest
	}{
		{name: "銘柄が指定がnilなら空配列に変えてから処理する",
			request: MarketPriceRequest{
				IssueCodes: nil,
				Columns:    AllMarketPriceColumns,
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 7, 12, 58, 0, 0, time.Local),
			want1: marketPriceRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 7, 12, 58, 0, 0, time.Local)},
					MessageType:    MessageTypeMarketPrice,
					ResponseFormat: commonResponseFormat,
				},
				IssueCodes: "",
				Columns:    "xLISS,pDPP,tDPP:T,pDPG,pDYWP,pDYRP,pDOP,tDOP:T,pDHP,tDHP:T,pDLP,tDLP:T,pDV,pQAS,pQAP,pAV,pQBS,pQBP,pBV,xDVES,xDCFS,pDHF,pDLF,pDJ,pAAV,pABV,pQOV,pGAV10,pGAP10,pGAV9,pGAP9,pGAV8,pGAP8,pGAV7,pGAP7,pGAV6,pGAP6,pGAV5,pGAP5,pGAV4,pGAP4,pGAV3,pGAP3,pGAV2,pGAP2,pGAV1,pGAP1,pGBV1,pGBP1,pGBV2,pGBP2,pGBV3,pGBP3,pGBV4,pGBP4,pGBV5,pGBP5,pGBV6,pGBP6,pGBV7,pGBP7,pGBV8,pGBP8,pGBV9,pGBP9,pGBV10,pGBP10,pQUV,pVWAP,pPRP",
			}},
		{name: "カラムが指定されていれば、指定されたカラムを登録する",
			request: MarketPriceRequest{
				IssueCodes: []string{"1475", "1476"},
				Columns:    []MarketPriceColumn{MarketPriceColumnCurrentPrice, MarketPriceColumnCurrentPriceTime},
			},
			arg1: 1234,
			arg2: time.Date(2022, 3, 7, 12, 58, 0, 0, time.Local),
			want1: marketPriceRequest{
				commonRequest: commonRequest{
					No:             1234,
					SendDate:       RequestTime{Time: time.Date(2022, 3, 7, 12, 58, 0, 0, time.Local)},
					MessageType:    MessageTypeMarketPrice,
					ResponseFormat: commonResponseFormat,
				},
				IssueCodes: "1475,1476",
				Columns:    "pDPP,tDPP:T",
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

func Test_client_MarketPrice(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		clock     *testClock
		requester *testRequester
		arg1      context.Context
		arg2      *Session
		arg3      MarketPriceRequest
		want1     *MarketPriceResponse
		want2     error
	}{
		{name: "全項目含む正常レスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 55, 45, 49, 50, 58, 51, 50, 58, 50, 55, 46, 56, 56, 57, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 55, 45, 49, 50, 58, 51, 50, 58, 50, 55, 46, 56, 53, 51, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 77, 102, 100, 115, 71, 101, 116, 77, 97, 114, 107, 101, 116, 80, 114, 105, 99, 101, 34, 44, 34, 97, 67, 76, 77, 77, 102, 100, 115, 77, 97, 114, 107, 101, 116, 80, 114, 105, 99, 101, 34, 58, 91, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 53, 34, 44, 34, 120, 76, 73, 83, 83, 34, 58, 34, 34, 44, 34, 112, 68, 80, 80, 34, 58, 34, 49, 56, 50, 55, 34, 44, 34, 116, 68, 80, 80, 58, 84, 34, 58, 34, 49, 50, 58, 51, 50, 34, 44, 34, 112, 68, 80, 71, 34, 58, 34, 48, 48, 53, 55, 34, 44, 34, 112, 68, 89, 87, 80, 34, 58, 34, 45, 53, 52, 34, 44, 34, 112, 68, 89, 82, 80, 34, 58, 34, 45, 50, 46, 56, 55, 34, 44, 34, 112, 68, 79, 80, 34, 58, 34, 49, 56, 52, 52, 34, 44, 34, 116, 68, 79, 80, 58, 84, 34, 58, 34, 48, 57, 58, 48, 48, 34, 44, 34, 112, 68, 72, 80, 34, 58, 34, 49, 56, 52, 54, 34, 44, 34, 116, 68, 72, 80, 58, 84, 34, 58, 34, 48, 57, 58, 48, 53, 34, 44, 34, 112, 68, 76, 80, 34, 58, 34, 49, 56, 49, 52, 34, 44, 34, 116, 68, 76, 80, 58, 84, 34, 58, 34, 49, 48, 58, 52, 53, 34, 44, 34, 112, 68, 86, 34, 58, 34, 49, 55, 55, 52, 54, 53, 34, 44, 34, 112, 81, 65, 83, 34, 58, 34, 48, 49, 48, 49, 34, 44, 34, 112, 81, 65, 80, 34, 58, 34, 49, 56, 50, 55, 34, 44, 34, 112, 65, 86, 34, 58, 34, 57, 50, 56, 56, 54, 34, 44, 34, 112, 81, 66, 83, 34, 58, 34, 48, 49, 48, 49, 34, 44, 34, 112, 81, 66, 80, 34, 58, 34, 49, 56, 50, 53, 34, 44, 34, 112, 66, 86, 34, 58, 34, 49, 49, 50, 56, 57, 51, 34, 44, 34, 120, 68, 86, 69, 83, 34, 58, 34, 34, 44, 34, 120, 68, 67, 70, 83, 34, 58, 34, 34, 44, 34, 112, 68, 72, 70, 34, 58, 34, 48, 48, 48, 48, 34, 44, 34, 112, 68, 76, 70, 34, 58, 34, 48, 48, 48, 48, 34, 44, 34, 112, 68, 74, 34, 58, 34, 51, 50, 51, 57, 51, 56, 54, 51, 50, 34, 44, 34, 112, 65, 65, 86, 34, 58, 34, 34, 44, 34, 112, 65, 66, 86, 34, 58, 34, 34, 44, 34, 112, 81, 79, 86, 34, 58, 34, 51, 53, 54, 51, 56, 34, 44, 34, 112, 71, 65, 86, 49, 48, 34, 58, 34, 51, 34, 44, 34, 112, 71, 65, 80, 49, 48, 34, 58, 34, 49, 56, 51, 54, 34, 44, 34, 112, 71, 65, 86, 57, 34, 58, 34, 51, 34, 44, 34, 112, 71, 65, 80, 57, 34, 58, 34, 49, 56, 51, 53, 34, 44, 34, 112, 71, 65, 86, 56, 34, 58, 34, 49, 55, 34, 44, 34, 112, 71, 65, 80, 56, 34, 58, 34, 49, 56, 51, 52, 34, 44, 34, 112, 71, 65, 86, 55, 34, 58, 34, 54, 34, 44, 34, 112, 71, 65, 80, 55, 34, 58, 34, 49, 56, 51, 51, 34, 44, 34, 112, 71, 65, 86, 54, 34, 58, 34, 53, 34, 44, 34, 112, 71, 65, 80, 54, 34, 58, 34, 49, 56, 51, 50, 34, 44, 34, 112, 71, 65, 86, 53, 34, 58, 34, 50, 48, 48, 48, 52, 34, 44, 34, 112, 71, 65, 80, 53, 34, 58, 34, 49, 56, 51, 49, 34, 44, 34, 112, 71, 65, 86, 52, 34, 58, 34, 53, 48, 48, 48, 55, 34, 44, 34, 112, 71, 65, 80, 52, 34, 58, 34, 49, 56, 51, 48, 34, 44, 34, 112, 71, 65, 86, 51, 34, 58, 34, 57, 49, 56, 57, 50, 34, 44, 34, 112, 71, 65, 80, 51, 34, 58, 34, 49, 56, 50, 57, 34, 44, 34, 112, 71, 65, 86, 50, 34, 58, 34, 49, 51, 49, 56, 57, 50, 34, 44, 34, 112, 71, 65, 80, 50, 34, 58, 34, 49, 56, 50, 56, 34, 44, 34, 112, 71, 65, 86, 49, 34, 58, 34, 57, 50, 56, 56, 54, 34, 44, 34, 112, 71, 65, 80, 49, 34, 58, 34, 49, 56, 50, 55, 34, 44, 34, 112, 71, 66, 86, 49, 34, 58, 34, 49, 49, 50, 56, 57, 51, 34, 44, 34, 112, 71, 66, 80, 49, 34, 58, 34, 49, 56, 50, 53, 34, 44, 34, 112, 71, 66, 86, 50, 34, 58, 34, 56, 57, 56, 51, 49, 34, 44, 34, 112, 71, 66, 80, 50, 34, 58, 34, 49, 56, 50, 52, 34, 44, 34, 112, 71, 66, 86, 51, 34, 58, 34, 53, 57, 57, 50, 54, 34, 44, 34, 112, 71, 66, 80, 51, 34, 58, 34, 49, 56, 50, 51, 34, 44, 34, 112, 71, 66, 86, 52, 34, 58, 34, 50, 57, 56, 50, 51, 34, 44, 34, 112, 71, 66, 80, 52, 34, 58, 34, 49, 56, 50, 50, 34, 44, 34, 112, 71, 66, 86, 53, 34, 58, 34, 50, 57, 57, 48, 54, 34, 44, 34, 112, 71, 66, 80, 53, 34, 58, 34, 49, 56, 50, 49, 34, 44, 34, 112, 71, 66, 86, 54, 34, 58, 34, 49, 48, 51, 51, 52, 34, 44, 34, 112, 71, 66, 80, 54, 34, 58, 34, 49, 56, 50, 48, 34, 44, 34, 112, 71, 66, 86, 55, 34, 58, 34, 49, 50, 34, 44, 34, 112, 71, 66, 80, 55, 34, 58, 34, 49, 56, 49, 57, 34, 44, 34, 112, 71, 66, 86, 56, 34, 58, 34, 49, 48, 52, 34, 44, 34, 112, 71, 66, 80, 56, 34, 58, 34, 49, 56, 49, 56, 34, 44, 34, 112, 71, 66, 86, 57, 34, 58, 34, 51, 53, 34, 44, 34, 112, 71, 66, 80, 57, 34, 58, 34, 49, 56, 49, 55, 34, 44, 34, 112, 71, 66, 86, 49, 48, 34, 58, 34, 51, 54, 50, 34, 44, 34, 112, 71, 66, 80, 49, 48, 34, 58, 34, 49, 56, 49, 54, 34, 44, 34, 112, 81, 85, 86, 34, 58, 34, 50, 56, 57, 53, 53, 34, 44, 34, 112, 86, 87, 65, 80, 34, 58, 34, 49, 56, 50, 53, 46, 51, 54, 54, 51, 34, 44, 34, 112, 80, 82, 80, 34, 58, 34, 49, 56, 56, 49, 34, 125, 44, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 54, 34, 44, 34, 120, 76, 73, 83, 83, 34, 58, 34, 34, 44, 34, 112, 68, 80, 80, 34, 58, 34, 49, 57, 53, 48, 34, 44, 34, 116, 68, 80, 80, 58, 84, 34, 58, 34, 49, 50, 58, 51, 48, 34, 44, 34, 112, 68, 80, 71, 34, 58, 34, 48, 48, 53, 55, 34, 44, 34, 112, 68, 89, 87, 80, 34, 58, 34, 45, 52, 34, 44, 34, 112, 68, 89, 82, 80, 34, 58, 34, 45, 48, 46, 50, 48, 34, 44, 34, 112, 68, 79, 80, 34, 58, 34, 49, 57, 52, 51, 34, 44, 34, 116, 68, 79, 80, 58, 84, 34, 58, 34, 48, 57, 58, 48, 48, 34, 44, 34, 112, 68, 72, 80, 34, 58, 34, 49, 57, 53, 48, 34, 44, 34, 116, 68, 72, 80, 58, 84, 34, 58, 34, 49, 50, 58, 51, 48, 34, 44, 34, 112, 68, 76, 80, 34, 58, 34, 49, 57, 50, 57, 34, 44, 34, 116, 68, 76, 80, 58, 84, 34, 58, 34, 49, 48, 58, 48, 53, 34, 44, 34, 112, 68, 86, 34, 58, 34, 50, 56, 48, 55, 53, 34, 44, 34, 112, 81, 65, 83, 34, 58, 34, 48, 49, 48, 49, 34, 44, 34, 112, 81, 65, 80, 34, 58, 34, 49, 57, 53, 48, 34, 44, 34, 112, 65, 86, 34, 58, 34, 52, 52, 52, 51, 55, 34, 44, 34, 112, 81, 66, 83, 34, 58, 34, 48, 49, 48, 49, 34, 44, 34, 112, 81, 66, 80, 34, 58, 34, 49, 57, 52, 55, 34, 44, 34, 112, 66, 86, 34, 58, 34, 52, 49, 57, 52, 55, 34, 44, 34, 120, 68, 86, 69, 83, 34, 58, 34, 34, 44, 34, 120, 68, 67, 70, 83, 34, 58, 34, 34, 44, 34, 112, 68, 72, 70, 34, 58, 34, 48, 48, 48, 48, 34, 44, 34, 112, 68, 76, 70, 34, 58, 34, 48, 48, 48, 48, 34, 44, 34, 112, 68, 74, 34, 58, 34, 53, 52, 51, 56, 50, 53, 54, 50, 34, 44, 34, 112, 65, 65, 86, 34, 58, 34, 34, 44, 34, 112, 65, 66, 86, 34, 58, 34, 34, 44, 34, 112, 81, 79, 86, 34, 58, 34, 49, 48, 48, 50, 48, 34, 44, 34, 112, 71, 65, 86, 49, 48, 34, 58, 34, 50, 51, 34, 44, 34, 112, 71, 65, 80, 49, 48, 34, 58, 34, 49, 57, 53, 57, 34, 44, 34, 112, 71, 65, 86, 57, 34, 58, 34, 51, 53, 48, 34, 44, 34, 112, 71, 65, 80, 57, 34, 58, 34, 49, 57, 53, 56, 34, 44, 34, 112, 71, 65, 86, 56, 34, 58, 34, 49, 34, 44, 34, 112, 71, 65, 80, 56, 34, 58, 34, 49, 57, 53, 55, 34, 44, 34, 112, 71, 65, 86, 55, 34, 58, 34, 49, 54, 51, 53, 48, 34, 44, 34, 112, 71, 65, 80, 55, 34, 58, 34, 49, 57, 53, 54, 34, 44, 34, 112, 71, 65, 86, 54, 34, 58, 34, 49, 54, 48, 50, 51, 34, 44, 34, 112, 71, 65, 80, 54, 34, 58, 34, 49, 57, 53, 53, 34, 44, 34, 112, 71, 65, 86, 53, 34, 58, 34, 49, 54, 50, 54, 55, 34, 44, 34, 112, 71, 65, 80, 53, 34, 58, 34, 49, 57, 53, 52, 34, 44, 34, 112, 71, 65, 86, 52, 34, 58, 34, 52, 54, 54, 57, 53, 34, 44, 34, 112, 71, 65, 80, 52, 34, 58, 34, 49, 57, 53, 51, 34, 44, 34, 112, 71, 65, 86, 51, 34, 58, 34, 54, 53, 56, 51, 55, 34, 44, 34, 112, 71, 65, 80, 51, 34, 58, 34, 49, 57, 53, 50, 34, 44, 34, 112, 71, 65, 86, 50, 34, 58, 34, 52, 51, 53, 57, 48, 34, 44, 34, 112, 71, 65, 80, 50, 34, 58, 34, 49, 57, 53, 49, 34, 44, 34, 112, 71, 65, 86, 49, 34, 58, 34, 52, 52, 52, 51, 55, 34, 44, 34, 112, 71, 65, 80, 49, 34, 58, 34, 49, 57, 53, 48, 34, 44, 34, 112, 71, 66, 86, 49, 34, 58, 34, 52, 49, 57, 52, 55, 34, 44, 34, 112, 71, 66, 80, 49, 34, 58, 34, 49, 57, 52, 55, 34, 44, 34, 112, 71, 66, 86, 50, 34, 58, 34, 52, 50, 50, 57, 55, 34, 44, 34, 112, 71, 66, 80, 50, 34, 58, 34, 49, 57, 52, 54, 34, 44, 34, 112, 71, 66, 86, 51, 34, 58, 34, 52, 55, 56, 50, 50, 34, 44, 34, 112, 71, 66, 80, 51, 34, 58, 34, 49, 57, 52, 53, 34, 44, 34, 112, 71, 66, 86, 52, 34, 58, 34, 52, 55, 52, 55, 53, 34, 44, 34, 112, 71, 66, 80, 52, 34, 58, 34, 49, 57, 52, 52, 34, 44, 34, 112, 71, 66, 86, 53, 34, 58, 34, 49, 56, 48, 50, 54, 34, 44, 34, 112, 71, 66, 80, 53, 34, 58, 34, 49, 57, 52, 51, 34, 44, 34, 112, 71, 66, 86, 54, 34, 58, 34, 49, 56, 51, 53, 48, 34, 44, 34, 112, 71, 66, 80, 54, 34, 58, 34, 49, 57, 52, 50, 34, 44, 34, 112, 71, 66, 86, 55, 34, 58, 34, 49, 54, 51, 50, 53, 34, 44, 34, 112, 71, 66, 80, 55, 34, 58, 34, 49, 57, 52, 49, 34, 44, 34, 112, 71, 66, 86, 56, 34, 58, 34, 49, 54, 51, 53, 51, 34, 44, 34, 112, 71, 66, 80, 56, 34, 58, 34, 49, 57, 52, 48, 34, 44, 34, 112, 71, 66, 86, 57, 34, 58, 34, 50, 48, 34, 44, 34, 112, 71, 66, 80, 57, 34, 58, 34, 49, 57, 51, 57, 34, 44, 34, 112, 71, 66, 86, 49, 48, 34, 58, 34, 51, 54, 48, 34, 44, 34, 112, 71, 66, 80, 49, 48, 34, 58, 34, 49, 57, 51, 56, 34, 44, 34, 112, 81, 85, 86, 34, 58, 34, 50, 48, 56, 51, 56, 34, 44, 34, 112, 86, 87, 65, 80, 34, 58, 34, 49, 57, 51, 55, 46, 48, 52, 53, 56, 34, 44, 34, 112, 80, 82, 80, 34, 58, 34, 49, 57, 53, 52, 34, 125, 93, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      MarketPriceRequest{},
			want1: &MarketPriceResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 7, 12, 32, 27, 889000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 7, 12, 32, 27, 853000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeMarketPrice,
				},
				MarketPrices: []MarketPrice{
					{
						IssueCode:         "1475",
						Section:           "",
						CurrentPrice:      1827,
						CurrentPriceTime:  time.Date(0, 1, 1, 12, 32, 0, 0, time.Local),
						ChangePriceType:   ChangePriceTypeRise,
						PrevDayRatio:      -54,
						PrevDayPercent:    -2.87,
						OpenPrice:         1844,
						OpenPriceTime:     time.Date(0, 1, 1, 9, 0, 0, 0, time.Local),
						HighPrice:         1846,
						HighPriceTime:     time.Date(0, 1, 1, 9, 5, 0, 0, time.Local),
						LowPrice:          1814,
						LowPriceTime:      time.Date(0, 1, 1, 10, 45, 0, 0, time.Local),
						Volume:            177_465,
						AskSign:           IndicationPriceTypeGeneral,
						AskPrice:          1827,
						AskQuantity:       92_886,
						BidSign:           IndicationPriceTypeGeneral,
						BidPrice:          1825,
						BidQuantity:       112_893,
						ExRightType:       "",
						DiscontinuityType: "",
						StopHigh:          CurrentPriceTypeNoChange,
						StopLow:           CurrentPriceTypeNoChange,
						TradingAmount:     323_938_632,
						AskQuantityMarket: 0,
						BidQuantityMarket: 0,
						AskQuantityOver:   35_638,
						AskQuantity10:     3,
						AskPrice10:        1836,
						AskQuantity9:      3,
						AskPrice9:         1835,
						AskQuantity8:      17,
						AskPrice8:         1834,
						AskQuantity7:      6,
						AskPrice7:         1833,
						AskQuantity6:      5,
						AskPrice6:         1832,
						AskQuantity5:      20004,
						AskPrice5:         1831,
						AskQuantity4:      50007,
						AskPrice4:         1830,
						AskQuantity3:      91892,
						AskPrice3:         1829,
						AskQuantity2:      131892,
						AskPrice2:         1828,
						AskQuantity1:      92886,
						AskPrice1:         1827,
						BidQuantity1:      112893,
						BidPrice1:         1825,
						BidQuantity2:      89831,
						BidPrice2:         1824,
						BidQuantity3:      59926,
						BidPrice3:         1823,
						BidQuantity4:      29823,
						BidPrice4:         1822,
						BidQuantity5:      29906,
						BidPrice5:         1821,
						BidQuantity6:      10334,
						BidPrice6:         1820,
						BidQuantity7:      12,
						BidPrice7:         1819,
						BidQuantity8:      104,
						BidPrice8:         1818,
						BidQuantity9:      35,
						BidPrice9:         1817,
						BidQuantity10:     362,
						BidPrice10:        1816,
						BidQuantityUnder:  28955,
						VWAP:              1825.3663,
						PRP:               1881,
					}, {
						IssueCode:         "1476",
						Section:           "",
						CurrentPrice:      1950,
						CurrentPriceTime:  time.Date(0, 1, 1, 12, 30, 0, 0, time.Local),
						ChangePriceType:   ChangePriceTypeRise,
						PrevDayRatio:      -4,
						PrevDayPercent:    -0.2,
						OpenPrice:         1943,
						OpenPriceTime:     time.Date(0, 1, 1, 9, 0, 0, 0, time.Local),
						HighPrice:         1950,
						HighPriceTime:     time.Date(0, 1, 1, 12, 30, 0, 0, time.Local),
						LowPrice:          1929,
						LowPriceTime:      time.Date(0, 1, 1, 10, 5, 0, 0, time.Local),
						Volume:            28075,
						AskSign:           IndicationPriceTypeGeneral,
						AskPrice:          1950,
						AskQuantity:       44437,
						BidSign:           IndicationPriceTypeGeneral,
						BidPrice:          1947,
						BidQuantity:       41947,
						ExRightType:       "",
						DiscontinuityType: "",
						StopHigh:          CurrentPriceTypeNoChange,
						StopLow:           CurrentPriceTypeNoChange,
						TradingAmount:     54382562,
						AskQuantityMarket: 0,
						BidQuantityMarket: 0,
						AskQuantityOver:   10020,
						AskQuantity10:     23,
						AskPrice10:        1959,
						AskQuantity9:      350,
						AskPrice9:         1958,
						AskQuantity8:      1,
						AskPrice8:         1957,
						AskQuantity7:      16350,
						AskPrice7:         1956,
						AskQuantity6:      16023,
						AskPrice6:         1955,
						AskQuantity5:      16267,
						AskPrice5:         1954,
						AskQuantity4:      46695,
						AskPrice4:         1953,
						AskQuantity3:      65837,
						AskPrice3:         1952,
						AskQuantity2:      43590,
						AskPrice2:         1951,
						AskQuantity1:      44437,
						AskPrice1:         1950,
						BidQuantity1:      41947,
						BidPrice1:         1947,
						BidQuantity2:      42297,
						BidPrice2:         1946,
						BidQuantity3:      47822,
						BidPrice3:         1945,
						BidQuantity4:      47475,
						BidPrice4:         1944,
						BidQuantity5:      18026,
						BidPrice5:         1943,
						BidQuantity6:      18350,
						BidPrice6:         1942,
						BidQuantity7:      16325,
						BidPrice7:         1941,
						BidQuantity8:      16353,
						BidPrice8:         1940,
						BidQuantity9:      20,
						BidPrice9:         1939,
						BidQuantity10:     360,
						BidPrice10:        1938,
						BidQuantityUnder:  20838,
						VWAP:              1937.0458,
						PRP:               1954,
					},
				},
			},
			want2: nil},
		{name: "銘柄指定なしのレスポンスパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 55, 45, 49, 50, 58, 51, 50, 58, 49, 52, 46, 56, 51, 50, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 55, 45, 49, 50, 58, 51, 50, 58, 49, 52, 46, 55, 55, 55, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 77, 102, 100, 115, 71, 101, 116, 77, 97, 114, 107, 101, 116, 80, 114, 105, 99, 101, 34, 44, 34, 97, 67, 76, 77, 77, 102, 100, 115, 77, 97, 114, 107, 101, 116, 80, 114, 105, 99, 101, 34, 58, 91, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 34, 44, 34, 112, 68, 80, 80, 34, 58, 34, 34, 44, 34, 116, 68, 80, 80, 58, 84, 34, 58, 34, 34, 125, 93, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      MarketPriceRequest{},
			want1: &MarketPriceResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 7, 12, 32, 14, 832000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 7, 12, 32, 14, 777000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeMarketPrice,
				},
				MarketPrices: []MarketPrice{},
			},
			want2: nil},
		{name: "存在しない銘柄の指定のレスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 55, 45, 49, 50, 58, 51, 50, 58, 48, 51, 46, 49, 54, 52, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 55, 45, 49, 50, 58, 51, 50, 58, 48, 51, 46, 49, 50, 51, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 77, 102, 100, 115, 71, 101, 116, 77, 97, 114, 107, 101, 116, 80, 114, 105, 99, 101, 34, 44, 34, 97, 67, 76, 77, 77, 102, 100, 115, 77, 97, 114, 107, 101, 116, 80, 114, 105, 99, 101, 34, 58, 91, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 42, 34, 44, 34, 112, 68, 80, 80, 34, 58, 34, 34, 44, 34, 116, 68, 80, 80, 58, 84, 34, 58, 34, 34, 125, 93, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      MarketPriceRequest{},
			want1: &MarketPriceResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 7, 12, 32, 3, 164000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 7, 12, 32, 3, 123000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeMarketPrice,
				},
				MarketPrices: []MarketPrice{{IssueCode: "*"}},
			},
			want2: nil},
		{name: "項目を制限して取得した場合のレスポンスをパースして返せる",
			clock:     &testClock{Now1: time.Date(2022, 3, 6, 11, 11, 0, 0, time.Local)},
			requester: &testRequester{get1: []byte{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 55, 45, 49, 50, 58, 51, 49, 58, 51, 56, 46, 48, 51, 48, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 48, 55, 45, 49, 50, 58, 51, 49, 58, 51, 56, 46, 48, 49, 52, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 77, 102, 100, 115, 71, 101, 116, 77, 97, 114, 107, 101, 116, 80, 114, 105, 99, 101, 34, 44, 34, 97, 67, 76, 77, 77, 102, 100, 115, 77, 97, 114, 107, 101, 116, 80, 114, 105, 99, 101, 34, 58, 91, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 53, 34, 44, 34, 112, 68, 80, 80, 34, 58, 34, 49, 56, 50, 54, 34, 44, 34, 116, 68, 80, 80, 58, 84, 34, 58, 34, 49, 50, 58, 51, 48, 34, 125, 44, 123, 34, 115, 73, 115, 115, 117, 101, 67, 111, 100, 101, 34, 58, 34, 49, 52, 55, 54, 34, 44, 34, 112, 68, 80, 80, 34, 58, 34, 49, 57, 53, 48, 34, 44, 34, 116, 68, 80, 80, 58, 84, 34, 58, 34, 49, 50, 58, 51, 48, 34, 125, 93, 125, 10}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      MarketPriceRequest{},
			want1: &MarketPriceResponse{
				CommonResponse: CommonResponse{
					No:           2,
					SendDate:     time.Date(2022, 3, 7, 12, 31, 38, 30000000, time.Local),
					ReceiveDate:  time.Date(2022, 3, 7, 12, 31, 38, 14000000, time.Local),
					ErrorNo:      ErrorNoProblem,
					ErrorMessage: "",
					MessageType:  MessageTypeMarketPrice,
				},
				MarketPrices: []MarketPrice{
					{IssueCode: "1475", CurrentPrice: 1826, CurrentPriceTime: time.Date(0, 1, 1, 12, 30, 0, 0, time.Local)},
					{IssueCode: "1476", CurrentPrice: 1950, CurrentPriceTime: time.Date(0, 1, 1, 12, 30, 0, 0, time.Local)},
				},
			},
			want2: nil},
		{name: "sessionがnilならエラー",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			arg1:  context.Background(),
			arg2:  nil,
			arg3:  MarketPriceRequest{},
			want1: nil,
			want2: NilArgumentErr},
		{name: "リクエストでエラーが返されたらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get2: StatusNotOkErr},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      MarketPriceRequest{},
			want1:     nil,
			want2:     StatusNotOkErr},
		{name: "レスポンスのUnmarshalに失敗したらエラーを返す",
			clock:     &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			requester: &testRequester{get1: []byte{}},
			arg1:      context.Background(),
			arg2:      &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:      MarketPriceRequest{},
			want1:     nil,
			want2:     UnmarshalFailedErr},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := &client{clock: test.clock, requester: test.requester}
			got1, got2 := client.MarketPrice(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), test.want1, test.want2, got1, got2)
			}
		})
	}
}

func Test_client_MarketPrice_Execute(t *testing.T) {
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

	got3, got4 := client.MarketPrice(context.Background(), session, MarketPriceRequest{
		IssueCodes: []string{"1475", "1476"},
		Columns:    []MarketPriceColumn{},
	})
	log.Printf("%+v, %+v\n", got3, got4)
}
