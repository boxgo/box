package mongodb

import (
	"context"

	"github.com/boxgo/box/v2/config"
	"github.com/boxgo/box/v2/logger"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	Mongo struct {
		client  *mongo.Client
		monitor Monitor
		cfg     *Config
	}

	Monitor interface {
		Setup(*mongo.Client)
		Serve()
		Shutdown()
		Started(context.Context, *event.CommandStartedEvent)
		Succeeded(context.Context, *event.CommandSucceededEvent)
		Failed(context.Context, *event.CommandFailedEvent)
		Event(*event.PoolEvent)
	}
)

func newMongo(cfg *Config) *Mongo {
	clientOptions := options.Client()
	clientOptions.ApplyURI(cfg.URI)
	clientOptions.SetAppName(config.ServiceName())

	var mo Monitor
	if cfg.monitor == nil && (cfg.EnableCommandMonitor || cfg.EnablePoolMonitor) {
		mo = newMonitor()
	}

	if cfg.EnableCommandMonitor {
		clientOptions.Monitor = &event.CommandMonitor{
			Started:   mo.Started,
			Succeeded: mo.Succeeded,
			Failed:    mo.Failed,
		}
	}
	if cfg.EnablePoolMonitor {
		clientOptions.PoolMonitor = &event.PoolMonitor{Event: mo.Event}
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		logger.Panicf("new mongodb client error: %s", err)
	}

	if mo != nil {
		mo.Setup(client)
	}

	mgo := &Mongo{
		client:  client,
		cfg:     cfg,
		monitor: mo,
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
		mgo.monitor.Serve()
	}

	return mgo.client.Ping(ctx, readpref.Primary())
}

func (mgo *Mongo) Shutdown(ctx context.Context) error {
	if mgo.monitor != nil {
		mgo.monitor.Shutdown()
	}

	return mgo.client.Disconnect(ctx)
}

func (mgo *Mongo) Client() *mongo.Client {
	return mgo.client
}

func (mgo *Mongo) DB(db string, opts ...*options.DatabaseOptions) *mongo.Database {
	return mgo.client.Database(db, opts...)
}
