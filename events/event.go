package events

import (
	"github.com/Shopify/sarama"
	"github.com/xgxw/foundation-go"
)

type kafkaEventOptions struct {
	*sarama.ProducerMessage
}

func (opts *kafkaEventOptions) Partition() int32 {
	return opts.ProducerMessage.Partition
}

func (opts *kafkaEventOptions) SetPartition(partition int32) {
	opts.ProducerMessage.Partition = partition
}

// SetPartition is action to specify a kafka partition to publish message
func SetPartition(partition int32) foundation.EventPublishOption {
	return func(options foundation.EventOptions) {
		options.SetPartition(partition)
	}
}
