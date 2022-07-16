package tachibana

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"testing"
	"time"
)

func Test_StreamRequest_Query(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		req   StreamRequest
		want1 []byte
	}{
		{name: "空リクエストをバイト列に出来る",
			req:   StreamRequest{},
			want1: []byte("p_rid=22&p_board_no=1000&p_eno=0&p_evt_cmd=")},
		{name: "各種設定をバイト列に出来る",
			req: StreamRequest{
				StartStreamNumber: 1000,
				StreamEventTypes: []EventType{
					EventTypeErrorStatus,
					EventTypeKeepAlive,
					EventTypeContract,
					EventTypeSystemStatus,
					EventTypeOperationStatus,
				},
			},
			want1: []byte("p_rid=22&p_board_no=1000&p_eno=1000&p_evt_cmd=ST,KP,EC,SS,US")},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got1 := test.req.Query()
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), string(test.want1), string(got1))
			}
		})
	}
}

func Test_client_streamResponseToMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  []byte
		want1 map[string][]string
	}{
		{name: "=がない場合は空文字を含んだ配列にしておく",
			arg1: []byte{112, 95, 110, 111, 2, 50, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 53, 49, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 51, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 67, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 53, 51, 48, 48, 50, 1, 112, 95, 76, 75, 2, 49, 1, 112, 95, 83, 83, 2, 49},
			want1: map[string][]string{
				"p_ALT":   {"1"},
				"p_CT":    {"20220712053002"},
				"p_ENO":   {"3"},
				"p_LK":    {"1"},
				"p_PV":    {"MSGSV"},
				"p_SS":    {"1"},
				"p_cmd":   {""},
				"p_date":  {"2022.07.12-21:01:56.551"},
				"p_err":   {""},
				"p_errno": {"0"},
				"p_no":    {"2"},
			}},
		{name: "SSのbodyをmapに変換できる",
			arg1: []byte{112, 95, 110, 111, 2, 50, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 53, 49, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 83, 83, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 51, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 67, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 53, 51, 48, 48, 50, 1, 112, 95, 76, 75, 2, 49, 1, 112, 95, 83, 83, 2, 49},
			want1: map[string][]string{
				"p_ALT":   {"1"},
				"p_CT":    {"20220712053002"},
				"p_ENO":   {"3"},
				"p_LK":    {"1"},
				"p_PV":    {"MSGSV"},
				"p_SS":    {"1"},
				"p_cmd":   {"SS"},
				"p_date":  {"2022.07.12-21:01:56.551"},
				"p_err":   {""},
				"p_errno": {"0"},
				"p_no":    {"2"},
			}},
		{name: "USのbodyをmapに変換できる",
			arg1: []byte{112, 95, 110, 111, 2, 51, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 53, 52, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 85, 83, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 49, 51, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 67, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 55, 49, 48, 49, 56, 1, 112, 95, 77, 67, 2, 48, 49, 1, 112, 95, 71, 83, 67, 68, 2, 49, 48, 49, 1, 112, 95, 83, 72, 83, 66, 2, 48, 52, 1, 112, 95, 85, 67, 2, 48, 50, 1, 112, 95, 85, 85, 2, 48, 50, 48, 50, 1, 112, 95, 69, 68, 75, 2, 48, 1, 112, 95, 85, 83, 2, 48, 53, 48},
			want1: map[string][]string{
				"p_ALT":   {"1"},
				"p_CT":    {"20220712071018"},
				"p_EDK":   {"0"},
				"p_ENO":   {"13"},
				"p_GSCD":  {"101"},
				"p_MC":    {"01"},
				"p_PV":    {"MSGSV"},
				"p_SHSB":  {"04"},
				"p_UC":    {"02"},
				"p_US":    {"050"},
				"p_UU":    {"0202"},
				"p_cmd":   {"US"},
				"p_date":  {"2022.07.12-21:01:56.554"},
				"p_err":   {""},
				"p_errno": {"0"},
				"p_no":    {"3"},
			}},
		{name: "ECのbodyをmapに変換できる",
			arg1: []byte{112, 95, 110, 111, 2, 56, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 57, 52, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 69, 67, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 49, 54, 50, 48, 48, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 78, 84, 2, 49, 48, 48, 1, 112, 95, 79, 78, 2, 49, 50, 48, 48, 52, 56, 53, 48, 1, 112, 95, 69, 68, 2, 50, 48, 50, 50, 48, 55, 49, 50, 1, 112, 95, 79, 79, 78, 2, 48, 1, 112, 95, 79, 84, 2, 49, 1, 112, 95, 83, 84, 2, 49, 1, 112, 95, 73, 67, 2, 49, 52, 55, 53, 1, 112, 95, 77, 67, 2, 48, 48, 1, 112, 95, 66, 66, 75, 66, 2, 51, 1, 112, 95, 84, 72, 75, 66, 2, 50, 1, 112, 95, 67, 82, 83, 74, 2, 48, 1, 112, 95, 67, 82, 80, 82, 75, 66, 2, 49, 1, 112, 95, 67, 82, 80, 82, 2, 48, 46, 48, 48, 48, 48, 48, 48, 1, 112, 95, 67, 82, 83, 82, 2, 51, 1, 112, 95, 67, 82, 84, 75, 83, 82, 2, 48, 1, 112, 95, 67, 82, 69, 80, 83, 82, 2, 48, 1, 112, 95, 67, 82, 69, 88, 83, 82, 2, 48, 1, 112, 95, 79, 68, 83, 84, 2, 48, 1, 112, 95, 75, 79, 70, 71, 2, 48, 1, 112, 95, 84, 84, 83, 84, 2, 48, 1, 112, 95, 69, 88, 83, 84, 2, 48, 1, 112, 95, 76, 77, 73, 84, 2, 48, 48, 48, 48, 48, 48, 48, 48, 1, 112, 95, 69, 80, 82, 67, 2, 1, 112, 95, 69, 88, 80, 82, 2, 48, 46, 48, 48, 48, 48, 48, 48, 1, 112, 95, 69, 88, 83, 82, 2, 48, 1, 112, 95, 69, 88, 82, 67, 2, 1, 112, 95, 69, 88, 68, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 56, 53, 57, 48, 51, 1, 112, 95, 73, 78, 2, 239, 189, 137, 227, 130, 183, 227, 130, 167, 227, 130, 162, 227, 131, 188, 227, 130, 186, 239, 188, 180, 239, 188, 175, 239, 188, 176, 239, 188, 169, 239, 188, 184, 1, 112, 95, 85, 80, 83, 74, 2, 1, 112, 95, 85, 80, 69, 88, 83, 82, 2, 1, 112, 95, 85, 80, 80, 82, 75, 66, 2, 1, 112, 95, 85, 80, 80, 82, 2, 1, 112, 95, 85, 80, 83, 82, 2, 1, 112, 95, 85, 80, 76, 77, 73, 84, 2, 1, 112, 95, 85, 80, 71, 75, 67, 68, 80, 82, 2, 1, 112, 95, 85, 80, 71, 75, 80, 82, 75, 66, 2, 1, 112, 95, 85, 80, 71, 75, 80, 82, 2},
			want1: map[string][]string{
				"p_ALT":      {"1"},
				"p_BBKB":     {"3"},
				"p_CREPSR":   {"0"},
				"p_CREXSR":   {"0"},
				"p_CRPR":     {"0.000000"},
				"p_CRPRKB":   {"1"},
				"p_CRSJ":     {"0"},
				"p_CRSR":     {"3"},
				"p_CRTKSR":   {"0"},
				"p_ED":       {"20220712"},
				"p_ENO":      {"16200"},
				"p_EPRC":     {""},
				"p_EXDT":     {"20220712085903"},
				"p_EXPR":     {"0.000000"},
				"p_EXRC":     {""},
				"p_EXSR":     {"0"},
				"p_EXST":     {"0"},
				"p_IC":       {"1475"},
				"p_IN":       {"ｉシェアーズＴＯＰＩＸ"},
				"p_KOFG":     {"0"},
				"p_LMIT":     {"00000000"},
				"p_MC":       {"00"},
				"p_NT":       {"100"},
				"p_ODST":     {"0"},
				"p_ON":       {"12004850"},
				"p_OON":      {"0"},
				"p_OT":       {"1"},
				"p_PV":       {"MSGSV"},
				"p_ST":       {"1"},
				"p_THKB":     {"2"},
				"p_TTST":     {"0"},
				"p_UPEXSR":   {""},
				"p_UPGKCDPR": {""},
				"p_UPGKPR":   {""},
				"p_UPGKPRKB": {""},
				"p_UPLMIT":   {""},
				"p_UPPR":     {""},
				"p_UPPRKB":   {""},
				"p_UPSJ":     {""},
				"p_UPSR":     {""},
				"p_cmd":      {"EC"},
				"p_date":     {"2022.07.12-21:01:56.594"},
				"p_err":      {""},
				"p_errno":    {"0"},
				"p_no":       {"8"},
			}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			client := &client{}
			got1 := client.streamResponseToMap(test.arg1)
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}

func Test_StreamResponse_getFromMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  map[string][]string
		arg2  string
		want1 []string
	}{
		{name: "マップにある値を取ればそのまま返される",
			arg1:  map[string][]string{"foo": {"bar"}},
			arg2:  "foo",
			want1: []string{"bar"}},
		{name: "マップにない値を取ろうとすると、空文字が入っている配列が返される",
			arg1:  map[string][]string{"foo": {"bar"}},
			arg2:  "bar",
			want1: []string{""}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			res := CommonStreamResponse{}
			got1 := res.getFromMap(test.arg1, test.arg2)
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}

func Test_StreamResponse_GetEventType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		res   *CommonStreamResponse
		want1 EventType
	}{
		{name: "イベントタイプが取れる",
			res:   &CommonStreamResponse{EventType: EventTypeContract},
			want1: EventTypeContract},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got1 := test.res.GetEventType()
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}

func Test_StreamResponse_GetErrorNo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		res   *CommonStreamResponse
		want1 ErrorNo
	}{
		{name: "イベントタイプが取れる",
			res:   &CommonStreamResponse{ErrorNo: ErrorNoProblem},
			want1: ErrorNoProblem},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got1 := test.res.GetErrorNo()
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}

func Test_StreamResponse_GetErrorText(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		res   *CommonStreamResponse
		want1 string
	}{
		{name: "イベントタイプが取れる",
			res:   &CommonStreamResponse{ErrorText: "parameter error."},
			want1: "parameter error."},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got1 := test.res.GetErrorText()
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}

func Test_SystemStatusStreamResponse_parse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  map[string][]string
		arg2  []byte
		want1 SystemStatusStreamResponse
	}{
		{name: "マップの値を使って構造体に反映できる",
			arg1: map[string][]string{
				"p_ALT":   {"1"},
				"p_CT":    {"20220712053002"},
				"p_ENO":   {"3"},
				"p_LK":    {"1"},
				"p_PV":    {"MSGSV"},
				"p_SS":    {"1"},
				"p_cmd":   {"SS"},
				"p_date":  {"2022.07.12-21:01:56.551"},
				"p_err":   {""},
				"p_errno": {"0"},
				"p_no":    {"2"},
			},
			arg2: []byte{112, 95, 110, 111, 2, 50, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 53, 49, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 83, 83, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 51, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 67, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 53, 51, 48, 48, 50, 1, 112, 95, 76, 75, 2, 49, 1, 112, 95, 83, 83, 2, 49},
			want1: SystemStatusStreamResponse{
				CommonStreamResponse: CommonStreamResponse{
					EventType:      EventTypeSystemStatus,
					StreamNumber:   2,
					StreamDateTime: time.Date(2022, 7, 12, 21, 1, 56, 551000000, time.Local),
					ErrorNo:        ErrorNoProblem,
					ErrorText:      "",
					Body:           []byte{112, 95, 110, 111, 2, 50, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 53, 49, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 83, 83, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 51, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 67, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 53, 51, 48, 48, 50, 1, 112, 95, 76, 75, 2, 49, 1, 112, 95, 83, 83, 2, 49},
				},
				Provider:       "MSGSV",
				EventNo:        3,
				FirstTime:      true,
				UpdateDateTime: time.Date(2022, 7, 12, 5, 30, 2, 0, time.Local),
				ApprovalLogin:  ApprovalLoginApproval,
				SystemStatus:   SystemStatusOpening,
			}},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var res SystemStatusStreamResponse
			res.parse(test.arg1, test.arg2)
			if !reflect.DeepEqual(test.want1, res) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, res)
			}
		})
	}
}

func Test_OperationStatusStreamResponse_parse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  map[string][]string
		arg2  []byte
		want1 OperationStatusStreamResponse
	}{
		{name: "マップの値を使って構造体に反映できる",
			arg1: map[string][]string{
				"p_ALT":   {"1"},
				"p_CT":    {"20220712071018"},
				"p_EDK":   {"0"},
				"p_ENO":   {"13"},
				"p_GSCD":  {"101"},
				"p_MC":    {"01"},
				"p_PV":    {"MSGSV"},
				"p_SHSB":  {"04"},
				"p_UC":    {"02"},
				"p_US":    {"050"},
				"p_UU":    {"0202"},
				"p_cmd":   {"US"},
				"p_date":  {"2022.07.12-21:01:56.554"},
				"p_err":   {""},
				"p_errno": {"0"},
				"p_no":    {"3"},
			},
			arg2: []byte{112, 95, 110, 111, 2, 50, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 53, 49, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 83, 83, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 51, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 67, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 53, 51, 48, 48, 50, 1, 112, 95, 76, 75, 2, 49, 1, 112, 95, 83, 83, 2, 49},
			want1: OperationStatusStreamResponse{
				CommonStreamResponse: CommonStreamResponse{
					EventType:      EventTypeOperationStatus,
					StreamNumber:   3,
					StreamDateTime: time.Date(2022, 7, 12, 21, 1, 56, 554000000, time.Local),
					ErrorNo:        ErrorNoProblem,
					ErrorText:      "",
					Body:           []byte{112, 95, 110, 111, 2, 50, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 53, 49, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 83, 83, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 51, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 67, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 53, 51, 48, 48, 50, 1, 112, 95, 76, 75, 2, 49, 1, 112, 95, 83, 83, 2, 49},
				},
				Provider:          "MSGSV",
				EventNo:           13,
				FirstTime:         true,
				UpdateDateTime:    time.Date(2022, 7, 12, 7, 10, 18, 0, time.Local),
				Exchange:          ExchangeDaishou,
				AssetCode:         "101",
				ProductType:       "04",
				OperationCategory: "02",
				OperationUnit:     "0202",
				BusinessDayType:   "0",
				OperationStatus:   "050",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var res OperationStatusStreamResponse
			res.parse(test.arg1, test.arg2)
			if !reflect.DeepEqual(test.want1, res) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, res)
			}
		})
	}
}

func Test_ContractStreamResponse_parse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  map[string][]string
		arg2  []byte
		want1 ContractStreamResponse
	}{
		{name: "注文受付の値を使って構造体に反映できる",
			arg1: map[string][]string{
				"p_ALT":      {"1"},
				"p_BBKB":     {"3"},
				"p_CREPSR":   {"0"},
				"p_CREXSR":   {"0"},
				"p_CRPR":     {"0.000000"},
				"p_CRPRKB":   {"1"},
				"p_CRSJ":     {"0"},
				"p_CRSR":     {"3"},
				"p_CRTKSR":   {"0"},
				"p_ED":       {"20220712"},
				"p_ENO":      {"16200"},
				"p_EPRC":     {""},
				"p_EXDT":     {"20220712085903"},
				"p_EXPR":     {"0.000000"},
				"p_EXRC":     {""},
				"p_EXSR":     {"0"},
				"p_EXST":     {"0"},
				"p_IC":       {"1475"},
				"p_IN":       {"ｉシェアーズＴＯＰＩＸ"},
				"p_KOFG":     {"0"},
				"p_LMIT":     {"00000000"},
				"p_MC":       {"00"},
				"p_NT":       {"100"},
				"p_ODST":     {"0"},
				"p_ON":       {"12004850"},
				"p_OON":      {"0"},
				"p_OT":       {"1"},
				"p_PV":       {"MSGSV"},
				"p_ST":       {"1"},
				"p_THKB":     {"2"},
				"p_TTST":     {"0"},
				"p_UPEXSR":   {""},
				"p_UPGKCDPR": {""},
				"p_UPGKPR":   {""},
				"p_UPGKPRKB": {""},
				"p_UPLMIT":   {""},
				"p_UPPR":     {""},
				"p_UPPRKB":   {""},
				"p_UPSJ":     {""},
				"p_UPSR":     {""},
				"p_cmd":      {"EC"},
				"p_date":     {"2022.07.12-21:01:56.594"},
				"p_err":      {""},
				"p_errno":    {"0"},
				"p_no":       {"8"},
			},
			want1: ContractStreamResponse{
				CommonStreamResponse: CommonStreamResponse{
					EventType:      EventTypeContract,
					StreamNumber:   8,
					StreamDateTime: time.Date(2022, 7, 12, 21, 1, 56, 594000000, time.Local),
					ErrorNo:        ErrorNoProblem,
					ErrorText:      "",
				},
				Provider:                 "MSGSV",
				EventNo:                  16200,
				FirstTime:                true,
				StreamOrderType:          StreamOrderTypeReceived,
				OrderNumber:              "12004850",
				ExecutionDate:            time.Date(2022, 7, 12, 0, 0, 0, 0, time.Local),
				ParentOrderNumber:        "0",
				ParentOrder:              true,
				ProductType:              ProductTypeStock,
				IssueCode:                "1475",
				Exchange:                 ExchangeToushou,
				Side:                     SideBuy,
				TradeType:                TradeTypeStandardEntry,
				ExecutionTiming:          ExecutionTimingNormal,
				ExecutionType:            ExecutionTypeMarket,
				Price:                    0,
				Quantity:                 3,
				CancelQuantity:           0,
				ExpireQuantity:           0,
				ContractQuantity:         0,
				StreamOrderStatus:        StreamOrderStatusNew,
				CarryOverType:            CarryOverTypeToday,
				CancelOrderStatus:        CancelOrderStatusNoCorrect,
				ContractStatus:           ContractStatusInOrder,
				ExpireDate:               time.Date(2022, 7, 12, 0, 0, 0, 0, time.Local),
				SecurityExpireReason:     "",
				SecurityContractPrice:    0,
				SecurityContractQuantity: 0,
				SecurityError:            "",
				NotifyDateTime:           time.Date(2022, 7, 12, 8, 59, 3, 0, time.Local),
				IssueName:                "ｉシェアーズＴＯＰＩＸ",
				CorrectExecutionTiming:   "",
				CorrectContractQuantity:  0,
				CorrectExecutionType:     "",
				CorrectPrice:             0,
				CorrectQuantity:          0,
				CorrectExpireDate:        time.Time{},
				CorrectStopOrderType:     "",
				CorrectTriggerPrice:      0,
				CorrectStopOrderPrice:    0,
			},
		},
		{name: "約定の値を使って構造体に反映できる",
			arg1: map[string][]string{
				"p_ALT":      {"1"},
				"p_BBKB":     {"3"},
				"p_CREPSR":   {"0"},
				"p_CREXSR":   {"3"},
				"p_CRPR":     {"0.000000"},
				"p_CRPRKB":   {"1"},
				"p_CRSJ":     {"0"},
				"p_CRSR":     {"0"},
				"p_CRTKSR":   {"0"},
				"p_ED":       {"20220712"},
				"p_ENO":      {"17392"},
				"p_EPRC":     {"0000"},
				"p_EXDT":     {"20220712090000"},
				"p_EXPR":     {"1966.000000"},
				"p_EXRC":     {""},
				"p_EXSR":     {"3"},
				"p_EXST":     {"2"},
				"p_IC":       {"1475"},
				"p_IN":       {"ｉシェアーズＴＯＰＩＸ"},
				"p_KOFG":     {"0"},
				"p_LMIT":     {"00000000"},
				"p_MC":       {"00"},
				"p_NT":       {"12"},
				"p_ODST":     {"1"},
				"p_ON":       {"12004850"},
				"p_OON":      {"0"},
				"p_OT":       {"1"},
				"p_PV":       {"MSGSV"},
				"p_ST":       {"1"},
				"p_THKB":     {"2"},
				"p_TTST":     {"0"},
				"p_UPEXSR":   {""},
				"p_UPGKCDPR": {""},
				"p_UPGKPR":   {""},
				"p_UPGKPRKB": {""},
				"p_UPLMIT":   {""},
				"p_UPPR":     {""},
				"p_UPPRKB":   {""},
				"p_UPSJ":     {""},
				"p_UPSR":     {""},
				"p_cmd":      {"EC"},
				"p_date":     {"2022.07.12-21:01:56.627"},
				"p_err":      {""},
				"p_errno":    {"0"},
				"p_no":       {"14"},
			},
			want1: ContractStreamResponse{
				CommonStreamResponse: CommonStreamResponse{
					EventType:      EventTypeContract,
					StreamNumber:   14,
					StreamDateTime: time.Date(2022, 7, 12, 21, 1, 56, 627000000, time.Local),
					ErrorNo:        ErrorNoProblem,
					ErrorText:      "",
				},
				Provider:                 "MSGSV",
				EventNo:                  17392,
				FirstTime:                true,
				StreamOrderType:          StreamOrderTypeContract,
				OrderNumber:              "12004850",
				ExecutionDate:            time.Date(2022, 7, 12, 0, 0, 0, 0, time.Local),
				ParentOrderNumber:        "0",
				ParentOrder:              true,
				ProductType:              ProductTypeStock,
				IssueCode:                "1475",
				Exchange:                 ExchangeToushou,
				Side:                     SideBuy,
				TradeType:                TradeTypeStandardEntry,
				ExecutionTiming:          ExecutionTimingNormal,
				ExecutionType:            ExecutionTypeMarket,
				Price:                    0,
				Quantity:                 0,
				CancelQuantity:           0,
				ExpireQuantity:           0,
				ContractQuantity:         3,
				StreamOrderStatus:        StreamOrderStatusReceived,
				CarryOverType:            CarryOverTypeToday,
				CancelOrderStatus:        CancelOrderStatusNoCorrect,
				ContractStatus:           ContractStatusDone,
				ExpireDate:               time.Date(2022, 7, 12, 0, 0, 0, 0, time.Local),
				SecurityExpireReason:     "0000",
				SecurityContractPrice:    1966,
				SecurityContractQuantity: 3,
				SecurityError:            "",
				NotifyDateTime:           time.Date(2022, 7, 12, 9, 0, 0, 0, time.Local),
				IssueName:                "ｉシェアーズＴＯＰＩＸ",
				CorrectExecutionTiming:   "",
				CorrectContractQuantity:  0,
				CorrectExecutionType:     "",
				CorrectPrice:             0,
				CorrectQuantity:          0,
				CorrectExpireDate:        time.Time{},
				CorrectStopOrderType:     "",
				CorrectTriggerPrice:      0,
				CorrectStopOrderPrice:    0,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var res ContractStreamResponse
			res.parse(test.arg1, test.arg2)
			if !reflect.DeepEqual(test.want1, res) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, res)
			}
		})
	}
}

func Test_client_Stream(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		stream func(stream1 chan<- []byte, stream2 chan<- error)
		arg1   context.Context
		arg2   *Session
		arg3   StreamRequest
		want1  []StreamResponse
		want2  error
	}{
		{name: "sessionがnilならエラー",
			stream: func(stream1 chan<- []byte, stream2 chan<- error) {
				defer close(stream1)
				defer close(stream2)
			},
			arg1:  context.Background(),
			arg2:  nil,
			arg3:  StreamRequest{},
			want1: nil,
			want2: NilArgumentErr},
		{name: "エラーが返されたらエラー",
			stream: func(stream1 chan<- []byte, stream2 chan<- error) {
				defer close(stream1)
				defer close(stream2)
				stream2 <- StatusNotOkErr
			},
			arg1:  context.Background(),
			arg2:  &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:  StreamRequest{},
			want1: nil,
			want2: StatusNotOkErr},
		{name: "失敗レスポンスが返されたらエラー",
			stream: func(stream1 chan<- []byte, stream2 chan<- error) {
				defer close(stream1)
				defer close(stream2)
				stream1 <- []byte{112, 95, 110, 111, 2, 49, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 52, 45, 48, 53, 58, 51, 55, 58, 48, 54, 46, 51, 57, 50, 1, 112, 95, 101, 114, 114, 110, 111, 2, 45, 49, 1, 112, 95, 101, 114, 114, 2, 112, 97, 114, 97, 109, 101, 116, 101, 114, 32, 101, 114, 114, 111, 114, 46, 1, 112, 95, 99, 109, 100, 2, 83, 84}
			},
			arg1:  context.Background(),
			arg2:  &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:  StreamRequest{},
			want1: nil,
			want2: StreamError},
		{name: "Keep Aliveは通知しない",
			stream: func(stream1 chan<- []byte, stream2 chan<- error) {
				defer close(stream1)
				defer close(stream2)
				stream1 <- []byte{112, 95, 110, 111, 2, 53, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 54, 45, 49, 54, 58, 52, 55, 58, 50, 56, 46, 49, 55, 52, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 75, 80}
			},
			arg1:  context.Background(),
			arg2:  &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3:  StreamRequest{},
			want1: nil,
			want2: nil},
		{name: "p_cmdがレスポンスに含まれていなかったらCommonStreamResponseとしてパース",
			stream: func(stream1 chan<- []byte, stream2 chan<- error) {
				defer close(stream1)
				defer close(stream2)
				stream1 <- []byte{112, 95, 110, 111, 2, 53, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 54, 45, 49, 54, 58, 52, 55, 58, 50, 56, 46, 49, 55, 52, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2}
			},
			arg1: context.Background(),
			arg2: &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3: StreamRequest{},
			want1: []StreamResponse{&CommonStreamResponse{
				EventType:      EventTypeUnspecified,
				StreamNumber:   5,
				StreamDateTime: time.Date(2022, 7, 16, 16, 47, 28, 174000000, time.Local),
				ErrorNo:        ErrorNoProblem,
				ErrorText:      "",
				Body:           []byte{112, 95, 110, 111, 2, 53, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 54, 45, 49, 54, 58, 52, 55, 58, 50, 56, 46, 49, 55, 52, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2},
			}},
			want2: nil},
		{name: "各種イベントを通知できる",
			stream: func(stream1 chan<- []byte, stream2 chan<- error) {
				defer close(stream1)
				defer close(stream2)
				stream1 <- []byte{112, 95, 110, 111, 2, 50, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 53, 49, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 83, 83, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 51, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 67, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 53, 51, 48, 48, 50, 1, 112, 95, 76, 75, 2, 49, 1, 112, 95, 83, 83, 2, 49}
				stream1 <- []byte{112, 95, 110, 111, 2, 51, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 53, 52, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 85, 83, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 49, 51, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 67, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 55, 49, 48, 49, 56, 1, 112, 95, 77, 67, 2, 48, 49, 1, 112, 95, 71, 83, 67, 68, 2, 49, 48, 49, 1, 112, 95, 83, 72, 83, 66, 2, 48, 52, 1, 112, 95, 85, 67, 2, 48, 50, 1, 112, 95, 85, 85, 2, 48, 50, 48, 50, 1, 112, 95, 69, 68, 75, 2, 48, 1, 112, 95, 85, 83, 2, 48, 53, 48}
				stream1 <- []byte{112, 95, 110, 111, 2, 56, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 57, 52, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 69, 67, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 49, 54, 50, 48, 48, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 78, 84, 2, 49, 48, 48, 1, 112, 95, 79, 78, 2, 49, 50, 48, 48, 52, 56, 53, 48, 1, 112, 95, 69, 68, 2, 50, 48, 50, 50, 48, 55, 49, 50, 1, 112, 95, 79, 79, 78, 2, 48, 1, 112, 95, 79, 84, 2, 49, 1, 112, 95, 83, 84, 2, 49, 1, 112, 95, 73, 67, 2, 49, 52, 55, 53, 1, 112, 95, 77, 67, 2, 48, 48, 1, 112, 95, 66, 66, 75, 66, 2, 51, 1, 112, 95, 84, 72, 75, 66, 2, 50, 1, 112, 95, 67, 82, 83, 74, 2, 48, 1, 112, 95, 67, 82, 80, 82, 75, 66, 2, 49, 1, 112, 95, 67, 82, 80, 82, 2, 48, 46, 48, 48, 48, 48, 48, 48, 1, 112, 95, 67, 82, 83, 82, 2, 51, 1, 112, 95, 67, 82, 84, 75, 83, 82, 2, 48, 1, 112, 95, 67, 82, 69, 80, 83, 82, 2, 48, 1, 112, 95, 67, 82, 69, 88, 83, 82, 2, 48, 1, 112, 95, 79, 68, 83, 84, 2, 48, 1, 112, 95, 75, 79, 70, 71, 2, 48, 1, 112, 95, 84, 84, 83, 84, 2, 48, 1, 112, 95, 69, 88, 83, 84, 2, 48, 1, 112, 95, 76, 77, 73, 84, 2, 48, 48, 48, 48, 48, 48, 48, 48, 1, 112, 95, 69, 80, 82, 67, 2, 1, 112, 95, 69, 88, 80, 82, 2, 48, 46, 48, 48, 48, 48, 48, 48, 1, 112, 95, 69, 88, 83, 82, 2, 48, 1, 112, 95, 69, 88, 82, 67, 2, 1, 112, 95, 69, 88, 68, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 56, 53, 57, 48, 51, 1, 112, 95, 73, 78, 2, 239, 189, 137, 227, 130, 183, 227, 130, 167, 227, 130, 162, 227, 131, 188, 227, 130, 186, 239, 188, 180, 239, 188, 175, 239, 188, 176, 239, 188, 169, 239, 188, 184, 1, 112, 95, 85, 80, 83, 74, 2, 1, 112, 95, 85, 80, 69, 88, 83, 82, 2, 1, 112, 95, 85, 80, 80, 82, 75, 66, 2, 1, 112, 95, 85, 80, 80, 82, 2, 1, 112, 95, 85, 80, 83, 82, 2, 1, 112, 95, 85, 80, 76, 77, 73, 84, 2, 1, 112, 95, 85, 80, 71, 75, 67, 68, 80, 82, 2, 1, 112, 95, 85, 80, 71, 75, 80, 82, 75, 66, 2, 1, 112, 95, 85, 80, 71, 75, 80, 82, 2}
			},
			arg1: context.Background(),
			arg2: &Session{lastRequestNo: 1, RequestURL: "", EventURL: ""},
			arg3: StreamRequest{},
			want1: []StreamResponse{
				&SystemStatusStreamResponse{
					CommonStreamResponse: CommonStreamResponse{
						EventType:      EventTypeSystemStatus,
						StreamNumber:   2,
						StreamDateTime: time.Date(2022, 7, 12, 21, 1, 56, 551000000, time.Local),
						ErrorNo:        ErrorNoProblem,
						ErrorText:      "",
						Body:           []byte{112, 95, 110, 111, 2, 50, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 53, 49, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 83, 83, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 51, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 67, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 53, 51, 48, 48, 50, 1, 112, 95, 76, 75, 2, 49, 1, 112, 95, 83, 83, 2, 49},
					},
					Provider:       "MSGSV",
					EventNo:        3,
					FirstTime:      true,
					UpdateDateTime: time.Date(2022, 7, 12, 5, 30, 2, 0, time.Local),
					ApprovalLogin:  ApprovalLoginApproval,
					SystemStatus:   SystemStatusOpening,
				},
				&OperationStatusStreamResponse{
					CommonStreamResponse: CommonStreamResponse{
						EventType:      EventTypeOperationStatus,
						StreamNumber:   3,
						StreamDateTime: time.Date(2022, 7, 12, 21, 1, 56, 554000000, time.Local),
						ErrorNo:        ErrorNoProblem,
						ErrorText:      "",
						Body:           []byte{112, 95, 110, 111, 2, 51, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 53, 52, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 85, 83, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 49, 51, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 67, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 55, 49, 48, 49, 56, 1, 112, 95, 77, 67, 2, 48, 49, 1, 112, 95, 71, 83, 67, 68, 2, 49, 48, 49, 1, 112, 95, 83, 72, 83, 66, 2, 48, 52, 1, 112, 95, 85, 67, 2, 48, 50, 1, 112, 95, 85, 85, 2, 48, 50, 48, 50, 1, 112, 95, 69, 68, 75, 2, 48, 1, 112, 95, 85, 83, 2, 48, 53, 48},
					},
					Provider:          "MSGSV",
					EventNo:           13,
					FirstTime:         true,
					UpdateDateTime:    time.Date(2022, 7, 12, 7, 10, 18, 0, time.Local),
					Exchange:          ExchangeDaishou,
					AssetCode:         "101",
					ProductType:       "04",
					OperationCategory: "02",
					OperationUnit:     "0202",
					BusinessDayType:   "0",
					OperationStatus:   "050",
				},
				&ContractStreamResponse{
					CommonStreamResponse: CommonStreamResponse{
						EventType:      EventTypeContract,
						StreamNumber:   8,
						StreamDateTime: time.Date(2022, 7, 12, 21, 1, 56, 594000000, time.Local),
						ErrorNo:        ErrorNoProblem,
						ErrorText:      "",
						Body:           []byte{112, 95, 110, 111, 2, 56, 1, 112, 95, 100, 97, 116, 101, 2, 50, 48, 50, 50, 46, 48, 55, 46, 49, 50, 45, 50, 49, 58, 48, 49, 58, 53, 54, 46, 53, 57, 52, 1, 112, 95, 101, 114, 114, 110, 111, 2, 48, 1, 112, 95, 101, 114, 114, 2, 1, 112, 95, 99, 109, 100, 2, 69, 67, 1, 112, 95, 80, 86, 2, 77, 83, 71, 83, 86, 1, 112, 95, 69, 78, 79, 2, 49, 54, 50, 48, 48, 1, 112, 95, 65, 76, 84, 2, 49, 1, 112, 95, 78, 84, 2, 49, 48, 48, 1, 112, 95, 79, 78, 2, 49, 50, 48, 48, 52, 56, 53, 48, 1, 112, 95, 69, 68, 2, 50, 48, 50, 50, 48, 55, 49, 50, 1, 112, 95, 79, 79, 78, 2, 48, 1, 112, 95, 79, 84, 2, 49, 1, 112, 95, 83, 84, 2, 49, 1, 112, 95, 73, 67, 2, 49, 52, 55, 53, 1, 112, 95, 77, 67, 2, 48, 48, 1, 112, 95, 66, 66, 75, 66, 2, 51, 1, 112, 95, 84, 72, 75, 66, 2, 50, 1, 112, 95, 67, 82, 83, 74, 2, 48, 1, 112, 95, 67, 82, 80, 82, 75, 66, 2, 49, 1, 112, 95, 67, 82, 80, 82, 2, 48, 46, 48, 48, 48, 48, 48, 48, 1, 112, 95, 67, 82, 83, 82, 2, 51, 1, 112, 95, 67, 82, 84, 75, 83, 82, 2, 48, 1, 112, 95, 67, 82, 69, 80, 83, 82, 2, 48, 1, 112, 95, 67, 82, 69, 88, 83, 82, 2, 48, 1, 112, 95, 79, 68, 83, 84, 2, 48, 1, 112, 95, 75, 79, 70, 71, 2, 48, 1, 112, 95, 84, 84, 83, 84, 2, 48, 1, 112, 95, 69, 88, 83, 84, 2, 48, 1, 112, 95, 76, 77, 73, 84, 2, 48, 48, 48, 48, 48, 48, 48, 48, 1, 112, 95, 69, 80, 82, 67, 2, 1, 112, 95, 69, 88, 80, 82, 2, 48, 46, 48, 48, 48, 48, 48, 48, 1, 112, 95, 69, 88, 83, 82, 2, 48, 1, 112, 95, 69, 88, 82, 67, 2, 1, 112, 95, 69, 88, 68, 84, 2, 50, 48, 50, 50, 48, 55, 49, 50, 48, 56, 53, 57, 48, 51, 1, 112, 95, 73, 78, 2, 239, 189, 137, 227, 130, 183, 227, 130, 167, 227, 130, 162, 227, 131, 188, 227, 130, 186, 239, 188, 180, 239, 188, 175, 239, 188, 176, 239, 188, 169, 239, 188, 184, 1, 112, 95, 85, 80, 83, 74, 2, 1, 112, 95, 85, 80, 69, 88, 83, 82, 2, 1, 112, 95, 85, 80, 80, 82, 75, 66, 2, 1, 112, 95, 85, 80, 80, 82, 2, 1, 112, 95, 85, 80, 83, 82, 2, 1, 112, 95, 85, 80, 76, 77, 73, 84, 2, 1, 112, 95, 85, 80, 71, 75, 67, 68, 80, 82, 2, 1, 112, 95, 85, 80, 71, 75, 80, 82, 75, 66, 2, 1, 112, 95, 85, 80, 71, 75, 80, 82, 2},
					},
					Provider:                 "MSGSV",
					EventNo:                  16200,
					FirstTime:                true,
					StreamOrderType:          StreamOrderTypeReceived,
					OrderNumber:              "12004850",
					ExecutionDate:            time.Date(2022, 7, 12, 0, 0, 0, 0, time.Local),
					ParentOrderNumber:        "0",
					ParentOrder:              true,
					ProductType:              ProductTypeStock,
					IssueCode:                "1475",
					Exchange:                 ExchangeToushou,
					Side:                     SideBuy,
					TradeType:                TradeTypeStandardEntry,
					ExecutionTiming:          ExecutionTimingNormal,
					ExecutionType:            ExecutionTypeMarket,
					Price:                    0,
					Quantity:                 3,
					CancelQuantity:           0,
					ExpireQuantity:           0,
					ContractQuantity:         0,
					StreamOrderStatus:        StreamOrderStatusNew,
					CarryOverType:            CarryOverTypeToday,
					CancelOrderStatus:        CancelOrderStatusNoCorrect,
					ContractStatus:           ContractStatusInOrder,
					ExpireDate:               time.Date(2022, 7, 12, 0, 0, 0, 0, time.Local),
					SecurityExpireReason:     "",
					SecurityContractPrice:    0,
					SecurityContractQuantity: 0,
					SecurityError:            "",
					NotifyDateTime:           time.Date(2022, 7, 12, 8, 59, 3, 0, time.Local),
					IssueName:                "ｉシェアーズＴＯＰＩＸ",
					CorrectExecutionTiming:   "",
					CorrectContractQuantity:  0,
					CorrectExecutionType:     "",
					CorrectPrice:             0,
					CorrectQuantity:          0,
					CorrectExpireDate:        time.Time{},
					CorrectStopOrderType:     "",
					CorrectTriggerPrice:      0,
					CorrectStopOrderPrice:    0,
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
			client := &client{requester: &testRequester{stream1: stream1, stream2: stream2}}
			ch1, ch2 := client.Stream(test.arg1, test.arg2, test.arg3)

			var got1 []StreamResponse
			var got2 error
			func() {
				for {
					select {
					case err, ok := <-ch2:
						if ok {
							got2 = err
							return
						}
					case b, ok := <-ch1:
						if !ok {
							return
						}
						got1 = append(got1, b)
					}
				}
			}()

			if !reflect.DeepEqual(test.want1, got1) || !errors.Is(got2, test.want2) {
				_want1, _ := json.Marshal(test.want1)
				_got1, _ := json.Marshal(got1)
				t.Errorf("%s error\nwant: %+v, %+v\ngot: %+v, %+v\n", t.Name(), string(_want1), test.want2, string(_got1), got2)
			}
		})
	}
}

func Test_client_Stream_Execute(t *testing.T) {
	t.Skip("実際にAPIを叩くテストのため、通常はスキップ")
	t.Parallel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	userId := "user-id"
	password := "password"

	client := NewClient(EnvironmentProduction, ApiVersionLatest)
	got1, got2 := client.Login(context.Background(), LoginRequest{
		UserId:   userId,
		Password: password,
	})
	log.Printf("%+v, %+v\n", got1, got2)

	session, err := got1.Session()
	if err != nil {
		t.Errorf("session: %+v\n", err)
	}

	got3, got4 := client.Stream(context.Background(), session, StreamRequest{
		StreamEventTypes: []EventType{
			EventTypeErrorStatus,
			EventTypeKeepAlive,
			//EventTypeMarketPrice,
			EventTypeContract,
			EventTypeNews,
			EventTypeSystemStatus,
			EventTypeOperationStatus,
		}})
	for {
		select {
		case res, ok := <-got3:
			if !ok {
				log.Println("got3 closed")
				return
			}
			log.Printf("%+v\n", res)
		case err, ok := <-got4:
			if !ok {
				log.Println("got4 closed")
				return
			}
			log.Println(err)
		}
	}
}
