package tachibana

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type StreamRequest struct {
	ColumnNumber      []int       // 株価ボード専用 行番号
	IssueCodes        []string    // 株価ボード専用 銘柄コード
	MarketCodes       []Exchange  // 株価ボード専用 市場コード
	StartStreamNumber int64       // 配信開始イベント通知番号
	StreamEventTypes  []EventType // 通知種別
}

func (r *StreamRequest) Query() []byte {
	var res []string

	res = append(res, "p_rid=22")
	res = append(res, "p_board_no=1000")

	if len(r.ColumnNumber) > 0 {
		columns := make([]string, len(r.ColumnNumber))
		for i, n := range r.ColumnNumber {
			columns[i] = strconv.Itoa(n)
		}
		res = append(res, "p_gyou_no="+strings.Join(columns, ","))
	}

	if len(r.IssueCodes) > 0 {
		res = append(res, "p_issue_code="+strings.Join(r.IssueCodes, ","))
	}

	if len(r.MarketCodes) > 0 {
		exchanges := make([]string, len(r.MarketCodes))
		for i, e := range r.MarketCodes {
			exchanges[i] = string(e)
		}
		res = append(res, "p_mkt_code="+strings.Join(exchanges, ","))
	}

	if r.StartStreamNumber <= 0 {
		res = append(res, "p_eno=0")
	} else {
		res = append(res, fmt.Sprintf("p_eno=%d", r.StartStreamNumber))
	}

	eventList := make([]string, len(r.StreamEventTypes))
	for i, e := range r.StreamEventTypes {
		eventList[i] = string(e)
	}
	res = append(res, "p_evt_cmd="+strings.Join(eventList, ","))

	return []byte(strings.Join(res, "&"))
}

type StreamResponse interface {
	GetEventType() EventType
	GetErrorNo() ErrorNo
	GetErrorText() string
	parse(m map[string][]string, b []byte)
}

type CommonStreamResponse struct {
	EventType      EventType
	StreamNumber   int64
	StreamDateTime time.Time
	ErrorNo        ErrorNo
	ErrorText      string
	Body           []byte
}

func (r *CommonStreamResponse) GetEventType() EventType {
	return r.EventType
}

func (r *CommonStreamResponse) GetErrorNo() ErrorNo {
	return r.ErrorNo
}

func (r *CommonStreamResponse) GetErrorText() string {
	return r.ErrorText
}

func (r *CommonStreamResponse) getFromMap(m map[string][]string, key string) []string {
	if s, ok := m[key]; ok {
		return s
	} else {
		return []string{""}
	}
}

func (r *CommonStreamResponse) parse(m map[string][]string, b []byte) {
	r.EventType = EventType(r.getFromMap(m, "p_cmd")[0])
	r.StreamNumber, _ = strconv.ParseInt(r.getFromMap(m, "p_no")[0], 10, 64)
	r.StreamDateTime, _ = time.ParseInLocation("2006.01.02-15:04:05.000", r.getFromMap(m, "p_date")[0], time.Local)
	r.ErrorNo = ErrorNo(r.getFromMap(m, "p_errno")[0])
	r.ErrorText = r.getFromMap(m, "p_err")[0]
	r.Body = b
}

type MarketPriceStreamResponse struct {
	CommonStreamResponse
	ColumnNumber      int                 // 行番号
	AskQuantityMarket float64             // 売数量(成行)
	BidQuantityMarket float64             // 買数量(成行)
	AskQuantity       float64             // 売気配数量
	BidQuantity       float64             // 買気配数量
	DiscontinuityType string              // 不連続要因銘柄区分
	StopHigh          CurrentPriceType    // 日通し高値フラグ
	HighPrice         float64             // 高値
	HighPriceTime     time.Time           // 高値時刻
	TradingAmount     float64             // 売買代金
	StopLow           CurrentPriceType    // 日通し安値フラグ
	LowPrice          float64             // 安値
	LowPriceTime      time.Time           // 安値時刻
	OpenPrice         float64             // 始値
	OpenPriceTime     time.Time           // 始値時刻
	ChangePriceType   ChangePriceType     // 現値前値比較
	CurrentPrice      float64             // 現在値
	CurrentPriceTime  time.Time           // 現在値時刻
	Volume            float64             // 出来高
	ExRightType       string              // 配当落銘柄区分
	PrevDayPercent    float64             // 騰落率
	PrevDayRatio      float64             // 前日比
	AskQuantity10     float64             // 売-10-数量
	AskPrice10        float64             // 売-10-値段
	AskQuantity9      float64             // 売-9-数量
	AskPrice9         float64             // 売-9-値段
	AskQuantity8      float64             // 売-8-数量
	AskPrice8         float64             // 売-8-値段
	AskQuantity7      float64             // 売-7-数量
	AskPrice7         float64             // 売-7-値段
	AskQuantity6      float64             // 売-6-数量
	AskPrice6         float64             // 売-6-値段
	AskQuantity5      float64             // 売-5-数量
	AskPrice5         float64             // 売-5-値段
	AskQuantity4      float64             // 売-4-数量
	AskPrice4         float64             // 売-4-値段
	AskQuantity3      float64             // 売-3-数量
	AskPrice3         float64             // 売-3-値段
	AskQuantity2      float64             // 売-2-数量
	AskPrice2         float64             // 売-2-値段
	AskQuantity1      float64             // 売-1-数量
	AskPrice1         float64             // 売-1-値段
	BidQuantity1      float64             // 買-1-数量
	BidPrice1         float64             // 買-1-値段
	BidQuantity2      float64             // 買-2-数量
	BidPrice2         float64             // 買-2-値段
	BidQuantity3      float64             // 買-3-数量
	BidPrice3         float64             // 買-3-値段
	BidQuantity4      float64             // 買-4-数量
	BidPrice4         float64             // 買-4-値段
	BidQuantity5      float64             // 買-5-数量
	BidPrice5         float64             // 買-5-値段
	BidQuantity6      float64             // 買-6-数量
	BidPrice6         float64             // 買-6-値段
	BidQuantity7      float64             // 買-7-数量
	BidPrice7         float64             // 買-7-値段
	BidQuantity8      float64             // 買-8-数量
	BidPrice8         float64             // 買-8-値段
	BidQuantity9      float64             // 買-9-数量
	BidPrice9         float64             // 買-9-値段
	BidQuantity10     float64             // 買-10-数量
	BidPrice10        float64             // 買-10-値段
	Section           string              // 所属
	PRP               float64             // 前日終値
	AskPrice          float64             // 売気配値
	AskSign           IndicationPriceType // 売気配値種類
	BidPrice          float64             // 買気配値
	BidSign           IndicationPriceType // 買気配値種類
	AskQuantityOver   float64             // 売-OVER
	BidQuantityUnder  float64             // 買-UNDER
	VWAP              float64             // VWAP
}

func (r *MarketPriceStreamResponse) getColumnNumber(m map[string][]string) int {
	for k := range m {
		match := regexp.MustCompile(`^p_(\d+)_.+$`).FindAllStringSubmatch(k, -1)
		if len(match) == 0 {
			continue
		}

		i, _ := strconv.Atoi(match[0][1])
		return i
	}
	return 0
}

func (r *MarketPriceStreamResponse) parse(m map[string][]string, b []byte) {
	r.CommonStreamResponse.parse(m, b)
	r.ColumnNumber = r.getColumnNumber(m)
	r.AskQuantityMarket, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_AAV", r.ColumnNumber))[0], 64)
	r.BidQuantityMarket, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_ABV", r.ColumnNumber))[0], 64)
	r.AskQuantity, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_AV", r.ColumnNumber))[0], 64)
	r.BidQuantity, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_BV", r.ColumnNumber))[0], 64)
	r.DiscontinuityType = r.getFromMap(m, fmt.Sprintf("p_%d_DCFS", r.ColumnNumber))[0]
	r.StopHigh = CurrentPriceType(r.getFromMap(m, fmt.Sprintf("p_%d_DHF", r.ColumnNumber))[0])
	r.HighPrice, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_DHP", r.ColumnNumber))[0], 64)
	r.HighPriceTime, _ = time.ParseInLocation("15:04", r.getFromMap(m, fmt.Sprintf("p_%d_DHP:T", r.ColumnNumber))[0], time.Local)
	r.TradingAmount, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_DJ", r.ColumnNumber))[0], 64)
	r.StopLow = CurrentPriceType(r.getFromMap(m, fmt.Sprintf("p_%d_DLF", r.ColumnNumber))[0])
	r.LowPrice, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_DLP", r.ColumnNumber))[0], 64)
	r.LowPriceTime, _ = time.ParseInLocation("15:04", r.getFromMap(m, fmt.Sprintf("p_%d_DLP:T", r.ColumnNumber))[0], time.Local)
	r.OpenPrice, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_DOP", r.ColumnNumber))[0], 64)
	r.OpenPriceTime, _ = time.ParseInLocation("15:04", r.getFromMap(m, fmt.Sprintf("p_%d_DOP:T", r.ColumnNumber))[0], time.Local)
	r.ChangePriceType = ChangePriceType(r.getFromMap(m, fmt.Sprintf("p_%d_DPG", r.ColumnNumber))[0])
	r.CurrentPrice, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_DPP", r.ColumnNumber))[0], 64)
	r.CurrentPriceTime, _ = time.ParseInLocation("15:04", r.getFromMap(m, fmt.Sprintf("p_%d_DPP:T", r.ColumnNumber))[0], time.Local)
	r.Volume, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_DV", r.ColumnNumber))[0], 64)
	r.ExRightType = r.getFromMap(m, fmt.Sprintf("p_%d_DVES", r.ColumnNumber))[0]
	r.PrevDayPercent, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_DYRP", r.ColumnNumber))[0], 64)
	r.PrevDayRatio, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_DYWP", r.ColumnNumber))[0], 64)
	r.AskQuantity10, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAV10", r.ColumnNumber))[0], 64)
	r.AskPrice10, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAP10", r.ColumnNumber))[0], 64)
	r.AskQuantity9, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAV9", r.ColumnNumber))[0], 64)
	r.AskPrice9, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAP9", r.ColumnNumber))[0], 64)
	r.AskQuantity8, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAV8", r.ColumnNumber))[0], 64)
	r.AskPrice8, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAP8", r.ColumnNumber))[0], 64)
	r.AskQuantity7, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAV7", r.ColumnNumber))[0], 64)
	r.AskPrice7, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAP7", r.ColumnNumber))[0], 64)
	r.AskQuantity6, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAV6", r.ColumnNumber))[0], 64)
	r.AskPrice6, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAP6", r.ColumnNumber))[0], 64)
	r.AskQuantity5, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAV5", r.ColumnNumber))[0], 64)
	r.AskPrice5, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAP5", r.ColumnNumber))[0], 64)
	r.AskQuantity4, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAV4", r.ColumnNumber))[0], 64)
	r.AskPrice4, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAP4", r.ColumnNumber))[0], 64)
	r.AskQuantity3, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAV3", r.ColumnNumber))[0], 64)
	r.AskPrice3, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAP3", r.ColumnNumber))[0], 64)
	r.AskQuantity2, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAV2", r.ColumnNumber))[0], 64)
	r.AskPrice2, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAP2", r.ColumnNumber))[0], 64)
	r.AskQuantity1, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAV1", r.ColumnNumber))[0], 64)
	r.AskPrice1, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GAP1", r.ColumnNumber))[0], 64)
	r.BidQuantity1, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBV1", r.ColumnNumber))[0], 64)
	r.BidPrice1, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBP1", r.ColumnNumber))[0], 64)
	r.BidQuantity2, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBV2", r.ColumnNumber))[0], 64)
	r.BidPrice2, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBP2", r.ColumnNumber))[0], 64)
	r.BidQuantity3, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBV3", r.ColumnNumber))[0], 64)
	r.BidPrice3, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBP3", r.ColumnNumber))[0], 64)
	r.BidQuantity4, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBV4", r.ColumnNumber))[0], 64)
	r.BidPrice4, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBP4", r.ColumnNumber))[0], 64)
	r.BidQuantity5, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBV5", r.ColumnNumber))[0], 64)
	r.BidPrice5, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBP5", r.ColumnNumber))[0], 64)
	r.BidQuantity6, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBV6", r.ColumnNumber))[0], 64)
	r.BidPrice6, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBP6", r.ColumnNumber))[0], 64)
	r.BidQuantity7, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBV7", r.ColumnNumber))[0], 64)
	r.BidPrice7, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBP7", r.ColumnNumber))[0], 64)
	r.BidQuantity8, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBV8", r.ColumnNumber))[0], 64)
	r.BidPrice8, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBP8", r.ColumnNumber))[0], 64)
	r.BidQuantity9, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBV9", r.ColumnNumber))[0], 64)
	r.BidPrice9, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBP9", r.ColumnNumber))[0], 64)
	r.BidQuantity10, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBV10", r.ColumnNumber))[0], 64)
	r.BidPrice10, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_GBP10", r.ColumnNumber))[0], 64)
	r.Section = r.getFromMap(m, fmt.Sprintf("p_%d_LISS", r.ColumnNumber))[0]
	r.PRP, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_PRP", r.ColumnNumber))[0], 64)
	r.AskPrice, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_QAP", r.ColumnNumber))[0], 64)
	r.AskSign = IndicationPriceType(r.getFromMap(m, fmt.Sprintf("p_%d_QAS", r.ColumnNumber))[0])
	r.BidPrice, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_QBP", r.ColumnNumber))[0], 64)
	r.BidSign = IndicationPriceType(r.getFromMap(m, fmt.Sprintf("p_%d_QBS", r.ColumnNumber))[0])
	r.AskQuantityOver, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_QOV", r.ColumnNumber))[0], 64)
	r.BidQuantityUnder, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_QUV", r.ColumnNumber))[0], 64)
	r.VWAP, _ = strconv.ParseFloat(r.getFromMap(m, fmt.Sprintf("p_%d_VWAP", r.ColumnNumber))[0], 64)
}

type ContractStreamResponse struct {
	CommonStreamResponse
	Provider                 string            // プロバイダ(情報提供元)
	EventNo                  int64             // イベント番号
	FirstTime                bool              // アラートフラグ
	StreamOrderType          StreamOrderType   // 通知種別
	OrderNumber              string            // 注文番号
	ExecutionDate            time.Time         // 営業日
	ParentOrderNumber        string            // 親注文番号
	ParentOrder              bool              // 注文種別
	ProductType              ProductType       // 商品種別
	IssueCode                string            // 銘柄コード
	Exchange                 Exchange          // 市場コード
	Side                     Side              // 売買区分
	TradeType                TradeType         // 取引区分
	ExecutionTiming          ExecutionTiming   // 執行条件
	ExecutionType            ExecutionType     // 注文値段区分
	Price                    float64           // 注文値段
	Quantity                 float64           // 注文数量
	CancelQuantity           float64           // 取消数量
	ExpireQuantity           float64           // 失効数量
	ContractQuantity         float64           // 約定済数量
	StreamOrderStatus        StreamOrderStatus // 注文ステータス
	CarryOverType            CarryOverType     // 繰越フラグ
	CancelOrderStatus        CancelOrderStatus // 訂正取消ステータス
	ContractStatus           ContractStatus    // 約定ステータス
	ExpireDate               time.Time         // 有効期限
	SecurityExpireReason     string            // 失効理由コード
	SecurityContractPrice    float64           // 約定値段
	SecurityContractQuantity float64           // 約定数量
	SecurityError            string            // 取引所エラーコード
	NotifyDateTime           time.Time         // 通知日時
	IssueName                string            // 銘柄名称
	CorrectExecutionTiming   ExecutionTiming   // 訂正執行条件
	CorrectContractQuantity  float64           // 訂正執行数量
	CorrectExecutionType     ExecutionType     // 訂正注文値段区分
	CorrectPrice             float64           // 訂正注文値段
	CorrectQuantity          float64           // 訂正注文数量
	CorrectExpireDate        time.Time         // 訂正注文期限
	CorrectStopOrderType     StopOrderType     // 訂正逆指値条件
	CorrectTriggerPrice      float64           // 訂正逆指値段区分
	CorrectStopOrderPrice    float64           // 訂正逆指値段
}

func (r *ContractStreamResponse) parse(m map[string][]string, b []byte) {
	r.CommonStreamResponse.parse(m, b)

	r.Provider = r.getFromMap(m, "p_PV")[0]
	r.EventNo, _ = strconv.ParseInt(r.getFromMap(m, "p_ENO")[0], 10, 64)
	r.FirstTime = r.getFromMap(m, "p_ALT")[0] != "0"
	r.StreamOrderType = StreamOrderType(r.getFromMap(m, "p_NT")[0])
	r.OrderNumber = r.getFromMap(m, "p_ON")[0]
	r.ExecutionDate, _ = time.ParseInLocation("20060102", r.getFromMap(m, "p_ED")[0], time.Local)
	r.ParentOrderNumber = r.getFromMap(m, "p_OON")[0]
	r.ParentOrder = r.getFromMap(m, "p_OT")[0] == "1"
	r.ProductType = ProductType(r.getFromMap(m, "p_ST")[0])
	r.IssueCode = r.getFromMap(m, "p_IC")[0]
	r.Exchange = Exchange(r.getFromMap(m, "p_MC")[0])
	r.Side = Side(r.getFromMap(m, "p_BBKB")[0])
	r.TradeType = TradeType(r.getFromMap(m, "p_THKB")[0])
	r.ExecutionTiming = ExecutionTiming(r.getFromMap(m, "p_CRSJ")[0])
	r.ExecutionType = ExecutionType(r.getFromMap(m, "p_CRPRKB")[0])
	r.Price, _ = strconv.ParseFloat(r.getFromMap(m, "p_CRPR")[0], 64)
	r.Quantity, _ = strconv.ParseFloat(r.getFromMap(m, "p_CRSR")[0], 64)
	r.CancelQuantity, _ = strconv.ParseFloat(r.getFromMap(m, "p_CRTKSR")[0], 64)
	r.ExpireQuantity, _ = strconv.ParseFloat(r.getFromMap(m, "p_CREPSR")[0], 64)
	r.ContractQuantity, _ = strconv.ParseFloat(r.getFromMap(m, "p_CREXSR")[0], 64)
	r.StreamOrderStatus = StreamOrderStatus(r.getFromMap(m, "p_ODST")[0])
	r.CarryOverType = CarryOverType(r.getFromMap(m, "p_KOFG")[0])
	r.CancelOrderStatus = CancelOrderStatus(r.getFromMap(m, "p_TTST")[0])
	r.ContractStatus = ContractStatus(r.getFromMap(m, "p_EXST")[0])
	if r.getFromMap(m, "p_LMIT")[0] == "00000000" { // 当日
		r.ExpireDate = r.ExecutionDate
	} else {
		r.ExpireDate, _ = time.ParseInLocation("20060102", r.getFromMap(m, "p_LMIT")[0], time.Local)
	}
	r.SecurityExpireReason = r.getFromMap(m, "p_EPRC")[0]
	r.SecurityContractPrice, _ = strconv.ParseFloat(r.getFromMap(m, "p_EXPR")[0], 64)
	r.SecurityContractQuantity, _ = strconv.ParseFloat(r.getFromMap(m, "p_EXSR")[0], 64)
	r.SecurityError = r.getFromMap(m, "p_EXRC")[0]
	r.NotifyDateTime, _ = time.ParseInLocation("20060102150405", r.getFromMap(m, "p_EXDT")[0], time.Local)
	r.IssueName = r.getFromMap(m, "p_IN")[0]
	r.CorrectExecutionTiming = ExecutionTiming(r.getFromMap(m, "p_UPSJ")[0])
	if r.getFromMap(m, "p_UPEXSR")[0] != "" {
		r.CorrectContractQuantity, _ = strconv.ParseFloat(r.getFromMap(m, "p_UPEXSR")[0], 64)
	}
	r.CorrectExecutionType = ExecutionType(r.getFromMap(m, "p_UPPRKB")[0])
	if r.getFromMap(m, "p_UPPR")[0] != "" {
		r.CorrectPrice, _ = strconv.ParseFloat(r.getFromMap(m, "p_UPPR")[0], 64)
	}
	if r.getFromMap(m, "p_UPSR")[0] != "" {
		r.CorrectQuantity, _ = strconv.ParseFloat(r.getFromMap(m, "p_UPSR")[0], 64)
	}
	if r.getFromMap(m, "p_UPLMIT")[0] != "" {
		r.CorrectExpireDate, _ = time.ParseInLocation("20060102", r.getFromMap(m, "p_UPLMIT")[0], time.Local)
	}
	r.CorrectStopOrderType = StopOrderType(r.getFromMap(m, "p_UPGKCDPR")[0])
	if r.getFromMap(m, "p_UPGKPRKB")[0] != "" {
		r.CorrectTriggerPrice, _ = strconv.ParseFloat(r.getFromMap(m, "p_UPGKPRKB")[0], 64)
	}
	if r.getFromMap(m, "p_UPGKPR")[0] != "" {
		r.CorrectStopOrderPrice, _ = strconv.ParseFloat(r.getFromMap(m, "p_UPGKPR")[0], 64)
	}
}

type NewsStreamResponse struct {
	CommonStreamResponse
	Provider      string    // プロバイダ(情報提供元)
	EventNo       int64     // イベント番号
	FirstTime     bool      // アラートフラグ
	NewsId        string    // ニュースID
	NewsDateTime  time.Time // ニュース日時
	NumOfCategory int       // ニュースカテゴリ数
	Categories    []string  // ニュースカテゴリリスト
	NumOfGenre    int       // ニュースジャンル数
	Genres        []string  // ニュースジャンルリスト
	NumOfIssue    int       // 関連銘柄コードリスト
	Issues        []string  // 関連銘柄コードリスト
	Title         string    // ニュースタイトル
	Content       string    // ニュース本文
}

func (r *NewsStreamResponse) parse(m map[string][]string, b []byte) {
	r.CommonStreamResponse.parse(m, b)

	r.Provider = r.getFromMap(m, "p_PV")[0]
	r.EventNo, _ = strconv.ParseInt(r.getFromMap(m, "p_ENO")[0], 10, 64)
	r.FirstTime = r.getFromMap(m, "p_ALT")[0] != "0"
	r.NewsId = r.getFromMap(m, "p_ID")[0]
	r.NewsDateTime, _ = time.ParseInLocation("20060102150405", r.getFromMap(m, "p_DT")[0]+r.getFromMap(m, "p_TM")[0], time.Local)
	r.NumOfCategory, _ = strconv.Atoi(r.getFromMap(m, "p_CGN")[0])
	r.Categories = r.getFromMap(m, "p_CGL")
	r.NumOfGenre, _ = strconv.Atoi(r.getFromMap(m, "p_GRN")[0])
	r.Genres = r.getFromMap(m, "p_GRL")
	r.NumOfIssue, _ = strconv.Atoi(r.getFromMap(m, "p_ISN")[0])
	r.Issues = r.getFromMap(m, "p_ISL")
	r.Title = r.getFromMap(m, "p_HDL")[0]
	r.Content = r.getFromMap(m, "p_TX")[0]
}

type SystemStatusStreamResponse struct {
	CommonStreamResponse
	Provider       string        // プロバイダ(情報提供元)
	EventNo        int64         // イベント番号
	FirstTime      bool          // アラートフラグ
	UpdateDateTime time.Time     // 情報更新時間
	ApprovalLogin  ApprovalLogin // ログイン許可区分
	SystemStatus   SystemStatus  // システムステータス
}

func (r *SystemStatusStreamResponse) parse(m map[string][]string, b []byte) {
	r.CommonStreamResponse.parse(m, b)

	r.Provider = r.getFromMap(m, "p_PV")[0]
	r.EventNo, _ = strconv.ParseInt(r.getFromMap(m, "p_ENO")[0], 10, 64)
	r.FirstTime = r.getFromMap(m, "p_ALT")[0] != "0"
	r.UpdateDateTime, _ = time.ParseInLocation("20060102150405", r.getFromMap(m, "p_CT")[0], time.Local)
	r.ApprovalLogin = ApprovalLogin(r.getFromMap(m, "p_LK")[0])
	r.SystemStatus = SystemStatus(r.getFromMap(m, "p_SS")[0])
}

type OperationStatusStreamResponse struct {
	CommonStreamResponse
	Provider          string    // プロバイダ(情報提供元)
	EventNo           int64     // イベント番号
	FirstTime         bool      // アラートフラグ
	UpdateDateTime    time.Time // 情報更新時間
	Exchange          Exchange  // 市場コード
	AssetCode         string    // 原資産コード
	ProductType       string    // 商品種別
	OperationCategory string    // 運用カテゴリー
	OperationUnit     string    // 運用ユニット
	BusinessDayType   string    // 営業日区分
	OperationStatus   string    // 運用ステータス
}

func (r *OperationStatusStreamResponse) parse(m map[string][]string, b []byte) {
	r.CommonStreamResponse.parse(m, b)

	r.Provider = r.getFromMap(m, "p_PV")[0]
	r.EventNo, _ = strconv.ParseInt(r.getFromMap(m, "p_ENO")[0], 10, 64)
	r.FirstTime = r.getFromMap(m, "p_ALT")[0] != "0"
	r.UpdateDateTime, _ = time.ParseInLocation("20060102150405", r.getFromMap(m, "p_CT")[0], time.Local)
	r.Exchange = Exchange(r.getFromMap(m, "p_MC")[0])
	r.AssetCode = r.getFromMap(m, "p_GSCD")[0]
	r.ProductType = r.getFromMap(m, "p_SHSB")[0]
	r.OperationCategory = r.getFromMap(m, "p_UC")[0]
	r.OperationUnit = r.getFromMap(m, "p_UU")[0]
	r.BusinessDayType = r.getFromMap(m, "p_EDK")[0]
	r.OperationStatus = r.getFromMap(m, "p_US")[0]
}

func (c *client) Stream(ctx context.Context, session *Session, req StreamRequest) (<-chan StreamResponse, <-chan error) {
	eventCh := make(chan StreamResponse)
	errCh := make(chan error)

	go func() {
		defer time.Sleep(1 * time.Second)
		defer close(eventCh)
		defer close(errCh)

		if session == nil {
			errCh <- NilArgumentErr
			return
		}

		ch1, ch2 := c.requester.stream(ctx, session.EventURL, req)
		for {
			select {
			case err, ok := <-ch2:
				if ok {
					errCh <- err
					return
				}
			case b, ok := <-ch1:
				// chanがcloseされたら抜ける
				if !ok {
					return
				}

				// ここでcommonResponseに一回変換する
				m := c.streamResponseToMap(b)
				t, ok := m["p_cmd"]
				if !ok {
					t = []string{""}
				}
				var res StreamResponse
				switch EventType(t[0]) {
				case EventTypeKeepAlive: // keep aliveは通知しない
					continue
				case EventTypeMarketPrice:
					res = new(MarketPriceStreamResponse)
				case EventTypeContract:
					res = new(ContractStreamResponse)
				case EventTypeNews:
					res = new(NewsStreamResponse)
				case EventTypeSystemStatus:
					res = new(SystemStatusStreamResponse)
				case EventTypeOperationStatus:
					res = new(OperationStatusStreamResponse)
				default:
					res = new(CommonStreamResponse)
				}
				res.parse(m, b)

				if res.GetErrorNo() != ErrorNoProblem {
					errCh <- fmt.Errorf("%w: %s(%s)", StreamError, res.GetErrorText(), res.GetErrorNo())
					return
				}

				eventCh <- res
			}
		}
	}()

	return eventCh, errCh
}

func (c *client) streamResponseToMap(b []byte) map[string][]string {
	res := make(map[string][]string)
	for _, b := range bytes.Split(b, []byte{1}) {
		bs := bytes.Split(b, []byte{2})
		switch len(bs) {
		case 1:
			res[string(bs[0])] = []string{""}
		case 2:
			res[string(bs[0])] = strings.Split(string(bs[1]), string([]byte{2}))
		}
	}
	return res
}
