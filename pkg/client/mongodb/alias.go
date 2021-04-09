package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	A          = primitive.A
	D          = primitive.D
	E          = primitive.E
	M          = primitive.M
	ObjectID   = primitive.ObjectID
	Binary     = primitive.Binary
	DateTime   = primitive.DateTime
	Decimal128 = primitive.Decimal128
	JavaScript = primitive.JavaScript
	MaxKey     = primitive.MaxKey
	MinKey     = primitive.MinKey
	Null       = primitive.Null
	Regex      = primitive.Regex
	Symbol     = primitive.Symbol
	Timestamp  = primitive.Timestamp
	Undefined  = primitive.Undefined
)

type (
	Database                      = mongo.Database
	Collection                    = mongo.Collection
	Cursor                        = mongo.Cursor
	AggregateOptions              = options.AggregateOptions
	ArrayFilters                  = options.ArrayFilters
	AutoEncryptionOptions         = options.AutoEncryptionOptions
	BucketOptions                 = options.BucketOptions
	BulkWriteOptions              = options.BulkWriteOptions
	ChangeStreamOptions           = options.ChangeStreamOptions
	ClientEncryptionOptions       = options.ClientEncryptionOptions
	ClientOptions                 = options.ClientOptions
	Collation                     = options.Collation
	CollectionOptions             = options.CollectionOptions
	ContextDialer                 = options.ContextDialer
	CountOptions                  = options.CountOptions
	CreateCollectionOptions       = options.CreateCollectionOptions
	CreateIndexesOptions          = options.CreateIndexesOptions
	CreateViewOptions             = options.CreateViewOptions
	Credential                    = options.Credential
	CursorType                    = options.CursorType
	DataKeyOptions                = options.DataKeyOptions
	DatabaseOptions               = options.DatabaseOptions
	DefaultIndexOptions           = options.DefaultIndexOptions
	DeleteOptions                 = options.DeleteOptions
	DistinctOptions               = options.DistinctOptions
	DropIndexesOptions            = options.DropIndexesOptions
	EncryptOptions                = options.EncryptOptions
	EstimatedDocumentCountOptions = options.EstimatedDocumentCountOptions
	FindOneAndDeleteOptions       = options.FindOneAndDeleteOptions
	FindOneAndReplaceOptions      = options.FindOneAndReplaceOptions
	FindOneAndUpdateOptions       = options.FindOneAndUpdateOptions
	FindOneOptions                = options.FindOneOptions
	FindOptions                   = options.FindOptions
	FullDocument                  = options.FullDocument
	GridFSFindOptions             = options.GridFSFindOptions
	IndexOptions                  = options.IndexOptions
	InsertManyOptions             = options.InsertManyOptions
	InsertOneOptions              = options.InsertOneOptions
	ListCollectionsOptions        = options.ListCollectionsOptions
	ListDatabasesOptions          = options.ListDatabasesOptions
	ListIndexesOptions            = options.ListIndexesOptions
	MarshalError                  = options.MarshalError
	NameOptions                   = options.NameOptions
	ReplaceOptions                = options.ReplaceOptions
	ReturnDocument                = options.ReturnDocument
	RunCmdOptions                 = options.RunCmdOptions
	SessionOptions                = options.SessionOptions
	TransactionOptions            = options.TransactionOptions
	UpdateOptions                 = options.UpdateOptions
	UploadOptions                 = options.UploadOptions
)

var (
	ErrDecodeToNil         = bson.ErrDecodeToNil
	ErrNilContext          = bson.ErrNilContext
	ErrNilReader           = bson.ErrNilReader
	ErrNilRegistry         = bson.ErrNilRegistry
	ErrClientDisconnected  = mongo.ErrClientDisconnected
	ErrEmptySlice          = mongo.ErrEmptySlice
	ErrInvalidIndexValue   = mongo.ErrInvalidIndexValue
	ErrMissingResumeToken  = mongo.ErrMissingResumeToken
	ErrMultipleIndexDrop   = mongo.ErrMultipleIndexDrop
	ErrNilCursor           = mongo.ErrNilCursor
	ErrNilDocument         = mongo.ErrNilDocument
	ErrNoDocuments         = mongo.ErrNoDocuments
	ErrNonStringIndexName  = mongo.ErrNonStringIndexName
	ErrUnacknowledgedWrite = mongo.ErrUnacknowledgedWrite
	ErrWrongClient         = mongo.ErrWrongClient
	ErrParseNaN            = primitive.ErrParseNaN
	ErrParseInf            = primitive.ErrParseInf
	ErrParseNegInf         = primitive.ErrParseNegInf
)

var (
	CompareTimestamp          = primitive.CompareTimestamp
	NewDateTimeFromTime       = primitive.NewDateTimeFromTime
	NewDecimal128             = primitive.NewDecimal128
	ParseDecimal128           = primitive.ParseDecimal128
	ParseDecimal128FromBigInt = primitive.ParseDecimal128FromBigInt
	NewObjectID               = primitive.NewObjectID
	NewObjectIDFromTimestamp  = primitive.NewObjectIDFromTimestamp
	ObjectIDFromHex           = primitive.ObjectIDFromHex
)
