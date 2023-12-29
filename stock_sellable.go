package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// StockSellableRequest - 売却可能数量リクエスト
type StockSellableRequest struct {
	IssueCode string // 銘柄コード
}

func (r *StockSellableRequest) request(no int64, now time.Time) stockSellableRequest {
	return stockSellableRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			MessageType:    MessageTypeStockSellable,
			ResponseFormat: commonResponseFormat,
		},
		IssueCode: r.IssueCode,
	}
}

type stockSellableRequest struct {
	commonRequest
	IssueCode string `json:"sIssueCode"` // 銘柄コード
}

type stockSellableResponse struct {
	commonResponse
	IssueCode        string  `json:"sIssueCode"`                            // 銘柄コード
	ResultCode       string  `json:"sResultCode"`                           // 結果コード
	ResultText       string  `json:"sResultText"`                           // 結果テキスト
	WarningCode      string  `json:"sWarningCode"`                          // 警告コード
	WarningText      string  `json:"sWarningText"`                          // 警告テキスト
	UpdateDateTime   YmdHm   `json:"sSummaryUpdate"`                        // 更新日時
	GeneralQuantity  float64 `json:"sZanKabuSuryouUriKanouIppan,string"`    // 売付可能株数(一般)
	SpecificQuantity float64 `json:"sZanKabuSuryouUriKanouTokutei,string"`  // 売付可能株数(特定)
	NisaQuantity     float64 `json:"sZanKabuSuryouUriKanouNisa,string"`     // 売付可能株数(NISA)
	GrowthQuantity   float64 `json:"sZanKabuSuryouUriKanouNseityou,string"` // 売付可能株数(N成長)
}

func (r *stockSellableResponse) UnmarshalJSON(b []byte) error {
	// 文字列でないところに空文字を返されることがあるので、置換しておく
	replaced := b
	replaces := map[string]string{
		`"sZanKabuSuryouUriKanouIppan":""`:   `"sZanKabuSuryouUriKanouIppan":"0"`,
		`"sZanKabuSuryouUriKanouTokutei":""`: `"sZanKabuSuryouUriKanouTokutei":"0"`,
		`"sZanKabuSuryouUriKanouNisa":""`:    `"sZanKabuSuryouUriKanouNisa":"0"`,
	}
	for o, n := range replaces {
		replaced = bytes.Replace(replaced, []byte(o), []byte(n), -1)
	}

	type alias stockSellableResponse
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *stockSellableResponse) response() StockSellableResponse {
	return StockSellableResponse{
		CommonResponse:   r.commonResponse.response(),
		IssueCode:        r.IssueCode,
		ResultCode:       r.ResultCode,
		ResultText:       r.ResultText,
		WarningCode:      r.WarningCode,
		WarningText:      r.WarningText,
		UpdateDateTime:   r.UpdateDateTime.Time,
		GeneralQuantity:  r.GeneralQuantity,
		SpecificQuantity: r.SpecificQuantity,
		NisaQuantity:     r.NisaQuantity,
	}
}

// StockSellableResponse - 売却可能数量レスポンス
type StockSellableResponse struct {
	CommonResponse
	IssueCode        string    // 銘柄コード
	ResultCode       string    // 結果コード
	ResultText       string    // 結果テキスト
	WarningCode      string    // 警告コード
	WarningText      string    // 警告テキスト
	UpdateDateTime   time.Time // 更新日時
	GeneralQuantity  float64   // 売付可能株数(一般)
	SpecificQuantity float64   // 売付可能株数(特定)
	NisaQuantity     float64   // 売付可能株数(NISA)
}

// StockSellable - 売却可能数量
func (c *client) StockSellable(ctx context.Context, session *Session, req StockSellableRequest) (*StockSellableResponse, error) {
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
	var res stockSellableResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, fmt.Errorf("%s: %w", err, UnmarshalFailedErr)
	}

	Res := res.response()
	return &Res, nil
}
