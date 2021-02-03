package kafka

import (
	"github.com/Shopify/sarama"
)

type (
	SyncProducer      sarama.SyncProducer
	AsyncSyncProducer sarama.AsyncProducer
	Consumer          sarama.Consumer
	ConsumerGroup     sarama.ConsumerGroup
	// ProduceRequest                sarama.ProduceRequest
	// ProduceResponse               sarama.ProduceResponse
	// ProduceResponseBlock          sarama.ProduceResponseBlock
	// ProducerError                 sarama.ProducerError
	// ProducerErrors                sarama.ProducerErrors
	// ProducerInterceptor           sarama.ProducerInterceptor
	// ProducerMessage               sarama.ProducerMessage
	// InitProducerIDRequest         sarama.InitProducerIDRequest
	// InitProducerIDResponse        sarama.InitProducerIDResponse
	// ConsumerError                 sarama.ConsumerError
	// ConsumerErrors                sarama.ConsumerErrors
	// ConsumerGroupClaim            sarama.ConsumerGroupClaim
	// ConsumerGroupHandler          sarama.ConsumerGroupHandler
	// ConsumerGroupMemberAssignment sarama.ConsumerGroupMemberAssignment
	// ConsumerGroupMemberMetadata   sarama.ConsumerGroupMemberMetadata
	// ConsumerGroupSession          sarama.ConsumerGroupSession
	// ConsumerInterceptor           sarama.ConsumerInterceptor
	// ConsumerMessage               sarama.ConsumerMessage
	// ConsumerMetadataRequest       sarama.ConsumerMetadataRequest
	// ConsumerMetadataResponse      sarama.ConsumerMetadataResponse
)
