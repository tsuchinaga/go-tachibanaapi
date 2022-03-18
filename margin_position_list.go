package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// MarginPositionListRequest - 信用建玉一覧リクエスト
type MarginPositionListRequest struct {
	IssueCode string // 銘柄コード
}

func (r *MarginPositionListRequest) request(no int64, now time.Time) marginPositionListRequest {
	return marginPositionListRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			FeatureType:    FeatureTypeMarginPositionList,
			ResponseFormat: commonResponseFormat,
		},
		IssueCode: r.IssueCode,
	}
}

type marginPositionListRequest struct {
	commonRequest
	IssueCode string `json:"sIssueCode,omitempty"` // 銘柄コード
}

type marginPositionListResponse struct {
	commonResponse
	IssueCode             string           `json:"sIssueCode"`                        // 銘柄コード
	ResultCode            string           `json:"sResultCode"`                       // 結果コード
	ResultText            string           `json:"sResultText"`                       // 結果テキスト
	WarningCode           string           `json:"sWarningCode"`                      // 警告コード
	WarningText           string           `json:"sWarningText"`                      // 警告テキスト
	TotalSellAmount       float64          `json:"sUritateDaikin,string"`             // 売建代金合計
	TotalBuyAmount        float64          `json:"sKaitateDaikin,string"`             // 買建代金合計
	TotalAmount           float64          `json:"sTotalDaikin,string"`               // 総代金合計
	TotalSellProfit       float64          `json:"sHyoukaSonekiGoukeiUridate,string"` // 評価損益合計_売建
	TotalBuyProfit        float64          `json:"sHyoukaSonekiGoukeiKaidate,string"` // 評価損益合計_買建
	TotalProfit           float64          `json:"sTotalHyoukaSonekiGoukei,string"`   // 総評価損益合計
	SpecificAccountProfit float64          `json:"sTokuteiHyoukaSonekiGoukei,string"` // 特定口座残高評価損益合計
	GeneralAccountProfit  float64          `json:"sIppanHyoukaSonekiGoukei,string"`   // 一般口座残高評価損益合計
	Positions             []marginPosition `json:"aShinyouTategyokuList"`             // 信用建玉リスト
}

func (r *marginPositionListResponse) UnmarshalJSON(b []byte) error {
	// 文字列でないところに空文字を返されることがあるので、置換しておく
	replaced := b
	replaces := map[string]string{
		`"aShinyouTategyokuList":""`: `"aShinyouTategyokuList":[]`,
	}
	for o, n := range replaces {
		replaced = bytes.Replace(replaced, []byte(o), []byte(n), -1)
	}

	type alias marginPositionListResponse
	ra := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}

	return json.Unmarshal(replaced, ra)
}

func (r *marginPositionListResponse) response() MarginPositionListResponse {
	positions := make([]MarginPosition, len(r.Positions))
	for i, p := range r.Positions {
		positions[i] = p.response()
	}

	return MarginPositionListResponse{
		CommonResponse:        r.commonResponse.response(),
		IssueCode:             r.IssueCode,
		ResultCode:            r.ResultCode,
		ResultText:            r.ResultText,
		WarningCode:           r.WarningCode,
		WarningText:           r.WarningText,
		TotalSellAmount:       r.TotalSellAmount,
		TotalBuyAmount:        r.TotalBuyAmount,
		TotalAmount:           r.TotalAmount,
		TotalSellProfit:       r.TotalSellProfit,
		TotalBuyProfit:        r.TotalBuyProfit,
		TotalProfit:           r.TotalProfit,
		SpecificAccountProfit: r.SpecificAccountProfit,
		GeneralAccountProfit:  r.GeneralAccountProfit,
		Positions:             positions,
	}
}

type marginPosition struct {
	WarningCode        string             `json:"sOrderWarningCode"`                   // 警告コード
	WarningText        string             `json:"sOrderWarningText"`                   // 警告テキスト
	PositionNumber     string             `json:"sOrderTategyokuNumber"`               // 建玉番号
	IssueCode          string             `json:"sOrderIssueCode"`                     // 銘柄コード
	Exchange           Exchange           `json:"sOrderSizyouC"`                       // 市場
	Side               Side               `json:"sOrderBaibaiKubun"`                   // 売買区分
	ExitTermType       ExitTermType       `json:"sOrderBensaiKubun"`                   // 弁済区分
	AccountType        AccountType        `json:"sOrderZyoutoekiKazeiC"`               // 譲渡益課税区分
	OrderQuantity      float64            `json:"sOrderTategyokuSuryou,string"`        // 建株数
	UnitPrice          float64            `json:"sOrderTategyokuTanka,string"`         // 建単価
	CurrentPrice       float64            `json:"sOrderHyoukaTanka,string"`            // 評価単価
	Profit             float64            `json:"sOrderGaisanHyoukaSoneki,string"`     // 評価損益
	ProfitRatio        float64            `json:"sOrderGaisanHyoukaSonekiRitu,string"` // 評価損益率
	TotalPrice         float64            `json:"sTategyokuDaikin,string"`             // 建玉代金
	Commission         float64            `json:"sOrderTateTesuryou,string"`           // 建手数料
	Interest           float64            `json:"sOrderZyunHibu,string"`               // 順日歩
	Premiums           float64            `json:"sOrderGyakuhibu,string"`              // 逆日歩
	RewritingFee       float64            `json:"sOrderKakikaeryou,string"`            // 書換料
	ManagementFee      float64            `json:"sOrderKanrihi,string"`                // 管理費
	LendingFee         float64            `json:"sOrderKasikaburyou,string"`           // 貸株料
	OtherFee           float64            `json:"sOrderSonota,string"`                 // その他
	ContractDate       Ymd                `json:"sOrderTategyokuDay"`                  // 建日
	ExitTerm           Ymd                `json:"sOrderTategyokuKizituDay"`            // 建玉期日日
	OwnedQuantity      float64            `json:"sTategyokuSuryou,string"`             // 建玉数量
	ExitQuantity       float64            `json:"sOrderYakuzyouHensaiKabusu,string"`   // 約定返済株数
	DeliveryQuantity   float64            `json:"sOrderGenbikiGenwatasiKabusu,string"` // 現引現渡株数
	HoldQuantity       float64            `json:"sOrderOrderSuryou,string"`            // 注文中数量
	ReturnableQuantity float64            `json:"sOrderHensaiKanouSuryou,string"`      // 返済可能数量
	PrevClosePrice     float64            `json:"sSyuzituOwarine,string"`              // 前日終値
	PrevCloseRatio     float64            `json:"sZenzituHi,string"`                   // 前日比
	PrevClosePercent   float64            `json:"sZenzituHiPer,string"`                // 前日比(%)
	PrevCloseRatioType PrevCloseRatioType `json:"sUpDownFlag"`                         // 騰落率Flag
}

func (r *marginPosition) response() MarginPosition {
	return MarginPosition{
		WarningCode:        r.WarningCode,
		WarningText:        r.WarningText,
		PositionNumber:     r.PositionNumber,
		IssueCode:          r.IssueCode,
		Exchange:           r.Exchange,
		Side:               r.Side,
		ExitTermType:       r.ExitTermType,
		AccountType:        r.AccountType,
		OrderQuantity:      r.OrderQuantity,
		UnitPrice:          r.UnitPrice,
		CurrentPrice:       r.CurrentPrice,
		Profit:             r.Profit,
		ProfitRatio:        r.ProfitRatio,
		TotalPrice:         r.TotalPrice,
		Commission:         r.Commission,
		Interest:           r.Interest,
		Premiums:           r.Premiums,
		RewritingFee:       r.RewritingFee,
		ManagementFee:      r.ManagementFee,
		LendingFee:         r.LendingFee,
		OtherFee:           r.OtherFee,
		ContractDate:       r.ContractDate.Time,
		ExitTerm:           r.ExitTerm.Time,
		OwnedQuantity:      r.OwnedQuantity,
		ExitQuantity:       r.ExitQuantity,
		DeliveryQuantity:   r.DeliveryQuantity,
		HoldQuantity:       r.HoldQuantity,
		ReturnableQuantity: r.ReturnableQuantity,
		PrevClosePrice:     r.PrevClosePrice,
		PrevCloseRatio:     r.PrevCloseRatio,
		PrevClosePercent:   r.PrevClosePercent,
		PrevCloseRatioType: r.PrevCloseRatioType,
	}
}

// MarginPositionListResponse - 信用建玉一覧レスポンス
type MarginPositionListResponse struct {
	CommonResponse
	IssueCode             string           // 銘柄コード
	ResultCode            string           // 結果コード
	ResultText            string           // 結果テキスト
	WarningCode           string           // 警告コード
	WarningText           string           // 警告テキスト
	TotalSellAmount       float64          // 売建代金合計
	TotalBuyAmount        float64          // 買建代金合計
	TotalAmount           float64          // 総代金合計
	TotalSellProfit       float64          // 評価損益合計_売建
	TotalBuyProfit        float64          // 評価損益合計_買建
	TotalProfit           float64          // 総評価損益合計
	SpecificAccountProfit float64          // 特定口座残高評価損益合計
	GeneralAccountProfit  float64          // 一般口座残高評価損益合計
	Positions             []MarginPosition // 信用建玉リスト
}

// MarginPosition - 信用建玉
type MarginPosition struct {
	WarningCode        string             // 警告コード
	WarningText        string             // 警告テキスト
	PositionNumber     string             // 建玉番号
	IssueCode          string             // 銘柄コード
	Exchange           Exchange           // 市場
	Side               Side               // 売買区分
	ExitTermType       ExitTermType       // 弁済区分
	AccountType        AccountType        // 譲渡益課税区分
	OrderQuantity      float64            // 建株数
	UnitPrice          float64            // 建単価
	CurrentPrice       float64            // 評価単価
	Profit             float64            // 評価損益
	ProfitRatio        float64            // 評価損益率
	TotalPrice         float64            // 建玉代金
	Commission         float64            // 建手数料
	Interest           float64            // 順日歩
	Premiums           float64            // 逆日歩
	RewritingFee       float64            // 書換料
	ManagementFee      float64            // 管理費
	LendingFee         float64            // 貸株料
	OtherFee           float64            // その他
	ContractDate       time.Time          // 建日
	ExitTerm           time.Time          // 建玉期日日
	OwnedQuantity      float64            // 建玉数量
	ExitQuantity       float64            // 約定返済株数
	DeliveryQuantity   float64            // 現引現渡株数
	HoldQuantity       float64            // 注文中数量
	ReturnableQuantity float64            // 返済可能数量
	PrevClosePrice     float64            // 前日終値
	PrevCloseRatio     float64            // 前日比
	PrevClosePercent   float64            // 前日比(%)
	PrevCloseRatioType PrevCloseRatioType // 騰落率Flag
}

// MarginPositionList - 信用建玉リスト
func (c *client) MarginPositionList(ctx context.Context, session *Session, req MarginPositionListRequest) (*MarginPositionListResponse, error) {
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
	var res marginPositionListResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, fmt.Errorf("%s: %w", err, UnmarshalFailedErr)
	}

	Res := res.response()
	return &Res, nil
}
