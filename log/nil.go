package log

import "io"

var NilWriter io.Writer = new(nilWriter)

type nilWriter struct{}

func (n *nilWriter) Write(b []byte) (int, error) {
	return len(b), nil
}
