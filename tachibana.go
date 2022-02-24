package tachibana

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/text/encoding/japanese"

	"golang.org/x/text/transform"
)

// NewClient - クライアントの生成
func NewClient(env Environment, ver ApiVersion) Client {
	client := &client{
		clock: newClock(),
		env:   env,
		ver:   ver,
	}
	client.auth = client.authURL(client.env, client.ver)

	return client
}

type Client interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, session *Session, req LogoutRequest) (*LogoutResponse, error)
}

type client struct {
	clock iClock
	env   Environment
	ver   ApiVersion
	auth  string
}

// encode - 文字コードの変換(UTF-8 -> Shift-JIS)と、URLエンコード
func (c *client) encode(str string) (string, error) {
	// utf-8 to shift-jis
	str, _, err := transform.String(japanese.ShiftJIS.NewEncoder(), str)
	if err != nil {
		return "", fmt.Errorf("%s: %w", err, EncodeErr)
	}

	// http encode
	return url.QueryEscape(str), nil
}

// encode - URLデコードと、文字コードの変換(Shift-JIS -> UTF-8)
func (c *client) decode(str string) (string, error) {
	// レスポンスはbodyにはいってくるのでhttp decodeが不要

	// shift-jis to utf-8
	//   基本的に Shift-JIS -> UTF-8ではエンコードに失敗しないはずなので、エラーを捨てる
	str, _, _ = transform.String(japanese.ShiftJIS.NewDecoder(), str)
	return str, nil
}

// authURL - ログインURLを返す
func (c *client) authURL(env Environment, ver ApiVersion) string {
	host := "kabuka.e-shiten.jp"
	if env == EnvironmentDemo {
		host = "demo-kabuka.e-shiten.jp"
	}

	path := "e_api_"
	switch ver {
	case ApiVersionV4R1, ApiVersionV4R2:
		path += string(ver)
	default:
		path += string(ApiVersionLatest) // latest
	}
	return fmt.Sprintf("https://%s/%s/auth/", host, path)
}

// get - GETリクエスト
func (c *client) get(ctx context.Context, uri string, request interface{}, response interface{}) error {
	rb, err := json.Marshal(request)
	if err != nil {
		return err
	}
	query, err := c.encode(string(rb))
	if err != nil {
		return err
	}

	u, _ := url.Parse(uri)
	u.RawQuery = query

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return err
	}

	// リクエスト送信
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Println(b, string(b))

	if res.StatusCode == http.StatusOK {
		if err := c.parseResponse(b, response); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("status is %d(body: %s): %w", res.StatusCode, string(b), StatusNotOkErr)
	}

	return nil
}

// parseResponse - レスポンスをパースする
func (c *client) parseResponse(body []byte, v interface{}) error {
	d, _ := c.decode(string(body)) // エラーが発生しないのでチェックせず捨てる

	if err := json.Unmarshal([]byte(d), v); err != nil {
		return err
	}
	return nil
}

// CommonRequest - リクエストの共通的な項目
type CommonRequest struct {
	No          int64       `json:"175,string"` // 送信通番
	SendDate    RequestTime `json:"177"`        // 送信日時
	FeatureType FeatureType `json:"192"`        // 機能ID
}

// CommonResponse - レスポンスの共通的な項目
type CommonResponse struct {
	No           int64       `json:"175,string"` // 送信通番
	SendDate     RequestTime `json:"177"`        // 送信日時
	ReceiveDate  RequestTime `json:"176"`        // 受信日時
	ErrorNo      ErrType     `json:"174"`        // エラー番号
	ErrorMessage string      `json:"173"`        // エラー文言
	FeatureType  FeatureType `json:"192"`        // 機能ID
}
