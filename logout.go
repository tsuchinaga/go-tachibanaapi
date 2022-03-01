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
		commonRequest: commonRequest{
			No:             no,
			SendDate:       RequestTime{Time: now},
			FeatureType:    FeatureTypeLogoutRequest,
			ResponseFormat: ResponseFormatReadable | ResponseFormatWrapped | ResponseFormatWordKey,
		},
	}
}

// logoutRequest - パース用ログアウトリクエスト
type logoutRequest struct {
	commonRequest
}

// logoutResponse - パース用ログアウトレスポンス
type logoutResponse struct {
	commonResponse
	ResultCode string `json:"sResultCode"` // 結果コード
	ResultText string `json:"sResultText"` // 結果テキスト
}

func (r *logoutResponse) response() LogoutResponse {
	return LogoutResponse{
		CommonResponse: r.commonResponse.response(),
		ResultCode:     r.ResultCode,
		ResultText:     r.ResultText,
	}
}

// LogoutResponse - ログアウトレスポンス
type LogoutResponse struct {
	CommonResponse
	ResultCode string // 結果コード
	ResultText string // 結果テキスト
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

	var res logoutResponse
	if err := c.get(ctx, session.RequestURL, r, &res); err != nil {
		return nil, err
	}

	Res := res.response()
	return &Res, nil
}
