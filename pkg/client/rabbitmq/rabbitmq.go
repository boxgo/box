// Package rabbitmq is an AMQP 0.9.1 client with RabbitMQ extensions in Go.
package rabbitmq

import (
	"context"

	"github.com/boxgo/box/pkg/logger"
	"github.com/streadway/amqp"
)

type (
	RabbitMQ struct {
		cfg  *Config
		conn *amqp.Connection
	}
)

func newRabbitMQ(c *Config) *RabbitMQ {
	conn, err := amqp.DialConfig(c.URI, amqp.Config{
		Vhost:      c.Vhost,
		ChannelMax: c.ChannelMax,
		FrameSize:  c.FrameSize,
		Heartbeat:  c.Heartbeat,
	})
	if err != nil {
		logger.Panicw("new rabbitmq error", "config", c, "err", err)
	}

	return &RabbitMQ{
		cfg:  c,
		conn: conn,
	}
}

func (mq *RabbitMQ) Serve(ctx context.Context) error {
	return nil
}

func (mq *RabbitMQ) Shutdown(ctx context.Context) error {
	if mq.conn != nil {
		return mq.conn.Close()
	}

	return nil
}

func (mq RabbitMQ) Channel() (*Channel, error) {
	return mq.conn.Channel()
}
