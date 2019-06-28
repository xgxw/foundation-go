package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func getOssOptions() (opts *OssOptions, err error) {
	var endpoint = "oss-cn-beijing.aliyuncs.com"
	var accessKeyID = os.Getenv("accrss_key_id")
	var accessKeySecret = os.Getenv("accrss_key_secret")
	if accessKeyID == "" || accessKeySecret == "" {
		return opts, errors.New("can't find accessKey")
	}
	return &OssOptions{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		Bucket:          "xgxw",
	}, nil
}

func Test_OssClienta_GetObject(t *testing.T) {
	opts, err := getOssOptions()
	if err != nil {
		return
	}
	Convey("Normal", t, func() {
		client, err := NewOssClient(opts)
		So(err, ShouldBeNil)
		buf, err := client.GetObject(context.Background(), "todo.md")
		So(err, ShouldBeNil)
		fmt.Println(string(buf))
	})
}

func Test_OssClienta_PutObject(t *testing.T) {
	opts, err := getOssOptions()
	if err != nil {
		return
	}
	Convey("Normal", t, func() {
		client, err := NewOssClient(opts)
		So(err, ShouldBeNil)
		err = client.PutObject(context.Background(), "todo.md", []byte("this is test content"))
		So(err, ShouldBeNil)
	})
}
