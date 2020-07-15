package foundation

import (
	"context"
)

type Event interface {
	Marshal() []byte
	Topic() string
}

type EventOptions interface {
	Partition() int32
	SetPartition(partition int32)
}

type EventPublishOption func(options EventOptions)

type EventHandler func(context.Context, ...[]byte) error

type EventManager interface {
	Publish(ctx context.Context, topic string, event Event, opts ...EventPublishOption) error
	Subscribe(ctx context.Context, group, topic string, handler EventHandler) error
	ProducerErrors() <-chan error
	ConsumerErrors() <-chan error
	Close() error
}

type EventInspector interface {
	RequeueEvent(ctx context.Context, targetType string, durationSec int64) error
	EventAck(ctx context.Context, e Event) error
}
