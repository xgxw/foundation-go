package storage

import (
	"context"
)

// ListOption 可选参数.
type ListOption int32

const (
	// ListOptionReverse is 表示是否递归显示目录下所有文件
	ListOptionReverse ListOption = 1 << 0
)

// ClientInterface 是为了对外界屏蔽storage内部具体的实现. 使用者通过使用 具体类(如ossClient) 初始化该接口,
// 各接口使用统一的标准. 好处是当更改storage选项时, 所造成的影响最小.
type ClientInterface interface {
	// GetObject is 以流的方式读取文件
	GetObject(ctx context.Context, fileID string) (buf []byte, err error)
	// PutObject is 上传文件, 可以新增/更新文件
	PutObject(ctx context.Context, fileID string, buf []byte) (err error)
	// DelObject is 删除文件, 文件夹也视为object, 在其后添加 `/` 标识文件夹
	DelObject(ctx context.Context, fileID string) (err error)
	// DelObjects is 删除多个文件
	DelObjects(ctx context.Context, fileIDs []string) (err error)
	// GetList is 获取文件列表, 返回文件目录的json格式. ops 使用位标识配置, 将需要的配置 或运算 即可.
	GetList(ctx context.Context, path string, ops ListOption) (buf []byte, err error)
}
