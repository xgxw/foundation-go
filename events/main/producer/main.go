package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/xgxw/foundation-go"
	"github.com/xgxw/foundation-go/events"
	"github.com/xgxw/foundation-go/log"
)

var topic = "tests2"
var host = "127.0.0.1:9092"

func main() {
	opts := events.KafkaManagerOptions{
		Host:               host,
		ErrChanSize:        10000,
		ThrottleBottleSize: 1,
	}
	opts.Producer.ReturnSuccess = true
	opts.Producer.ReturnErrors = true
	opts.Consumer.ReturnErrors = true
	opts.Consumer.GroupNotification = true

	logger := log.NewLogger(log.Options{
		Mode: "debug",
	}, os.Stdout)
	em, err := events.NewKafkaEventManager(opts, logger)
	if err != nil {
		logger.Fatalf("create kafka EventManager error. error:%w", err)
	}
	logger.Info("create kafka EventManager success")

	msg := &Event{"this is test msg"}
	err = em.Publish(context.Background(), msg.Topic(), msg)
	if err != nil {
		logger.Fatalf("publish message error. error:%w", err)
	}
	logger.Info("send message success")

	// 等待内部程序将消息推送出去
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

type Event struct {
	Msg string `json:"msg"`
}

var _ foundation.Event = new(Event)

func (e *Event) Marshal() []byte {
	b, _ := json.Marshal(e)
	return b
}
func (e *Event) Topic() string {
	return topic
}
