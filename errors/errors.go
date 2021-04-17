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
	// IError is 自定义错误接口
	IError interface {
		Error() string
		Code() int
	}
	// Error is 自定义错误信息
	Error struct {
		code       int
		innerError error
	}
)

var _ error = new(Error)
var _ IError = new(Error)

// Error is 转换为字符串
func (e Error) Error() string {
	return fmt.Sprintf("code: %d, err: %s", e.code, e.innerError.Error())
}

// Code is 转换为字符串
func (e Error) Code() int {
	return e.code
}

// New 创建一个内部error
func New(code int, err error) Error {
	return Error{
		code:       code,
		innerError: err,
	}
}

// Code 如果error是IError类型, 则获取error
func Code(err error) int {
	if e, ok := err.(Error); ok {
		return e.code
	}
	return Unkown
}
