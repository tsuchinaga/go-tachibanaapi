package tachibana

import (
	"context"
	"time"
)

// LogoutRequest - ログアウトリクエスト
type LogoutRequest struct{}

// request - 送信できるログアウトリクエストを取得
func (r *LogoutRequest) request(no int64, now time.Time) logoutRequest {
	return logoutRequest{
		CommonRequest: CommonRequest{
			No:          no,
			SendDate:    RequestTime{Time: now},
			FeatureType: FeatureTypeLogoutRequest,
		},
	}
}

// logoutRequest - パース用ログアウトリクエスト
type logoutRequest struct {
	CommonRequest
}

// LogoutResponse - ログアウトレスポンス
type LogoutResponse struct {
	CommonResponse
	ResultCode string `json:"534"` // 結果コード
	ResultText string `json:"535"` // 結果テキスト
}

// Logout - ログアウト
func (c *client) Logout(ctx context.Context, session *Session, req LogoutRequest) (*LogoutResponse, error) {
	if session == nil {
		return nil, NilArgumentErr
	}
	session.mtx.Lock()
	defer session.mtx.Unlock()

	session.lastRequestNo++
	r := req.request(session.lastRequestNo, c.clock.Now())

	var res LogoutResponse
	if err := c.get(ctx, session.RequestURL, r, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
