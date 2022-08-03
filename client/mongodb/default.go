package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Default = StdConfig("default").Build()
)

func Client() *mongo.Client {
	return Default.Client()
}

func DB(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	return Default.DB(name, opts...)
}
