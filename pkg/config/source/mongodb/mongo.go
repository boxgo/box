package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/boxgo/box/pkg/config/source"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	mongoSource struct {
		err        error
		opts       source.Options
		client     *mongo.Client
		db         string
		collection string
	}

	Config struct {
		Format string
		Config string
	}
)

func NewSource(opts ...source.Option) source.Source {
	var (
		sOpts      = source.NewOptions(opts...)
		client     *mongo.Client
		db         string
		collection string
	)

	if val, ok := sOpts.Context.Value(mongoURIKey{}).(string); !ok {
		log.Panic("config source mongo is not set.")
	} else {
		clientOptions := options.Client()
		clientOptions.ApplyURI(val)
		if cli, err := mongo.NewClient(clientOptions); err != nil {
			log.Panicf("build mongo error: %s", err)
		} else {
			client = cli
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.Panicf("config connnect to mongodb error %s", err)
	}

	if val, ok := sOpts.Context.Value(mongoDBKey{}).(string); ok && val != "" {
		db = val
	}

	if val, ok := sOpts.Context.Value(mongoCollectionKey{}).(string); ok && val != "" {
		collection = val
	}

	return &mongoSource{
		err:        nil,
		opts:       sOpts,
		client:     client,
		db:         db,
		collection: collection,
	}
}

func (rs *mongoSource) Read() (*source.ChangeSet, error) {
	if rs.err != nil {
		return nil, rs.err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	cfg := &Config{}
	if err := rs.client.
		Database(rs.db).
		Collection(rs.collection).
		FindOne(ctx, bson.D{}).
		Decode(cfg); err != nil {
		return nil, err
	}

	format := cfg.Format
	if format == "" {
		format = rs.opts.Encoder.String()
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Source:    rs.String(),
		Data:      []byte(cfg.Config),
		Format:    format,
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}
func (rs *mongoSource) Watch() (source.Watcher, error) {
	if rs.err != nil {
		return nil, rs.err
	}

	return newWatcher(rs.db, rs.collection, rs.client, rs.opts)
}

func (rs *mongoSource) String() string {
	return "mongodb"
}
