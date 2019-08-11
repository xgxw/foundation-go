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
	var accessKeyID = os.Getenv("access_key_id")
	var accessKeySecret = os.Getenv("access_key_secret")
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

func Test_OssClient_GetObject(t *testing.T) {
	opts, err := getOssOptions()
	if err != nil {
		fmt.Println(err)
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

func Test_OssClient_PutObject(t *testing.T) {
	opts, err := getOssOptions()
	if err != nil {
		fmt.Println(err)
		return
	}
	Convey("Normal", t, func() {
		client, err := NewOssClient(opts)
		So(err, ShouldBeNil)
		err = client.PutObject(context.Background(), "test.md", []byte("this is test content"))
		So(err, ShouldBeNil)
	})
}

func Test_OssClient_DelObject(t *testing.T) {
	opts, err := getOssOptions()
	if err != nil {
		fmt.Println(err)
		return
	}
	Convey("Normal", t, func() {
		client, err := NewOssClient(opts)
		So(err, ShouldBeNil)
		err = client.DelObject(context.Background(), "test/")
		So(err, ShouldBeNil)
	})
}

func Test_OssClient_DelObjects(t *testing.T) {
	opts, err := getOssOptions()
	if err != nil {
		fmt.Println(err)
		return
	}
	Convey("Normal", t, func() {
		client, err := NewOssClient(opts)
		So(err, ShouldBeNil)
		err = client.DelObjects(context.Background(), []string{"test/"})
		So(err, ShouldBeNil)
	})
}

func Test_OssClient_GetList(t *testing.T) {
	opts, err := getOssOptions()
	if err != nil {
		fmt.Println(err)
		return
	}
	Convey("Normal", t, func() {
		client, err := NewOssClient(opts)
		So(err, ShouldBeNil)
		files, err := client.GetList(context.Background(), "", ListOptionReverse)
		fmt.Println(string(files))
		So(err, ShouldBeNil)
	})
}
