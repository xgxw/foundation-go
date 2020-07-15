package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToContext(t *testing.T) {
	entry := NewEntry(DefaultLogger)

	ctx := ToContext(context.Background(), entry)
	actual := ctx.Value(ctxLoggerKey)

	assert.Equalf(t, entry, actual, "logrus entry should match")
}

func TestExtract(t *testing.T) {
	ctx := context.Background()
	got := Extract(ctx)

	assert.NotNil(t, got, "extract should give a valid entry")

	entry := NewEntry(DefaultLogger)
	ctx = ToContext(ctx, entry)
	actual := Extract(ctx)

	assert.Equalf(t, entry, actual, "logrus entry should match")
}
