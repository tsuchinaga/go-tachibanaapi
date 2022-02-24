package tachibana

import "sync"

// Session - リクエストセッション
type Session struct {
	lastRequestNo int64
	RequestURL    string
	EventURL      string
	mtx           sync.Mutex
}
