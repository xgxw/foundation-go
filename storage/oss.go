package storage

import (
	"bytes"
	"context"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssOptions struct {
	Endpoint        string `yaml:"endpoint" mapstructure:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id" mapstructure:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret" mapstructure:"access_key_secret"`
	Bucket          string `yaml:"bucket" mapstructure:"bucket"`
}

type OssClient struct {
	*oss.Client
	bucket *oss.Bucket
}

var _ StorageClientInterface = &OssClient{}

func NewOssClient(opts *OssOptions) (client StorageClientInterface, err error) {
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

func (o *OssClient) PutObject(ctx context.Context, fileID string, buf []byte) (err error) {
	err = o.bucket.PutObject(fileID, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	return nil
}
