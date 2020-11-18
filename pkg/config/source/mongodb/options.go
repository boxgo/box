package mongodb

import (
	"context"
	"encoding/json"
	"log"

	"github.com/boxgo/box/pkg/config/source"
)

type (
	mongoURIKey        struct{}
	mongoDBKey         struct{}
	mongoCollectionKey struct{}
	mongoServiceKey    struct{}

	mongoConfig struct {
		URI        string `json:"uri" desc:"mongodb connection uri."`
		DB         string `json:"db"`
		Collection string `json:"collection"`
		Service    string `json:"service"`
	}
)

func WithConfig(data []byte) []source.Option {
	v := &mongoConfig{}

	if err := json.Unmarshal(data, v); err != nil {
		log.Panicf("service mongo build error: %#v", err)
	}

	return []source.Option{
		WithURI(v.URI),
		WithDB(v.DB),
		WithCollection(v.Collection),
		WithService(v.Service),
	}
}

func WithURI(uri string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, mongoURIKey{}, uri)
	}
}

func WithDB(db string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, mongoDBKey{}, db)
	}
}

func WithCollection(collection string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, mongoCollectionKey{}, collection)
	}
}

func WithService(service string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, mongoServiceKey{}, service)
	}
}
