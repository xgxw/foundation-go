package main

import (
	"context"
	"fmt"
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
		ThrottleBottleSize: 2,
	}
	opts.Producer.ReturnSuccess = true
	opts.Producer.ReturnErrors = true
	opts.Consumer.ReturnErrors = true
	opts.Consumer.GroupNotification = false

	logger := log.NewLogger(log.Options{
		Mode: "debug",
	}, os.Stdout)
	em, err := events.NewKafkaEventManager(opts, logger)
	if err != nil {
		logger.Fatalf("create kafka EventManager error. error:%w", err)
	}
	logger.Info("create kafka EventManager success")

	err = em.Subscribe(context.Background(), "asda", topic, getEventHandler())
	if err != nil {
		logger.Fatalf("subscribe message error. error:%+v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	if err := em.Close(); err != nil {
		logger.Errorf("failed to close kafka manager error: %s", err)
	}
}

func getEventHandler() foundation.EventHandler {
	return func(ctx context.Context, msgs ...[]byte) error {
		fmt.Println("into handler")
		for _, msg := range msgs {
			fmt.Println("received message")
			fmt.Println(msg)
		}
		return nil
	}
}
