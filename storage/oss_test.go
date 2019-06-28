package storage

import (
	"context"
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_OssClienta_GetObject(t *testing.T) {
	var endpoint = "oss-cn-beijing.aliyuncs.com"
	var accessKeyID = os.Getenv("accrss_key_id")
	var accessKeySecret = os.Getenv("accrss_key_secret")
	if accessKeyID == "" || accessKeySecret == "" {
		return
	}
	Convey("Normal", t, func() {
		client, err := NewOssClient(endpoint, accessKeyID, accessKeySecret, "xgxw")
		So(err, ShouldBeNil)
		buf, err := client.GetObject(context.Background(), "todo.md")
		So(err, ShouldBeNil)
		fmt.Println(string(buf))
	})
}

func Test_OssClienta_PutObject(t *testing.T) {
	var endpoint = "oss-cn-beijing.aliyuncs.com"
	var accessKeyID = os.Getenv("accrss_key_id")
	var accessKeySecret = os.Getenv("accrss_key_secret")
	if accessKeyID == "" || accessKeySecret == "" {
		return
	}
	Convey("Normal", t, func() {
		client, err := NewOssClient(endpoint, accessKeyID, accessKeySecret, "xgxw")
		So(err, ShouldBeNil)
		err = client.PutObject(context.Background(), "todo.md", []byte("this is test content"))
		So(err, ShouldBeNil)
	})
}
