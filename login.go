package tachibana

import (
	"context"
	"time"
)

// LoginRequest - ログインリクエスト
type LoginRequest struct {
	UserId   string // ログインユーザー
	Password string // ログインパスワード
}

// request - 送信できるログインリクエストを取得
func (r *LoginRequest) request(no int64, now time.Time) loginRequest {
	return loginRequest{
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			FeatureType:    FeatureTypeLoginRequest,
			ResponseFormat: ResponseFormatReadable | ResponseFormatWrapped | ResponseFormatWordKey,
		},
		UserId:   r.UserId,
		Password: r.Password,
	}
}

// loginRequest - パース用ログインリクエスト
type loginRequest struct {
	commonRequest
	UserId   string `json:"sUserId"`   // ログインユーザー
	Password string `json:"sPassword"` // ログインパスワード
}

// loginResponse - ログインレスポンス
type loginResponse struct {
	commonResponse
	ResultCode                string              `json:"sResultCode"`               // 結果コード
	ResultText                string              `json:"sResultText"`               // 結果テキスト
	AccountType               AccountType         `json:"sZyoutoekiKazeiC"`          // 譲渡益課税区分
	SecondPasswordOmit        NumberBool          `json:"sSecondPasswordOmit"`       // 暗証番号省略有無C
	LastLoginDateTime         YmdHms              `json:"sLastLoginDate"`            // 最終ログイン日時
	GeneralAccount            NumberBool          `json:"sSogoKouzaKubun"`           // 総合口座開設区分
	SafekeepingAccount        NumberBool          `json:"sHogoAdukariKouzaKubun"`    // 保護預り口座開設区分
	TransferAccount           NumberBool          `json:"sFurikaeKouzaKubun"`        // 振替決済口座開設区分
	ForeignAccount            NumberBool          `json:"sGaikokuKouzaKubun"`        // 外国口座開設区分
	MRFAccount                NumberBool          `json:"sMRFKouzaKubun"`            // MRF口座開設区分
	StockSpecificAccount      SpecificAccountType `json:"sTokuteiKouzaKubunGenbutu"` // 特定口座区分現物
	MarginSpecificAccount     SpecificAccountType `json:"sTokuteiKouzaKubunSinyou"`  // 特定口座区分信用
	InvestmentSpecificAccount SpecificAccountType `json:"sTokuteiKouzaKubunTousin"`  // 特定口座区分投信
	DividendAccount           NumberBool          `json:"sTokuteiHaitouKouzaKubun"`  // 配当特定口座区分
	SpecificAccount           NumberBool          `json:"sTokuteiKanriKouzaKubun"`   // 特定管理口座開設区分
	MarginAccount             NumberBool          `json:"sSinyouKouzaKubun"`         // 信用取引口座開設区分
	FutureOptionAccount       NumberBool          `json:"sSakopKouzaKubun"`          // 先物OP口座開設区分
	MMFAccount                NumberBool          `json:"sMMFKouzaKubun"`            // MMF口座開設区分
	ChinaForeignAccount       NumberBool          `json:"sTyukokufKouzaKubun"`       // 中国F口座開設区分
	FXAccount                 NumberBool          `json:"sKawaseKouzaKubun"`         // 為替保証金口座開設区分
	NISAAccount               NumberBool          `json:"sHikazeiKouzaKubun"`        // 非課税口座開設区分
	UnreadDocument            NumberBool          `json:"sKinsyouhouMidokuFlg"`      // 金商法交付書面未読フラグ
	RequestURL                string              `json:"sUrlRequest"`               // 仮想URL(REQUEST)
	EventURL                  string              `json:"sUrlEvent"`                 // 仮想URL(EVENT)
}

func (r *loginResponse) response() LoginResponse {
	return LoginResponse{
		CommonResponse:            r.commonResponse.response(),
		ResultCode:                r.ResultCode,
		ResultText:                r.ResultText,
		AccountType:               r.AccountType,
		SecondPasswordOmit:        r.SecondPasswordOmit.Bool(),
		LastLoginDateTime:         r.LastLoginDateTime.Time,
		GeneralAccount:            r.GeneralAccount.Bool(),
		SafekeepingAccount:        r.SafekeepingAccount.Bool(),
		TransferAccount:           r.TransferAccount.Bool(),
		ForeignAccount:            r.ForeignAccount.Bool(),
		MRFAccount:                r.MRFAccount.Bool(),
		StockSpecificAccount:      r.StockSpecificAccount,
		MarginSpecificAccount:     r.MarginSpecificAccount,
		InvestmentSpecificAccount: r.InvestmentSpecificAccount,
		DividendAccount:           r.DividendAccount.Bool(),
		SpecificAccount:           r.SpecificAccount.Bool(),
		MarginAccount:             r.MarginAccount.Bool(),
		FutureOptionAccount:       r.FutureOptionAccount.Bool(),
		MMFAccount:                r.MMFAccount.Bool(),
		ChinaForeignAccount:       r.ChinaForeignAccount.Bool(),
		FXAccount:                 r.FXAccount.Bool(),
		NISAAccount:               r.NISAAccount.Bool(),
		UnreadDocument:            r.UnreadDocument.Bool(),
		RequestURL:                r.RequestURL,
		EventURL:                  r.EventURL,
	}
}

// LoginResponse - ログインレスポンス
type LoginResponse struct {
	CommonResponse
	ResultCode                string              // 結果コード
	ResultText                string              // 結果テキスト
	AccountType               AccountType         // 譲渡益課税区分
	SecondPasswordOmit        bool                // 暗証番号省略有無C
	LastLoginDateTime         time.Time           // 最終ログイン日時
	GeneralAccount            bool                // 総合口座開設区分
	SafekeepingAccount        bool                // 保護預り口座開設区分
	TransferAccount           bool                // 振替決済口座開設区分
	ForeignAccount            bool                // 外国口座開設区分
	MRFAccount                bool                // MRF口座開設区分
	StockSpecificAccount      SpecificAccountType // 特定口座区分現物
	MarginSpecificAccount     SpecificAccountType // 特定口座区分信用
	InvestmentSpecificAccount SpecificAccountType // 特定口座区分投信
	DividendAccount           bool                // 配当特定口座区分
	SpecificAccount           bool                // 特定管理口座開設区分
	MarginAccount             bool                // 信用取引口座開設区分
	FutureOptionAccount       bool                // 先物OP口座開設区分
	MMFAccount                bool                // MMF口座開設区分
	ChinaForeignAccount       bool                // 中国F口座開設区分
	FXAccount                 bool                // 為替保証金口座開設区分
	NISAAccount               bool                // 非課税口座開設区分
	UnreadDocument            bool                // 金商法交付書面未読フラグ
	RequestURL                string              // 仮想URL(REQUEST)
	EventURL                  string              // 仮想URL(EVENT)
}

// Session - ログインレスポンスからセッションを取り出す
func (r *LoginResponse) Session() (*Session, error) {
	if r.ErrorNo != ErrorNoProblem || r.FeatureType != FeatureTypeLoginResponse || r.ResultCode != "0" {
		return nil, CanNotCreateSessionErr
	}

	return &Session{
		lastRequestNo: r.No,
		RequestURL:    r.RequestURL,
		EventURL:      r.EventURL,
	}, nil
}

// Login - ログイン
func (c *client) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	r := req.request(1, c.clock.Now())

	var res loginResponse
	if err := c.get(ctx, c.auth, r, &res); err != nil {
		return nil, err
	}

	Res := res.response()
	return &Res, nil
}
