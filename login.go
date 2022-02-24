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
		CommonRequest: CommonRequest{
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
	CommonRequest
	UserId   string `json:"690"` // ログインユーザー
	Password string `json:"531"` // ログインパスワード
}

// LoginResponse - ログインレスポンス
type LoginResponse struct {
	CommonResponse
	ResultCode            string              `json:"534"` // 結果コード
	ResultText            string              `json:"535"` // 結果テキスト
	AccountType           AccountType         `json:"744"` // 譲渡益課税区分
	LastLoginDateTime     LoginDateTime       `json:"401"` // 最終ログイン日時
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

// Session - ログインレスポンスからセッションを取り出す
func (r *LoginResponse) Session() (*Session, error) {
	if r.ErrorNo != ErrTypeNoProblem && r.FeatureType != FeatureTypeLoginResponse && r.ResultCode != "0" {
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

	var res LoginResponse
	if err := c.get(ctx, c.auth, r, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
