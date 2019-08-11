package storage

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/everywan/foundation-go/internal/utils"
)

// OssOptions is Oss 存储所需的配置想
type OssOptions struct {
	Endpoint        string `yaml:"endpoint" mapstructure:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id" mapstructure:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret" mapstructure:"access_key_secret"`
	Bucket          string `yaml:"bucket" mapstructure:"bucket"`
}

// OssClient oss 客户端
type OssClient struct {
	*oss.Client
	bucket *oss.Bucket
}

var _ ClientInterface = &OssClient{}

// NewOssClient is 创建oss客户端
func NewOssClient(opts *OssOptions) (client ClientInterface, err error) {
	ossClient, err := oss.New(opts.Endpoint, opts.AccessKeyID, opts.AccessKeySecret)
	if err != nil {
		return client, err
	}
	bucket, err := ossClient.Bucket(opts.Bucket)
	if err != nil {
		return client, err
	}
	return &OssClient{
		Client: ossClient,
		bucket: bucket,
	}, nil
}

// GetObject is 根据fileid从oss获取文件
func (o *OssClient) GetObject(ctx context.Context, fileID string) (buf []byte, err error) {
	buf = make([]byte, 0)
	body, err := o.bucket.GetObject(fileID)
	if err != nil {
		return buf, err
	}
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, body)
	buf = buffer.Bytes()
	return buf, err
}

// PutObject is 新建或上传文件到oss
func (o *OssClient) PutObject(ctx context.Context, fileID string, buf []byte) (err error) {
	return o.bucket.PutObject(fileID, bytes.NewReader(buf))
}

// DelObject is 删除 oss 上的一个文件
func (o *OssClient) DelObject(ctx context.Context, fileID string) (err error) {
	return o.bucket.DeleteObject(fileID)
}

// DelObjects is 删除 oss 上的多个文件
func (o *OssClient) DelObjects(ctx context.Context, fileIDs []string) (err error) {
	_, err = o.bucket.DeleteObjects(fileIDs)
	return err
}

// GetList is 获取指定文件夹下所有的文件列表, ListOption 参考 storage 中定义,
// 采用位运算添加配置. 返回结果为 json 格式
func (o *OssClient) GetList(ctx context.Context, path string, ops ListOption) (buf []byte, err error) {
	delimter := "/"
	if ops&ListOptionReverse > 0 {
		delimter = ""
	}
	lsRes, err := o.bucket.ListObjects(oss.Prefix(path), oss.Delimiter(delimter))
	if err != nil {
		return nil, err
	}
	files := make([]string, len(lsRes.Objects))
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	for i, file := range lsRes.Objects {
		files[i] = strings.Replace(file.Key, path, "", 1)
	}
	return utils.ParseOssLsPaths(files, delimter)
}
