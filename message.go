package foundation

import (
	"context"
)

type (
	MessageTemplate struct {
		To            string                 `json:"to"`
		SignName      string                 `json:"sign_name"`
		TemplateCode  string                 `json:"template_name"`
		TemplateParam map[string]interface{} `json:"template_param"`
	}
)

type Message interface {
	Send(context.Context, *MessageTemplate) error
}
