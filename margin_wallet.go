package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"time"
)

// MarginWalletRequest - 建余力&本日維持率リクエスト
type MarginWalletRequest struct {
	SymbolCode string   // 銘柄コード
	Exchange   Exchange // 市場
}

func (r *MarginWalletRequest) request(no int64, now time.Time) marginWalletRequest {
	return marginWalletRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			FeatureType:    FeatureTypeMarginWallet,
			ResponseFormat: commonResponseFormat,
		},
		SymbolCode: r.SymbolCode,
		Exchange:   r.Exchange,
	}
}

type marginWalletRequest struct {
	commonRequest
	SymbolCode string   `json:"sIssueCode"` // 銘柄コード
	Exchange   Exchange `json:"sSizyouC"`   // 市場
}

type marginWalletResponse struct {
	commonResponse
	SymbolCode     string     `json:"sIssueCode"`                     // 銘柄コード
	Exchange       Exchange   `json:"sSizyouC"`                       // 市場
	ResultCode     string     `json:"sResultCode"`                    // 結果コード
	ResultText     string     `json:"sResultText"`                    // 結果テキスト
	WarningCode    string     `json:"sWarningCode"`                   // 警告コード
	WarningText    string     `json:"sWarningText"`                   // 警告テキスト
	UpdateDateTime YmdHm      `json:"sSummaryUpdate"`                 // 更新日時
	MarginWallet   float64    `json:"sSummarySinyouSinkidate,string"` // 信用新規建可能額
	DepositRate    float64    `json:"sItakuhosyoukin,string"`         // 委託保証金率
	Shortage       NumberBool `json:"sOisyouKakuteiFlg"`              // 追証フラグ
}

func (r *marginWalletResponse) UnmarshalJSON(b []byte) error {
	// 文字列でないところに空文字を返されることがあるので、置換しておく
	replaced := b
	replaces := map[string]string{
		`"sSummarySinyouSinkidate":""`: `"sSummarySinyouSinkidate":"0"`,
		`"sItakuhosyoukin":""`:         `"sItakuhosyoukin":"0"`,
	}
	for o, n := range replaces {
		replaced = bytes.Replace(replaced, []byte(o), []byte(n), -1)
	}

	type alias marginWalletResponse
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *marginWalletResponse) response() MarginWalletResponse {
	return MarginWalletResponse{
		CommonResponse: r.commonResponse.response(),
		SymbolCode:     r.SymbolCode,
		Exchange:       r.Exchange,
		ResultCode:     r.ResultCode,
		ResultText:     r.ResultText,
		WarningCode:    r.WarningCode,
		WarningText:    r.WarningText,
		UpdateDateTime: r.UpdateDateTime.Time,
		MarginWallet:   r.MarginWallet,
		DepositRate:    r.DepositRate,
		Shortage:       r.Shortage.Bool(),
	}
}

// MarginWalletResponse - 建余力&本日維持率レスポンス
type MarginWalletResponse struct {
	CommonResponse
	SymbolCode     string    // 銘柄コード
	Exchange       Exchange  // 市場
	ResultCode     string    // 結果コード
	ResultText     string    // 結果テキスト
	WarningCode    string    // 警告コード
	WarningText    string    // 警告テキスト
	UpdateDateTime time.Time // 更新日時
	MarginWallet   float64   // 信用新規建可能額
	DepositRate    float64   // 委託保証金率
	Shortage       bool      // 追証フラグ
}

// MarginWallet - 建余力&本日維持率
func (c *client) MarginWallet(ctx context.Context, session *Session, req MarginWalletRequest) (*MarginWalletResponse, error) {
	if session == nil {
		return nil, NilArgumentErr
	}
	session.mtx.Lock()
	defer session.mtx.Unlock()

	session.lastRequestNo++
	r := req.request(session.lastRequestNo, c.clock.Now())

	var res marginWalletResponse
	if err := c.get(ctx, session.RequestURL, r, &res); err != nil {
		return nil, err
	}

	Res := res.response()
	return &Res, nil
}
