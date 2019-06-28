package storage

import (
	"bytes"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssClient struct {
	*oss.Client
	endpoint        string
	accessKeyID     string
	accessKeySecret string
}

type OssClientService interface {
	// GetObject is 以流的方式读取文件
	GetObject(bucketID, fileID string) (buf []byte, err error)
	// PutObject is 上传文件
	PutObject(bucketID, fileID string, buf []byte) (err error)
}

var _ OssClientService = &OssClient{}

func NewOssClient(endpoint, accessKeyID, accessKeySecret string) (client OssClientService, err error) {
	ossClient, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return client, err
	}
	return &OssClient{
		Client:          ossClient,
		endpoint:        endpoint,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
	}, nil
}

func (o *OssClient) getBucket(bucketID string) (bucket *oss.Bucket, err error) {
	bucket, err = o.Bucket(bucketID)
	return bucket, err
}

func (o *OssClient) GetObject(bucketID, fileID string) (buf []byte, err error) {
	buf = make([]byte, 0)
	bucket, err := o.getBucket(bucketID)
	if err != nil {
		return buf, err
	}
	body, err := bucket.GetObject(fileID)
	if err != nil {
		return buf, err
	}
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, body)
	buf = buffer.Bytes()
	return buf, err
}

func (o *OssClient) PutObject(bucketID, fileID string, buf []byte) (err error) {
	bucket, err := o.getBucket(bucketID)
	if err != nil {
		return err
	}
	err = bucket.PutObject(fileID, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	return nil
}
