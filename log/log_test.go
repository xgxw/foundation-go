package log

import (
	"bytes"
	"testing"

	assertpkg "github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	assert := assertpkg.New(t)

	testCases := []struct {
		Level string
		Msg   string
	}{
		{
			Level: "warn",
		},
	}
	for _, tc := range testCases {
		opts := Options{
			Level: tc.Level,
		}
		output := bytes.NewBuffer([]byte{})
		l := NewLogger(opts, output)
		if !assert.NotNil(l, "failed to create logger") {
			t.FailNow()
		}
		l.Debug("debug")
		assert.Empty(output.String(), "unexpected logger level output")
		output.Reset()

		l.Warn("warn")
		assert.NotEmpty(output.String(), "unexpected logger level output")
		output.Reset()

		l.Warn("error")
		assert.NotEmpty(output.String(), "unexpected logger level output")
		output.Reset()
	}
}
