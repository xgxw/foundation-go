package messages

import (
	"context"
	"os"
	"testing"

	"github.com/xgxw/foundation-go"
	"github.com/xgxw/foundation-go/utils"
)

func TestAliSendManual(t *testing.T) {
	if utils.SkipTest() {
		t.SkipNow()
		return
	}
	client, err := NewAliClient(&AliMessageOptions{
		AccessKeyId:     os.Getenv("access_key_id"),
		AccessKeySecret: os.Getenv("access_key_secret"),
	})
	if err != nil {
		t.Errorf("new client error. err: %v", err)
		return
	}
	err = client.Send(context.Background(), &foundation.MessageTemplate{
		To:           "13260105235",
		SignName:     "读书行路",
		TemplateCode: "SMS_215340597",
		TemplateParam: map[string]interface{}{
			"code": 123456,
		},
	})
	if err != nil {
		t.Errorf("new client error. err: %v", err)
		return
	}
	t.Log("send success")
}
