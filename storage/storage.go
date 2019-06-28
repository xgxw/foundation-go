package storage

import "context"

// StorageClientInterface 是为了对外界屏蔽storage内部具体的实现. 使用者通过使用 具体类(如ossClient) 初始化该接口,
// 各接口使用统一的标准. 好处是当更改storage选项时, 所造成的影响最小.
type StorageClientInterface interface {
	// GetObject is 以流的方式读取文件
	GetObject(ctx context.Context, fileID string) (buf []byte, err error)
	// PutObject is 上传文件
	PutObject(ctx context.Context, fileID string, buf []byte) (err error)
}
