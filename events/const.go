package events

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

var (
	requiredACKs = map[int]sarama.RequiredAcks{
		1:  sarama.WaitForLocal,
		-1: sarama.WaitForAll,
	}

	compressionOption = map[int]sarama.CompressionCodec{
		0: sarama.CompressionNone,
		1: sarama.CompressionGZIP,
		2: sarama.CompressionSnappy,
		3: sarama.CompressionLZ4,
	}

	consumerGroupMode = map[int]cluster.ConsumerMode{
		0: cluster.ConsumerModeMultiplex,
		1: cluster.ConsumerModePartitions,
	}
	consumerOffsetInitial = map[int]int64{
		-1: sarama.OffsetNewest,
		-2: sarama.OffsetOldest,
	}
)
