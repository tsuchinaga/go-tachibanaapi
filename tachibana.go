package tachibana

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

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
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)                                                          // ログイン
	Logout(ctx context.Context, session *Session, req LogoutRequest) (*LogoutResponse, error)                                     // ログアウト
	NewOrder(ctx context.Context, session *Session, req NewOrderRequest) (*NewOrderResponse, error)                               // 新規注文
	CorrectOrder(ctx context.Context, session *Session, req CorrectOrderRequest) (*CorrectOrderResponse, error)                   // 訂正注文
	CancelOrder(ctx context.Context, session *Session, req CancelOrderRequest) (*CancelOrderResponse, error)                      // 取消注文
	OrderList(ctx context.Context, session *Session, req OrderListRequest) (*OrderListResponse, error)                            // 注文一覧
	OrderListDetail(ctx context.Context, session *Session, req OrderListDetailRequest) (*OrderListDetailResponse, error)          // 注文一覧(詳細)
	StockPositionList(ctx context.Context, session *Session, req StockPositionListRequest) (*StockPositionListResponse, error)    // 現物株リスト
	MarginPositionList(ctx context.Context, session *Session, req MarginPositionListRequest) (*MarginPositionListResponse, error) // 信用建玉リスト
	StockWallet(ctx context.Context, session *Session, req StockWalletRequest) (*StockWalletResponse, error)                      // 買余力
	MarginWallet(ctx context.Context, session *Session, req MarginWalletRequest) (*MarginWalletResponse, error)                   // 建余力&本日維持率
	StockMaster(ctx context.Context, session *Session, req StockMasterRequest) (*StockMasterResponse, error)                      // 株式銘柄マスタ
	MarketPrice(ctx context.Context, session *Session, req MarketPriceRequest) (*MarketPriceResponse, error)                      // 時価関連情報
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

// host - ホスト
func (c *client) host(env Environment) string {
	host := "kabuka.e-shiten.jp"
	if env == EnvironmentDemo {
		host = "demo-kabuka.e-shiten.jp"
	}
	return host
}

// authURL - ログインURLを返す
func (c *client) authURL(env Environment, ver ApiVersion) string {
	path := "e_api_"
	switch ver {
	case ApiVersionV4R1, ApiVersionV4R2:
		path += string(ver)
	default:
		path += string(ApiVersionLatest) // latest
	}
	return fmt.Sprintf("https://%s/%s/auth/", c.host(env), path)
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

// commonResponseFormat - 共通レスポンスフォーマット
var commonResponseFormat ResponseFormat = ResponseFormatWrapped | ResponseFormatWordKey

// commonRequest - リクエストの共通的な項目
type commonRequest struct {
	No             int64          `json:"p_no,string"`      // 送信通番
	SendDate       RequestTime    `json:"p_sd_date"`        // 送信日時
	FeatureType    FeatureType    `json:"sCLMID"`           // 機能ID
	ResponseFormat ResponseFormat `json:"sJsonOfmt,string"` // レスポンスフォーマット
}

// commonResponse - パース用レスポンスの共通的な項目
type commonResponse struct {
	No           int64       `json:"p_no,string"` // 送信通番
	SendDate     RequestTime `json:"p_sd_date"`   // 送信日時
	ReceiveDate  RequestTime `json:"p_rv_date"`   // 受信日時
	ErrorNo      ErrorNo     `json:"p_errno"`     // エラー番号
	ErrorMessage string      `json:"p_err"`       // エラー文言
	FeatureType  FeatureType `json:"sCLMID"`      // 機能ID
}

// response - パース用レスポンスから使いやすい形のレスポンスに変換して返す
func (r *commonResponse) response() CommonResponse {
	return CommonResponse{
		No:           r.No,
		SendDate:     r.SendDate.Time,
		ReceiveDate:  r.ReceiveDate.Time,
		ErrorNo:      r.ErrorNo,
		ErrorMessage: r.ErrorMessage,
		FeatureType:  r.FeatureType,
	}
}

// CommonResponse - レスポンスの共通的な項目
type CommonResponse struct {
	No           int64       // 送信通番
	SendDate     time.Time   // 送信日時
	ReceiveDate  time.Time   // 受信日時
	ErrorNo      ErrorNo     // エラー番号
	ErrorMessage string      // エラー文言
	FeatureType  FeatureType // 機能ID
}
