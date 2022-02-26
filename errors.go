package tachibana

import "errors"

var (
	NilArgumentErr         = errors.New("nil argument")
	EncodeErr              = errors.New("encode error")
	StatusNotOkErr         = errors.New("status not ok")
	CanNotCreateSessionErr = errors.New("cannot create session")
)