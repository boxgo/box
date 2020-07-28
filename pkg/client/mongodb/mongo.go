package mongodb

import (
	"context"

	"github.com/boxgo/box/pkg/config"
	"github.com/boxgo/box/pkg/dummybox"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	Mongo struct {
		dummybox.DummyBox
		name   string
		cfg    config.SubConfigurator
		client *mongo.Client
		uri    *config.Field
	}

	Options struct {
		name string
		cfg  config.SubConfigurator
	}

	OptionFunc func(*Options)
)

func WithName(name string) OptionFunc {
	return func(options *Options) {
		options.name = name
	}
}

func WithConfigurator(cfg config.SubConfigurator) OptionFunc {
	return func(options *Options) {
		options.cfg = cfg
	}
}

func New(optionFunc ...OptionFunc) (*Mongo, error) {
	opts := &Options{}
	for _, fn := range optionFunc {
		fn(opts)
	}

	if opts.name == "" {
		opts.name = "mongo.default"
	} else {
		opts.name = "mongo." + opts.name
	}
	if opts.cfg == nil {
		opts.cfg = config.Default
	}

	uriField := config.NewField(opts.name, "uri", "mongodb connection uri. https://docs.mongodb.com/manual/reference/connection-string", "mongodb://127.0.0.1:27017")

	// set override options
	clientOptions := options.Client()
	clientOptions.Monitor = newCommandMonitor()
	clientOptions.PoolMonitor = newPoolMonitor()
	clientOptions.ApplyURI(opts.cfg.GetString(uriField))
	clientOptions.SetAppName(opts.cfg.GetBoxName())

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	mgo := &Mongo{
		name:   opts.name,
		client: client,
		cfg:    opts.cfg,
		uri:    uriField,
	}

	opts.cfg.Mount(mgo.uri)

	return mgo, nil
}

func (mgo *Mongo) Serve(ctx context.Context) error {
	if err := mgo.client.Connect(ctx); err != nil {
		return err
	}

	return mgo.client.Ping(ctx, readpref.Primary())
}

func (mgo *Mongo) Shutdown(ctx context.Context) error {
	return mgo.client.Disconnect(ctx)
}

func (mgo *Mongo) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	return mgo.client.Database(name, opts...)
}

func (mgo *Mongo) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	return mgo.client.StartSession(opts...)
}

func (mgo *Mongo) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	return mgo.client.UseSession(ctx, fn)
}

func (mgo *Mongo) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error {
	return mgo.client.UseSessionWithOptions(ctx, opts, fn)
}
