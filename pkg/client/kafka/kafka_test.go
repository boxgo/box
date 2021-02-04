package kafka_test

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
	"github.com/boxgo/box/pkg/client/kafka"
)

const (
	testTopic = "wechat_event"
)

func Example() {
	kfk := kafka.StdConfig("default").Build()

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

	partitionConsumer, err := consumer.ConsumePartition(testTopic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	cnt := 0
	go func() {
		for {
			select {
			case <-partitionConsumer.Messages():
				cnt++
			case <-signals:
				break
			}
		}
	}()

	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: testTopic,
		Value: sarama.StringEncoder("hi"),
	})
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)

	fmt.Println(offset > 0, partition == 0, cnt > 0)
	// Output: true true true
}
