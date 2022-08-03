package mongodb

import (
	"context"
	"log"
	"time"

	source2 "github.com/boxgo/box/v2/config/source"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	mongoSource struct {
		err    error
		opts   source2.Options
		client *mongo.Client
		db         string
		collection string
		service    string
	}

	Config struct {
		Service string
		Format  string
		Config  string
	}
)

func NewSource(opts ...source2.Option) source2.Source {
	var (
		sOpts      = source2.NewOptions(opts...)
		client     *mongo.Client
		db         string
		collection string
		service    string
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
		log.Panicf("service connnect to mongodb error %s", err)
	}

	if val, ok := sOpts.Context.Value(mongoDBKey{}).(string); ok && val != "" {
		db = val
	}
	if val, ok := sOpts.Context.Value(mongoCollectionKey{}).(string); ok && val != "" {
		collection = val
	}
	if val, ok := sOpts.Context.Value(mongoServiceKey{}).(string); ok && val != "" {
		service = val
	}

	return &mongoSource{
		err:        nil,
		opts:       sOpts,
		client:     client,
		db:         db,
		collection: collection,
		service:    service,
	}
}

func (rs *mongoSource) Read() (*source2.ChangeSet, error) {
	if rs.err != nil {
		return nil, rs.err
	}

	cfg, err := loadConfig(rs.client, rs.db, rs.collection, rs.service)
	if err != nil {
		return nil, err
	}

	format := cfg.Format
	if format == "" {
		format = rs.opts.Encoder.String()
	}

	cs := &source2.ChangeSet{
		Timestamp: time.Now(),
		Source:    rs.String(),
		Data:      []byte(cfg.Config),
		Format:    format,
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}
func (rs *mongoSource) Watch() (source2.Watcher, error) {
	if rs.err != nil {
		return nil, rs.err
	}

	return newWatcher(rs)
}

func (rs *mongoSource) String() string {
	return "mongodb"
}

func loadConfig(client *mongo.Client, db, col, svc string) (*Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	cfg := &Config{}
	err := client.
		Database(db).
		Collection(col).
		FindOne(ctx, bson.M{"service": svc}).
		Decode(cfg)

	return cfg, err
}
