package storages

import (
	"context"
	"io"
	"net/url"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
	"github.com/xgxw/foundation-go"
)

// OSSOptions is Oss 存储所需的配置项
type OSSOptions struct {
	Endpoint        string `yaml:"endpoint" mapstructure:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id" mapstructure:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret" mapstructure:"access_key_secret"`
	Bucket          string `yaml:"bucket" mapstructure:"bucket"`
	UseHTTPS        bool   `yaml:"use_https" mapstructure:"use_https"`
	Domain          string `yaml:"domain" mapstructure:"domain"`
}

// OSSBucketInterface 是为了脱离bucket带来的影响.
// 如mock时, 如果使用 oss.Client 则无法mock(Client is struct)
type OSSBucketInterface interface {
	PutObject(objectKey string, reader io.Reader, options ...oss.Option) error
	GetObject(objectKey string, options ...oss.Option) (io.ReadCloser, error)
	DeleteObject(objectKey string, options ...oss.Option) error
}

// OSSStorage oss 客户端
type OSSStorage struct {
	name      string
	bucket    OSSBucketInterface
	accessURL *url.URL
}

var _ foundation.Storage = &OSSStorage{}

// NewOSSStorage is 创建oss客户端
func NewOSSStorage(opts *OSSOptions) (*OSSStorage, error) {
	ossClient, err := oss.New(opts.Endpoint, opts.AccessKeyID, opts.AccessKeySecret)
	if err != nil {
		return &OSSStorage{}, errors.Wrap(err, "create oss client error")
	}

	bucket, err := ossClient.Bucket(opts.Bucket)
	if err != nil {
		return &OSSStorage{}, errors.Wrap(err, "conn oss bucket error")
	}

	var buf strings.Builder
	buf.WriteString("oss:")
	buf.WriteString(opts.Endpoint)
	buf.WriteString(":")
	buf.WriteString(opts.Bucket)

	accessURL := initAccessURL(opts)

	return &OSSStorage{
		name:      buf.String(),
		bucket:    bucket,
		accessURL: accessURL,
	}, nil
}

func initAccessURL(opts *OSSOptions) *url.URL {
	// scheme
	var scheme = "http"
	if opts.UseHTTPS {
		scheme = "https"
	}

	// domain
	var domainBuf strings.Builder
	if len(opts.Domain) > 0 {
		domainBuf.WriteString(opts.Domain)
	} else {
		domainBuf.WriteString(opts.Bucket)
		domainBuf.WriteString(".")
		domainBuf.WriteString(opts.Endpoint)
	}

	return &url.URL{
		Scheme: scheme,
		Host:   domainBuf.String(),
	}
}

func (s *OSSStorage) Name() string {
	return s.name
}

func (s *OSSStorage) Save(ctx context.Context, fileID string, reader io.Reader) error {
	return s.bucket.PutObject(fileID, reader)
}

// Writer 可以直接将数据写入传入的结构体. 如传入 bytes.Buffer, 数据直接可写入buffer.
func (s *OSSStorage) Fetch(ctx context.Context, fileID string, writer io.Writer) error {
	reader, err := s.bucket.GetObject(fileID)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, reader)
	return err
}

func (s *OSSStorage) Del(ctx context.Context, fileID string) error {
	return s.bucket.DeleteObject(fileID)
}

func (s *OSSStorage) URL(ctx context.Context, fileID string) (*url.URL, error) {
	var accessURL = &url.URL{
		Scheme: s.accessURL.Scheme,
		Host:   s.accessURL.Host,
		Path:   fileID,
	}
	return accessURL, nil
}
