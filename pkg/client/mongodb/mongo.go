package mongodb

import (
	"context"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	Mongo struct {
		client  *mongo.Client
		monitor *Monitor
		cfg     *Config
	}
)

func newMongo(cfg *Config) *Mongo {
	clientOptions := options.Client()
	clientOptions.ApplyURI(cfg.URI)
	clientOptions.SetAppName(config.ServiceName())

	if cfg.commandMonitor == nil {
		cfg.commandMonitor = defaultMonitor.CommandMonitor()
	}
	if cfg.poolMonitor == nil {
		cfg.poolMonitor = defaultMonitor.PoolEventMonitor()
	}

	if cfg.EnableCommandMonitor {
		clientOptions.Monitor = cfg.commandMonitor
	}
	if cfg.EnablePoolMonitor {
		clientOptions.PoolMonitor = cfg.poolMonitor
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		logger.Panicf("new mongodb client error: %s", err)
	}

	mgo := &Mongo{
		client:  client,
		cfg:     cfg,
		monitor: defaultMonitor,
	}

	return mgo
}

func (mgo *Mongo) Name() string {
	return "mongo"
}

func (mgo *Mongo) Serve(ctx context.Context) error {
	if err := mgo.client.Connect(ctx); err != nil {
		return err
	}

	if mgo.monitor != nil {
		mgo.monitor.watch(mgo.client)
	}

	return mgo.client.Ping(ctx, readpref.Primary())
}

func (mgo *Mongo) Shutdown(ctx context.Context) error {
	if mgo.monitor != nil {
		mgo.monitor.shutdown()
	}

	return mgo.client.Disconnect(ctx)
}

func (mgo *Mongo) Client() *mongo.Client {
	return mgo.client
}

func (mgo *Mongo) DB(db string, opts ...*options.DatabaseOptions) *mongo.Database {
	return mgo.client.Database(db, opts...)
}
