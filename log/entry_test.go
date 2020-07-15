package log

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewEntry(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	l := NewLogger(Options{Level: "debug"}, output)

	logger := NewEntry(l).
		WithScope("foundation").
		WithContext(ContextFields{
			RequestID:  "bkti6vtcdckp19rehim0",
			TargetID:   "460",
			TargetType: "merchant",
		})

	// passing logger with context.Context
	ctx := context.Background()
	ctx = ToContext(ctx, logger)

	logger = Extract(ctx).WithClient(ProtocolHTTP, ClientFields{
		URI:    "http://gateway.lehuipay.com/v1/merchants",
		Method: "POST",
		Kind:   "client",
	})

	logger.WithRequest(ProtocolHTTP, PayloadFields{
		Content: `{"name":"Test Merchant"}`,
	}).Debug("sending http request")

	start := time.Now()
	// do the HTTP request
	end := time.Now()

	// record some metrics
	logger.WithMetrics(MetricsFields{
		Name:  "internal_rpc_call",
		Value: end.UnixNano() - start.UnixNano(),
	}).Debug("internal call finished")

	logger.WithResponse(ProtocolHTTP, PayloadFields{
		Content: `{"id":"xxx","name":"Test Merchant"}`,
	}).Debug("http response received")

	fmt.Println(output.String())
}

func TestEntry_ExtractContext(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	l := NewLogger(Options{Level: "debug"}, output)
	entry := l.NewEntry()
	ctx := ContextFields{
		RequestID:  "bkti6vtcdckp19rehim0",
		TargetID:   "460",
		TargetType: "merchant",
	}

	require.NotEqual(t, ctx, entry.ExtractContext(), "extracted empty context should not match")

	entry = entry.WithContext(ctx)
	require.Equal(t, ctx, entry.ExtractContext(), "extracted context should match")
}

func TestEntry_ExtractServer(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	l := NewLogger(Options{Level: "debug"}, output)
	entry := l.NewEntry()
	server := ServerFields{
		URI:  "/v1/merchants",
		Host: "m.lehuipay.com",
	}

	require.NotEqual(t, server, entry.ExtractServer(ProtocolHTTP), "extracted empty server should not match")

	entry = entry.WithServer(ProtocolHTTP, server)

	// WithServer sets Kind to "server"
	server.Kind = "server"

	require.Equal(t, server, entry.ExtractServer(ProtocolHTTP), "extracted server should match")
	require.NotEqual(t, server, entry.ExtractServer(ProtocolTCP), "extracted server which was not set should not match")
}

func TestEntry_ExtractClient(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	l := NewLogger(Options{Level: "debug"}, output)
	entry := l.NewEntry()
	client := ClientFields{
		URI:  "/v1/merchants",
		Host: "m.lehuipay.com",
	}

	require.NotEqual(t, client, entry.ExtractClient(ProtocolHTTP), "extracted empty client should not match")

	entry = entry.WithClient(ProtocolHTTP, client)

	// WithClient sets Kind to "client"
	client.Kind = "client"

	require.Equal(t, client, entry.ExtractClient(ProtocolHTTP), "extracted client should match")
	require.NotEqual(t, client, entry.ExtractClient(ProtocolTCP), "extracted client which was not set should not match")
}
