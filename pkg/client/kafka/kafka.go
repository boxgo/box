package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/boxgo/box/pkg/logger"
)

type (
	Kafka struct {
		cfg    *Config
		client sarama.Client
	}
)

func newKafka(c *Config) *Kafka {
	client, err := sarama.NewClient(c.Addrs, c.kfkCfg)
	if err != nil {
		logger.Panicw("kafka build error", "err", err)
	}

	return &Kafka{
		cfg:    c,
		client: client,
	}
}

func (kfk Kafka) NewSyncProducer() (SyncProducer, error) {
	return sarama.NewSyncProducerFromClient(kfk.client)
}

func (kfk Kafka) NewAsyncProducer() (AsyncSyncProducer, error) {
	return sarama.NewAsyncProducerFromClient(kfk.client)
}

func (kfk Kafka) NewConsumer() (Consumer, error) {
	return sarama.NewConsumerFromClient(kfk.client)
}

func (kfk Kafka) NewConsumerGroup(groupID string) (ConsumerGroup, error) {
	return sarama.NewConsumerGroupFromClient(groupID, kfk.client)
}
