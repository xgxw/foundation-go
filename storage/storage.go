package storage

import (
	"context"
)

// ListOption 可选参数.
type ListOption int32

const (
	ListOption_Reverse ListOption = 1 << 0
)

// StorageClientInterface 是为了对外界屏蔽storage内部具体的实现. 使用者通过使用 具体类(如ossClient) 初始化该接口,
// 各接口使用统一的标准. 好处是当更改storage选项时, 所造成的影响最小.
type StorageClientInterface interface {
	// GetObject is 以流的方式读取文件
	GetObject(ctx context.Context, fileID string) (buf []byte, err error)
	// PutObject is 上传文件
	PutObject(ctx context.Context, fileID string, buf []byte) (err error)
	// GetList is 获取文件列表, 返回文件目录的json格式. ops 使用位标识配置, 将需要的配置 或运算 即可.
	GetList(ctx context.Context, path string, ops ListOption) (buf []byte, err error)
}
