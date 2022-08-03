package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/boxgo/box/v2/config"
	"github.com/boxgo/box/v2/logger"
)

type (
	// Config 配置
	Config struct {
		path              string
		kfkCfg            *sarama.Config
		Addrs    []string       `config:"addrs"`
		Net      Net            `config:"net"`
		Metadata Metadata       `config:"metadata"`
		Producer ProducerConfig `config:"producer"`
		Consumer ConsumerConfig `config:"consumer"`
		ClientID string         `config:"clientId"`
		ChannelBufferSize int                 `config:"channelBufferSize"`
		Version           sarama.KafkaVersion `config:"version"`
	}

	Net struct {
		MaxOpenRequests int           `config:"maxOpenRequests"`
		DialTimeout     time.Duration `config:"dialTimeout"`
		ReadTimeout     time.Duration `config:"readTimeout"`
		WriteTimeout    time.Duration `config:"writeTimeout"`
		KeepAlive       time.Duration `config:"keepAlive"`
	}

	Metadata struct {
		RetryMax         int                                         `config:"retryMax"`
		RetryBackoff     time.Duration                               `config:"retryBackoff"`
		RetryBackoffFunc func(retries, maxRetries int) time.Duration `config:"-" json:"-"`
		RefreshFrequency time.Duration                               `config:"refreshFrequency"`
		Full             bool                                        `config:"full"`
		Timeout          time.Duration                               `config:"timeout"`
	}

	ProducerConfig struct {
		MaxMessageBytes  int                                         `config:"maxMessageBytes"`
		RequiredAcks     sarama.RequiredAcks                         `config:"requiredAcks"`
		Timeout          time.Duration                               `config:"timeout"`
		Compression      sarama.CompressionCodec                     `config:"compression"`
		CompressionLevel int                                         `config:"compressionLevel"`
		Partitioner      sarama.PartitionerConstructor               `config:"-" json:"-"`
		Idempotent       bool                                        `config:"idempotent"`
		ReturnSuccesses  bool                                        `config:"returnSuccesses"`
		ReturnErrors     bool                                        `config:"returnErrors"`
		FlushBytes       int                                         `config:"flushBytes"`
		FlushMessages    int                                         `config:"flushMessages"`
		FlushFrequency   time.Duration                               `config:"flushFrequency"`
		FlushMaxMessages int                                         `config:"FlushMaxMessages"`
		RetryMax         int                                         `config:"retryMax"`
		RetryBackoff     time.Duration                               `config:"retryBackoff"`
		RetryBackoffFunc func(retries, maxRetries int) time.Duration `config:"-" json:"-"`
		Interceptors     []sarama.ProducerInterceptor                `config:"-" json:"-"`
	}

	ConsumerConfig struct {
		GroupSessionTimeout        time.Duration                   `config:"groupSessionTimeout"`
		GroupHeartbeatInterval     time.Duration                   `config:"groupHeartbeatInterval"`
		GroupRebalanceStrategy     sarama.BalanceStrategy          `config:"groupRebalanceStrategy" json:"-"`
		GroupRebalanceTimeout      time.Duration                   `config:"groupRebalanceTimeout"`
		GroupRebalanceRetryMax     int                             `config:"groupRebalanceRetryMax"`
		GroupRebalanceRetryBackoff time.Duration                   `config:"groupRebalanceRetryBackoff"`
		GroupMemberUserData        []byte                          `config:"groupMemberUserData"`
		RetryBackoff               time.Duration                   `config:"retryBackoff"`
		RetryBackoffFunc           func(retries int) time.Duration `config:"-" json:"-"`
		FetchMin                   int32                           `config:"fetchMin"`
		FetchMax                   int32                           `config:"fetchMax"`
		FetchDefault               int32                           `config:"fetchDefault"`
		MaxWaitTime                time.Duration                   `config:"maxWaitTime"`
		MaxProcessingTime          time.Duration                   `config:"maxProcessingTime"`
		ReturnErrors               bool                            `config:"returnErrors"`
		OffsetsCommitInterval      time.Duration                   `config:"offsetsCommitInterval"`
		OffsetsInitial             int64                           `config:"offsetsInitial"`
		OffsetsRetention           time.Duration                   `config:"offsetsRetention"`
		OffsetRetryMax             int                             `config:"offsetRetryMax"`
		OffsetAutoCommitEnable     bool                            `config:"offsetAutoCommitEnable"`
		OffsetAutoCommitInterval   time.Duration                   `config:"offsetAutoCommitInterval"`
		IsolationLevel             sarama.IsolationLevel           `config:"isolationLevel"`
		Interceptors               []sarama.ConsumerInterceptor    `config:"-" json:"-"`
	}

	// OptionFunc 选项信息
	OptionFunc func(*Config)
)

// StdConfig 标准配置
func StdConfig(key string, optionFunc ...OptionFunc) *Config {
	cfg := DefaultConfig(key)
	for _, fn := range optionFunc {
		fn(cfg)
	}

	if err := config.Scan(cfg); err != nil {
		logger.Panicf("Kafka load config error: %s", err)
	}

	cfg.kfkCfg.Net.MaxOpenRequests = cfg.Net.MaxOpenRequests
	cfg.kfkCfg.Net.DialTimeout = cfg.Net.DialTimeout
	cfg.kfkCfg.Net.ReadTimeout = cfg.Net.ReadTimeout
	cfg.kfkCfg.Net.WriteTimeout = cfg.Net.WriteTimeout
	cfg.kfkCfg.Net.KeepAlive = cfg.Net.KeepAlive
	cfg.kfkCfg.Metadata.Retry.Max = cfg.Metadata.RetryMax
	cfg.kfkCfg.Metadata.Retry.Backoff = cfg.Metadata.RetryBackoff
	cfg.kfkCfg.Metadata.Retry.BackoffFunc = cfg.Metadata.RetryBackoffFunc
	cfg.kfkCfg.Metadata.RefreshFrequency = cfg.Metadata.RefreshFrequency
	cfg.kfkCfg.Metadata.Full = cfg.Metadata.Full
	cfg.kfkCfg.Metadata.Timeout = cfg.Metadata.Timeout
	cfg.kfkCfg.Producer.MaxMessageBytes = cfg.Producer.MaxMessageBytes
	cfg.kfkCfg.Producer.RequiredAcks = cfg.Producer.RequiredAcks
	cfg.kfkCfg.Producer.Timeout = cfg.Producer.Timeout
	cfg.kfkCfg.Producer.Compression = cfg.Producer.Compression
	cfg.kfkCfg.Producer.CompressionLevel = cfg.Producer.CompressionLevel
	cfg.kfkCfg.Producer.Partitioner = cfg.Producer.Partitioner
	cfg.kfkCfg.Producer.Idempotent = cfg.Producer.Idempotent
	cfg.kfkCfg.Producer.Return.Successes = cfg.Producer.ReturnSuccesses
	cfg.kfkCfg.Producer.Return.Errors = cfg.Producer.ReturnErrors
	cfg.kfkCfg.Producer.Flush.Bytes = cfg.Producer.FlushBytes
	cfg.kfkCfg.Producer.Flush.Messages = cfg.Producer.FlushMessages
	cfg.kfkCfg.Producer.Flush.Frequency = cfg.Producer.FlushFrequency
	cfg.kfkCfg.Producer.Flush.MaxMessages = cfg.Producer.FlushMaxMessages
	cfg.kfkCfg.Producer.Retry.Max = cfg.Producer.RetryMax
	cfg.kfkCfg.Producer.Retry.Backoff = cfg.Producer.RetryBackoff
	cfg.kfkCfg.Producer.Retry.BackoffFunc = cfg.Producer.RetryBackoffFunc
	cfg.kfkCfg.Producer.Interceptors = cfg.Producer.Interceptors
	cfg.kfkCfg.Consumer.Group.Session.Timeout = cfg.Consumer.GroupSessionTimeout
	cfg.kfkCfg.Consumer.Group.Heartbeat.Interval = cfg.Consumer.GroupHeartbeatInterval
	cfg.kfkCfg.Consumer.Group.Rebalance.Strategy = cfg.Consumer.GroupRebalanceStrategy
	cfg.kfkCfg.Consumer.Group.Rebalance.Timeout = cfg.Consumer.GroupRebalanceTimeout
	cfg.kfkCfg.Consumer.Group.Rebalance.Retry.Max = cfg.Consumer.GroupRebalanceRetryMax
	cfg.kfkCfg.Consumer.Group.Rebalance.Retry.Backoff = cfg.Consumer.GroupRebalanceRetryBackoff
	cfg.kfkCfg.Consumer.Retry.Backoff = cfg.Consumer.RetryBackoff
	cfg.kfkCfg.Consumer.Retry.BackoffFunc = cfg.Consumer.RetryBackoffFunc
	cfg.kfkCfg.Consumer.Fetch.Min = cfg.Consumer.FetchMin
	cfg.kfkCfg.Consumer.Fetch.Max = cfg.Consumer.FetchMax
	cfg.kfkCfg.Consumer.Fetch.Default = cfg.Consumer.FetchDefault
	cfg.kfkCfg.Consumer.MaxWaitTime = cfg.Consumer.MaxWaitTime
	cfg.kfkCfg.Consumer.MaxProcessingTime = cfg.Consumer.MaxProcessingTime
	cfg.kfkCfg.Consumer.Return.Errors = cfg.Consumer.ReturnErrors
	cfg.kfkCfg.Consumer.Offsets.Initial = cfg.Consumer.OffsetsInitial
	cfg.kfkCfg.Consumer.Offsets.Retention = cfg.Consumer.OffsetsRetention
	cfg.kfkCfg.Consumer.Offsets.Retry.Max = cfg.Consumer.OffsetRetryMax
	cfg.kfkCfg.Consumer.Offsets.AutoCommit.Enable = cfg.Consumer.OffsetAutoCommitEnable
	cfg.kfkCfg.Consumer.Offsets.AutoCommit.Interval = cfg.Consumer.OffsetAutoCommitInterval
	cfg.kfkCfg.Consumer.IsolationLevel = cfg.Consumer.IsolationLevel
	cfg.kfkCfg.Consumer.Interceptors = cfg.Consumer.Interceptors
	cfg.kfkCfg.ClientID = cfg.ClientID
	cfg.kfkCfg.ChannelBufferSize = cfg.ChannelBufferSize
	cfg.kfkCfg.Version = cfg.Version

	if err := cfg.kfkCfg.Validate(); err != nil {
		logger.Panicf("Kafka config invalid error: %s", err)
	}

	return cfg
}

// DefaultConfig 默认配置
func DefaultConfig(key string) *Config {
	kfkCfg := sarama.NewConfig()

	return &Config{
		path:  "kafka." + key,
		Addrs: []string{},
		Net: Net{
			MaxOpenRequests: kfkCfg.Net.MaxOpenRequests,
			DialTimeout:     kfkCfg.Net.DialTimeout,
			ReadTimeout:     kfkCfg.Net.ReadTimeout,
			WriteTimeout:    kfkCfg.Net.WriteTimeout,
			KeepAlive:       kfkCfg.Net.KeepAlive,
		},
		Metadata: Metadata{
			RetryMax:         kfkCfg.Metadata.Retry.Max,
			RetryBackoff:     kfkCfg.Metadata.Retry.Backoff,
			RetryBackoffFunc: kfkCfg.Metadata.Retry.BackoffFunc,
			RefreshFrequency: kfkCfg.Metadata.RefreshFrequency,
			Full:             kfkCfg.Metadata.Full,
			Timeout:          kfkCfg.Metadata.Timeout,
		},
		Producer: ProducerConfig{
			MaxMessageBytes:  kfkCfg.Producer.MaxMessageBytes,
			RequiredAcks:     kfkCfg.Producer.RequiredAcks,
			Timeout:          kfkCfg.Producer.Timeout,
			Compression:      kfkCfg.Producer.Compression,
			CompressionLevel: kfkCfg.Producer.CompressionLevel,
			Partitioner:      kfkCfg.Producer.Partitioner,
			Idempotent:       kfkCfg.Producer.Idempotent,
			ReturnSuccesses:  kfkCfg.Producer.Return.Successes,
			ReturnErrors:     kfkCfg.Producer.Return.Errors,
			FlushBytes:       kfkCfg.Producer.Flush.Bytes,
			FlushMessages:    kfkCfg.Producer.Flush.Messages,
			FlushFrequency:   kfkCfg.Producer.Flush.Frequency,
			FlushMaxMessages: kfkCfg.Producer.Flush.MaxMessages,
			RetryMax:         kfkCfg.Producer.Retry.Max,
			RetryBackoff:     kfkCfg.Producer.Retry.Backoff,
			RetryBackoffFunc: kfkCfg.Producer.Retry.BackoffFunc,
			Interceptors:     kfkCfg.Producer.Interceptors,
		},
		Consumer: ConsumerConfig{
			GroupSessionTimeout:        kfkCfg.Consumer.Group.Session.Timeout,
			GroupHeartbeatInterval:     kfkCfg.Consumer.Group.Heartbeat.Interval,
			GroupRebalanceStrategy:     kfkCfg.Consumer.Group.Rebalance.Strategy,
			GroupRebalanceTimeout:      kfkCfg.Consumer.Group.Rebalance.Timeout,
			GroupRebalanceRetryMax:     kfkCfg.Consumer.Group.Rebalance.Retry.Max,
			GroupRebalanceRetryBackoff: kfkCfg.Consumer.Group.Rebalance.Retry.Backoff,
			GroupMemberUserData:        kfkCfg.Consumer.Group.Member.UserData,
			RetryBackoff:               kfkCfg.Consumer.Retry.Backoff,
			RetryBackoffFunc:           kfkCfg.Consumer.Retry.BackoffFunc,
			FetchMin:                   kfkCfg.Consumer.Fetch.Min,
			FetchMax:                   kfkCfg.Consumer.Fetch.Max,
			FetchDefault:               kfkCfg.Consumer.Fetch.Default,
			MaxWaitTime:                kfkCfg.Consumer.MaxWaitTime,
			MaxProcessingTime:          kfkCfg.Consumer.MaxProcessingTime,
			ReturnErrors:               kfkCfg.Consumer.Return.Errors,
			OffsetsInitial:             kfkCfg.Consumer.Offsets.Initial,
			OffsetsRetention:           kfkCfg.Consumer.Offsets.Retention,
			OffsetRetryMax:             kfkCfg.Consumer.Offsets.Retry.Max,
			OffsetAutoCommitEnable:     kfkCfg.Consumer.Offsets.AutoCommit.Enable,
			OffsetAutoCommitInterval:   kfkCfg.Consumer.Offsets.AutoCommit.Interval,
			IsolationLevel:             kfkCfg.Consumer.IsolationLevel,
			Interceptors:               kfkCfg.Consumer.Interceptors,
		},
		ClientID:          kfkCfg.ClientID,
		ChannelBufferSize: kfkCfg.ChannelBufferSize,
		Version:           kfkCfg.Version,
		kfkCfg:            kfkCfg,
	}
}

// Build 构建实例
func (c *Config) Build() *Kafka {
	return newKafka(c)
}

// Path 实例配置目录
func (c *Config) Path() string {
	return c.path
}
