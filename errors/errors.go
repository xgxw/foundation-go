package errors

import (
	"errors"
	"fmt"
)

// 错误类型
var (
	Unkown      = 0
	UnkownError = New(Unkown, errors.New("unkown error"))
)

type (
	// Error is 自定义错误信息
	Error struct {
		Code       int   `json:"code"`
		InnerError error `json:"-"`
	}
)

// Error is 转换为字符串
func (e Error) Error() string {
	return fmt.Sprintf("code: %d, err: %s", e.Code, e.InnerError.Error())
}

func New(code int, err error) Error {
	return Error{
		Code:       code,
		InnerError: err,
	}
}

func Code(err error) int {
	if e, ok := err.(Error); ok {
		return e.Code
	}
	return Unkown
}
