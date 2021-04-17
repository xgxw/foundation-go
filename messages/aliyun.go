package messages

import (
	"context"
	"encoding/json"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/pkg/errors"
	"github.com/xgxw/foundation-go"
	"github.com/xgxw/foundation-go/utils"
)

const BaseURL = "dysmsapi.aliyuncs.com"

type (
	AliClient struct {
		*dysmsapi20170525.Client
	}
	AliMessageOptions struct {
		AccessKeyId     string `json:"access_key_id"`
		AccessKeySecret string `json:"access_key_secret"`
	}
)

func NewAliClient(opts *AliMessageOptions) (*AliClient, error) {
	baseURL := BaseURL
	config := &openapi.Config{
		AccessKeyId:     &opts.AccessKeyId,
		AccessKeySecret: &opts.AccessKeySecret,
		Endpoint:        &baseURL,
	}
	client, err := dysmsapi20170525.NewClient(config)
	return &AliClient{
		Client: client,
	}, err
}

func (c *AliClient) Send(ctx context.Context, msg *foundation.MessageTemplate) error {
	params, err := json.Marshal(msg.TemplateParam)
	if err != nil {
		return errors.Wrap(err, "template_param must be json")
	}
	paramsStr := utils.Bytes2String(params)
	request := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  &msg.To,
		SignName:      &msg.SignName,
		TemplateCode:  &msg.TemplateCode,
		TemplateParam: &paramsStr,
	}
	_, err = c.Client.SendSms(request)
	if err != nil {
		return errors.Wrap(err, "send sms error")
	}
	return nil
}
