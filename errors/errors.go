package errors

import "fmt"

var (
	InvalidSource    = "invalid_source"
	InvalidSourceErr = &Error{Code: InvalidSource, Msg: InvalidSource}
)

type (
	Error struct {
		Code       string `json:"code"`
		Msg        string `json:"msg,omitempty"`
		InnerError error  `json:"-"`
	}
)

func (e Error) Error() string {
	return fmt.Sprintf("code: %s, msg: %s", e.Code, e.Msg)
}
