package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

// StockWalletRequest - 買余力リクエスト
type StockWalletRequest struct {
	SymbolCode string   // 銘柄コード
	Exchange   Exchange // 市場
}

func (r *StockWalletRequest) request(no int64, now time.Time) stockWalletRequest {
	return stockWalletRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			FeatureType:    FeatureTypeStockWallet,
			ResponseFormat: commonResponseFormat,
		},
		SymbolCode: r.SymbolCode,
		Exchange:   r.Exchange,
	}
}

type stockWalletRequest struct {
	commonRequest
	SymbolCode string   `json:"sIssueCode"` // 銘柄コード
	Exchange   Exchange `json:"sSizyouC"`   // 市場
}

type stockWalletResponse struct {
	commonResponse
	SymbolCode     string     `json:"sIssueCode"`                          // 銘柄コード
	Exchange       Exchange   `json:"sSizyouC"`                            // 市場
	ResultCode     string     `json:"sResultCode"`                         // 結果コード
	ResultText     string     `json:"sResultText"`                         // 結果テキスト
	WarningCode    string     `json:"sWarningCode"`                        // 警告コード
	WarningText    string     `json:"sWarningText"`                        // 警告テキスト
	UpdateDateTime YmdHm      `json:"sSummaryUpdate"`                      // 更新日時
	StockWallet    float64    `json:"sSummaryGenkabuKaituke,string"`       // 株式現物買付可能額
	NisaWallet     float64    `json:"sSummaryNisaKaitukeKanougaku,string"` // NISA口座買付可能額
	Shortage       NumberBool `json:"sHusokukinHasseiFlg"`                 // 不足金発生フラグ
}

func (r *stockWalletResponse) UnmarshalJSON(b []byte) error {
	// 文字列でないところに空文字を返されることがあるので、置換しておく
	replaced := b
	replaces := map[string]string{
		`"sSummaryGenkabuKaituke":""`:       `"sSummaryGenkabuKaituke":"0"`,
		`"sSummaryNisaKaitukeKanougaku":""`: `"sSummaryNisaKaitukeKanougaku":"0"`,
	}
	for o, n := range replaces {
		replaced = bytes.Replace(replaced, []byte(o), []byte(n), -1)
	}

	type alias stockWalletResponse
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *stockWalletResponse) response() StockWalletResponse {
	return StockWalletResponse{
		CommonResponse: r.commonResponse.response(),
		SymbolCode:     r.SymbolCode,
		Exchange:       r.Exchange,
		ResultCode:     r.ResultCode,
		ResultText:     r.ResultText,
		WarningCode:    r.WarningCode,
		WarningText:    r.WarningText,
		UpdateDateTime: r.UpdateDateTime.Time,
		StockWallet:    r.StockWallet,
		NisaWallet:     r.NisaWallet,
		Shortage:       r.Shortage.Bool(),
	}
}

// StockWalletResponse - 買余力レスポンス
type StockWalletResponse struct {
	CommonResponse
	SymbolCode     string    // 銘柄コード
	Exchange       Exchange  // 市場
	ResultCode     string    // 結果コード
	ResultText     string    // 結果テキスト
	WarningCode    string    // 警告コード
	WarningText    string    // 警告テキスト
	UpdateDateTime time.Time // 更新日時
	StockWallet    float64   // 株式現物買付可能額
	NisaWallet     float64   // NISA口座買付可能額
	Shortage       bool      // 不足金発生フラグ
}

// StockWallet - 買余力
func (c *client) StockWallet(ctx context.Context, session *Session, req StockWalletRequest) (*StockWalletResponse, error) {
	if session == nil {
		return nil, NilArgumentErr
	}
	session.mtx.Lock()
	defer session.mtx.Unlock()

	session.lastRequestNo++
	r := req.request(session.lastRequestNo, c.clock.Now())

	var res stockWalletResponse
	if err := c.get(ctx, session.RequestURL, r, &res); err != nil {
		return nil, err
	}

	Res := res.response()
	return &Res, nil
}
