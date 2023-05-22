package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// StockExchangeMasterColumn - 株式銘柄市場マスタカラム
type StockExchangeMasterColumn string

const (
	StockExchangeMasterColumnIssueCode                   StockExchangeMasterColumn = "sIssueCode"                // 銘柄コード
	StockExchangeMasterColumnExchange                    StockExchangeMasterColumn = "sZyouzyouSizyou"           // 上場市場
	StockExchangeMasterColumnStockSystemType             StockExchangeMasterColumn = "sSystemC"                  // システムC
	StockExchangeMasterColumnUnderLimitPrice             StockExchangeMasterColumn = "sNehabaMin"                // 値幅下限
	StockExchangeMasterColumnUpperLimitPrice             StockExchangeMasterColumn = "sNehabaMax"                // 値幅上限
	StockExchangeMasterColumnSymbolCategory              StockExchangeMasterColumn = "sIssueKubunC"              // 銘柄区分C
	StockExchangeMasterColumnLimitPriceExchange          StockExchangeMasterColumn = "sNehabaSizyouC"            // 値幅市場C
	StockExchangeMasterColumnMarginType                  StockExchangeMasterColumn = "sSinyouC"                  // 信用C
	StockExchangeMasterColumnListingDate                 StockExchangeMasterColumn = "sSinkiZyouzyouDay"         // 新規上場日
	StockExchangeMasterColumnLimitPriceDate              StockExchangeMasterColumn = "sNehabaKigenDay"           // 値幅期限日
	StockExchangeMasterColumnLimitPriceCategory          StockExchangeMasterColumn = "sNehabaKiseiC"             // 値幅規制C
	StockExchangeMasterColumnLimitPriceValue             StockExchangeMasterColumn = "sNehabaKiseiTi"            // 値幅規制値
	StockExchangeMasterColumnConfirmLimitPrice           StockExchangeMasterColumn = "sNehabaCheckKahiC"         // 値幅チェック可否C
	StockExchangeMasterColumnSection                     StockExchangeMasterColumn = "sIssueBubetuC"             // 銘柄部別C
	StockExchangeMasterColumnPrevClosePrice              StockExchangeMasterColumn = "sZenzituOwarine"           // 前日終値
	StockExchangeMasterColumnCalculateLimitPriceExchange StockExchangeMasterColumn = "sNehabaSansyutuSizyouC"    // 値幅算出市場C
	StockExchangeMasterColumnRegulation1                 StockExchangeMasterColumn = "sIssueKisei1C"             // 銘柄規制1C
	StockExchangeMasterColumnRegulation2                 StockExchangeMasterColumn = "sIssueKisei2C"             // 銘柄規制2C
	StockExchangeMasterColumnSectionType                 StockExchangeMasterColumn = "sZyouzyouKubun"            // 上場区分
	StockExchangeMasterColumnDelistingDate               StockExchangeMasterColumn = "sZyouzyouHaisiDay"         // 上場廃止日
	StockExchangeMasterColumnTradingUnit                 StockExchangeMasterColumn = "sSizyoubetuBaibaiTani"     // 売買単位
	StockExchangeMasterColumnNextTradingUnit             StockExchangeMasterColumn = "sSizyoubetuBaibaiTaniYoku" // 売買単位(翌営業日)
	StockExchangeMasterColumnTickGroupType               StockExchangeMasterColumn = "sYobineTaniNumber"         // 呼値の単位番号
	StockExchangeMasterColumnNextTickGroupType           StockExchangeMasterColumn = "sYobineTaniNumberYoku"     // 呼値の単位番号(翌営業日)
	StockExchangeMasterColumnInformationSource           StockExchangeMasterColumn = "sZyouhouSource"            // 情報系ソース
	StockExchangeMasterColumnInformationCode             StockExchangeMasterColumn = "sZyouhouCode"              // 情報系コード
	StockExchangeMasterColumnOfferPrice                  StockExchangeMasterColumn = "sKouboPrice"               // 公募価格
	StockExchangeMasterColumnCreateDateTime              StockExchangeMasterColumn = "sCreateDate"               // 作成日時
	StockExchangeMasterColumnUpdateDateTime              StockExchangeMasterColumn = "sUpdateDate"               // 更新日時
	StockExchangeMasterColumnUpdateNumber                StockExchangeMasterColumn = "sUpdateNumber"             // 更新通番
)

// StockExchangeMasterRequest - 株式銘柄市場マスタリクエスト
type StockExchangeMasterRequest struct {
	Columns []StockExchangeMasterColumn // 取得したい情報
}

func (r *StockExchangeMasterRequest) request(no int64, now time.Time) stockExchangeMasterRequest {
	columns := make([]string, len(r.Columns))
	for i, column := range r.Columns {
		columns[i] = string(column)
	}

	return stockExchangeMasterRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			MessageType:    MessageTypeMasterData,
			ResponseFormat: commonResponseFormat,
		},
		TargetFeature: string(MessageTypeStockExchangeMaster),
		Columns:       strings.Join(columns, ","),
	}
}

type stockExchangeMasterRequest struct {
	commonRequest
	TargetFeature string `json:"sTargetCLMID,omitempty"`  // 取得したいマスタデータ
	Columns       string `json:"sTargetColumn,omitempty"` // 取得したい情報
}

type stockExchangeMasterResponse struct {
	commonResponse
	StockExchangeMasters []stockExchangeMaster `json:"CLMIssueSizyouMstKabu"` // 株式銘柄市場マスタ
}

func (r *stockExchangeMasterResponse) response() StockExchangeMasterResponse {
	stockExchangeMasters := make([]StockExchangeMaster, len(r.StockExchangeMasters))
	for i, stock := range r.StockExchangeMasters {
		stockExchangeMasters[i] = stock.response()
	}
	return StockExchangeMasterResponse{
		CommonResponse:       r.commonResponse.response(),
		StockExchangeMasters: stockExchangeMasters,
	}
}

type stockExchangeMaster struct {
	IssueCode                   string        `json:"sIssueCode"`                       // 銘柄コード
	Exchange                    Exchange      `json:"sZyouzyouSizyou"`                  // 上場市場
	StockSystemType             string        `json:"sSystemC"`                         // システムC
	UnderLimitPrice             float64       `json:"sNehabaMin,string"`                // 値幅下限
	UpperLimitPrice             float64       `json:"sNehabaMax,string"`                // 値幅上限
	SymbolCategory              string        `json:"sIssueKubunC"`                     // 銘柄区分C
	LimitPriceExchange          Exchange      `json:"sNehabaSizyouC"`                   // 値幅市場C
	MarginType                  MarginType    `json:"sSinyouC"`                         // 信用C
	ListingDate                 Ymd           `json:"sSinkiZyouzyouDay"`                // 新規上場日
	LimitPriceDate              Ymd           `json:"sNehabaKigenDay"`                  // 値幅期限日
	LimitPriceCategory          string        `json:"sNehabaKiseiC"`                    // 値幅規制C
	LimitPriceValue             float64       `json:"sNehabaKiseiTi,string"`            // 値幅規制値
	ConfirmLimitPrice           NumberBool    `json:"sNehabaCheckKahiC"`                // 値幅チェック可否C
	Section                     string        `json:"sIssueBubetuC"`                    // 銘柄部別C
	PrevClosePrice              float64       `json:"sZenzituOwarine,string"`           // 前日終値
	CalculateLimitPriceExchange Exchange      `json:"sNehabaSansyutuSizyouC"`           // 値幅算出市場C
	Regulation1                 string        `json:"sIssueKisei1C"`                    // 銘柄規制1C
	Regulation2                 string        `json:"sIssueKisei2C"`                    // 銘柄規制2C
	SectionType                 string        `json:"sZyouzyouKubun"`                   // 上場区分
	DelistingDate               Ymd           `json:"sZyouzyouHaisiDay"`                // 上場廃止日
	TradingUnit                 float64       `json:"sSizyoubetuBaibaiTan,stringi"`     // 売買単位
	NextTradingUnit             float64       `json:"sSizyoubetuBaibaiTaniYoku,string"` // 売買単位(翌営業日)
	TickGroupType               TickGroupType `json:"sYobineTaniNumber"`                // 呼値の単位番号
	NextTickGroupType           TickGroupType `json:"sYobineTaniNumberYoku"`            // 呼値の単位番号(翌営業日)
	InformationSource           string        `json:"sZyouhouSource"`                   // 情報系ソース
	InformationCode             string        `json:"sZyouhouCode"`                     // 情報系コード
	OfferPrice                  float64       `json:"sKouboPrice,string"`               // 公募価格
	CreateDateTime              YmdHms        `json:"sCreateDate"`                      // 作成日時
	UpdateDateTime              YmdHms        `json:"sUpdateDate"`                      // 更新日時
	UpdateNumber                string        `json:"sUpdateNumber"`                    // 更新通番
}

func (r *stockExchangeMaster) UnmarshalJSON(b []byte) error {
	replaced := b
	replaces := map[string]string{
		`"sNehabaMin":""`:      `"sNehabaMin":"0"`,
		`"sNehabaMax":""`:      `"sNehabaMax":"0"`,
		`"sZenzituOwarine":""`: `"sZenzituOwarine":"0"`,
	}
	for o, n := range replaces {
		replaced = bytes.Replace(replaced, []byte(o), []byte(n), -1)
	}

	type alias stockExchangeMaster
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *stockExchangeMaster) response() StockExchangeMaster {
	return StockExchangeMaster{
		IssueCode:                   r.IssueCode,
		Exchange:                    r.Exchange,
		StockSystemType:             r.StockSystemType,
		UnderLimitPrice:             r.UnderLimitPrice,
		UpperLimitPrice:             r.UpperLimitPrice,
		SymbolCategory:              r.SymbolCategory,
		LimitPriceExchange:          r.LimitPriceExchange,
		MarginType:                  r.MarginType,
		ListingDate:                 r.ListingDate.Time,
		LimitPriceDate:              r.LimitPriceDate.Time,
		LimitPriceCategory:          r.LimitPriceCategory,
		LimitPriceValue:             r.LimitPriceValue,
		ConfirmLimitPrice:           r.ConfirmLimitPrice.Bool(),
		Section:                     r.Section,
		PrevClosePrice:              r.PrevClosePrice,
		CalculateLimitPriceExchange: r.CalculateLimitPriceExchange,
		Regulation1:                 r.Regulation1,
		Regulation2:                 r.Regulation2,
		SectionType:                 r.SectionType,
		DelistingDate:               r.DelistingDate.Time,
		TradingUnit:                 r.TradingUnit,
		NextTradingUnit:             r.NextTradingUnit,
		TickGroupType:               r.TickGroupType,
		NextTickGroupType:           r.NextTickGroupType,
		InformationSource:           r.InformationSource,
		InformationCode:             r.InformationCode,
		OfferPrice:                  r.OfferPrice,
		CreateDateTime:              r.CreateDateTime.Time,
		UpdateDateTime:              r.UpdateDateTime.Time,
		UpdateNumber:                r.UpdateNumber,
	}
}

// StockExchangeMasterResponse - 株式銘柄市場マスタレスポンス
type StockExchangeMasterResponse struct {
	CommonResponse
	StockExchangeMasters []StockExchangeMaster // 株式銘柄市場マスタ
}

// StockExchangeMaster - 株式銘柄市場マスタ
type StockExchangeMaster struct {
	IssueCode                   string        // 銘柄コード
	Exchange                    Exchange      // 上場市場
	StockSystemType             string        // システムC
	UnderLimitPrice             float64       // 値幅下限
	UpperLimitPrice             float64       // 値幅上限
	SymbolCategory              string        // 銘柄区分C
	LimitPriceExchange          Exchange      // 値幅市場C
	MarginType                  MarginType    // 信用C
	ListingDate                 time.Time     // 新規上場日
	LimitPriceDate              time.Time     // 値幅期限日
	LimitPriceCategory          string        // 値幅規制C
	LimitPriceValue             float64       // 値幅規制値
	ConfirmLimitPrice           bool          // 値幅チェック可否C
	Section                     string        // 銘柄部別C
	PrevClosePrice              float64       // 前日終値
	CalculateLimitPriceExchange Exchange      // 値幅算出市場C
	Regulation1                 string        // 銘柄規制1C
	Regulation2                 string        // 銘柄規制2C
	SectionType                 string        // 上場区分
	DelistingDate               time.Time     // 上場廃止日
	TradingUnit                 float64       // 売買単位
	NextTradingUnit             float64       // 売買単位(翌営業日)
	TickGroupType               TickGroupType // 呼値の単位番号
	NextTickGroupType           TickGroupType // 呼値の単位番号(翌営業日)
	InformationSource           string        // 情報系ソース
	InformationCode             string        // 情報系コード
	OfferPrice                  float64       // 公募価格
	CreateDateTime              time.Time     // 作成日時
	UpdateDateTime              time.Time     // 更新日時
	UpdateNumber                string        // 更新通番
}

// StockExchangeMaster - 株式銘柄市場マスタ
func (c *client) StockExchangeMaster(ctx context.Context, session *Session, req StockExchangeMasterRequest) (*StockExchangeMasterResponse, error) {
	if session == nil {
		return nil, NilArgumentErr
	}
	session.mtx.Lock()
	defer session.mtx.Unlock()

	session.lastRequestNo++
	r := req.request(session.lastRequestNo, c.clock.Now())

	b, err := c.requester.get(ctx, session.MasterURL, r)
	if err != nil {
		return nil, err
	}
	var res stockExchangeMasterResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, fmt.Errorf("%s: %w", err, UnmarshalFailedErr)
	}

	Res := res.response()
	return &Res, nil
}
