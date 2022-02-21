package postman

import (
	"context"
	"crypto/tls"
	"net/smtp"
	"time"

	"github.com/boxgo/box/pkg/logger"
	"github.com/jordan-wright/email"
)

type (
	PostMan struct {
		cfg    *Config
		client *email.Pool
	}

	Email      = email.Email
	Attachment = email.Attachment
)

var (
	NewEmail = email.NewEmail
)

var (
	Default = StdConfig("default").Build()
)

func newPostMan(cfg *Config) *PostMan {
	client, err := email.NewPool(
		cfg.Address,
		4,
		smtp.PlainAuth(cfg.Identity, cfg.Username, cfg.Password, cfg.Host),
	)

	if err != nil {
		logger.Panicw("Email.NewPool", "config", cfg, "err", err)
	}

	return &PostMan{
		cfg:    cfg,
		client: client,
	}
}

func (pm PostMan) Serve(ctx context.Context) error {
	pm.client.SetHelloHostname("")
	return nil
}

func (pm PostMan) Shutdown(ctx context.Context) error {
	pm.client.Close()
	return nil
}

func (pm PostMan) Send(e *Email, t ...time.Duration) error {
	if e == nil {
		return nil
	}

	timeout := pm.cfg.Timeout
	if len(t) != 0 {
		timeout = t[0]
	}

	if e.From == "" {
		e.From = pm.cfg.From
	}
	if e.From == "" {
		e.From = pm.cfg.Username
	}

	if pm.cfg.SSL {
		return e.SendWithTLS(
			pm.cfg.Address,
			smtp.PlainAuth(pm.cfg.Identity, pm.cfg.Username, pm.cfg.Password, pm.cfg.Host),
			&tls.Config{ServerName: pm.cfg.Host},
		)
	}

	return pm.client.Send(e, timeout)
}
