package storage

import (
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_OssClienta_GetObject(t *testing.T) {
	var endpoint = "oss-cn-beijing.aliyuncs.com"
	var accessKeyID = os.Getenv("accrss_key_id")
	var accessKeySecret = os.Getenv("accrss_key_secret")
	Convey("Normal", t, func() {
		client, err := NewOssClient(endpoint, accessKeyID, accessKeySecret)
		So(err, ShouldBeNil)
		buf, err := client.GetObject("xgxw", "todo.md")
		So(err, ShouldBeNil)
		fmt.Println(string(buf))
	})
}
