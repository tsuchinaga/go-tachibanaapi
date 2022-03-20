package tachibana

import (
	"context"
	"errors"
	"log"
	"reflect"
	"testing"
	"time"
)

func Test_client_BusinessDay(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		clock  iClock
		stream func(stream1 chan<- []byte, stream2 chan<- error)
		arg1   context.Context
		arg2   *Session
		arg3   BusinessDayRequest
		want1  []*BusinessDayResponse
		want2  error
	}{
		{name: "sessionがnilならエラー",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			stream: func(stream1 chan<- []byte, stream2 chan<- error) {
				defer close(stream1)
				defer close(stream2)
			},
			arg1:  context.Background(),
			arg2:  nil,
			arg3:  BusinessDayRequest{},
			want1: nil,
			want2: NilArgumentErr},
		{name: "エラーが返されたらエラーを返す",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			stream: func(stream1 chan<- []byte, stream2 chan<- error) {
				defer close(stream1)
				defer close(stream2)
				stream2 <- StatusNotOkErr
			},
			arg1:  context.Background(),
			arg2:  &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:  BusinessDayRequest{},
			want1: nil,
			want2: StatusNotOkErr},
		{name: "パースできなければエラーを返す",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			stream: func(stream1 chan<- []byte, stream2 chan<- error) {
				defer close(stream1)
				defer close(stream2)
				stream1 <- nil
			},
			arg1:  context.Background(),
			arg2:  &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:  BusinessDayRequest{},
			want1: nil,
			want2: UnmarshalFailedErr},
		{name: "初期ダウンロード終了通知がきたらそこで終わる",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			stream: func(stream1 chan<- []byte, stream2 chan<- error) {
				defer close(stream1)
				defer close(stream2)
				data := [][]byte{
					{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 50, 49, 45, 48, 53, 58, 51, 51, 58, 50, 55, 46, 50, 57, 51, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 76, 77, 68, 97, 116, 101, 90, 121, 111, 117, 104, 111, 117, 34, 44, 34, 115, 68, 97, 121, 75, 101, 121, 34, 58, 34, 48, 48, 49, 34, 44, 34, 115, 77, 97, 101, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 49, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 56, 34, 44, 34, 115, 77, 97, 101, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 50, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 55, 34, 44, 34, 115, 77, 97, 101, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 51, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 54, 34, 44, 34, 115, 84, 104, 101, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 50, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 49, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 51, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 50, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 52, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 51, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 53, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 52, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 56, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 53, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 57, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 54, 34, 58, 34, 50, 48, 50, 50, 48, 51, 51, 48, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 55, 34, 58, 34, 50, 48, 50, 50, 48, 51, 51, 49, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 56, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 49, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 57, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 52, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 49, 48, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 53, 34, 44, 34, 115, 75, 97, 98, 117, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 52, 34, 44, 34, 115, 75, 97, 98, 117, 75, 97, 114, 105, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 53, 34, 44, 34, 115, 66, 111, 110, 100, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 52, 34, 125},
					{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 50, 49, 45, 48, 53, 58, 51, 51, 58, 50, 55, 46, 50, 57, 51, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 76, 77, 68, 97, 116, 101, 90, 121, 111, 117, 104, 111, 117, 34, 44, 34, 115, 68, 97, 121, 75, 101, 121, 34, 58, 34, 48, 48, 50, 34, 44, 34, 115, 77, 97, 101, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 49, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 50, 34, 44, 34, 115, 77, 97, 101, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 50, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 56, 34, 44, 34, 115, 77, 97, 101, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 51, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 55, 34, 44, 34, 115, 84, 104, 101, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 51, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 49, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 52, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 50, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 53, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 51, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 56, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 52, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 57, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 53, 34, 58, 34, 50, 48, 50, 50, 48, 51, 51, 48, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 54, 34, 58, 34, 50, 48, 50, 50, 48, 51, 51, 49, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 55, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 49, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 56, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 52, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 57, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 53, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 49, 48, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 54, 34, 44, 34, 115, 75, 97, 98, 117, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 53, 34, 44, 34, 115, 75, 97, 98, 117, 75, 97, 114, 105, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 56, 34, 44, 34, 115, 66, 111, 110, 100, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 53, 34, 125},
					{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 50, 49, 45, 48, 53, 58, 51, 51, 58, 52, 53, 46, 52, 53, 57, 34, 44, 34, 112, 95, 110, 111, 34, 58, 34, 50, 34, 44, 34, 112, 95, 114, 118, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 50, 49, 45, 48, 53, 58, 51, 51, 58, 50, 55, 46, 50, 55, 54, 34, 44, 34, 112, 95, 101, 114, 114, 110, 111, 34, 58, 34, 48, 34, 44, 34, 112, 95, 101, 114, 114, 34, 58, 34, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 67, 76, 77, 69, 118, 101, 110, 116, 68, 111, 119, 110, 108, 111, 97, 100, 67, 111, 109, 112, 108, 101, 116, 101, 34, 125},
				}
				for _, d := range data {
					stream1 <- d
				}
			},
			arg1: context.Background(),
			arg2: &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3: BusinessDayRequest{},
			want1: []*BusinessDayResponse{
				{
					CommonResponse: CommonResponse{
						SendDate:    time.Date(2022, 3, 21, 5, 33, 27, 293000000, time.Local),
						FeatureType: "LMDateZyouhou",
					},
					DayKey:                 DayKeyToday,
					PrevDay1:               time.Date(2022, 3, 18, 0, 0, 0, 0, time.Local),
					PrevDay2:               time.Date(2022, 3, 17, 0, 0, 0, 0, time.Local),
					PrevDay3:               time.Date(2022, 3, 16, 0, 0, 0, 0, time.Local),
					Today:                  time.Date(2022, 3, 22, 0, 0, 0, 0, time.Local),
					NextDay1:               time.Date(2022, 3, 23, 0, 0, 0, 0, time.Local),
					NextDay2:               time.Date(2022, 3, 24, 0, 0, 0, 0, time.Local),
					NextDay3:               time.Date(2022, 3, 25, 0, 0, 0, 0, time.Local),
					NextDay4:               time.Date(2022, 3, 28, 0, 0, 0, 0, time.Local),
					NextDay5:               time.Date(2022, 3, 29, 0, 0, 0, 0, time.Local),
					NextDay6:               time.Date(2022, 3, 30, 0, 0, 0, 0, time.Local),
					NextDay7:               time.Date(2022, 3, 31, 0, 0, 0, 0, time.Local),
					NextDay8:               time.Date(2022, 4, 1, 0, 0, 0, 0, time.Local),
					NextDay9:               time.Date(2022, 4, 4, 0, 0, 0, 0, time.Local),
					NextDay10:              time.Date(2022, 4, 5, 0, 0, 0, 0, time.Local),
					DeliveryDay:            time.Date(2022, 3, 24, 0, 0, 0, 0, time.Local),
					ProvisionalDeliveryDay: time.Date(2022, 3, 25, 0, 0, 0, 0, time.Local),
					BondDeliveryDay:        time.Date(2022, 3, 24, 0, 0, 0, 0, time.Local),
				},
				{
					CommonResponse: CommonResponse{
						SendDate:    time.Date(2022, 3, 21, 5, 33, 27, 293000000, time.Local),
						FeatureType: "LMDateZyouhou",
					},
					DayKey:                 DayKeyNextDay,
					PrevDay2:               time.Date(2022, 3, 18, 0, 0, 0, 0, time.Local),
					PrevDay3:               time.Date(2022, 3, 17, 0, 0, 0, 0, time.Local),
					PrevDay1:               time.Date(2022, 3, 22, 0, 0, 0, 0, time.Local),
					Today:                  time.Date(2022, 3, 23, 0, 0, 0, 0, time.Local),
					NextDay1:               time.Date(2022, 3, 24, 0, 0, 0, 0, time.Local),
					NextDay2:               time.Date(2022, 3, 25, 0, 0, 0, 0, time.Local),
					NextDay3:               time.Date(2022, 3, 28, 0, 0, 0, 0, time.Local),
					NextDay4:               time.Date(2022, 3, 29, 0, 0, 0, 0, time.Local),
					NextDay5:               time.Date(2022, 3, 30, 0, 0, 0, 0, time.Local),
					NextDay6:               time.Date(2022, 3, 31, 0, 0, 0, 0, time.Local),
					NextDay7:               time.Date(2022, 4, 1, 0, 0, 0, 0, time.Local),
					NextDay8:               time.Date(2022, 4, 4, 0, 0, 0, 0, time.Local),
					NextDay9:               time.Date(2022, 4, 5, 0, 0, 0, 0, time.Local),
					NextDay10:              time.Date(2022, 4, 6, 0, 0, 0, 0, time.Local),
					DeliveryDay:            time.Date(2022, 3, 25, 0, 0, 0, 0, time.Local),
					ProvisionalDeliveryDay: time.Date(2022, 3, 28, 0, 0, 0, 0, time.Local),
					BondDeliveryDay:        time.Date(2022, 3, 25, 0, 0, 0, 0, time.Local),
				},
			},
			want2: nil},
		{name: "初期ダウンロード終了通知がこなくても、chanがcloseされたら終わる",
			clock: &testClock{Now1: time.Date(2022, 2, 24, 21, 2, 23, 365000000, time.Local)},
			stream: func(stream1 chan<- []byte, stream2 chan<- error) {
				defer close(stream1)
				defer close(stream2)
				data := [][]byte{
					{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 50, 49, 45, 48, 53, 58, 51, 51, 58, 50, 55, 46, 50, 57, 51, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 76, 77, 68, 97, 116, 101, 90, 121, 111, 117, 104, 111, 117, 34, 44, 34, 115, 68, 97, 121, 75, 101, 121, 34, 58, 34, 48, 48, 49, 34, 44, 34, 115, 77, 97, 101, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 49, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 56, 34, 44, 34, 115, 77, 97, 101, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 50, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 55, 34, 44, 34, 115, 77, 97, 101, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 51, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 54, 34, 44, 34, 115, 84, 104, 101, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 50, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 49, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 51, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 50, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 52, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 51, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 53, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 52, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 56, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 53, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 57, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 54, 34, 58, 34, 50, 48, 50, 50, 48, 51, 51, 48, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 55, 34, 58, 34, 50, 48, 50, 50, 48, 51, 51, 49, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 56, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 49, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 57, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 52, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 49, 48, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 53, 34, 44, 34, 115, 75, 97, 98, 117, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 52, 34, 44, 34, 115, 75, 97, 98, 117, 75, 97, 114, 105, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 53, 34, 44, 34, 115, 66, 111, 110, 100, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 52, 34, 125},
					{123, 34, 112, 95, 115, 100, 95, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 46, 48, 51, 46, 50, 49, 45, 48, 53, 58, 51, 51, 58, 50, 55, 46, 50, 57, 51, 34, 44, 34, 115, 67, 76, 77, 73, 68, 34, 58, 34, 76, 77, 68, 97, 116, 101, 90, 121, 111, 117, 104, 111, 117, 34, 44, 34, 115, 68, 97, 121, 75, 101, 121, 34, 58, 34, 48, 48, 50, 34, 44, 34, 115, 77, 97, 101, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 49, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 50, 34, 44, 34, 115, 77, 97, 101, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 50, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 56, 34, 44, 34, 115, 77, 97, 101, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 51, 34, 58, 34, 50, 48, 50, 50, 48, 51, 49, 55, 34, 44, 34, 115, 84, 104, 101, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 51, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 49, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 52, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 50, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 53, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 51, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 56, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 52, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 57, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 53, 34, 58, 34, 50, 48, 50, 50, 48, 51, 51, 48, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 54, 34, 58, 34, 50, 48, 50, 50, 48, 51, 51, 49, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 55, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 49, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 56, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 52, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 57, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 53, 34, 44, 34, 115, 89, 111, 107, 117, 69, 105, 103, 121, 111, 117, 68, 97, 121, 95, 49, 48, 34, 58, 34, 50, 48, 50, 50, 48, 52, 48, 54, 34, 44, 34, 115, 75, 97, 98, 117, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 53, 34, 44, 34, 115, 75, 97, 98, 117, 75, 97, 114, 105, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 56, 34, 44, 34, 115, 66, 111, 110, 100, 85, 107, 101, 119, 97, 116, 97, 115, 105, 68, 97, 121, 34, 58, 34, 50, 48, 50, 50, 48, 51, 50, 53, 34, 125},
				}
				for _, d := range data {
					stream1 <- d
				}
			},
			arg1: context.Background(),
			arg2: &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3: BusinessDayRequest{},
			want1: []*BusinessDayResponse{
				{
					CommonResponse: CommonResponse{
						SendDate:    time.Date(2022, 3, 21, 5, 33, 27, 293000000, time.Local),
						FeatureType: "LMDateZyouhou",
					},
					DayKey:                 DayKeyToday,
					PrevDay1:               time.Date(2022, 3, 18, 0, 0, 0, 0, time.Local),
					PrevDay2:               time.Date(2022, 3, 17, 0, 0, 0, 0, time.Local),
					PrevDay3:               time.Date(2022, 3, 16, 0, 0, 0, 0, time.Local),
					Today:                  time.Date(2022, 3, 22, 0, 0, 0, 0, time.Local),
					NextDay1:               time.Date(2022, 3, 23, 0, 0, 0, 0, time.Local),
					NextDay2:               time.Date(2022, 3, 24, 0, 0, 0, 0, time.Local),
					NextDay3:               time.Date(2022, 3, 25, 0, 0, 0, 0, time.Local),
					NextDay4:               time.Date(2022, 3, 28, 0, 0, 0, 0, time.Local),
					NextDay5:               time.Date(2022, 3, 29, 0, 0, 0, 0, time.Local),
					NextDay6:               time.Date(2022, 3, 30, 0, 0, 0, 0, time.Local),
					NextDay7:               time.Date(2022, 3, 31, 0, 0, 0, 0, time.Local),
					NextDay8:               time.Date(2022, 4, 1, 0, 0, 0, 0, time.Local),
					NextDay9:               time.Date(2022, 4, 4, 0, 0, 0, 0, time.Local),
					NextDay10:              time.Date(2022, 4, 5, 0, 0, 0, 0, time.Local),
					DeliveryDay:            time.Date(2022, 3, 24, 0, 0, 0, 0, time.Local),
					ProvisionalDeliveryDay: time.Date(2022, 3, 25, 0, 0, 0, 0, time.Local),
					BondDeliveryDay:        time.Date(2022, 3, 24, 0, 0, 0, 0, time.Local),
				},
				{
					CommonResponse: CommonResponse{
						SendDate:    time.Date(2022, 3, 21, 5, 33, 27, 293000000, time.Local),
						FeatureType: "LMDateZyouhou",
					},
					DayKey:                 DayKeyNextDay,
					PrevDay2:               time.Date(2022, 3, 18, 0, 0, 0, 0, time.Local),
					PrevDay3:               time.Date(2022, 3, 17, 0, 0, 0, 0, time.Local),
					PrevDay1:               time.Date(2022, 3, 22, 0, 0, 0, 0, time.Local),
					Today:                  time.Date(2022, 3, 23, 0, 0, 0, 0, time.Local),
					NextDay1:               time.Date(2022, 3, 24, 0, 0, 0, 0, time.Local),
					NextDay2:               time.Date(2022, 3, 25, 0, 0, 0, 0, time.Local),
					NextDay3:               time.Date(2022, 3, 28, 0, 0, 0, 0, time.Local),
					NextDay4:               time.Date(2022, 3, 29, 0, 0, 0, 0, time.Local),
					NextDay5:               time.Date(2022, 3, 30, 0, 0, 0, 0, time.Local),
					NextDay6:               time.Date(2022, 3, 31, 0, 0, 0, 0, time.Local),
					NextDay7:               time.Date(2022, 4, 1, 0, 0, 0, 0, time.Local),
					NextDay8:               time.Date(2022, 4, 4, 0, 0, 0, 0, time.Local),
					NextDay9:               time.Date(2022, 4, 5, 0, 0, 0, 0, time.Local),
					NextDay10:              time.Date(2022, 4, 6, 0, 0, 0, 0, time.Local),
					DeliveryDay:            time.Date(2022, 3, 25, 0, 0, 0, 0, time.Local),
					ProvisionalDeliveryDay: time.Date(2022, 3, 28, 0, 0, 0, 0, time.Local),
					BondDeliveryDay:        time.Date(2022, 3, 25, 0, 0, 0, 0, time.Local),
				},
			},
			want2: nil},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			stream1 := make(chan []byte)
			stream2 := make(chan error)
			go test.stream(stream1, stream2)
			client := &client{clock: test.clock, requester: &testRequester{stream1: stream1, stream2: stream2}}
			got1, got2 := client.BusinessDay(test.arg1, test.arg2, test.arg3)

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				t.Errorf("%s error\nresult: %+v, %+v\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(),
					!reflect.DeepEqual(test.want1, got1),
					!errors.Is(got2, test.want2),
					test.want1, test.want2,
					got1, got2)
			}
		})
	}
}

func Test_client_BusinessDay_Execute(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	userId := "user-id"
	password := "password"

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
	got3, got4 := client.BusinessDay(ctx, session, BusinessDayRequest{})
	log.Printf("%+v, %+v\n", got3, got4)
	for _, b := range got3 {
		log.Printf("%+v\n", b)
	}
}
