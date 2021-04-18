package utils

import "errors"

// ValidatePhone 验证手机号格式
func ValidatePhone(phone string) error {
	if len(phone) != 11 {
		return errors.New("手机号长度应该是11位")
	}
	return nil
}
