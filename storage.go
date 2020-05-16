package foundation

import (
	"context"
	"io"
	"net/url"
)

type Storage interface {
	Name() string
	Save(ctx context.Context, fileID string, reader io.Reader) error
	Fetch(ctx context.Context, fileID string, writer io.Writer) error
	Del(ctx context.Context, fileID string) error
	URL(ctx context.Context, fileID string) (*url.URL, error)
}
