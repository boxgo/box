package kafka_test

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	kafka2 "github.com/boxgo/box/v2/client/kafka"
)

const (
	testTopic = "test"
)

func Example() {
	kfk := kafka2.StdConfig("default").Build()

	producer, err := kfk.NewSyncProducer()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	consumer, err := kfk.NewConsumer()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			panic(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(testTopic, 0, kafka2.OffsetNewest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var cnt int32

	go func() {
		for {
			select {
			case <-partitionConsumer.Messages():
				atomic.AddInt32(&cnt, 1)
			case <-signals:
				break
			}
		}
	}()

	partition, offset, err := producer.SendMessage(&kafka2.ProducerMessage{
		Topic: testTopic,
		Value: kafka2.StringEncoder("hi"),
	})
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)

	fmt.Println(offset >= 0, partition == 0, atomic.LoadInt32(&cnt) > 0)
	// Output: true true true
}
