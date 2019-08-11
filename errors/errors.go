package errors

import "fmt"

// 错误类型
var (
	InvalidSource    = "invalid_source"
	InvalidSourceErr = &Error{Code: InvalidSource, Msg: InvalidSource}
)

type (
	// Error is 自定义错误信息
	Error struct {
		Code       string `json:"code"`
		Msg        string `json:"msg,omitempty"`
		InnerError error  `json:"-"`
	}
)

// Error is 转换为字符串
func (e Error) Error() string {
	return fmt.Sprintf("code: %s, msg: %s", e.Code, e.Msg)
}
