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
				ColumnNumber:      []int{1, 2, 3},
				IssueCodes:        []string{"1111", "2222", "3333"},
				MarketCodes:       []Exchange{ExchangeToushou, ExchangeMeishou, ExchangeSatsushou},
				StartStreamNumber: 1000,
				StreamEventTypes: []EventType{
					EventTypeErrorStatus,
					EventTypeKeepAlive,
					EventTypeContract,
					EventTypeSystemStatus,
					EventTypeOperationStatus,
				},
			},
			want1: []byte("p_rid=22&p_board_no=1000&p_gyou_no=1,2,3&p_issue_code=1111,2222,3333&p_mkt_code=00,02,07&p_eno=1000&p_evt_cmd=ST,KP,EC,SS,US")},
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
		{name: "NSのbodyをmapに変換できる",
			arg1: []byte("p_no\x024864\x01p_date\x022022.07.25-17:49:57.424\x01p_errno\x020\x01p_err\x02\x01p_cmd\x02NS\x01p_PV\x02QNSD\x01p_ENO\x0295371\x01p_ALT\x021\x01p_ID\x0220220725174500_NOV6627\x01p_DT\x0220220725\x01p_TM\x02174500\x01p_CGN\x021\x01p_CGL\x02129\x01p_GRN\x021\x01p_GRL\x0262199\x01p_ISN\x021\x01p_ISL\x023494\x01p_SKF\x021\x01p_UPF\x02\x01p_HLD\x02<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\x01p_TX\x02<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\\r\\r\\r\\r開示会社：マリオン(3494)\\r\\r開示書類：i-Bond(アイボンド）サイトリニューアルのお知らせ\\r\\r開示日時：2022/07/25 17:45\\r\\r\\r\\r＜引用＞\\r\\r不動産賃貸事業、不動産証券化事業を展開する株式会社マリオン(本社:東京都新宿区、代表取締役社長福田敬司)は、2022年7月25日(月)、お金第3の置き場「i-Bond」のサービスサイトをリニューアルいたしました。\\r\\r\\r"),
			want1: map[string][]string{
				"p_no":    {"4864"},
				"p_date":  {"2022.07.25-17:49:57.424"},
				"p_errno": {"0"},
				"p_err":   {""},
				"p_cmd":   {"NS"},
				"p_PV":    {"QNSD"},
				"p_ENO":   {"95371"},
				"p_ALT":   {"1"},
				"p_ID":    {"20220725174500_NOV6627"},
				"p_DT":    {"20220725"},
				"p_TM":    {"174500"},
				"p_CGN":   {"1"},
				"p_CGL":   {"129"},
				"p_GRN":   {"1"},
				"p_GRL":   {"62199"},
				"p_ISN":   {"1"},
				"p_ISL":   {"3494"},
				"p_SKF":   {"1"},
				"p_UPF":   {""},
				"p_HLD":   {"<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ"},
				"p_TX":    {`<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\r\r\r\r開示会社：マリオン(3494)\r\r開示書類：i-Bond(アイボンド）サイトリニューアルのお知らせ\r\r開示日時：2022/07/25 17:45\r\r\r\r＜引用＞\r\r不動産賃貸事業、不動産証券化事業を展開する株式会社マリオン(本社:東京都新宿区、代表取締役社長福田敬司)は、2022年7月25日(月)、お金第3の置き場「i-Bond」のサービスサイトをリニューアルいたしました。\r\r\r`},
			}},
		{name: "NSのbodyをmapに変換できる",
			arg1: []byte("p_no\x023\x01p_date\x022022.07.26-20:04:48.809\x01p_errno\x020\x01p_err\x02\x01p_cmd\x02FD\x01p_5_AV\x02196970\x01p_5_BV\x02140000\x01p_5_DHF\x020000\x01p_5_DHP\x02378.6\x01p_5_DHP:T\x0214:59\x01p_5_DJ\x021168368891\x01p_5_DLF\x020000\x01p_5_DLP\x02368.5\x01p_5_DLP:T\x0209:50\x01p_5_DOP\x02369.7\x01p_5_DOP:T\x0209:03\x01p_5_DPG\x020057\x01p_5_DPP\x02378.6\x01p_5_DPP:T\x0215:00\x01p_5_DV\x023120080\x01p_5_DYRP\x024.87\x01p_5_DYWP\x0217.6\x01p_5_GAP1\x02378.6\x01p_5_GAP10\x02379.5\x01p_5_GAP2\x02378.7\x01p_5_GAP3\x02378.8\x01p_5_GAP4\x02378.9\x01p_5_GAP5\x02379.0\x01p_5_GAP6\x02379.1\x01p_5_GAP7\x02379.2\x01p_5_GAP8\x02379.3\x01p_5_GAP9\x02379.4\x01p_5_GAV1\x02196970\x01p_5_GAV10\x027510\x01p_5_GAV2\x02192610\x01p_5_GAV3\x02162960\x01p_5_GAV4\x0262870\x01p_5_GAV5\x0283130\x01p_5_GAV6\x0221220\x01p_5_GAV7\x0220\x01p_5_GAV8\x02110\x01p_5_GAV9\x0210\x01p_5_GBP1\x02378.3\x01p_5_GBP10\x02377.4\x01p_5_GBP2\x02378.2\x01p_5_GBP3\x02378.1\x01p_5_GBP4\x02378.0\x01p_5_GBP5\x02377.9\x01p_5_GBP6\x02377.8\x01p_5_GBP7\x02377.7\x01p_5_GBP8\x02377.6\x01p_5_GBP9\x02377.5\x01p_5_GBV1\x02140000\x01p_5_GBV10\x027010\x01p_5_GBV2\x02208910\x01p_5_GBV3\x02172860\x01p_5_GBV4\x0262900\x01p_5_GBV5\x0263300\x01p_5_GBV6\x0277910\x01p_5_GBV7\x0277920\x01p_5_GBV8\x0269940\x01p_5_GBV9\x027000\x01p_5_PRP\x02361.0\x01p_5_QAP\x02378.6\x01p_5_QAS\x020101\x01p_5_QBP\x02378.3\x01p_5_QBS\x020101\x01p_5_QOV\x021020310\x01p_5_QUV\x02781260\x01p_5_VWAP\x02374.4676"),
			want1: map[string][]string{
				"p_no":      {"3"},
				"p_date":    {"2022.07.26-20:04:48.809"},
				"p_errno":   {"0"},
				"p_err":     {""},
				"p_cmd":     {"FD"},
				"p_5_AV":    {"196970"},
				"p_5_BV":    {"140000"},
				"p_5_DHF":   {"0000"},
				"p_5_DHP":   {"378.6"},
				"p_5_DHP:T": {"14:59"},
				"p_5_DJ":    {"1168368891"},
				"p_5_DLF":   {"0000"},
				"p_5_DLP":   {"368.5"},
				"p_5_DLP:T": {"09:50"},
				"p_5_DOP":   {"369.7"},
				"p_5_DOP:T": {"09:03"},
				"p_5_DPG":   {"0057"},
				"p_5_DPP":   {"378.6"},
				"p_5_DPP:T": {"15:00"},
				"p_5_DV":    {"3120080"},
				"p_5_DYRP":  {"4.87"},
				"p_5_DYWP":  {"17.6"},
				"p_5_GAP1":  {"378.6"},
				"p_5_GAP10": {"379.5"},
				"p_5_GAP2":  {"378.7"},
				"p_5_GAP3":  {"378.8"},
				"p_5_GAP4":  {"378.9"},
				"p_5_GAP5":  {"379.0"},
				"p_5_GAP6":  {"379.1"},
				"p_5_GAP7":  {"379.2"},
				"p_5_GAP8":  {"379.3"},
				"p_5_GAP9":  {"379.4"},
				"p_5_GAV1":  {"196970"},
				"p_5_GAV10": {"7510"},
				"p_5_GAV2":  {"192610"},
				"p_5_GAV3":  {"162960"},
				"p_5_GAV4":  {"62870"},
				"p_5_GAV5":  {"83130"},
				"p_5_GAV6":  {"21220"},
				"p_5_GAV7":  {"20"},
				"p_5_GAV8":  {"110"},
				"p_5_GAV9":  {"10"},
				"p_5_GBP1":  {"378.3"},
				"p_5_GBP10": {"377.4"},
				"p_5_GBP2":  {"378.2"},
				"p_5_GBP3":  {"378.1"},
				"p_5_GBP4":  {"378.0"},
				"p_5_GBP5":  {"377.9"},
				"p_5_GBP6":  {"377.8"},
				"p_5_GBP7":  {"377.7"},
				"p_5_GBP8":  {"377.6"},
				"p_5_GBP9":  {"377.5"},
				"p_5_GBV1":  {"140000"},
				"p_5_GBV10": {"7010"},
				"p_5_GBV2":  {"208910"},
				"p_5_GBV3":  {"172860"},
				"p_5_GBV4":  {"62900"},
				"p_5_GBV5":  {"63300"},
				"p_5_GBV6":  {"77910"},
				"p_5_GBV7":  {"77920"},
				"p_5_GBV8":  {"69940"},
				"p_5_GBV9":  {"7000"},
				"p_5_PRP":   {"361.0"},
				"p_5_QAP":   {"378.6"},
				"p_5_QAS":   {"0101"},
				"p_5_QBP":   {"378.3"},
				"p_5_QBS":   {"0101"},
				"p_5_QOV":   {"1020310"},
				"p_5_QUV":   {"781260"},
				"p_5_VWAP":  {"374.4676"},
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

func Test_NewsStreamResponse_parse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  map[string][]string
		arg2  []byte
		want1 NewsStreamResponse
	}{
		{name: "マップの値を使って構造体に反映できる",
			arg1: map[string][]string{
				"p_no":    {"4864"},
				"p_date":  {"2022.07.25-17:49:57.424"},
				"p_errno": {"0"},
				"p_err":   {""},
				"p_cmd":   {"NS"},
				"p_PV":    {"QNSD"},
				"p_ENO":   {"95371"},
				"p_ALT":   {"1"},
				"p_ID":    {"20220725174500_NOV6627"},
				"p_DT":    {"20220725"},
				"p_TM":    {"174500"},
				"p_CGN":   {"1"},
				"p_CGL":   {"129"},
				"p_GRN":   {"1"},
				"p_GRL":   {"62199"},
				"p_ISN":   {"1"},
				"p_ISL":   {"3494"},
				"p_SKF":   {"1"},
				"p_UPF":   {""},
				"p_HLD":   {"<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ"},
				"p_TX":    {`<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\r\r\r\r開示会社：マリオン(3494)\r\r開示書類：i-Bond(アイボンド）サイトリニューアルのお知らせ\r\r開示日時：2022/07/25 17:45\r\r\r\r＜引用＞\r\r不動産賃貸事業、不動産証券化事業を展開する株式会社マリオン(本社:東京都新宿区、代表取締役社長福田敬司)は、2022年7月25日(月)、お金第3の置き場「i-Bond」のサービスサイトをリニューアルいたしました。\r\r\r`},
			},
			arg2: []byte("p_no\x024864\x01p_date\x022022.07.25-17:49:57.424\x01p_errno\x020\x01p_err\x02\x01p_cmd\x02NS\x01p_PV\x02QNSD\x01p_ENO\x0295371\x01p_ALT\x021\x01p_ID\x0220220725174500_NOV6627\x01p_DT\x0220220725\x01p_TM\x02174500\x01p_CGN\x021\x01p_CGL\x02129\x01p_GRN\x021\x01p_GRL\x0262199\x01p_ISN\x021\x01p_ISL\x023494\x01p_SKF\x021\x01p_UPF\x02\x01p_HLD\x02<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\x01p_TX\x02<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\\r\\r\\r\\r開示会社：マリオン(3494)\\r\\r開示書類：i-Bond(アイボンド）サイトリニューアルのお知らせ\\r\\r開示日時：2022/07/25 17:45\\r\\r\\r\\r＜引用＞\\r\\r不動産賃貸事業、不動産証券化事業を展開する株式会社マリオン(本社:東京都新宿区、代表取締役社長福田敬司)は、2022年7月25日(月)、お金第3の置き場「i-Bond」のサービスサイトをリニューアルいたしました。\\r\\r\\r"),
			want1: NewsStreamResponse{
				CommonStreamResponse: CommonStreamResponse{
					EventType:      EventTypeNews,
					StreamNumber:   4864,
					StreamDateTime: time.Date(2022, 7, 25, 17, 49, 57, 424000000, time.Local),
					ErrorNo:        ErrorNoProblem,
					ErrorText:      "",
					Body:           []byte("p_no\x024864\x01p_date\x022022.07.25-17:49:57.424\x01p_errno\x020\x01p_err\x02\x01p_cmd\x02NS\x01p_PV\x02QNSD\x01p_ENO\x0295371\x01p_ALT\x021\x01p_ID\x0220220725174500_NOV6627\x01p_DT\x0220220725\x01p_TM\x02174500\x01p_CGN\x021\x01p_CGL\x02129\x01p_GRN\x021\x01p_GRL\x0262199\x01p_ISN\x021\x01p_ISL\x023494\x01p_SKF\x021\x01p_UPF\x02\x01p_HLD\x02<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\x01p_TX\x02<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\\r\\r\\r\\r開示会社：マリオン(3494)\\r\\r開示書類：i-Bond(アイボンド）サイトリニューアルのお知らせ\\r\\r開示日時：2022/07/25 17:45\\r\\r\\r\\r＜引用＞\\r\\r不動産賃貸事業、不動産証券化事業を展開する株式会社マリオン(本社:東京都新宿区、代表取締役社長福田敬司)は、2022年7月25日(月)、お金第3の置き場「i-Bond」のサービスサイトをリニューアルいたしました。\\r\\r\\r"),
				},
				Provider:      "QNSD",
				EventNo:       95371,
				FirstTime:     true,
				NewsId:        "20220725174500_NOV6627",
				NewsDateTime:  time.Date(2022, 7, 25, 17, 45, 0, 0, time.Local),
				NumOfCategory: 1,
				Categories:    []string{"129"},
				NumOfGenre:    1,
				Genres:        []string{"62199"},
				NumOfIssue:    1,
				Issues:        []string{"3494"},
				Title:         "<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ",
				Content:       `<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\r\r\r\r開示会社：マリオン(3494)\r\r開示書類：i-Bond(アイボンド）サイトリニューアルのお知らせ\r\r開示日時：2022/07/25 17:45\r\r\r\r＜引用＞\r\r不動産賃貸事業、不動産証券化事業を展開する株式会社マリオン(本社:東京都新宿区、代表取締役社長福田敬司)は、2022年7月25日(月)、お金第3の置き場「i-Bond」のサービスサイトをリニューアルいたしました。\r\r\r`,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var res NewsStreamResponse
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

func Test_MarketPriceStreamResponse_getColumnNumber(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  map[string][]string
		want1 int
	}{
		{name: "空っぽのmapなら0が返される",
			arg1:  map[string][]string{},
			want1: 0},
		{name: "抽出できるキーがなければ0が返される",
			arg1: map[string][]string{
				"p_no":    {"6"},
				"p_date":  {"2022.07.26-20:03:17.899"},
				"p_errno": {"0"},
				"p_err":   {""},
				"p_cmd":   {"FD"},
			},
			want1: 0},
		{name: "抽出できるなら数値を返す(1)",
			arg1: map[string][]string{
				"p_no":    {"6"},
				"p_date":  {"2022.07.26-20:03:17.899"},
				"p_errno": {"0"},
				"p_err":   {""},
				"p_cmd":   {"FD"},
				"p_1_AV":  {"57974"},
			},
			want1: 1},
		{name: "抽出できるなら数値を返す(120)",
			arg1: map[string][]string{
				"p_no":     {"6"},
				"p_date":   {"2022.07.26-20:03:17.899"},
				"p_errno":  {"0"},
				"p_err":    {""},
				"p_cmd":    {"FD"},
				"p_120_AV": {"57974"},
			},
			want1: 120},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			res := &MarketPriceStreamResponse{}
			got1 := res.getColumnNumber(test.arg1)
			if !reflect.DeepEqual(test.want1, got1) {
				t.Errorf("%s error\nwant: %+v\ngot: %+v\n", t.Name(), test.want1, got1)
			}
		})
	}
}

func Test_MarketPriceStreamResponse_parse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		arg1  map[string][]string
		arg2  []byte
		want1 MarketPriceStreamResponse
	}{
		{name: "マップの値を使って構造体に反映できる",
			arg1: map[string][]string{
				"p_no":      {"3"},
				"p_date":    {"2022.07.26-20:04:48.809"},
				"p_errno":   {"0"},
				"p_err":     {""},
				"p_cmd":     {"FD"},
				"p_5_AV":    {"196970"},
				"p_5_BV":    {"140000"},
				"p_5_DHF":   {"0000"},
				"p_5_DHP":   {"378.6"},
				"p_5_DHP:T": {"14:59"},
				"p_5_DJ":    {"1168368891"},
				"p_5_DLF":   {"0000"},
				"p_5_DLP":   {"368.5"},
				"p_5_DLP:T": {"09:50"},
				"p_5_DOP":   {"369.7"},
				"p_5_DOP:T": {"09:03"},
				"p_5_DPG":   {"0057"},
				"p_5_DPP":   {"378.6"},
				"p_5_DPP:T": {"15:00"},
				"p_5_DV":    {"3120080"},
				"p_5_DYRP":  {"4.87"},
				"p_5_DYWP":  {"17.6"},
				"p_5_GAP1":  {"378.6"},
				"p_5_GAP10": {"379.5"},
				"p_5_GAP2":  {"378.7"},
				"p_5_GAP3":  {"378.8"},
				"p_5_GAP4":  {"378.9"},
				"p_5_GAP5":  {"379.0"},
				"p_5_GAP6":  {"379.1"},
				"p_5_GAP7":  {"379.2"},
				"p_5_GAP8":  {"379.3"},
				"p_5_GAP9":  {"379.4"},
				"p_5_GAV1":  {"196970"},
				"p_5_GAV10": {"7510"},
				"p_5_GAV2":  {"192610"},
				"p_5_GAV3":  {"162960"},
				"p_5_GAV4":  {"62870"},
				"p_5_GAV5":  {"83130"},
				"p_5_GAV6":  {"21220"},
				"p_5_GAV7":  {"20"},
				"p_5_GAV8":  {"110"},
				"p_5_GAV9":  {"10"},
				"p_5_GBP1":  {"378.3"},
				"p_5_GBP10": {"377.4"},
				"p_5_GBP2":  {"378.2"},
				"p_5_GBP3":  {"378.1"},
				"p_5_GBP4":  {"378.0"},
				"p_5_GBP5":  {"377.9"},
				"p_5_GBP6":  {"377.8"},
				"p_5_GBP7":  {"377.7"},
				"p_5_GBP8":  {"377.6"},
				"p_5_GBP9":  {"377.5"},
				"p_5_GBV1":  {"140000"},
				"p_5_GBV10": {"7010"},
				"p_5_GBV2":  {"208910"},
				"p_5_GBV3":  {"172860"},
				"p_5_GBV4":  {"62900"},
				"p_5_GBV5":  {"63300"},
				"p_5_GBV6":  {"77910"},
				"p_5_GBV7":  {"77920"},
				"p_5_GBV8":  {"69940"},
				"p_5_GBV9":  {"7000"},
				"p_5_PRP":   {"361.0"},
				"p_5_QAP":   {"378.6"},
				"p_5_QAS":   {"0101"},
				"p_5_QBP":   {"378.3"},
				"p_5_QBS":   {"0101"},
				"p_5_QOV":   {"1020310"},
				"p_5_QUV":   {"781260"},
				"p_5_VWAP":  {"374.4676"},
			},
			arg2: []byte("p_no\x023\x01p_date\x022022.07.26-20:04:48.809\x01p_errno\x020\x01p_err\x02\x01p_cmd\x02FD\x01p_5_AV\x02196970\x01p_5_BV\x02140000\x01p_5_DHF\x020000\x01p_5_DHP\x02378.6\x01p_5_DHP:T\x0214:59\x01p_5_DJ\x021168368891\x01p_5_DLF\x020000\x01p_5_DLP\x02368.5\x01p_5_DLP:T\x0209:50\x01p_5_DOP\x02369.7\x01p_5_DOP:T\x0209:03\x01p_5_DPG\x020057\x01p_5_DPP\x02378.6\x01p_5_DPP:T\x0215:00\x01p_5_DV\x023120080\x01p_5_DYRP\x024.87\x01p_5_DYWP\x0217.6\x01p_5_GAP1\x02378.6\x01p_5_GAP10\x02379.5\x01p_5_GAP2\x02378.7\x01p_5_GAP3\x02378.8\x01p_5_GAP4\x02378.9\x01p_5_GAP5\x02379.0\x01p_5_GAP6\x02379.1\x01p_5_GAP7\x02379.2\x01p_5_GAP8\x02379.3\x01p_5_GAP9\x02379.4\x01p_5_GAV1\x02196970\x01p_5_GAV10\x027510\x01p_5_GAV2\x02192610\x01p_5_GAV3\x02162960\x01p_5_GAV4\x0262870\x01p_5_GAV5\x0283130\x01p_5_GAV6\x0221220\x01p_5_GAV7\x0220\x01p_5_GAV8\x02110\x01p_5_GAV9\x0210\x01p_5_GBP1\x02378.3\x01p_5_GBP10\x02377.4\x01p_5_GBP2\x02378.2\x01p_5_GBP3\x02378.1\x01p_5_GBP4\x02378.0\x01p_5_GBP5\x02377.9\x01p_5_GBP6\x02377.8\x01p_5_GBP7\x02377.7\x01p_5_GBP8\x02377.6\x01p_5_GBP9\x02377.5\x01p_5_GBV1\x02140000\x01p_5_GBV10\x027010\x01p_5_GBV2\x02208910\x01p_5_GBV3\x02172860\x01p_5_GBV4\x0262900\x01p_5_GBV5\x0263300\x01p_5_GBV6\x0277910\x01p_5_GBV7\x0277920\x01p_5_GBV8\x0269940\x01p_5_GBV9\x027000\x01p_5_PRP\x02361.0\x01p_5_QAP\x02378.6\x01p_5_QAS\x020101\x01p_5_QBP\x02378.3\x01p_5_QBS\x020101\x01p_5_QOV\x021020310\x01p_5_QUV\x02781260\x01p_5_VWAP\x02374.4676"),
			want1: MarketPriceStreamResponse{
				CommonStreamResponse: CommonStreamResponse{
					EventType:      EventTypeMarketPrice,
					StreamNumber:   3,
					StreamDateTime: time.Date(2022, 7, 26, 20, 4, 48, 809000000, time.Local),
					ErrorNo:        ErrorNoProblem,
					ErrorText:      "",
					Body:           []byte("p_no\x023\x01p_date\x022022.07.26-20:04:48.809\x01p_errno\x020\x01p_err\x02\x01p_cmd\x02FD\x01p_5_AV\x02196970\x01p_5_BV\x02140000\x01p_5_DHF\x020000\x01p_5_DHP\x02378.6\x01p_5_DHP:T\x0214:59\x01p_5_DJ\x021168368891\x01p_5_DLF\x020000\x01p_5_DLP\x02368.5\x01p_5_DLP:T\x0209:50\x01p_5_DOP\x02369.7\x01p_5_DOP:T\x0209:03\x01p_5_DPG\x020057\x01p_5_DPP\x02378.6\x01p_5_DPP:T\x0215:00\x01p_5_DV\x023120080\x01p_5_DYRP\x024.87\x01p_5_DYWP\x0217.6\x01p_5_GAP1\x02378.6\x01p_5_GAP10\x02379.5\x01p_5_GAP2\x02378.7\x01p_5_GAP3\x02378.8\x01p_5_GAP4\x02378.9\x01p_5_GAP5\x02379.0\x01p_5_GAP6\x02379.1\x01p_5_GAP7\x02379.2\x01p_5_GAP8\x02379.3\x01p_5_GAP9\x02379.4\x01p_5_GAV1\x02196970\x01p_5_GAV10\x027510\x01p_5_GAV2\x02192610\x01p_5_GAV3\x02162960\x01p_5_GAV4\x0262870\x01p_5_GAV5\x0283130\x01p_5_GAV6\x0221220\x01p_5_GAV7\x0220\x01p_5_GAV8\x02110\x01p_5_GAV9\x0210\x01p_5_GBP1\x02378.3\x01p_5_GBP10\x02377.4\x01p_5_GBP2\x02378.2\x01p_5_GBP3\x02378.1\x01p_5_GBP4\x02378.0\x01p_5_GBP5\x02377.9\x01p_5_GBP6\x02377.8\x01p_5_GBP7\x02377.7\x01p_5_GBP8\x02377.6\x01p_5_GBP9\x02377.5\x01p_5_GBV1\x02140000\x01p_5_GBV10\x027010\x01p_5_GBV2\x02208910\x01p_5_GBV3\x02172860\x01p_5_GBV4\x0262900\x01p_5_GBV5\x0263300\x01p_5_GBV6\x0277910\x01p_5_GBV7\x0277920\x01p_5_GBV8\x0269940\x01p_5_GBV9\x027000\x01p_5_PRP\x02361.0\x01p_5_QAP\x02378.6\x01p_5_QAS\x020101\x01p_5_QBP\x02378.3\x01p_5_QBS\x020101\x01p_5_QOV\x021020310\x01p_5_QUV\x02781260\x01p_5_VWAP\x02374.4676"),
				},
				ColumnNumber:      5,
				AskQuantityMarket: 0,
				BidQuantityMarket: 0,
				AskQuantity:       196970,
				BidQuantity:       140000,
				DiscontinuityType: "",
				StopHigh:          CurrentPriceTypeNoChange,
				HighPrice:         378.6,
				HighPriceTime:     time.Date(0, 1, 1, 14, 59, 0, 0, time.Local),
				TradingAmount:     1168368891,
				StopLow:           CurrentPriceTypeNoChange,
				LowPrice:          368.5,
				LowPriceTime:      time.Date(0, 1, 1, 9, 50, 0, 0, time.Local),
				OpenPrice:         369.7,
				OpenPriceTime:     time.Date(0, 1, 1, 9, 3, 0, 0, time.Local),
				ChangePriceType:   ChangePriceTypeRise,
				CurrentPrice:      378.6,
				CurrentPriceTime:  time.Date(0, 1, 1, 15, 0, 0, 0, time.Local),
				Volume:            3120080,
				ExRightType:       "",
				PrevDayPercent:    4.87,
				PrevDayRatio:      17.6,
				AskQuantity10:     7510,
				AskPrice10:        379.5,
				AskQuantity9:      10,
				AskPrice9:         379.4,
				AskQuantity8:      110,
				AskPrice8:         379.3,
				AskQuantity7:      20,
				AskPrice7:         379.2,
				AskQuantity6:      21220,
				AskPrice6:         379.1,
				AskQuantity5:      83130,
				AskPrice5:         379.0,
				AskQuantity4:      62870,
				AskPrice4:         378.9,
				AskQuantity3:      162960,
				AskPrice3:         378.8,
				AskQuantity2:      192610,
				AskPrice2:         378.7,
				AskQuantity1:      196970,
				AskPrice1:         378.6,
				BidQuantity1:      140000,
				BidPrice1:         378.3,
				BidQuantity2:      208910,
				BidPrice2:         378.2,
				BidQuantity3:      172860,
				BidPrice3:         378.1,
				BidQuantity4:      62900,
				BidPrice4:         378.0,
				BidQuantity5:      63300,
				BidPrice5:         377.9,
				BidQuantity6:      77910,
				BidPrice6:         377.8,
				BidQuantity7:      77920,
				BidPrice7:         377.7,
				BidQuantity8:      69940,
				BidPrice8:         377.6,
				BidQuantity9:      7000,
				BidPrice9:         377.5,
				BidQuantity10:     7010,
				BidPrice10:        377.4,
				Section:           "",
				PRP:               361.0,
				AskPrice:          378.6,
				AskSign:           IndicationPriceTypeGeneral,
				BidPrice:          378.3,
				BidSign:           IndicationPriceTypeGeneral,
				AskQuantityOver:   1020310,
				BidQuantityUnder:  781260,
				VWAP:              374.4676,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var res MarketPriceStreamResponse
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
				stream1 <- []byte("p_no\x024864\x01p_date\x022022.07.25-17:49:57.424\x01p_errno\x020\x01p_err\x02\x01p_cmd\x02NS\x01p_PV\x02QNSD\x01p_ENO\x0295371\x01p_ALT\x021\x01p_ID\x0220220725174500_NOV6627\x01p_DT\x0220220725\x01p_TM\x02174500\x01p_CGN\x021\x01p_CGL\x02129\x01p_GRN\x021\x01p_GRL\x0262199\x01p_ISN\x021\x01p_ISL\x023494\x01p_SKF\x021\x01p_UPF\x02\x01p_HLD\x02<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\x01p_TX\x02<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\\r\\r\\r\\r開示会社：マリオン(3494)\\r\\r開示書類：i-Bond(アイボンド）サイトリニューアルのお知らせ\\r\\r開示日時：2022/07/25 17:45\\r\\r\\r\\r＜引用＞\\r\\r不動産賃貸事業、不動産証券化事業を展開する株式会社マリオン(本社:東京都新宿区、代表取締役社長福田敬司)は、2022年7月25日(月)、お金第3の置き場「i-Bond」のサービスサイトをリニューアルいたしました。\\r\\r\\r")
				stream1 <- []byte("p_no\x023\x01p_date\x022022.07.26-20:04:48.809\x01p_errno\x020\x01p_err\x02\x01p_cmd\x02FD\x01p_5_AV\x02196970\x01p_5_BV\x02140000\x01p_5_DHF\x020000\x01p_5_DHP\x02378.6\x01p_5_DHP:T\x0214:59\x01p_5_DJ\x021168368891\x01p_5_DLF\x020000\x01p_5_DLP\x02368.5\x01p_5_DLP:T\x0209:50\x01p_5_DOP\x02369.7\x01p_5_DOP:T\x0209:03\x01p_5_DPG\x020057\x01p_5_DPP\x02378.6\x01p_5_DPP:T\x0215:00\x01p_5_DV\x023120080\x01p_5_DYRP\x024.87\x01p_5_DYWP\x0217.6\x01p_5_GAP1\x02378.6\x01p_5_GAP10\x02379.5\x01p_5_GAP2\x02378.7\x01p_5_GAP3\x02378.8\x01p_5_GAP4\x02378.9\x01p_5_GAP5\x02379.0\x01p_5_GAP6\x02379.1\x01p_5_GAP7\x02379.2\x01p_5_GAP8\x02379.3\x01p_5_GAP9\x02379.4\x01p_5_GAV1\x02196970\x01p_5_GAV10\x027510\x01p_5_GAV2\x02192610\x01p_5_GAV3\x02162960\x01p_5_GAV4\x0262870\x01p_5_GAV5\x0283130\x01p_5_GAV6\x0221220\x01p_5_GAV7\x0220\x01p_5_GAV8\x02110\x01p_5_GAV9\x0210\x01p_5_GBP1\x02378.3\x01p_5_GBP10\x02377.4\x01p_5_GBP2\x02378.2\x01p_5_GBP3\x02378.1\x01p_5_GBP4\x02378.0\x01p_5_GBP5\x02377.9\x01p_5_GBP6\x02377.8\x01p_5_GBP7\x02377.7\x01p_5_GBP8\x02377.6\x01p_5_GBP9\x02377.5\x01p_5_GBV1\x02140000\x01p_5_GBV10\x027010\x01p_5_GBV2\x02208910\x01p_5_GBV3\x02172860\x01p_5_GBV4\x0262900\x01p_5_GBV5\x0263300\x01p_5_GBV6\x0277910\x01p_5_GBV7\x0277920\x01p_5_GBV8\x0269940\x01p_5_GBV9\x027000\x01p_5_PRP\x02361.0\x01p_5_QAP\x02378.6\x01p_5_QAS\x020101\x01p_5_QBP\x02378.3\x01p_5_QBS\x020101\x01p_5_QOV\x021020310\x01p_5_QUV\x02781260\x01p_5_VWAP\x02374.4676")
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
				&NewsStreamResponse{
					CommonStreamResponse: CommonStreamResponse{
						EventType:      EventTypeNews,
						StreamNumber:   4864,
						StreamDateTime: time.Date(2022, 7, 25, 17, 49, 57, 424000000, time.Local),
						ErrorNo:        ErrorNoProblem,
						ErrorText:      "",
						Body:           []byte("p_no\x024864\x01p_date\x022022.07.25-17:49:57.424\x01p_errno\x020\x01p_err\x02\x01p_cmd\x02NS\x01p_PV\x02QNSD\x01p_ENO\x0295371\x01p_ALT\x021\x01p_ID\x0220220725174500_NOV6627\x01p_DT\x0220220725\x01p_TM\x02174500\x01p_CGN\x021\x01p_CGL\x02129\x01p_GRN\x021\x01p_GRL\x0262199\x01p_ISN\x021\x01p_ISL\x023494\x01p_SKF\x021\x01p_UPF\x02\x01p_HLD\x02<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\x01p_TX\x02<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\\r\\r\\r\\r開示会社：マリオン(3494)\\r\\r開示書類：i-Bond(アイボンド）サイトリニューアルのお知らせ\\r\\r開示日時：2022/07/25 17:45\\r\\r\\r\\r＜引用＞\\r\\r不動産賃貸事業、不動産証券化事業を展開する株式会社マリオン(本社:東京都新宿区、代表取締役社長福田敬司)は、2022年7月25日(月)、お金第3の置き場「i-Bond」のサービスサイトをリニューアルいたしました。\\r\\r\\r"),
					},
					Provider:      "QNSD",
					EventNo:       95371,
					FirstTime:     true,
					NewsId:        "20220725174500_NOV6627",
					NewsDateTime:  time.Date(2022, 7, 25, 17, 45, 0, 0, time.Local),
					NumOfCategory: 1,
					Categories:    []string{"129"},
					NumOfGenre:    1,
					Genres:        []string{"62199"},
					NumOfIssue:    1,
					Issues:        []string{"3494"},
					Title:         "<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ",
					Content:       `<TDnet>AI: マリオン(3494) i-Bond(アイボンド）サイトリニューアルのお知らせ\r\r\r\r開示会社：マリオン(3494)\r\r開示書類：i-Bond(アイボンド）サイトリニューアルのお知らせ\r\r開示日時：2022/07/25 17:45\r\r\r\r＜引用＞\r\r不動産賃貸事業、不動産証券化事業を展開する株式会社マリオン(本社:東京都新宿区、代表取締役社長福田敬司)は、2022年7月25日(月)、お金第3の置き場「i-Bond」のサービスサイトをリニューアルいたしました。\r\r\r`,
				},
				&MarketPriceStreamResponse{
					CommonStreamResponse: CommonStreamResponse{
						EventType:      EventTypeMarketPrice,
						StreamNumber:   3,
						StreamDateTime: time.Date(2022, 7, 26, 20, 4, 48, 809000000, time.Local),
						ErrorNo:        ErrorNoProblem,
						ErrorText:      "",
						Body:           []byte("p_no\x023\x01p_date\x022022.07.26-20:04:48.809\x01p_errno\x020\x01p_err\x02\x01p_cmd\x02FD\x01p_5_AV\x02196970\x01p_5_BV\x02140000\x01p_5_DHF\x020000\x01p_5_DHP\x02378.6\x01p_5_DHP:T\x0214:59\x01p_5_DJ\x021168368891\x01p_5_DLF\x020000\x01p_5_DLP\x02368.5\x01p_5_DLP:T\x0209:50\x01p_5_DOP\x02369.7\x01p_5_DOP:T\x0209:03\x01p_5_DPG\x020057\x01p_5_DPP\x02378.6\x01p_5_DPP:T\x0215:00\x01p_5_DV\x023120080\x01p_5_DYRP\x024.87\x01p_5_DYWP\x0217.6\x01p_5_GAP1\x02378.6\x01p_5_GAP10\x02379.5\x01p_5_GAP2\x02378.7\x01p_5_GAP3\x02378.8\x01p_5_GAP4\x02378.9\x01p_5_GAP5\x02379.0\x01p_5_GAP6\x02379.1\x01p_5_GAP7\x02379.2\x01p_5_GAP8\x02379.3\x01p_5_GAP9\x02379.4\x01p_5_GAV1\x02196970\x01p_5_GAV10\x027510\x01p_5_GAV2\x02192610\x01p_5_GAV3\x02162960\x01p_5_GAV4\x0262870\x01p_5_GAV5\x0283130\x01p_5_GAV6\x0221220\x01p_5_GAV7\x0220\x01p_5_GAV8\x02110\x01p_5_GAV9\x0210\x01p_5_GBP1\x02378.3\x01p_5_GBP10\x02377.4\x01p_5_GBP2\x02378.2\x01p_5_GBP3\x02378.1\x01p_5_GBP4\x02378.0\x01p_5_GBP5\x02377.9\x01p_5_GBP6\x02377.8\x01p_5_GBP7\x02377.7\x01p_5_GBP8\x02377.6\x01p_5_GBP9\x02377.5\x01p_5_GBV1\x02140000\x01p_5_GBV10\x027010\x01p_5_GBV2\x02208910\x01p_5_GBV3\x02172860\x01p_5_GBV4\x0262900\x01p_5_GBV5\x0263300\x01p_5_GBV6\x0277910\x01p_5_GBV7\x0277920\x01p_5_GBV8\x0269940\x01p_5_GBV9\x027000\x01p_5_PRP\x02361.0\x01p_5_QAP\x02378.6\x01p_5_QAS\x020101\x01p_5_QBP\x02378.3\x01p_5_QBS\x020101\x01p_5_QOV\x021020310\x01p_5_QUV\x02781260\x01p_5_VWAP\x02374.4676"),
					},
					ColumnNumber:      5,
					AskQuantityMarket: 0,
					BidQuantityMarket: 0,
					AskQuantity:       196970,
					BidQuantity:       140000,
					DiscontinuityType: "",
					StopHigh:          CurrentPriceTypeNoChange,
					HighPrice:         378.6,
					HighPriceTime:     time.Date(0, 1, 1, 14, 59, 0, 0, time.Local),
					TradingAmount:     1168368891,
					StopLow:           CurrentPriceTypeNoChange,
					LowPrice:          368.5,
					LowPriceTime:      time.Date(0, 1, 1, 9, 50, 0, 0, time.Local),
					OpenPrice:         369.7,
					OpenPriceTime:     time.Date(0, 1, 1, 9, 3, 0, 0, time.Local),
					ChangePriceType:   ChangePriceTypeRise,
					CurrentPrice:      378.6,
					CurrentPriceTime:  time.Date(0, 1, 1, 15, 0, 0, 0, time.Local),
					Volume:            3120080,
					ExRightType:       "",
					PrevDayPercent:    4.87,
					PrevDayRatio:      17.6,
					AskQuantity10:     7510,
					AskPrice10:        379.5,
					AskQuantity9:      10,
					AskPrice9:         379.4,
					AskQuantity8:      110,
					AskPrice8:         379.3,
					AskQuantity7:      20,
					AskPrice7:         379.2,
					AskQuantity6:      21220,
					AskPrice6:         379.1,
					AskQuantity5:      83130,
					AskPrice5:         379.0,
					AskQuantity4:      62870,
					AskPrice4:         378.9,
					AskQuantity3:      162960,
					AskPrice3:         378.8,
					AskQuantity2:      192610,
					AskPrice2:         378.7,
					AskQuantity1:      196970,
					AskPrice1:         378.6,
					BidQuantity1:      140000,
					BidPrice1:         378.3,
					BidQuantity2:      208910,
					BidPrice2:         378.2,
					BidQuantity3:      172860,
					BidPrice3:         378.1,
					BidQuantity4:      62900,
					BidPrice4:         378.0,
					BidQuantity5:      63300,
					BidPrice5:         377.9,
					BidQuantity6:      77910,
					BidPrice6:         377.8,
					BidQuantity7:      77920,
					BidPrice7:         377.7,
					BidQuantity8:      69940,
					BidPrice8:         377.6,
					BidQuantity9:      7000,
					BidPrice9:         377.5,
					BidQuantity10:     7010,
					BidPrice10:        377.4,
					Section:           "",
					PRP:               361.0,
					AskPrice:          378.6,
					AskSign:           IndicationPriceTypeGeneral,
					BidPrice:          378.3,
					BidSign:           IndicationPriceTypeGeneral,
					AskQuantityOver:   1020310,
					BidQuantityUnder:  781260,
					VWAP:              374.4676,
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
