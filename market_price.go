package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// MarketPriceColumn - 時価関連情報カラム
type MarketPriceColumn string

// AllMarketPriceColumns - 時価関連情報の全カラム
var AllMarketPriceColumns = []MarketPriceColumn{
	MarketPriceColumnSection,
	MarketPriceColumnCurrentPrice,
	MarketPriceColumnCurrentPriceTime,
	MarketPriceColumnChangePriceType,
	MarketPriceColumnPrevDayRatio,
	MarketPriceColumnPrevDayPercent,
	MarketPriceColumnOpenPrice,
	MarketPriceColumnOpenPriceTime,
	MarketPriceColumnHighPrice,
	MarketPriceColumnHighPriceTime,
	MarketPriceColumnLowPrice,
	MarketPriceColumnLowPriceTime,
	MarketPriceColumnVolume,
	MarketPriceColumnAskSign,
	MarketPriceColumnAskPrice,
	MarketPriceColumnAskQuantity,
	MarketPriceColumnBidSign,
	MarketPriceColumnBidPrice,
	MarketPriceColumnBidQuantity,
	MarketPriceColumnExRightType,
	MarketPriceColumnDiscontinuityType,
	MarketPriceColumnStopHigh,
	MarketPriceColumnStopLow,
	MarketPriceColumnTradingAmount,
	MarketPriceColumnAskQuantityMarket,
	MarketPriceColumnBidQuantityMarket,
	MarketPriceColumnAskQuantityOver,
	MarketPriceColumnAskQuantity10,
	MarketPriceColumnAskPrice10,
	MarketPriceColumnAskQuantity9,
	MarketPriceColumnAskPrice9,
	MarketPriceColumnAskQuantity8,
	MarketPriceColumnAskPrice8,
	MarketPriceColumnAskQuantity7,
	MarketPriceColumnAskPrice7,
	MarketPriceColumnAskQuantity6,
	MarketPriceColumnAskPrice6,
	MarketPriceColumnAskQuantity5,
	MarketPriceColumnAskPrice5,
	MarketPriceColumnAskQuantity4,
	MarketPriceColumnAskPrice4,
	MarketPriceColumnAskQuantity3,
	MarketPriceColumnAskPrice3,
	MarketPriceColumnAskQuantity2,
	MarketPriceColumnAskPrice2,
	MarketPriceColumnAskQuantity1,
	MarketPriceColumnAskPrice1,
	MarketPriceColumnBidQuantity1,
	MarketPriceColumnBidPrice1,
	MarketPriceColumnBidQuantity2,
	MarketPriceColumnBidPrice2,
	MarketPriceColumnBidQuantity3,
	MarketPriceColumnBidPrice3,
	MarketPriceColumnBidQuantity4,
	MarketPriceColumnBidPrice4,
	MarketPriceColumnBidQuantity5,
	MarketPriceColumnBidPrice5,
	MarketPriceColumnBidQuantity6,
	MarketPriceColumnBidPrice6,
	MarketPriceColumnBidQuantity7,
	MarketPriceColumnBidPrice7,
	MarketPriceColumnBidQuantity8,
	MarketPriceColumnBidPrice8,
	MarketPriceColumnBidQuantity9,
	MarketPriceColumnBidPrice9,
	MarketPriceColumnBidQuantity10,
	MarketPriceColumnBidPrice10,
	MarketPriceColumnBidQuantityUnder,
	MarketPriceColumnVWAP,
	MarketPriceColumnPRP,
}

const (
	MarketPriceColumnSection           MarketPriceColumn = "xLISS"  // 所属
	MarketPriceColumnCurrentPrice      MarketPriceColumn = "pDPP"   // 現在値
	MarketPriceColumnCurrentPriceTime  MarketPriceColumn = "tDPP:T" // 現在値時刻
	MarketPriceColumnChangePriceType   MarketPriceColumn = "pDPG"   // 現値前値比較
	MarketPriceColumnPrevDayRatio      MarketPriceColumn = "pDYWP"  // 前日比
	MarketPriceColumnPrevDayPercent    MarketPriceColumn = "pDYRP"  // 騰落率
	MarketPriceColumnOpenPrice         MarketPriceColumn = "pDOP"   // 始値
	MarketPriceColumnOpenPriceTime     MarketPriceColumn = "tDOP:T" // 始値時刻
	MarketPriceColumnHighPrice         MarketPriceColumn = "pDHP"   // 高値
	MarketPriceColumnHighPriceTime     MarketPriceColumn = "tDHP:T" // 高値時刻
	MarketPriceColumnLowPrice          MarketPriceColumn = "pDLP"   // 安値
	MarketPriceColumnLowPriceTime      MarketPriceColumn = "tDLP:T" // 安値時刻
	MarketPriceColumnVolume            MarketPriceColumn = "pDV"    // 出来高
	MarketPriceColumnAskSign           MarketPriceColumn = "pQAS"   // 売気配値種類
	MarketPriceColumnAskPrice          MarketPriceColumn = "pQAP"   // 売気配値
	MarketPriceColumnAskQuantity       MarketPriceColumn = "pAV"    // 売気配数量
	MarketPriceColumnBidSign           MarketPriceColumn = "pQBS"   // 買気配値種類
	MarketPriceColumnBidPrice          MarketPriceColumn = "pQBP"   // 買気配値
	MarketPriceColumnBidQuantity       MarketPriceColumn = "pBV"    // 買気配数量
	MarketPriceColumnExRightType       MarketPriceColumn = "xDVES"  // 配当落銘柄区分
	MarketPriceColumnDiscontinuityType MarketPriceColumn = "xDCFS"  // 不連続要因銘柄区分
	MarketPriceColumnStopHigh          MarketPriceColumn = "pDHF"   // 日通し高値フラグ
	MarketPriceColumnStopLow           MarketPriceColumn = "pDLF"   // 日通し安値フラグ
	MarketPriceColumnTradingAmount     MarketPriceColumn = "pDJ"    // 売買代金
	MarketPriceColumnAskQuantityMarket MarketPriceColumn = "pAAV"   // 売数量(成行)
	MarketPriceColumnBidQuantityMarket MarketPriceColumn = "pABV"   // 買数量(成行)
	MarketPriceColumnAskQuantityOver   MarketPriceColumn = "pQOV"   // 売-OVER
	MarketPriceColumnAskQuantity10     MarketPriceColumn = "pGAV10" // 売-10-数量
	MarketPriceColumnAskPrice10        MarketPriceColumn = "pGAP10" // 売-10-値段
	MarketPriceColumnAskQuantity9      MarketPriceColumn = "pGAV9"  // 売-9-数量
	MarketPriceColumnAskPrice9         MarketPriceColumn = "pGAP9"  // 売-9-値段
	MarketPriceColumnAskQuantity8      MarketPriceColumn = "pGAV8"  // 売-8-数量
	MarketPriceColumnAskPrice8         MarketPriceColumn = "pGAP8"  // 売-8-値段
	MarketPriceColumnAskQuantity7      MarketPriceColumn = "pGAV7"  // 売-7-数量
	MarketPriceColumnAskPrice7         MarketPriceColumn = "pGAP7"  // 売-7-値段
	MarketPriceColumnAskQuantity6      MarketPriceColumn = "pGAV6"  // 売-6-数量
	MarketPriceColumnAskPrice6         MarketPriceColumn = "pGAP6"  // 売-6-値段
	MarketPriceColumnAskQuantity5      MarketPriceColumn = "pGAV5"  // 売-5-数量
	MarketPriceColumnAskPrice5         MarketPriceColumn = "pGAP5"  // 売-5-値段
	MarketPriceColumnAskQuantity4      MarketPriceColumn = "pGAV4"  // 売-4-数量
	MarketPriceColumnAskPrice4         MarketPriceColumn = "pGAP4"  // 売-4-値段
	MarketPriceColumnAskQuantity3      MarketPriceColumn = "pGAV3"  // 売-3-数量
	MarketPriceColumnAskPrice3         MarketPriceColumn = "pGAP3"  // 売-3-値段
	MarketPriceColumnAskQuantity2      MarketPriceColumn = "pGAV2"  // 売-2-数量
	MarketPriceColumnAskPrice2         MarketPriceColumn = "pGAP2"  // 売-2-値段
	MarketPriceColumnAskQuantity1      MarketPriceColumn = "pGAV1"  // 売-1-数量
	MarketPriceColumnAskPrice1         MarketPriceColumn = "pGAP1"  // 売-1-値段
	MarketPriceColumnBidQuantity1      MarketPriceColumn = "pGBV1"  // 買-1-数量
	MarketPriceColumnBidPrice1         MarketPriceColumn = "pGBP1"  // 買-1-値段
	MarketPriceColumnBidQuantity2      MarketPriceColumn = "pGBV2"  // 買-2-数量
	MarketPriceColumnBidPrice2         MarketPriceColumn = "pGBP2"  // 買-2-値段
	MarketPriceColumnBidQuantity3      MarketPriceColumn = "pGBV3"  // 買-3-数量
	MarketPriceColumnBidPrice3         MarketPriceColumn = "pGBP3"  // 買-3-値段
	MarketPriceColumnBidQuantity4      MarketPriceColumn = "pGBV4"  // 買-4-数量
	MarketPriceColumnBidPrice4         MarketPriceColumn = "pGBP4"  // 買-4-値段
	MarketPriceColumnBidQuantity5      MarketPriceColumn = "pGBV5"  // 買-5-数量
	MarketPriceColumnBidPrice5         MarketPriceColumn = "pGBP5"  // 買-5-値段
	MarketPriceColumnBidQuantity6      MarketPriceColumn = "pGBV6"  // 買-6-数量
	MarketPriceColumnBidPrice6         MarketPriceColumn = "pGBP6"  // 買-6-値段
	MarketPriceColumnBidQuantity7      MarketPriceColumn = "pGBV7"  // 買-7-数量
	MarketPriceColumnBidPrice7         MarketPriceColumn = "pGBP7"  // 買-7-値段
	MarketPriceColumnBidQuantity8      MarketPriceColumn = "pGBV8"  // 買-8-数量
	MarketPriceColumnBidPrice8         MarketPriceColumn = "pGBP8"  // 買-8-値段
	MarketPriceColumnBidQuantity9      MarketPriceColumn = "pGBV9"  // 買-9-数量
	MarketPriceColumnBidPrice9         MarketPriceColumn = "pGBP9"  // 買-9-値段
	MarketPriceColumnBidQuantity10     MarketPriceColumn = "pGBV10" // 買-10-数量
	MarketPriceColumnBidPrice10        MarketPriceColumn = "pGBP10" // 買-10-値段
	MarketPriceColumnBidQuantityUnder  MarketPriceColumn = "pQUV"   // 買-UNDER
	MarketPriceColumnVWAP              MarketPriceColumn = "pVWAP"  // VWAP
	MarketPriceColumnPRP               MarketPriceColumn = "pPRP"   // 前日終値
)

// MarketPriceRequest - 時価関連情報リクエスト
type MarketPriceRequest struct {
	IssueCodes []string            // 取得したい銘柄コード
	Columns    []MarketPriceColumn // 取得したい情報
}

func (r *MarketPriceRequest) request(no int64, now time.Time) marketPriceRequest {
	if r.IssueCodes == nil {
		r.IssueCodes = []string{}
	}

	if r.Columns == nil || len(r.Columns) == 0 {
		r.Columns = AllMarketPriceColumns
	}
	columns := make([]string, len(r.Columns))
	for i, column := range r.Columns {
		columns[i] = string(column)
	}

	return marketPriceRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			FeatureType:    FeatureTypeMarketPrice,
			ResponseFormat: commonResponseFormat,
		},
		IssueCodes: strings.Join(r.IssueCodes, ","),
		Columns:    strings.Join(columns, ","),
	}
}

type marketPriceRequest struct {
	commonRequest
	IssueCodes string `json:"sTargetIssueCode,omitempty"` // 取得したい銘柄コード
	Columns    string `json:"sTargetColumn,omitempty"`    // 取得したい情報
}

type marketPriceResponse struct {
	commonResponse
	MarketPrices []marketPrice `json:"aCLMMfdsMarketPrice"` // 時価情報
}

func (r *marketPriceResponse) UnmarshalJSON(b []byte) error {
	// 文字列でないところに空文字を返されることがあるので、置換しておく
	replaced := b
	replaces := map[string]string{
		`"aCLMMfdsMarketPrice":""`: `"aCLMMfdsMarketPrice":[]`,
	}
	for o, n := range replaces {
		replaced = bytes.Replace(replaced, []byte(o), []byte(n), -1)
	}

	type alias marketPriceResponse
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *marketPriceResponse) response() MarketPriceResponse {
	marketPrices := make([]MarketPrice, 0)
	for _, price := range r.MarketPrices {
		if price.IssueCode == "" {
			continue
		}
		marketPrices = append(marketPrices, price.response())
	}

	return MarketPriceResponse{
		CommonResponse: r.commonResponse.response(),
		MarketPrices:   marketPrices,
	}
}

type marketPrice struct {
	commonResponse
	IssueCode         string              `json:"sIssueCode"`    // 銘柄コード
	Section           string              `json:"xLISS"`         // 所属
	CurrentPrice      float64             `json:"pDPP,string"`   // 現在値
	CurrentPriceTime  Hm                  `json:"tDPP:T"`        // 現在値時刻
	ChangePriceType   ChangePriceType     `json:"pDPG"`          // 現値前値比較
	PrevDayRatio      float64             `json:"pDYWP,string"`  // 前日比
	PrevDayPercent    float64             `json:"pDYRP,string"`  // 騰落率
	OpenPrice         float64             `json:"pDOP,string"`   // 始値
	OpenPriceTime     Hm                  `json:"tDOP:T"`        // 始値時刻
	HighPrice         float64             `json:"pDHP,string"`   // 高値
	HighPriceTime     Hm                  `json:"tDHP:T"`        // 高値時刻
	LowPrice          float64             `json:"pDLP,string"`   // 安値
	LowPriceTime      Hm                  `json:"tDLP:T"`        // 安値時刻
	Volume            float64             `json:"pDV,string"`    // 出来高
	AskSign           IndicationPriceType `json:"pQAS"`          // 売気配値種類
	AskPrice          float64             `json:"pQAP,string"`   // 売気配値
	AskQuantity       float64             `json:"pAV,string"`    // 売気配数量
	BidSign           IndicationPriceType `json:"pQBS"`          // 買気配値種類
	BidPrice          float64             `json:"pQBP,string"`   // 買気配値
	BidQuantity       float64             `json:"pBV,string"`    // 買気配数量
	ExRightType       string              `json:"xDVES"`         // 配当落銘柄区分
	DiscontinuityType string              `json:"xDCFS"`         // 不連続要因銘柄区分
	StopHigh          CurrentPriceType    `json:"pDHF"`          // 日通し高値フラグ
	StopLow           CurrentPriceType    `json:"pDLF"`          // 日通し安値フラグ
	TradingAmount     float64             `json:"pDJ,string"`    // 売買代金
	AskQuantityMarket float64             `json:"pAAV,string"`   // 売数量(成行)
	BidQuantityMarket float64             `json:"pABV,string"`   // 買数量(成行)
	AskQuantityOver   float64             `json:"pQOV,string"`   // 売-OVER
	AskQuantity10     float64             `json:"pGAV10,string"` // 売-10-数量
	AskPrice10        float64             `json:"pGAP10,string"` // 売-10-値段
	AskQuantity9      float64             `json:"pGAV9,string"`  // 売-9-数量
	AskPrice9         float64             `json:"pGAP9,string"`  // 売-9-値段
	AskQuantity8      float64             `json:"pGAV8,string"`  // 売-8-数量
	AskPrice8         float64             `json:"pGAP8,string"`  // 売-8-値段
	AskQuantity7      float64             `json:"pGAV7,string"`  // 売-7-数量
	AskPrice7         float64             `json:"pGAP7,string"`  // 売-7-値段
	AskQuantity6      float64             `json:"pGAV6,string"`  // 売-6-数量
	AskPrice6         float64             `json:"pGAP6,string"`  // 売-6-値段
	AskQuantity5      float64             `json:"pGAV5,string"`  // 売-5-数量
	AskPrice5         float64             `json:"pGAP5,string"`  // 売-5-値段
	AskQuantity4      float64             `json:"pGAV4,string"`  // 売-4-数量
	AskPrice4         float64             `json:"pGAP4,string"`  // 売-4-値段
	AskQuantity3      float64             `json:"pGAV3,string"`  // 売-3-数量
	AskPrice3         float64             `json:"pGAP3,string"`  // 売-3-値段
	AskQuantity2      float64             `json:"pGAV2,string"`  // 売-2-数量
	AskPrice2         float64             `json:"pGAP2,string"`  // 売-2-値段
	AskQuantity1      float64             `json:"pGAV1,string"`  // 売-1-数量
	AskPrice1         float64             `json:"pGAP1,string"`  // 売-1-値段
	BidQuantity1      float64             `json:"pGBV1,string"`  // 買-1-数量
	BidPrice1         float64             `json:"pGBP1,string"`  // 買-1-値段
	BidQuantity2      float64             `json:"pGBV2,string"`  // 買-2-数量
	BidPrice2         float64             `json:"pGBP2,string"`  // 買-2-値段
	BidQuantity3      float64             `json:"pGBV3,string"`  // 買-3-数量
	BidPrice3         float64             `json:"pGBP3,string"`  // 買-3-値段
	BidQuantity4      float64             `json:"pGBV4,string"`  // 買-4-数量
	BidPrice4         float64             `json:"pGBP4,string"`  // 買-4-値段
	BidQuantity5      float64             `json:"pGBV5,string"`  // 買-5-数量
	BidPrice5         float64             `json:"pGBP5,string"`  // 買-5-値段
	BidQuantity6      float64             `json:"pGBV6,string"`  // 買-6-数量
	BidPrice6         float64             `json:"pGBP6,string"`  // 買-6-値段
	BidQuantity7      float64             `json:"pGBV7,string"`  // 買-7-数量
	BidPrice7         float64             `json:"pGBP7,string"`  // 買-7-値段
	BidQuantity8      float64             `json:"pGBV8,string"`  // 買-8-数量
	BidPrice8         float64             `json:"pGBP8,string"`  // 買-8-値段
	BidQuantity9      float64             `json:"pGBV9,string"`  // 買-9-数量
	BidPrice9         float64             `json:"pGBP9,string"`  // 買-9-値段
	BidQuantity10     float64             `json:"pGBV10,string"` // 買-10-数量
	BidPrice10        float64             `json:"pGBP10,string"` // 買-10-値段
	BidQuantityUnder  float64             `json:"pQUV,string"`   // 買-UNDER
	VWAP              float64             `json:"pVWAP,string"`  // VWAP
	PRP               float64             `json:"pPRP,string"`   // 前日終値
}

func (r *marketPrice) UnmarshalJSON(b []byte) error {
	// 文字列でないところに空文字を返されることがあるので、置換しておく
	replaced := b
	replaces := map[string]string{
		`"pDPP":""`:   `"pDPP":"0"`,
		`"pDYWP":""`:  `"pDYWP":"0"`,
		`"pDYRP":""`:  `"pDYRP":"0"`,
		`"pDOP":""`:   `"pDOP":"0"`,
		`"pDHP":""`:   `"pDHP":"0"`,
		`"pDLP":""`:   `"pDLP":"0"`,
		`"pDV":""`:    `"pDV":"0"`,
		`"pQAP":""`:   `"pQAP":"0"`,
		`"pAV":""`:    `"pAV":"0"`,
		`"pQBP":""`:   `"pQBP":"0"`,
		`"pBV":""`:    `"pBV":"0"`,
		`"pDJ":""`:    `"pDJ":"0"`,
		`"pAAV":""`:   `"pAAV":"0"`,
		`"pABV":""`:   `"pABV":"0"`,
		`"pQOV":""`:   `"pQOV":"0"`,
		`"pGAV10":""`: `"pGAV10":"0"`,
		`"pGAP10":""`: `"pGAP10":"0"`,
		`"pGAV9":""`:  `"pGAV9":"0"`,
		`"pGAP9":""`:  `"pGAP9":"0"`,
		`"pGAV8":""`:  `"pGAV8":"0"`,
		`"pGAP8":""`:  `"pGAP8":"0"`,
		`"pGAV7":""`:  `"pGAV7":"0"`,
		`"pGAP7":""`:  `"pGAP7":"0"`,
		`"pGAV6":""`:  `"pGAV6":"0"`,
		`"pGAP6":""`:  `"pGAP6":"0"`,
		`"pGAV5":""`:  `"pGAV5":"0"`,
		`"pGAP5":""`:  `"pGAP5":"0"`,
		`"pGAV4":""`:  `"pGAV4":"0"`,
		`"pGAP4":""`:  `"pGAP4":"0"`,
		`"pGAV3":""`:  `"pGAV3":"0"`,
		`"pGAP3":""`:  `"pGAP3":"0"`,
		`"pGAV2":""`:  `"pGAV2":"0"`,
		`"pGAP2":""`:  `"pGAP2":"0"`,
		`"pGAV1":""`:  `"pGAV1":"0"`,
		`"pGAP1":""`:  `"pGAP1":"0"`,
		`"pGBV1":""`:  `"pGBV1":"0"`,
		`"pGBP1":""`:  `"pGBP1":"0"`,
		`"pGBV2":""`:  `"pGBV2":"0"`,
		`"pGBP2":""`:  `"pGBP2":"0"`,
		`"pGBV3":""`:  `"pGBV3":"0"`,
		`"pGBP3":""`:  `"pGBP3":"0"`,
		`"pGBV4":""`:  `"pGBV4":"0"`,
		`"pGBP4":""`:  `"pGBP4":"0"`,
		`"pGBV5":""`:  `"pGBV5":"0"`,
		`"pGBP5":""`:  `"pGBP5":"0"`,
		`"pGBV6":""`:  `"pGBV6":"0"`,
		`"pGBP6":""`:  `"pGBP6":"0"`,
		`"pGBV7":""`:  `"pGBV7":"0"`,
		`"pGBP7":""`:  `"pGBP7":"0"`,
		`"pGBV8":""`:  `"pGBV8":"0"`,
		`"pGBP8":""`:  `"pGBP8":"0"`,
		`"pGBV9":""`:  `"pGBV9":"0"`,
		`"pGBP9":""`:  `"pGBP9":"0"`,
		`"pGBV10":""`: `"pGBV10":"0"`,
		`"pGBP10":""`: `"pGBP10":"0"`,
		`"pQUV":""`:   `"pQUV":"0"`,
		`"pVWAP":""`:  `"pVWAP":"0"`,
		`"pPRP":""`:   `"pPRP":"0"`,
	}
	for o, n := range replaces {
		replaced = bytes.Replace(replaced, []byte(o), []byte(n), -1)
	}

	type alias marketPrice
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *marketPrice) response() MarketPrice {
	return MarketPrice{
		IssueCode:         r.IssueCode,
		Section:           r.Section,
		CurrentPrice:      r.CurrentPrice,
		CurrentPriceTime:  r.CurrentPriceTime.Time,
		ChangePriceType:   r.ChangePriceType,
		PrevDayRatio:      r.PrevDayRatio,
		PrevDayPercent:    r.PrevDayPercent,
		OpenPrice:         r.OpenPrice,
		OpenPriceTime:     r.OpenPriceTime.Time,
		HighPrice:         r.HighPrice,
		HighPriceTime:     r.HighPriceTime.Time,
		LowPrice:          r.LowPrice,
		LowPriceTime:      r.LowPriceTime.Time,
		Volume:            r.Volume,
		AskSign:           r.AskSign,
		AskPrice:          r.AskPrice,
		AskQuantity:       r.AskQuantity,
		BidSign:           r.BidSign,
		BidPrice:          r.BidPrice,
		BidQuantity:       r.BidQuantity,
		ExRightType:       r.ExRightType,
		DiscontinuityType: r.DiscontinuityType,
		StopHigh:          r.StopHigh,
		StopLow:           r.StopLow,
		TradingAmount:     r.TradingAmount,
		AskQuantityMarket: r.AskQuantityMarket,
		BidQuantityMarket: r.BidQuantityMarket,
		AskQuantityOver:   r.AskQuantityOver,
		AskQuantity10:     r.AskQuantity10,
		AskPrice10:        r.AskPrice10,
		AskQuantity9:      r.AskQuantity9,
		AskPrice9:         r.AskPrice9,
		AskQuantity8:      r.AskQuantity8,
		AskPrice8:         r.AskPrice8,
		AskQuantity7:      r.AskQuantity7,
		AskPrice7:         r.AskPrice7,
		AskQuantity6:      r.AskQuantity6,
		AskPrice6:         r.AskPrice6,
		AskQuantity5:      r.AskQuantity5,
		AskPrice5:         r.AskPrice5,
		AskQuantity4:      r.AskQuantity4,
		AskPrice4:         r.AskPrice4,
		AskQuantity3:      r.AskQuantity3,
		AskPrice3:         r.AskPrice3,
		AskQuantity2:      r.AskQuantity2,
		AskPrice2:         r.AskPrice2,
		AskQuantity1:      r.AskQuantity1,
		AskPrice1:         r.AskPrice1,
		BidQuantity1:      r.BidQuantity1,
		BidPrice1:         r.BidPrice1,
		BidQuantity2:      r.BidQuantity2,
		BidPrice2:         r.BidPrice2,
		BidQuantity3:      r.BidQuantity3,
		BidPrice3:         r.BidPrice3,
		BidQuantity4:      r.BidQuantity4,
		BidPrice4:         r.BidPrice4,
		BidQuantity5:      r.BidQuantity5,
		BidPrice5:         r.BidPrice5,
		BidQuantity6:      r.BidQuantity6,
		BidPrice6:         r.BidPrice6,
		BidQuantity7:      r.BidQuantity7,
		BidPrice7:         r.BidPrice7,
		BidQuantity8:      r.BidQuantity8,
		BidPrice8:         r.BidPrice8,
		BidQuantity9:      r.BidQuantity9,
		BidPrice9:         r.BidPrice9,
		BidQuantity10:     r.BidQuantity10,
		BidPrice10:        r.BidPrice10,
		BidQuantityUnder:  r.BidQuantityUnder,
		VWAP:              r.VWAP,
		PRP:               r.PRP,
	}
}

// MarketPriceResponse - 時価関連情報レスポンス
type MarketPriceResponse struct {
	CommonResponse
	MarketPrices []MarketPrice // 時価情報
}

// MarketPrice - 時価関連情報
type MarketPrice struct {
	IssueCode         string              // 銘柄コード
	Section           string              // 所属
	CurrentPrice      float64             // 現在値
	CurrentPriceTime  time.Time           // 現在値時刻
	ChangePriceType   ChangePriceType     // 現値前値比較
	PrevDayRatio      float64             // 前日比
	PrevDayPercent    float64             // 騰落率
	OpenPrice         float64             // 始値
	OpenPriceTime     time.Time           // 始値時刻
	HighPrice         float64             // 高値
	HighPriceTime     time.Time           // 高値時刻
	LowPrice          float64             // 安値
	LowPriceTime      time.Time           // 安値時刻
	Volume            float64             // 出来高
	AskSign           IndicationPriceType // 売気配値種類
	AskPrice          float64             // 売気配値
	AskQuantity       float64             // 売気配数量
	BidSign           IndicationPriceType // 買気配値種類
	BidPrice          float64             // 買気配値
	BidQuantity       float64             // 買気配数量
	ExRightType       string              // 配当落銘柄区分
	DiscontinuityType string              // 不連続要因銘柄区分
	StopHigh          CurrentPriceType    // 日通し高値フラグ
	StopLow           CurrentPriceType    // 日通し安値フラグ
	TradingAmount     float64             // 売買代金
	AskQuantityMarket float64             // 売数量(成行)
	BidQuantityMarket float64             // 買数量(成行)
	AskQuantityOver   float64             // 売-OVER
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
	BidQuantityUnder  float64             // 買-UNDER
	VWAP              float64             // VWAP
	PRP               float64             // 前日終値
}

// MarketPrice - 時価関連情報
func (c *client) MarketPrice(ctx context.Context, session *Session, req MarketPriceRequest) (*MarketPriceResponse, error) {
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
	var res marketPriceResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, fmt.Errorf("%s: %w", err, UnmarshalFailedErr)
	}

	Res := res.response()
	return &Res, nil
}
