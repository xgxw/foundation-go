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

// KafkaConsumer implement EventConsumer interface
type KafkaConsumer struct {
	logger                  *log.Logger
	c                       *cluster.Consumer
	stopChan                chan bool
	batchConsumptionSize    int
	batchConsumptionTimeout int64 // millisecond
}

// Do implement EventConsumer interface
func (kc *KafkaConsumer) Do(ctx context.Context, wg *sync.WaitGroup, h foundation.EventHandler, errQueue chan error) {
	defer wg.Done()
	logger := log.Extract(ctx)
	msgArr := make([][]byte, kc.batchConsumptionSize)
	ticker := time.NewTicker(time.Duration(kc.batchConsumptionTimeout) * time.Millisecond)

	var (
		pendingCommitMessage *sarama.ConsumerMessage
		msgArrCurrCount      int
	)

	markOffset := func() {
		kc.c.MarkOffset(pendingCommitMessage, "")
		logger.Debug("consumer succeed to process kafka message")
		// clear
		for i := 0; i < len(msgArr); i++ {
			msgArr[i] = nil
		}
		msgArrCurrCount = 0
	}

loop:
	for {
		select {
		case <-kc.stopChan:
			if err := kc.c.Close(); err != nil {
				errQueue <- errors.Wrap(err, "failed to close consumer")
			}
			break loop
		case err := <-kc.c.Errors():
			logger.WithError(err).Error("received error message when consuming event")
			errQueue <- err
		case ntf := <-kc.c.Notifications():
			logger.WithField("notification type", ntf.Type.String()).Warnf("Notification: %+v", ntf)
		case msg, ok := <-kc.c.Messages():
			if ok {
				pendingCommitMessage = msg
				msgArr[msgArrCurrCount] = msg.Value
				msgArrCurrCount++

				// add msg to arr
				if msgArrCurrCount == kc.batchConsumptionSize {
					// trigger to handle a batch of message once in handler
					if err := h(ctx, msgArr...); err != nil {
						errQueue <- errors.Wrap(err, "event handler failed to processing event message")
						continue
					}
					markOffset()
				}

				// notify event status through channel
				if c, ok := EventNotifyFromContext(ctx); ok {
					c <- true
				}
			}
		case <-ticker.C:
			if msgArrCurrCount > 0 {
				// trigger to handle a batch of message once in handler
				if err := h(ctx, msgArr[0:msgArrCurrCount]...); err != nil {
					errQueue <- errors.Wrap(err, "event handler failed to processing event message")
					continue
				}
				markOffset()
			}
		}
	}
}

func (kc *KafkaConsumer) SetBatchConsumptionSize(size int) *KafkaConsumer {
	kc.batchConsumptionSize = size
	return kc
}

func (kc *KafkaConsumer) SetBatchConsumptionHandleTimeout(timeout int64) *KafkaConsumer {
	kc.batchConsumptionTimeout = timeout
	return kc
}

// Close kafka consumer
func (kc *KafkaConsumer) Close() {
	kc.stopChan <- true
	close(kc.stopChan)
}
