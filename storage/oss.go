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
	GetObject(name string) (err error)
	// PutObject is 上传文件
	PutObject(name, content string) (err error)
}

func NewOssClient(endpoint, accessKeyID, accessKeySecret string) (client *OssClient, err error) {
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

func (o *OssClient) getBucket(bucketName string) (bucket *oss.Bucket, err error) {
	bucket, err = o.Bucket(bucketName)
	return bucket, err
}

func (o *OssClient) GetObject(bucketName, fileName string) (buf []byte, err error) {
	buf = make([]byte, 0)
	bucket, err := o.getBucket(bucketName)
	if err != nil {
		return buf, err
	}
	body, err := bucket.GetObject(fileName)
	if err != nil {
		return buf, err
	}
	defer body.Close()
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, body)
	buf = buffer.Bytes()
	return buf, err
}

func (o *OssClient) PutObject(bucketName, fileName string, buf []byte) (err error) {
	bucket, err := o.getBucket(bucketName)
	if err != nil {
		return err
	}
	err = bucket.PutObject(fileName, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	return nil
}
