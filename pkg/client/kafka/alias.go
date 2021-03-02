package kafka

import (
	"github.com/Shopify/sarama"
)

type (
	SyncProducer                  = sarama.SyncProducer
	AsyncSyncProducer             = sarama.AsyncProducer
	Consumer                      = sarama.Consumer
	ConsumerGroup                 = sarama.ConsumerGroup
	ProduceRequest                = sarama.ProduceRequest
	ProduceResponse               = sarama.ProduceResponse
	ProduceResponseBlock          = sarama.ProduceResponseBlock
	ProducerError                 = sarama.ProducerError
	ProducerErrors                = sarama.ProducerErrors
	ProducerInterceptor           = sarama.ProducerInterceptor
	ProducerMessage               = sarama.ProducerMessage
	InitProducerIDRequest         = sarama.InitProducerIDRequest
	InitProducerIDResponse        = sarama.InitProducerIDResponse
	ConsumerError                 = sarama.ConsumerError
	ConsumerErrors                = sarama.ConsumerErrors
	ConsumerGroupClaim            = sarama.ConsumerGroupClaim
	ConsumerGroupHandler          = sarama.ConsumerGroupHandler
	ConsumerGroupMemberAssignment = sarama.ConsumerGroupMemberAssignment
	ConsumerGroupMemberMetadata   = sarama.ConsumerGroupMemberMetadata
	ConsumerGroupSession          = sarama.ConsumerGroupSession
	ConsumerInterceptor           = sarama.ConsumerInterceptor
	ConsumerMessage               = sarama.ConsumerMessage
	ConsumerMetadataRequest       = sarama.ConsumerMetadataRequest
	ConsumerMetadataResponse      = sarama.ConsumerMetadataResponse
)

const (
	// RangeBalanceStrategyName identifies strategies that use the range partition assignment strategy
	RangeBalanceStrategyName = sarama.RangeBalanceStrategyName

	// RoundRobinBalanceStrategyName identifies strategies that use the round-robin partition assignment strategy
	RoundRobinBalanceStrategyName = sarama.RoundRobinBalanceStrategyName

	// StickyBalanceStrategyName identifies strategies that use the sticky-partition assignment strategy
	StickyBalanceStrategyName = sarama.StickyBalanceStrategyName
)

const (
	// SASLTypeOAuth represents the SASL/OAUTHBEARER mechanism (Kafka 2.0.0+)
	SASLTypeOAuth = sarama.SASLTypeOAuth
	// SASLTypePlaintext represents the SASL/PLAIN mechanism
	SASLTypePlaintext = sarama.SASLTypePlaintext
	// SASLTypeSCRAMSHA256 represents the SCRAM-SHA-256 mechanism.
	SASLTypeSCRAMSHA256 = sarama.SASLTypeSCRAMSHA256
	// SASLTypeSCRAMSHA512 represents the SCRAM-SHA-512 mechanism.
	SASLTypeSCRAMSHA512 = sarama.SASLTypeSCRAMSHA512
	SASLTypeGSSAPI      = sarama.SASLTypeGSSAPI
	// SASLHandshakeV0 is v0 of the Kafka SASL handshake protocol. Client and
	// server negotiate SASL auth using opaque packets.
	SASLHandshakeV0 = sarama.SASLHandshakeV0
	// SASLHandshakeV1 is v1 of the Kafka SASL handshake protocol. Client and
	// server negotiate SASL by wrapping tokens with Kafka protocol headers.
	SASLHandshakeV1 = sarama.SASLHandshakeV1
	// SASLExtKeyAuth is the reserved extension key name sent as part of the
	// SASL/OAUTHBEARER initial client response
	SASLExtKeyAuth = sarama.SASLExtKeyAuth
)

const (
	// OffsetNewest stands for the log head offset, i.e. the offset that will be
	// assigned to the next message that will be produced to the partition. You
	// can send this to a client's GetOffset method to get this offset, or when
	// calling ConsumePartition to start consuming new messages.
	OffsetNewest = sarama.OffsetNewest
	// OffsetOldest stands for the oldest offset available on the broker for a
	// partition. You can send this to a client's GetOffset method to get this
	// offset, or when calling ConsumePartition to start consuming from the
	// oldest offset that is still available on the broker.
	OffsetOldest = sarama.OffsetOldest
)

const (
	TOK_ID_KRB_AP_REQ   = sarama.TOK_ID_KRB_AP_REQ
	GSS_API_GENERIC_TAG = sarama.GSS_API_GENERIC_TAG
	KRB5_USER_AUTH      = sarama.KRB5_USER_AUTH
	KRB5_KEYTAB_AUTH    = sarama.KRB5_KEYTAB_AUTH
	GSS_API_INITIAL     = sarama.GSS_API_INITIAL
	GSS_API_VERIFY      = sarama.GSS_API_VERIFY
	GSS_API_FINISH      = sarama.GSS_API_FINISH
)

const (
	// APIKeySASLAuth is the API key for the SaslAuthenticate Kafka API
	APIKeySASLAuth = sarama.APIKeySASLAuth

	// GroupGenerationUndefined is a special value for the group generation field of Offset Commit Requests that should be used when a consumer group does not rely on Kafka for partition management.
	GroupGenerationUndefined = sarama.GroupGenerationUndefined

	// ReceiveTime is a special value for the timestamp field of Offset Commit Requests which tells the broker to set the timestamp to the time at which the request was received. The timestamp is only used if message version 1 is used, which requires kafka 0.8.2.
	ReceiveTime int64 = sarama.ReceiveTime
)

// Numeric error codes returned by the Kafka server.
const (
	ErrNoError                            = sarama.ErrNoError
	ErrUnknown                            = sarama.ErrUnknown
	ErrOffsetOutOfRange                   = sarama.ErrOffsetOutOfRange
	ErrInvalidMessage                     = sarama.ErrInvalidMessage
	ErrUnknownTopicOrPartition            = sarama.ErrUnknownTopicOrPartition
	ErrInvalidMessageSize                 = sarama.ErrInvalidMessageSize
	ErrLeaderNotAvailable                 = sarama.ErrLeaderNotAvailable
	ErrNotLeaderForPartition              = sarama.ErrNotLeaderForPartition
	ErrRequestTimedOut                    = sarama.ErrRequestTimedOut
	ErrBrokerNotAvailable                 = sarama.ErrBrokerNotAvailable
	ErrReplicaNotAvailable                = sarama.ErrReplicaNotAvailable
	ErrMessageSizeTooLarge                = sarama.ErrMessageSizeTooLarge
	ErrStaleControllerEpochCode           = sarama.ErrStaleControllerEpochCode
	ErrOffsetMetadataTooLarge             = sarama.ErrOffsetMetadataTooLarge
	ErrNetworkException                   = sarama.ErrNetworkException
	ErrOffsetsLoadInProgress              = sarama.ErrOffsetsLoadInProgress
	ErrConsumerCoordinatorNotAvailable    = sarama.ErrConsumerCoordinatorNotAvailable
	ErrNotCoordinatorForConsumer          = sarama.ErrNotCoordinatorForConsumer
	ErrInvalidTopic                       = sarama.ErrInvalidTopic
	ErrMessageSetSizeTooLarge             = sarama.ErrMessageSetSizeTooLarge
	ErrNotEnoughReplicas                  = sarama.ErrNotEnoughReplicas
	ErrNotEnoughReplicasAfterAppend       = sarama.ErrNotEnoughReplicasAfterAppend
	ErrInvalidRequiredAcks                = sarama.ErrInvalidRequiredAcks
	ErrIllegalGeneration                  = sarama.ErrIllegalGeneration
	ErrInconsistentGroupProtocol          = sarama.ErrInconsistentGroupProtocol
	ErrInvalidGroupId                     = sarama.ErrInvalidGroupId
	ErrUnknownMemberId                    = sarama.ErrUnknownMemberId
	ErrInvalidSessionTimeout              = sarama.ErrInvalidSessionTimeout
	ErrRebalanceInProgress                = sarama.ErrRebalanceInProgress
	ErrInvalidCommitOffsetSize            = sarama.ErrInvalidCommitOffsetSize
	ErrTopicAuthorizationFailed           = sarama.ErrTopicAuthorizationFailed
	ErrGroupAuthorizationFailed           = sarama.ErrGroupAuthorizationFailed
	ErrClusterAuthorizationFailed         = sarama.ErrClusterAuthorizationFailed
	ErrInvalidTimestamp                   = sarama.ErrInvalidTimestamp
	ErrUnsupportedSASLMechanism           = sarama.ErrUnsupportedSASLMechanism
	ErrIllegalSASLState                   = sarama.ErrIllegalSASLState
	ErrUnsupportedVersion                 = sarama.ErrUnsupportedVersion
	ErrTopicAlreadyExists                 = sarama.ErrTopicAlreadyExists
	ErrInvalidPartitions                  = sarama.ErrInvalidPartitions
	ErrInvalidReplicationFactor           = sarama.ErrInvalidReplicationFactor
	ErrInvalidReplicaAssignment           = sarama.ErrInvalidReplicaAssignment
	ErrInvalidConfig                      = sarama.ErrInvalidConfig
	ErrNotController                      = sarama.ErrNotController
	ErrInvalidRequest                     = sarama.ErrInvalidRequest
	ErrUnsupportedForMessageFormat        = sarama.ErrUnsupportedForMessageFormat
	ErrPolicyViolation                    = sarama.ErrPolicyViolation
	ErrOutOfOrderSequenceNumber           = sarama.ErrOutOfOrderSequenceNumber
	ErrDuplicateSequenceNumber            = sarama.ErrDuplicateSequenceNumber
	ErrInvalidProducerEpoch               = sarama.ErrInvalidProducerEpoch
	ErrInvalidTxnState                    = sarama.ErrInvalidTxnState
	ErrInvalidProducerIDMapping           = sarama.ErrInvalidProducerIDMapping
	ErrInvalidTransactionTimeout          = sarama.ErrInvalidTransactionTimeout
	ErrConcurrentTransactions             = sarama.ErrConcurrentTransactions
	ErrTransactionCoordinatorFenced       = sarama.ErrTransactionCoordinatorFenced
	ErrTransactionalIDAuthorizationFailed = sarama.ErrTransactionalIDAuthorizationFailed
	ErrSecurityDisabled                   = sarama.ErrSecurityDisabled
	ErrOperationNotAttempted              = sarama.ErrOperationNotAttempted
	ErrKafkaStorageError                  = sarama.ErrKafkaStorageError
	ErrLogDirNotFound                     = sarama.ErrLogDirNotFound
	ErrSASLAuthenticationFailed           = sarama.ErrSASLAuthenticationFailed
	ErrUnknownProducerID                  = sarama.ErrUnknownProducerID
	ErrReassignmentInProgress             = sarama.ErrReassignmentInProgress
	ErrDelegationTokenAuthDisabled        = sarama.ErrDelegationTokenAuthDisabled
	ErrDelegationTokenNotFound            = sarama.ErrDelegationTokenNotFound
	ErrDelegationTokenOwnerMismatch       = sarama.ErrDelegationTokenOwnerMismatch
	ErrDelegationTokenRequestNotAllowed   = sarama.ErrDelegationTokenRequestNotAllowed
	ErrDelegationTokenAuthorizationFailed = sarama.ErrDelegationTokenAuthorizationFailed
	ErrDelegationTokenExpired             = sarama.ErrDelegationTokenExpired
	ErrInvalidPrincipalType               = sarama.ErrInvalidPrincipalType
	ErrNonEmptyGroup                      = sarama.ErrNonEmptyGroup
	ErrGroupIDNotFound                    = sarama.ErrGroupIDNotFound
	ErrFetchSessionIDNotFound             = sarama.ErrFetchSessionIDNotFound
	ErrInvalidFetchSessionEpoch           = sarama.ErrInvalidFetchSessionEpoch
	ErrListenerNotFound                   = sarama.ErrListenerNotFound
	ErrTopicDeletionDisabled              = sarama.ErrTopicDeletionDisabled
	ErrFencedLeaderEpoch                  = sarama.ErrFencedLeaderEpoch
	ErrUnknownLeaderEpoch                 = sarama.ErrUnknownLeaderEpoch
	ErrUnsupportedCompressionType         = sarama.ErrUnsupportedCompressionType
	ErrStaleBrokerEpoch                   = sarama.ErrStaleBrokerEpoch
	ErrOffsetNotAvailable                 = sarama.ErrOffsetNotAvailable
	ErrMemberIdRequired                   = sarama.ErrMemberIdRequired
	ErrPreferredLeaderNotAvailable        = sarama.ErrPreferredLeaderNotAvailable
	ErrGroupMaxSizeReached                = sarama.ErrGroupMaxSizeReached
	ErrFencedInstancedId                  = sarama.ErrFencedInstancedId
)

var (
	// MaxRequestSize is the maximum size (in bytes) of any request that Sarama will attempt to send. Trying
	// to send a request larger than this will result in an PacketEncodingError. The default of 100 MiB is aligned
	// with Kafka's default `socket.request.max.bytes`, which is the largest request the broker will attempt
	// to process.
	MaxRequestSize = sarama.MaxRequestSize

	// MaxResponseSize is the maximum size (in bytes) of any response that Sarama will attempt to parse. If
	// a broker returns a response message larger than this value, Sarama will return a PacketDecodingError to
	// protect the client from running out of memory. Please note that brokers do not have any natural limit on
	// the size of responses they send. In particular, they can send arbitrarily large fetch responses to consumers
	// (see https://issues.apache.org/jira/browse/KAFKA-2063).
	MaxResponseSize = sarama.MaxResponseSize
)

var (
	// ErrOutOfBrokers is the error returned when the client has run out of brokers to talk to because all of them errored
	// or otherwise failed to respond.
	ErrOutOfBrokers = sarama.ErrOutOfBrokers

	// ErrClosedClient is the error returned when a method is called on a client that has been closed.
	ErrClosedClient = sarama.ErrClosedClient

	// ErrIncompleteResponse is the error returned when the server returns a syntactically valid response, but it does
	// not contain the expected information.
	ErrIncompleteResponse = sarama.ErrIncompleteResponse

	// ErrInvalidPartition is the error returned when a partitioner returns an invalid partition index
	// (meaning one outside of the range [0...numPartitions-1]).
	ErrInvalidPartition = sarama.ErrInvalidPartition

	// ErrAlreadyConnected is the error returned when calling Open() on a Broker that is already connected or connecting.
	ErrAlreadyConnected = sarama.ErrAlreadyConnected

	// ErrNotConnected is the error returned when trying to send or call Close() on a Broker that is not connected.
	ErrNotConnected = sarama.ErrNotConnected

	// ErrInsufficientData is returned when decoding and the packet is truncated. This can be expected
	// when requesting messages, since as an optimization the server is allowed to return a partial message at the end
	// of the message set.
	ErrInsufficientData = sarama.ErrInsufficientData

	// ErrShuttingDown is returned when a producer receives a message during shutdown.
	ErrShuttingDown = sarama.ErrShuttingDown

	// ErrMessageTooLarge is returned when the next message to consume is larger than the configured Consumer.Fetch.Max
	ErrMessageTooLarge = sarama.ErrMessageTooLarge

	// ErrConsumerOffsetNotAdvanced is returned when a partition consumer didn't advance its offset after parsing
	// a RecordBatch.
	ErrConsumerOffsetNotAdvanced = sarama.ErrConsumerOffsetNotAdvanced

	// ErrControllerNotAvailable is returned when server didn't give correct controller id. May be kafka server's version
	// is lower than 0.10.0.0.
	ErrControllerNotAvailable = sarama.ErrControllerNotAvailable

	// ErrNoTopicsToUpdateMetadata is returned when Meta.Full is set to false but no specific topics were found to update
	// the metadata.
	ErrNoTopicsToUpdateMetadata = sarama.ErrNoTopicsToUpdateMetadata
)
