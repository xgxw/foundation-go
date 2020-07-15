package events

import (
	"context"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/pkg/errors"
	"github.com/xgxw/foundation-go"
	"github.com/xgxw/foundation-go/log"
)

// KafkaConsumerMaker is a kafka consumer factory
type KafkaConsumerMaker struct {
	host         string
	consumerConf *cluster.Config
}

// Consumer create a kafka cluster consumer by consumer maker
func (kcm *KafkaConsumerMaker) Consumer(group, topic string) (*KafkaConsumer, error) {
	c, err := cluster.NewConsumer([]string{kcm.host}, group, []string{topic}, kcm.consumerConf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a consumer")
	}
	return &KafkaConsumer{
		c:        c,
		stopChan: make(chan bool, 1),
	}, nil
}

type KafkaManagerOptions struct {
	// yaml
	Host               string `mapstructure:"host"`
	ErrChanSize        int    `mapstructure:"err_chan_size"`
	ThrottleBottleSize int    `mapstructure:"throttle_bottle_size"`
	Producer           struct {
		RequiredACKs    int  `mapstructure:"required_acks"`
		MaxMessageBytes int  `mapstructure:"max_message_bytes"`
		Compression     int  `mapstructure:"compression"`
		FlushFrequency  int  `mapstructure:"flush_frequency"`
		ReturnSuccess   bool `mapstructure:"return_success"`
		ReturnErrors    bool `mapstructure:"return_errors"`
	} `mapstructure:"producer"`

	Consumer struct {
		GroupMode               int   `mapstructure:"group_mode"`
		GroupNotification       bool  `mapstructure:"group_notification"`
		OffsetInitial           int   `mapstructure:"offset_initial"`
		ReturnErrors            bool  `mapstructure:"return_errors"`
		BatchConsumptionSize    int   `mapstructure:"batch_consumption_size"`
		BatchConsumptionTimeout int64 `mapstructure:"batch_consumption_timeout"`
	} `mapstructure:"consumer"`
}

// KafkaManager implement event manager interface
type KafkaManager struct {
	// common arguments
	logger *log.Logger
	host   string
	wg     *sync.WaitGroup

	// producer
	producer        sarama.AsyncProducer
	producerConf    *sarama.Config
	producerErrChan chan error
	throttleBottle  chan bool

	// consumer
	consumerMaker           *KafkaConsumerMaker
	consumerGroup           []*KafkaConsumer
	cgMutex                 *sync.Mutex
	consumerErrChan         chan error
	batchConsumptionSize    int
	batchConsumptionTimeout int64
}

// NewKafkaEventManager create a kafka event manager instance
func NewKafkaEventManager(opts KafkaManagerOptions, l *log.Logger) (foundation.EventManager, error) {
	// kafka producer
	config := sarama.NewConfig()
	config.Version = sarama.V0_10_0_0
	config.Producer.Return.Successes = opts.Producer.ReturnSuccess
	config.Producer.Return.Errors = opts.Producer.ReturnErrors
	config.Producer.Partitioner = newKafkaPartitioner
	if opt, ok := requiredACKs[opts.Producer.RequiredACKs]; ok {
		config.Producer.RequiredAcks = opt
	} else {
		config.Producer.RequiredAcks = sarama.WaitForLocal
	}
	if opts.Producer.MaxMessageBytes < 1048576 {
		config.Producer.MaxMessageBytes = 1048576
	} else {
		config.Producer.MaxMessageBytes = opts.Producer.MaxMessageBytes
	}
	if opt, ok := compressionOption[opts.Producer.Compression]; ok {
		config.Producer.Compression = opt
	} else {
		config.Producer.Compression = sarama.CompressionNone
	}
	if opts.Producer.FlushFrequency < 100 {
		config.Producer.Flush.Frequency = time.Millisecond * 100
	} else {
		config.Producer.Flush.Frequency = time.Millisecond * time.Duration(opts.Producer.FlushFrequency)
	}
	async, err := sarama.NewAsyncProducer([]string{opts.Host}, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a kafka producer")
	}

	// kafka consumer
	consumerConf := cluster.NewConfig()
	consumerConf.Version = sarama.V0_10_0_0
	consumerConf.Group.Return.Notifications = opts.Consumer.GroupNotification
	consumerConf.Consumer.Return.Errors = opts.Consumer.ReturnErrors

	if opt, ok := consumerGroupMode[opts.Consumer.GroupMode]; ok {
		consumerConf.Group.Mode = opt
	} else {
		consumerConf.Group.Mode = cluster.ConsumerModeMultiplex
	}

	if opt, ok := consumerOffsetInitial[opts.Consumer.OffsetInitial]; ok {
		consumerConf.Consumer.Offsets.Initial = opt
	} else {
		consumerConf.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	consumerConf.Consumer.Offsets.CommitInterval = time.Second

	if opts.Consumer.BatchConsumptionSize <= 0 {
		opts.Consumer.BatchConsumptionSize = 100
	}

	if opts.Consumer.BatchConsumptionTimeout <= 0 {
		opts.Consumer.BatchConsumptionTimeout = 500
	}

	km := &KafkaManager{
		logger: l,
		host:   opts.Host,
		wg:     &sync.WaitGroup{},
		// Producer
		producer:        async,
		producerConf:    config,
		producerErrChan: make(chan error, opts.ErrChanSize),
		throttleBottle:  make(chan bool, opts.ThrottleBottleSize),
		// Consumer
		consumerMaker: &KafkaConsumerMaker{
			host:         opts.Host,
			consumerConf: consumerConf,
		},
		cgMutex:                 &sync.Mutex{},
		consumerErrChan:         make(chan error, opts.ErrChanSize),
		batchConsumptionSize:    opts.Consumer.BatchConsumptionSize,
		batchConsumptionTimeout: opts.Consumer.BatchConsumptionTimeout,
	}
	km.producerEventHandler()

	return km, nil
}

func (km *KafkaManager) producerEventHandler() {
	km.wg.Add(1)
	go func() {
		defer km.wg.Done()
		for msg := range km.producer.Successes() {
			km.logger.WithField("topic", msg.Topic).
				WithField("partition", msg.Partition).
				WithField("offset", msg.Offset).
				Debugf("Publish success")
			_ = <-km.throttleBottle
		}
	}()

	km.wg.Add(1)
	go func() {
		defer km.wg.Done()
		for pErr := range km.producer.Errors() {
			km.logger.WithField("topic", pErr.Msg.Topic).WithError(pErr).Error("failed to publish")
			km.producerErrChan <- errors.Wrap(pErr.Err, "failed to publish")
			_ = <-km.throttleBottle
		}
	}()
}

// Publish a message to Kafka by specific topic
func (km *KafkaManager) Publish(ctx context.Context, topic string, event foundation.Event, opts ...foundation.EventPublishOption) error {
	km.throttleBottle <- true

	pMsg := &sarama.ProducerMessage{
		Topic:     event.Topic(),
		Partition: -1,
		Value:     sarama.ByteEncoder(event.Marshal()),
	}

	if len(opts) > 0 {
		options := &kafkaEventOptions{pMsg}
		for _, opt := range opts {
			opt(options)
		}
	}

	km.producer.Input() <- pMsg
	return nil
}

// Subscribe a topic from Kafka
func (km *KafkaManager) Subscribe(ctx context.Context, group, topic string, handler foundation.EventHandler) error {
	c, err := km.consumerMaker.Consumer(group, topic)
	if err != nil {
		return errors.Wrap(err, "failed to create event consumer")
	}
	c.SetBatchConsumptionSize(km.batchConsumptionSize)
	c.SetBatchConsumptionHandleTimeout(km.batchConsumptionTimeout)

	logger := km.logger.NewEntry().
		WithScope("event_consumer").
		WithContext(log.ContextFields{
			TargetID:   topic,
			TargetType: "topic",
			ObjectID:   group,
		})

	km.wg.Add(1)
	go c.Do(log.ToContext(ctx, logger), km.wg, handler, km.consumerErrChan)

	km.cgMutex.Lock()
	km.consumerGroup = append(km.consumerGroup, c)
	km.cgMutex.Unlock()
	return nil
}

func (km *KafkaManager) ProducerErrors() <-chan error {
	return km.producerErrChan
}

func (km *KafkaManager) ConsumerErrors() <-chan error {
	return km.consumerErrChan
}

// Close is that close message consumer and publisher safety
func (km *KafkaManager) Close() error {
	// stop producer && consumer group
	km.producer.AsyncClose()
	km.cgMutex.Lock()
	for _, c := range km.consumerGroup {
		c.Close()
	}
	km.cgMutex.Unlock()

	// binding waitGroup with a channel
	cb := make(chan bool, 1)
	go func() {
		km.wg.Wait()
		cb <- true
		close(cb)
	}()

	// waiting for producer and consumer finish their job in processing
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	select {
	case <-cb:
		err = nil
	case <-ctx.Done():
		err = errors.New("closing event manager timeout")
	}

	km.producer = nil
	km.consumerGroup = km.consumerGroup[:0]
	close(km.producerErrChan)
	close(km.consumerErrChan)
	return err
}

type kafkaPartitioner struct {
	roundRobin sarama.Partitioner
}

func newKafkaPartitioner(topic string) sarama.Partitioner {
	return &kafkaPartitioner{
		roundRobin: sarama.NewRoundRobinPartitioner(topic),
	}
}

func (kp *kafkaPartitioner) Partition(message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	if message.Partition >= 0 {
		return message.Partition, nil
	}
	return kp.roundRobin.Partition(message, numPartitions)
}

func (kp *kafkaPartitioner) RequiresConsistency() bool {
	return true
}

func (kp *kafkaPartitioner) MessageRequiresConsistency(message *sarama.ProducerMessage) bool {
	if message.Partition >= 0 {
		return true
	}
	return kp.roundRobin.RequiresConsistency()
}
