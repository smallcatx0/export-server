package exception

import "strings"

// Exception 包装响应体
type Exception struct {
	HTTPCode int
	Code     uint32
	Msg      string
	Data     interface{}
}

func (r *Exception) Error() string {
	return r.Msg
}

// NewException 客户端错误
func NewException(code uint32, msg ...string) error {
	err := &Exception{
		Code:     code,
		HTTPCode: 400,
	}
	if len(msg) == 0 {
		err.Msg = ErrNos[code]
	} else {
		err.Msg = strings.Join(msg, " ")
	}
	return err
}

// NewError 服务端错误
func NewError(code uint32, httpCode int, msg ...string) error {
	err := &Exception{
		Code:     code,
		HTTPCode: httpCode,
	}
	if len(msg) == 0 {
		err.Msg = ErrNos[code]
	} else {
		err.Msg = strings.Join(msg, " ")
	}
	return err
}
