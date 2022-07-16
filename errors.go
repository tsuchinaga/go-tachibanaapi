package tachibana

import "errors"

var (
	NilArgumentErr         = errors.New("nil argument")
	EncodeErr              = errors.New("encode error")
	StatusNotOkErr         = errors.New("status not ok")
	CanNotCreateSessionErr = errors.New("cannot create session")
	UnmarshalFailedErr     = errors.New("unmarshal failed")
	StreamError            = errors.New("stream error")
)
