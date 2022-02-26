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
			No:          no,
			SendDate:    RequestTime{Time: now},
			FeatureType: FeatureTypeLoginRequest,
		},
		UserId:   r.UserId,
		Password: r.Password,
	}
}

// loginRequest - パース用ログインリクエスト
type loginRequest struct {
	commonRequest
	UserId   string `json:"690"` // ログインユーザー
	Password string `json:"531"` // ログインパスワード
}

// loginResponse - ログインレスポンス
type loginResponse struct {
	commonResponse
	ResultCode            string              `json:"534"` // 結果コード
	ResultText            string              `json:"535"` // 結果テキスト
	AccountType           AccountType         `json:"744"` // 譲渡益課税区分
	LastLoginDateTime     YmdHms              `json:"401"` // 最終ログイン日時
	GeneralAccount        NumberBool          `json:"580"` // 総合口座開設区分
	SafekeepingAccount    NumberBool          `json:"287"` // 保護預り口座開設区分
	TransferAccount       NumberBool          `json:"232"` // 振替決済口座開設区分
	ForeignAccount        NumberBool          `json:"234"` // 外国口座開設区分
	MRFAccount            NumberBool          `json:"404"` // MRF口座開設区分
	StockSpecificAccount  SpecificAccountType `json:"645"` // 特定口座区分現物
	MarginSpecificAccount SpecificAccountType `json:"646"` // 特定口座区分信用
	DividendAccount       NumberBool          `json:"642"` // 配当特定口座区分
	SpecificAccount       NumberBool          `json:"644"` // 特定管理口座開設区分
	MarginAccount         NumberBool          `json:"565"` // 信用取引口座開設区分
	FutureOptionAccount   NumberBool          `json:"542"` // 先物OP口座開設区分
	MMFAccount            NumberBool          `json:"403"` // MMF口座開設区分
	ChinaForeignAccount   NumberBool          `json:"660"` // 中国F口座開設区分
	FXAccount             NumberBool          `json:"348"` // 為替保証金口座開設区分
	NISAAccount           NumberBool          `json:"281"` // 非課税口座開設区分
	UnreadDocument        NumberBool          `json:"374"` // 金商法交付書面未読フラグ
	RequestURL            string              `json:"689"` // 仮想URL(REQUEST)
	EventURL              string              `json:"688"` // 仮想URL(EVENT)
}

func (r *loginResponse) response() LoginResponse {
	return LoginResponse{
		CommonResponse:        r.commonResponse.response(),
		ResultCode:            r.ResultCode,
		ResultText:            r.ResultText,
		AccountType:           r.AccountType,
		LastLoginDateTime:     r.LastLoginDateTime.Time,
		GeneralAccount:        r.GeneralAccount.Bool(),
		SafekeepingAccount:    r.SafekeepingAccount.Bool(),
		TransferAccount:       r.TransferAccount.Bool(),
		ForeignAccount:        r.ForeignAccount.Bool(),
		MRFAccount:            r.MRFAccount.Bool(),
		StockSpecificAccount:  r.StockSpecificAccount,
		MarginSpecificAccount: r.MarginSpecificAccount,
		DividendAccount:       r.DividendAccount.Bool(),
		SpecificAccount:       r.SpecificAccount.Bool(),
		MarginAccount:         r.MarginAccount.Bool(),
		FutureOptionAccount:   r.FutureOptionAccount.Bool(),
		MMFAccount:            r.MMFAccount.Bool(),
		ChinaForeignAccount:   r.ChinaForeignAccount.Bool(),
		FXAccount:             r.FXAccount.Bool(),
		NISAAccount:           r.NISAAccount.Bool(),
		UnreadDocument:        r.UnreadDocument.Bool(),
		RequestURL:            r.RequestURL,
		EventURL:              r.EventURL,
	}
}

// LoginResponse - ログインレスポンス
type LoginResponse struct {
	CommonResponse
	ResultCode            string              // 結果コード
	ResultText            string              // 結果テキスト
	AccountType           AccountType         // 譲渡益課税区分
	LastLoginDateTime     time.Time           // 最終ログイン日時
	GeneralAccount        bool                // 総合口座開設区分
	SafekeepingAccount    bool                // 保護預り口座開設区分
	TransferAccount       bool                // 振替決済口座開設区分
	ForeignAccount        bool                // 外国口座開設区分
	MRFAccount            bool                // MRF口座開設区分
	StockSpecificAccount  SpecificAccountType // 特定口座区分現物
	MarginSpecificAccount SpecificAccountType // 特定口座区分信用
	DividendAccount       bool                // 配当特定口座区分
	SpecificAccount       bool                // 特定管理口座開設区分
	MarginAccount         bool                // 信用取引口座開設区分
	FutureOptionAccount   bool                // 先物OP口座開設区分
	MMFAccount            bool                // MMF口座開設区分
	ChinaForeignAccount   bool                // 中国F口座開設区分
	FXAccount             bool                // 為替保証金口座開設区分
	NISAAccount           bool                // 非課税口座開設区分
	UnreadDocument        bool                // 金商法交付書面未読フラグ
	RequestURL            string              // 仮想URL(REQUEST)
	EventURL              string              // 仮想URL(EVENT)
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
