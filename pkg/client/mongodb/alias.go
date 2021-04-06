package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
