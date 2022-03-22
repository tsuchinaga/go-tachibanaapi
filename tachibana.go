package tachibana

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"

	"golang.org/x/text/transform"
)

// NewClient - クライアントの生成
func NewClient(env Environment, ver ApiVersion) Client {
	client := &client{
		clock:     newClock(),
		env:       env,
		ver:       ver,
		requester: &requester{},
	}

	return client
}

type Client interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)                                                          // ログイン
	Logout(ctx context.Context, session *Session, req LogoutRequest) (*LogoutResponse, error)                                     // ログアウト
	NewOrder(ctx context.Context, session *Session, req NewOrderRequest) (*NewOrderResponse, error)                               // 新規注文
	CorrectOrder(ctx context.Context, session *Session, req CorrectOrderRequest) (*CorrectOrderResponse, error)                   // 訂正注文
	CancelOrder(ctx context.Context, session *Session, req CancelOrderRequest) (*CancelOrderResponse, error)                      // 取消注文
	StockWallet(ctx context.Context, session *Session, req StockWalletRequest) (*StockWalletResponse, error)                      // 買余力
	MarginWallet(ctx context.Context, session *Session, req MarginWalletRequest) (*MarginWalletResponse, error)                   // 建余力&本日維持率
	StockSellable(ctx context.Context, session *Session, req StockSellableRequest) (*StockSellableResponse, error)                // 売却可能数量
	OrderList(ctx context.Context, session *Session, req OrderListRequest) (*OrderListResponse, error)                            // 注文一覧
	OrderListDetail(ctx context.Context, session *Session, req OrderListDetailRequest) (*OrderListDetailResponse, error)          // 注文一覧(詳細)
	StockPositionList(ctx context.Context, session *Session, req StockPositionListRequest) (*StockPositionListResponse, error)    // 現物株リスト
	MarginPositionList(ctx context.Context, session *Session, req MarginPositionListRequest) (*MarginPositionListResponse, error) // 信用建玉リスト
	StockMaster(ctx context.Context, session *Session, req StockMasterRequest) (*StockMasterResponse, error)                      // 株式銘柄マスタ
	MarketPrice(ctx context.Context, session *Session, req MarketPriceRequest) (*MarketPriceResponse, error)                      // 時価関連情報
	BusinessDay(ctx context.Context, session *Session, req BusinessDayRequest) ([]*BusinessDayResponse, error)                    // 日付情報
	TickGroup(ctx context.Context, session *Session, req TickGroupRequest) ([]*TickGroupResponse, error)                          // 呼値
}

type client struct {
	clock     iClock
	env       Environment
	ver       ApiVersion
	requester iRequester
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

// commonResponseFormat - 共通レスポンスフォーマット
var commonResponseFormat = ResponseFormatWrapped | ResponseFormatWordKey

// commonRequest - リクエストの共通的な項目
type commonRequest struct {
	No             int64          `json:"p_no,string"`      // 送信通番
	SendDate       RequestTime    `json:"p_sd_date"`        // 送信日時
	MessageType    MessageType    `json:"sCLMID"`           // 機能ID
	ResponseFormat ResponseFormat `json:"sJsonOfmt,string"` // レスポンスフォーマット
}

// commonResponse - パース用レスポンスの共通的な項目
type commonResponse struct {
	No           int64       `json:"p_no,string"` // 送信通番
	SendDate     RequestTime `json:"p_sd_date"`   // 送信日時
	ReceiveDate  RequestTime `json:"p_rv_date"`   // 受信日時
	ErrorNo      ErrorNo     `json:"p_errno"`     // エラー番号
	ErrorMessage string      `json:"p_err"`       // エラー文言
	MessageType  MessageType `json:"sCLMID"`      // 機能ID
}

// response - パース用レスポンスから使いやすい形のレスポンスに変換して返す
func (r *commonResponse) response() CommonResponse {
	return CommonResponse{
		No:           r.No,
		SendDate:     r.SendDate.Time,
		ReceiveDate:  r.ReceiveDate.Time,
		ErrorNo:      r.ErrorNo,
		ErrorMessage: r.ErrorMessage,
		MessageType:  r.MessageType,
	}
}

// CommonResponse - レスポンスの共通的な項目
type CommonResponse struct {
	No           int64       // 送信通番
	SendDate     time.Time   // 送信日時
	ReceiveDate  time.Time   // 受信日時
	ErrorNo      ErrorNo     // エラー番号
	ErrorMessage string      // エラー文言
	MessageType  MessageType // 機能ID
}

type iRequester interface {
	get(ctx context.Context, uri string, request interface{}) ([]byte, error)
	stream(ctx context.Context, uri string, request interface{}) (<-chan []byte, <-chan error)
}

type requester struct {
	insecureSkipVerify bool
}

// encode - 文字コードの変換(UTF-8 -> Shift-JIS)と、URLエンコード
func (r *requester) encode(b []byte) ([]byte, error) {
	// utf-8 to shift-jis
	b, _, err := transform.Bytes(japanese.ShiftJIS.NewEncoder(), b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, EncodeErr)
	}

	// http encode
	return []byte(url.QueryEscape(string(b))), nil
}

// encode - URLデコードと、文字コードの変換(Shift-JIS -> UTF-8)
func (r *requester) decode(b []byte) ([]byte, error) {
	// レスポンスはbodyにはいってくるのでhttp decodeが不要

	// shift-jis to utf-8
	//   基本的に Shift-JIS -> UTF-8ではエンコードに失敗しないはずなので、エラーを捨てる
	b, _, _ = transform.Bytes(japanese.ShiftJIS.NewDecoder(), b)
	return b, nil
}

// get - GETリクエスト
func (r *requester) get(ctx context.Context, uri string, request interface{}) ([]byte, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	qb, err := r.encode(rb)
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(uri)
	u.RawQuery = string(qb)

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// リクエスト送信
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		return b, nil
	} else {
		return nil, fmt.Errorf("status is %d(body: %s): %w", res.StatusCode, string(b), StatusNotOkErr)
	}
}

// stream - chunked response リクエスト
func (r *requester) stream(ctx context.Context, uri string, request interface{}) (<-chan []byte, <-chan error) {
	ch := make(chan []byte)
	errCh := make(chan error)

	go func() {
		defer close(ch)
		defer close(errCh)

		rb, err := json.Marshal(request)
		if err != nil {
			errCh <- err
			return
		}

		query, err := r.encode(rb)
		if err != nil {
			errCh <- err
			return
		}
		u, _ := url.Parse(uri)
		u.RawQuery = string(query)

		// TCPソケットオープン
		dialer := &net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		}
		host := u.Host
		if !strings.Contains(host, ":") {
			host += ":443"
		}
		conn, err := tls.DialWithDialer(dialer, "tcp", host, &tls.Config{InsecureSkipVerify: r.insecureSkipVerify})
		if err != nil {
			errCh <- err
			return
		}
		defer conn.Close()

		// コネクションを通してリクエストの送信
		req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
		if err != nil {
			errCh <- err
			return
		}
		err = req.Write(conn)
		if err != nil {
			errCh <- err
			return
		}

		// レスポンスの確認
		reader := bufio.NewReader(conn)
		res, err := http.ReadResponse(reader, req)
		if err != nil {
			errCh <- err
			return
		}
		if len(res.TransferEncoding) < 1 || res.TransferEncoding[0] != "chunked" {
			errCh <- errors.New("response is not chunked")
			return
		}

		// レスポンスをチャネルに流していく
		sCh, sErrCh := r.scanChunkedResponse(reader)
		for {
			select {
			case <-ctx.Done():
				return
			case err, ok := <-sErrCh:
				if ok && err != nil {
					errCh <- err
				}
				return
			case b, ok := <-sCh:
				if !ok {
					return
				}

				d, _ := r.decode(b) // decodeでは失敗がおきないのでエラーを捨てる
				ch <- d
			}
		}
	}()

	return ch, errCh
}

func (r *requester) scanChunkedResponse(reader io.Reader) (<-chan []byte, <-chan error) {
	ch := make(chan []byte)
	errCh := make(chan error)

	go func() {
		defer close(ch)
		defer close(errCh)

		scanner := bufio.NewScanner(httputil.NewChunkedReader(reader))
		for scanner.Scan() {
			ch <- scanner.Bytes()
		}
		errCh <- scanner.Err()
	}()

	return ch, errCh
}
