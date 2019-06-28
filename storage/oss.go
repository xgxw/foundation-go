package storage

import (
	"bytes"
	"context"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssClient struct {
	*oss.Client
	bucket          *oss.Bucket
	endpoint        string
	accessKeyID     string
	accessKeySecret string
}

var _ StorageClientInterface = &OssClient{}

func NewOssClient(endpoint, accessKeyID, accessKeySecret, bucketID string) (client StorageClientInterface, err error) {
	ossClient, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return client, err
	}
	bucket, err := ossClient.Bucket(bucketID)
	if err != nil {
		return client, err
	}
	return &OssClient{
		Client:          ossClient,
		bucket:          bucket,
		endpoint:        endpoint,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
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
