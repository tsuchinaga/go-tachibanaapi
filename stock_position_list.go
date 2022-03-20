package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// StockPositionListRequest - 現物保有銘柄一覧リクエスト
type StockPositionListRequest struct {
	IssueCode string // 銘柄コード
}

func (r *StockPositionListRequest) request(no int64, now time.Time) stockPositionListRequest {
	return stockPositionListRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			MessageType:    MessageTypeStockPositionList,
			ResponseFormat: commonResponseFormat,
		},
		IssueCode: r.IssueCode,
	}
}

type stockPositionListRequest struct {
	commonRequest
	IssueCode string `json:"sIssueCode,omitempty"` // 銘柄コード
}

type stockPositionListResponse struct {
	commonResponse
	IssueCode      string          `json:"sIssueCode,omitempty"`                    // 銘柄コード
	ResultCode     string          `json:"sResultCode"`                             // 結果コード
	ResultText     string          `json:"sResultText"`                             // 結果テキスト
	WarningCode    string          `json:"sWarningCode"`                            // 警告コード
	WarningText    string          `json:"sWarningText"`                            // 警告テキスト
	SpecificAmount float64         `json:"sTokuteiGaisanHyoukagakuGoukei,string"`   // 概算評価額合計(特定口座残高)
	GeneralAmount  float64         `json:"sIppanGaisanHyoukagakuGoukei,string"`     // 概算評価額合計(一般口座残高)
	NisaAmount     float64         `json:"sNisaGaisanHyoukagakuGoukei,string"`      // 概算評価額合計(NISA口座残高)
	TotalAmount    float64         `json:"sTotalGaisanHyoukagakuGoukei,string"`     // 残高合計 概算評価額合計
	SpecificProfit float64         `json:"sTokuteiGaisanHyoukaSonekiGoukei,string"` // 概算評価損益合計(特定口座残高)
	GeneralProfit  float64         `json:"sIppanGaisanHyoukaSonekiGoukei,string"`   // 概算評価損益合計(一般口座残高)
	NisaProfit     float64         `json:"sNisaGaisanHyoukaSonekiGoukei,string"`    // 概算評価損益合計(NISA口座残高)
	TotalProfit    float64         `json:"sTotalGaisanHyoukaSonekiGoukei,string"`   // 概算評価損益合計(残高合計)
	Positions      []stockPosition `json:"aGenbutuKabuList"`                        // 現物株リスト
}

func (r *stockPositionListResponse) UnmarshalJSON(b []byte) error {
	// 文字列でないところに空文字を返されることがあるので、置換しておく
	replaced := b
	replaces := map[string]string{
		`"aGenbutuKabuList":""`: `"aGenbutuKabuList":[]`,
	}
	for o, n := range replaces {
		replaced = bytes.Replace(replaced, []byte(o), []byte(n), -1)
	}

	type alias stockPositionListResponse
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *stockPositionListResponse) response() StockPositionListResponse {
	positions := make([]StockPosition, len(r.Positions))
	for i, p := range r.Positions {
		positions[i] = p.response()
	}

	return StockPositionListResponse{
		CommonResponse: r.commonResponse.response(),
		IssueCode:      r.IssueCode,
		ResultCode:     r.ResultCode,
		ResultText:     r.ResultText,
		WarningCode:    r.WarningCode,
		WarningText:    r.WarningText,
		SpecificAmount: r.SpecificAmount,
		GeneralAmount:  r.GeneralAmount,
		NisaAmount:     r.NisaAmount,
		TotalAmount:    r.TotalAmount,
		SpecificProfit: r.SpecificProfit,
		GeneralProfit:  r.GeneralProfit,
		NisaProfit:     r.NisaProfit,
		TotalProfit:    r.TotalProfit,
		Positions:      positions,
	}
}

type stockPosition struct {
	WarningCode        string             `json:"sUriOrderWarningCode"`                   // 警告コード
	WarningText        string             `json:"sUriOrderWarningText"`                   // 警告テキスト
	IssueCode          string             `json:"sUriOrderIssueCode"`                     // 銘柄コード
	AccountType        AccountType        `json:"sUriOrderZyoutoekiKazeiC"`               // 口座
	OwnedQuantity      float64            `json:"sUriOrderZanKabuSuryou,string"`          // 残高株数
	UnHoldQuantity     float64            `json:"sUriOrderUritukeKanouSuryou,string"`     // 売付可能株数
	UnitValuation      float64            `json:"sUriOrderGaisanBokaTanka,string"`        // 評価単価
	BookValuation      float64            `json:"sUriOrderHyoukaTanka,string"`            // 概算簿価単価
	TotalValuation     float64            `json:"sUriOrderGaisanHyoukagaku,string"`       // 評価金額
	Profit             float64            `json:"sUriOrderGaisanHyoukaSoneki,string"`     // 評価損益
	ProfitRatio        float64            `json:"sUriOrderGaisanHyoukaSonekiRitu,string"` // 評価損益率
	PrevClosePrice     float64            `json:"sSyuzituOwarine,string"`                 // 前日終値
	PrevCloseRatio     float64            `json:"sZenzituHi,string"`                      // 前日比
	PrevClosePercent   float64            `json:"sZenzituHiPer,string"`                   // 前日比(%)
	PrevCloseRatioType PrevCloseRatioType `json:"sUpDownFlag"`                            // 騰落率Flag
	MarginBalance      float64            `json:"sNissyoukinKasikabuZan,string"`          // 証金貸株残
}

func (r *stockPosition) response() StockPosition {
	return StockPosition{
		WarningCode:        r.WarningCode,
		WarningText:        r.WarningText,
		IssueCode:          r.IssueCode,
		AccountType:        r.AccountType,
		OwnedQuantity:      r.OwnedQuantity,
		UnHoldQuantity:     r.UnHoldQuantity,
		UnitValuation:      r.UnitValuation,
		BookValuation:      r.BookValuation,
		TotalValuation:     r.TotalValuation,
		Profit:             r.Profit,
		ProfitRatio:        r.ProfitRatio,
		PrevClosePrice:     r.PrevClosePrice,
		PrevCloseRatio:     r.PrevCloseRatio,
		PrevClosePercent:   r.PrevClosePercent,
		PrevCloseRatioType: r.PrevCloseRatioType,
		MarginBalance:      r.MarginBalance,
	}
}

// StockPositionListResponse - 現物保有銘柄一覧レスポンス
type StockPositionListResponse struct {
	CommonResponse
	IssueCode      string          // 銘柄コード
	ResultCode     string          // 結果コード
	ResultText     string          // 結果テキスト
	WarningCode    string          // 警告コード
	WarningText    string          // 警告テキスト
	SpecificAmount float64         // 概算評価額合計(特定口座残高)
	GeneralAmount  float64         // 概算評価額合計(一般口座残高)
	NisaAmount     float64         // 概算評価額合計(NISA口座残高)
	TotalAmount    float64         // 残高合計 概算評価額合計
	SpecificProfit float64         // 概算評価損益合計(特定口座残高)
	GeneralProfit  float64         // 概算評価損益合計(一般口座残高)
	NisaProfit     float64         // 概算評価損益合計(NISA口座残高)
	TotalProfit    float64         // 概算評価損益合計(残高合計)
	Positions      []StockPosition // 現物株リスト
}

// StockPosition - 現物株リスト
type StockPosition struct {
	WarningCode        string             // 警告コード
	WarningText        string             // 警告テキスト
	IssueCode          string             // 銘柄コード
	AccountType        AccountType        // 口座
	OwnedQuantity      float64            // 残高株数
	UnHoldQuantity     float64            // 売付可能株数
	UnitValuation      float64            // 評価単価
	BookValuation      float64            // 概算簿価単価
	TotalValuation     float64            // 評価金額
	Profit             float64            // 評価損益
	ProfitRatio        float64            // 評価損益率
	PrevClosePrice     float64            // 前日終値
	PrevCloseRatio     float64            // 前日比
	PrevClosePercent   float64            // 前日比(%)
	PrevCloseRatioType PrevCloseRatioType // 騰落率Flag
	MarginBalance      float64            // 証金貸株残
}

// StockPositionList - 現物株リスト
func (c *client) StockPositionList(ctx context.Context, session *Session, req StockPositionListRequest) (*StockPositionListResponse, error) {
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
	var res stockPositionListResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, fmt.Errorf("%s: %w", err, UnmarshalFailedErr)
	}

	Res := res.response()
	return &Res, nil
}
