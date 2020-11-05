package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	A bson.A
	D bson.D
	E bson.E
	M bson.M
)

var (
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
)
