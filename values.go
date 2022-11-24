package tachibana

import (
	"math"
	"sync"
)

// NoChangeFloat - float64で変更しないことを指定
var NoChangeFloat float64 = math.Inf(-1)

// Session - リクエストセッション
type Session struct {
	lastRequestNo int64
	RequestURL    string
	MasterURL     string
	PriceURL      string
	EventURL      string
	mtx           sync.Mutex
}
