package tachibana

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// StockMasterColumn - 株式銘柄マスタカラム
type StockMasterColumn string

const (
	StockMasterColumnIssueCode            StockMasterColumn = "sIssueCode"             // 銘柄コード
	StockMasterColumnName                 StockMasterColumn = "sIssueName"             // 銘柄名
	StockMasterColumnShortName            StockMasterColumn = "sIssueNameRyaku"        // 銘柄名略称
	StockMasterColumnKana                 StockMasterColumn = "sIssueNameKana"         // 銘柄名(カナ)
	StockMasterColumnAlphabet             StockMasterColumn = "sIssueNameEizi"         // 銘柄名(英語表記)
	StockMasterColumnSpecificTarget       StockMasterColumn = "sTokuteiF"              // 特定口座対象C
	StockMasterColumnTaxFree              StockMasterColumn = "sHikazeiC"              // 非課税対象C
	StockMasterColumnSharedStocks         StockMasterColumn = "sZyouzyouHakkouKabusu"  // 上場発行株数
	StockMasterColumnExRight              StockMasterColumn = "sKenriotiFlag"          // 権利落ちフラグ
	StockMasterColumnLastRightDay         StockMasterColumn = "sKenritukiSaisyuDay"    // 権利付最終日
	StockMasterColumnListingType          StockMasterColumn = "sZyouzyouNyusatuC"      // 上場・入札C
	StockMasterColumnReleaseTradingDate   StockMasterColumn = "sNyusatuKaizyoDay"      // 入札解除日
	StockMasterColumnTradingDate          StockMasterColumn = "sNyusatuDay"            // 入札日
	StockMasterColumnTradingUnit          StockMasterColumn = "sBaibaiTani"            // 売買単位
	StockMasterColumnNextTradingUnit      StockMasterColumn = "sBaibaiTaniYoku"        // 売買単位(翌営業日)
	StockMasterColumnStopTradingType      StockMasterColumn = "sBaibaiTeisiC"          // 売買停止C
	StockMasterColumnStartPublicationDate StockMasterColumn = "sHakkouKaisiDay"        // 発行開始日
	StockMasterColumnLastPublicationDate  StockMasterColumn = "sHakkouSaisyuDay"       // 発行最終日
	StockMasterColumnSettlementType       StockMasterColumn = "sKessanC"               // 決算C
	StockMasterColumnSettlementDate       StockMasterColumn = "sKessanDay"             // 決算日
	StockMasterColumnListingDate          StockMasterColumn = "sZyouzyouOutouDay"      // 上場応答日
	StockMasterColumnExpireDate2Type      StockMasterColumn = "sNiruiKizituC"          // 二類期日C
	StockMasterColumnLargeUnit            StockMasterColumn = "sOogutiKabusu"          // 大口株数
	StockMasterColumnLargeAmount          StockMasterColumn = "sOogutiKingmaker"       // 大口金額
	StockMasterColumnOutputTicketType     StockMasterColumn = "sBadenpyouOutputYNC"    // 場伝票出力有無C
	StockMasterColumnDepositAmount        StockMasterColumn = "sHosyoukinDaiyouKakeme" // 保証金代用掛目
	StockMasterColumnDepositValuation     StockMasterColumn = "sDaiyouHyoukaTanka"     // 代用証券評価単価
	StockMasterColumnOrganizationType     StockMasterColumn = "sKikoSankaC"            // 機構参加C
	StockMasterColumnProvisionalType      StockMasterColumn = "sKarikessaiC"           // 仮決済C
	StockMasterColumnPrimaryExchange      StockMasterColumn = "sYusenSizyou"           // 優先市場
	StockMasterColumnIndefinitePeriodType StockMasterColumn = "sMukigenC"              // 無期限対象C
	StockMasterColumnIndustryCode         StockMasterColumn = "sGyousyuCode"           // 業種コード
	StockMasterColumnIndustryName         StockMasterColumn = "sGyousyuName"           // 業種コード名
	StockMasterColumnSORTargetType        StockMasterColumn = "sSorC"                  // SOR対象銘柄C
	StockMasterColumnCreateDateTime       StockMasterColumn = "sCreateDate"            // 作成日時
	StockMasterColumnUpdateDateTime       StockMasterColumn = "sUpdateDate"            // 更新日時
	StockMasterColumnUpdateNumber         StockMasterColumn = "sUpdateNumber"          // 更新通番
)

// StockMasterRequest - 株式銘柄マスタリクエスト
type StockMasterRequest struct {
	Columns []StockMasterColumn // 取得したい情報
}

func (r *StockMasterRequest) request(no int64, now time.Time) stockMasterRequest {
	columns := make([]string, len(r.Columns))
	for i, column := range r.Columns {
		columns[i] = string(column)
	}

	return stockMasterRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			MessageType:    MessageTypeMasterData,
			ResponseFormat: commonResponseFormat,
		},
		TargetFeature: string(MessageTypeStockMaster),
		Columns:       strings.Join(columns, ","),
	}
}

type stockMasterRequest struct {
	commonRequest
	TargetFeature string `json:"sTargetCLMID,omitempty"`  // 取得したいマスタデータ
	Columns       string `json:"sTargetColumn,omitempty"` // 取得したい情報
}

type stockMasterResponse struct {
	commonResponse
	StockMasters []stockMaster `json:"CLMIssueMstKabu"` // 株式銘柄マスタ
}

func (r *stockMasterResponse) response() StockMasterResponse {
	stockMasters := make([]StockMaster, len(r.StockMasters))
	for i, stock := range r.StockMasters {
		stockMasters[i] = stock.response()
	}
	return StockMasterResponse{
		CommonResponse: r.commonResponse.response(),
		StockMasters:   stockMasters,
	}
}

type stockMaster struct {
	Code                 string          `json:"sIssueCode"`                    // 銘柄コード
	Name                 string          `json:"sIssueName"`                    // 銘柄名
	ShortName            string          `json:"sIssueNameRyaku"`               // 銘柄名略称
	Kana                 string          `json:"sIssueNameKana"`                // 銘柄名(カナ)
	Alphabet             string          `json:"sIssueNameEizi"`                // 銘柄名(英語表記)
	SpecificTarget       NumberBool      `json:"sTokuteiF"`                     // 特定口座対象C
	TaxFree              TaxFree         `json:"sHikazeiC"`                     // 非課税対象C
	SharedStocks         int64           `json:"sZyouzyouHakkouKabusu,string"`  // 上場発行株数
	ExRight              ExRightType     `json:"sKenriotiFlag"`                 // 権利落ちフラグ
	LastRightDay         Ymd             `json:"sKenritukiSaisyuDay"`           // 権利付最終日
	ListingType          ListingType     `json:"sZyouzyouNyusatuC"`             // 上場・入札C
	ReleaseTradingDate   Ymd             `json:"sNyusatuKaizyoDay"`             // 入札解除日
	TradingDate          Ymd             `json:"sNyusatuDay"`                   // 入札日
	TradingUnit          float64         `json:"sBaibaiTani,string"`            // 売買単位
	NextTradingUnit      float64         `json:"sBaibaiTaniYoku,string"`        // 売買単位(翌営業日)
	StopTradingType      StopTradingType `json:"sBaibaiTeisiC"`                 // 売買停止C
	StartPublicationDate Ymd             `json:"sHakkouKaisiDay"`               // 発行開始日
	LastPublicationDate  Ymd             `json:"sHakkouSaisyuDay"`              // 発行最終日
	SettlementType       SettlementType  `json:"sKessanC"`                      // 決算C
	SettlementDate       Ymd             `json:"sKessanDay"`                    // 決算日
	ListingDate          Ymd             `json:"sZyouzyouOutouDay"`             // 上場応答日
	ExpireDate2Type      string          `json:"sNiruiKizituC"`                 // 二類期日C
	LargeUnit            float64         `json:"sOogutiKabusu,string"`          // 大口株数
	LargeAmount          float64         `json:"sOogutiKingmaker,string"`       // 大口金額
	OutputTicketType     string          `json:"sBadenpyouOutputYNC"`           // 場伝票出力有無C
	DepositAmount        float64         `json:"sHosyoukinDaiyouKakeme,string"` // 保証金代用掛目
	DepositValuation     float64         `json:"sDaiyouHyoukaTanka,string"`     // 代用証券評価単価
	OrganizationType     string          `json:"sKikoSankaC"`                   // 機構参加C
	ProvisionalType      string          `json:"sKarikessaiC"`                  // 仮決済C
	PrimaryExchange      Exchange        `json:"sYusenSizyou"`                  // 優先市場
	IndefinitePeriodType string          `json:"sMukigenC"`                     // 無期限対象C
	IndustryCode         string          `json:"sGyousyuCode"`                  // 業種コード
	IndustryName         string          `json:"sGyousyuName"`                  // 業種コード名
	SORTargetType        string          `json:"sSorC"`                         // SOR対象銘柄C
	CreateDateTime       YmdHms          `json:"sCreateDate"`                   // 作成日時
	UpdateDateTime       YmdHms          `json:"sUpdateDate"`                   // 更新日時
	UpdateNumber         string          `json:"sUpdateNumber"`                 // 更新通番
}

func (r *stockMaster) response() StockMaster {
	return StockMaster{
		IssueCode:            r.Code,
		Name:                 r.Name,
		ShortName:            r.ShortName,
		Kana:                 r.Kana,
		Alphabet:             r.Alphabet,
		SpecificTarget:       r.SpecificTarget.Bool(),
		TaxFree:              r.TaxFree,
		SharedStocks:         r.SharedStocks,
		ExRight:              r.ExRight,
		LastRightDay:         r.LastRightDay.Time,
		ListingType:          r.ListingType,
		ReleaseTradingDate:   r.ReleaseTradingDate.Time,
		TradingDate:          r.TradingDate.Time,
		TradingUnit:          r.TradingUnit,
		NextTradingUnit:      r.NextTradingUnit,
		StopTradingType:      r.StopTradingType,
		StartPublicationDate: r.StartPublicationDate.Time,
		LastPublicationDate:  r.LastPublicationDate.Time,
		SettlementType:       r.SettlementType,
		SettlementDate:       r.SettlementDate.Time,
		ListingDate:          r.ListingDate.Time,
		ExpireDate2Type:      r.ExpireDate2Type,
		LargeUnit:            r.LargeUnit,
		LargeAmount:          r.LargeAmount,
		OutputTicketType:     r.OutputTicketType,
		DepositAmount:        r.DepositAmount,
		DepositValuation:     r.DepositValuation,
		OrganizationType:     r.OrganizationType,
		ProvisionalType:      r.ProvisionalType,
		PrimaryExchange:      r.PrimaryExchange,
		IndefinitePeriodType: r.IndefinitePeriodType,
		IndustryCode:         r.IndustryCode,
		IndustryName:         r.IndustryName,
		SORTargetType:        r.SORTargetType,
		CreateDateTime:       r.CreateDateTime.Time,
		UpdateDateTime:       r.UpdateDateTime.Time,
		UpdateNumber:         r.UpdateNumber,
	}
}

// StockMasterResponse - 株式銘柄マスタレスポンス
type StockMasterResponse struct {
	CommonResponse
	StockMasters []StockMaster // 株式銘柄マスタ
}

// StockMaster - 株式銘柄
type StockMaster struct {
	IssueCode            string          `json:"sIssueCode"`                    // 銘柄コード
	Name                 string          `json:"sIssueName"`                    // 銘柄名
	ShortName            string          `json:"sIssueNameRyaku"`               // 銘柄名略称
	Kana                 string          `json:"sIssueNameKana"`                // 銘柄名(カナ)
	Alphabet             string          `json:"sIssueNameEizi"`                // 銘柄名(英語表記)
	SpecificTarget       bool            `json:"sTokuteiF"`                     // 特定口座対象C
	TaxFree              TaxFree         `json:"sHikazeiC"`                     // 非課税対象C
	SharedStocks         int64           `json:"sZyouzyouHakkouKabusu,string"`  // 上場発行株数
	ExRight              ExRightType     `json:"sKenriotiFlag"`                 // 権利落ちフラグ
	LastRightDay         time.Time       `json:"sKenritukiSaisyuDay"`           // 権利付最終日
	ListingType          ListingType     `json:"sZyouzyouNyusatuC"`             // 上場・入札C
	ReleaseTradingDate   time.Time       `json:"sNyusatuKaizyoDay"`             // 入札解除日
	TradingDate          time.Time       `json:"sNyusatuDay"`                   // 入札日
	TradingUnit          float64         `json:"sBaibaiTani,string"`            // 売買単位
	NextTradingUnit      float64         `json:"sBaibaiTaniYoku,string"`        // 売買単位(翌営業日)
	StopTradingType      StopTradingType `json:"sBaibaiTeisiC"`                 // 売買停止C
	StartPublicationDate time.Time       `json:"sHakkouKaisiDay"`               // 発行開始日
	LastPublicationDate  time.Time       `json:"sHakkouSaisyuDay"`              // 発行最終日
	SettlementType       SettlementType  `json:"sKessanC"`                      // 決算C
	SettlementDate       time.Time       `json:"sKessanDay"`                    // 決算日
	ListingDate          time.Time       `json:"sZyouzyouOutouDay"`             // 上場応答日
	ExpireDate2Type      string          `json:"sNiruiKizituC"`                 // 二類期日C
	LargeUnit            float64         `json:"sOogutiKabusu,string"`          // 大口株数
	LargeAmount          float64         `json:"sOogutiKingmaker,string"`       // 大口金額
	OutputTicketType     string          `json:"sBadenpyouOutputYNC"`           // 場伝票出力有無C
	DepositAmount        float64         `json:"sHosyoukinDaiyouKakeme,string"` // 保証金代用掛目
	DepositValuation     float64         `json:"sDaiyouHyoukaTanka,string"`     // 代用証券評価単価
	OrganizationType     string          `json:"sKikoSankaC"`                   // 機構参加C
	ProvisionalType      string          `json:"sKarikessaiC"`                  // 仮決済C
	PrimaryExchange      Exchange        `json:"sYusenSizyou"`                  // 優先市場
	IndefinitePeriodType string          `json:"sMukigenC"`                     // 無期限対象C
	IndustryCode         string          `json:"sGyousyuCode"`                  // 業種コード
	IndustryName         string          `json:"sGyousyuName"`                  // 業種コード名
	SORTargetType        string          `json:"sSorC"`                         // SOR対象銘柄C
	CreateDateTime       time.Time       `json:"sCreateDate"`                   // 作成日時
	UpdateDateTime       time.Time       `json:"sUpdateDate"`                   // 更新日時
	UpdateNumber         string          `json:"sUpdateNumber"`                 // 更新通番
}

// StockMaster - 株式銘柄マスタ
func (c *client) StockMaster(ctx context.Context, session *Session, req StockMasterRequest) (*StockMasterResponse, error) {
	if session == nil {
		return nil, NilArgumentErr
	}
	session.mtx.Lock()
	defer session.mtx.Unlock()

	session.lastRequestNo++
	r := req.request(session.lastRequestNo, c.clock.Now())

	b, err := c.requester.get(ctx, session.RequestURL, r)
	if err != nil {
		return nil, err
	}
	var res stockMasterResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, fmt.Errorf("%s: %w", err, UnmarshalFailedErr)
	}

	Res := res.response()
	return &Res, nil
}
