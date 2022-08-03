package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/boxgo/box/v2/logger"
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
	} else {
		logger.Debugw("kafka build config", "config", c)
	}

	return &Kafka{
		cfg:    c,
		client: client,
	}
}

func (kfk Kafka) NewSyncProducer() (SyncProducer, error) {
	return sarama.NewSyncProducerFromClient(kfk.client)
}

func (kfk Kafka) NewAsyncProducer() (AsyncProducer, error) {
	return sarama.NewAsyncProducerFromClient(kfk.client)
}

func (kfk Kafka) NewConsumer() (Consumer, error) {
	return sarama.NewConsumerFromClient(kfk.client)
}

func (kfk Kafka) NewConsumerGroup(groupID string) (ConsumerGroup, error) {
	return sarama.NewConsumerGroupFromClient(groupID, kfk.client)
}

// Config returns the Config struct of the client. This struct should not be
// altered after it has been created.
func (kfk Kafka) Config() *ConfigKafka {
	return kfk.client.Config()
}

// Controller returns the cluster controller broker. It will return a
// locally cached value if it's available. You can call RefreshController
// to update the cached value. Requires Kafka 0.10 or higher.
func (kfk Kafka) Controller() (*Broker, error) {
	return kfk.client.Controller()
}

// RefreshController retrieves the cluster controller from fresh metadata
// and stores it in the local cache. Requires Kafka 0.10 or higher.
func (kfk Kafka) RefreshController() (*Broker, error) {
	return kfk.client.RefreshController()
}

// Brokers returns the current set of active brokers as retrieved from cluster metadata.
func (kfk Kafka) Brokers() []*Broker {
	return kfk.client.Brokers()
}

// Topics returns the set of available topics as retrieved from cluster metadata.
func (kfk Kafka) Topics() ([]string, error) {
	return kfk.client.Topics()
}

// Partitions returns the sorted list of all partition IDs for the given topic.
func (kfk Kafka) Partitions(topic string) ([]int32, error) {
	return kfk.client.Partitions(topic)
}

// WritablePartitions returns the sorted list of all writable partition IDs for
// the given topic, where "writable" means "having a valid leader accepting
// writes".
func (kfk Kafka) WritablePartitions(topic string) ([]int32, error) {
	return kfk.client.WritablePartitions(topic)
}

// Leader returns the broker object that is the leader of the current
// topic/partition, as determined by querying the cluster metadata.
func (kfk Kafka) Leader(topic string, partitionID int32) (*Broker, error) {
	return kfk.client.Leader(topic, partitionID)
}

// Replicas returns the set of all replica IDs for the given partition.
func (kfk Kafka) Replicas(topic string, partitionID int32) ([]int32, error) {
	return kfk.client.Replicas(topic, partitionID)
}

// InSyncReplicas returns the set of all in-sync replica IDs for the given
// partition. In-sync replicas are replicas which are fully caught up with
// the partition leader.
func (kfk Kafka) InSyncReplicas(topic string, partitionID int32) ([]int32, error) {
	return kfk.client.InSyncReplicas(topic, partitionID)
}

// OfflineReplicas returns the set of all offline replica IDs for the given
// partition. Offline replicas are replicas which are offline
func (kfk Kafka) OfflineReplicas(topic string, partitionID int32) ([]int32, error) {
	return kfk.client.OfflineReplicas(topic, partitionID)
}

// RefreshBrokers takes a list of addresses to be used as seed brokers.
// Existing broker connections are closed and the updated list of seed brokers
// will be used for the next metadata fetch.
func (kfk Kafka) RefreshBrokers(addrs []string) error {
	return kfk.client.RefreshBrokers(addrs)
}

// RefreshMetadata takes a list of topics and queries the cluster to refresh the
// available metadata for those topics. If no topics are provided, it will refresh
// metadata for all topics.
func (kfk Kafka) RefreshMetadata(topics ...string) error {
	return kfk.client.RefreshMetadata(topics...)
}

// GetOffset queries the cluster to get the most recent available offset at the
// given time (in milliseconds) on the topic/partition combination.
// Time should be OffsetOldest for the earliest available offset,
// OffsetNewest for the offset of the message that will be produced next, or a time.
func (kfk Kafka) GetOffset(topic string, partitionID int32, time int64) (int64, error) {
	return kfk.client.GetOffset(topic, partitionID, time)
}

// Coordinator returns the coordinating broker for a consumer group. It will
// return a locally cached value if it's available. You can call
// RefreshCoordinator to update the cached value. This function only works on
// Kafka 0.8.2 and higher.
func (kfk Kafka) Coordinator(consumerGroup string) (*Broker, error) {
	return kfk.client.Coordinator(consumerGroup)
}

// RefreshCoordinator retrieves the coordinator for a consumer group and stores it
// in local cache. This function only works on Kafka 0.8.2 and higher.
func (kfk Kafka) RefreshCoordinator(consumerGroup string) error {
	return kfk.client.RefreshCoordinator(consumerGroup)
}

// InitProducerID retrieves information required for Idempotent Producer
func (kfk Kafka) InitProducerID() (*InitProducerIDResponse, error) {
	return kfk.client.InitProducerID()
}

// Close shuts down all broker connections managed by this client. It is required
// to call this function before a client object passes out of scope, as it will
// otherwise leak memory. You must close any Producers or Consumers using a client
// before you close the client.
func (kfk Kafka) Close() error {
	return kfk.client.Close()
}

// Closed returns true if the client has already had Close called on it
func (kfk Kafka) Closed() bool {
	return kfk.client.Closed()
}
