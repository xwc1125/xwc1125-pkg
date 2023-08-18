// Package email
package email

import (
	"crypto/tls"
	"net/smtp"
	"time"

	"github.com/chain5j/logger"
	"github.com/jordan-wright/email"
)

type Email interface {
	Send(e *email.Email, timeout time.Duration) (err error)
	Close()
}

var (
	_ Email = new(mail)
)

type mail struct {
	log  logger.Logger
	pool *email.Pool
}

func NewEmail(config Config) (Email, error) {
	return NewEmailByAuth(config, smtp.PlainAuth("", config.Username, config.Password,
		config.Addr))
}

func NewEmailByAuth(config Config, auth smtp.Auth, tlsConfig ...*tls.Config) (Email, error) {
	log := logger.Log("email")
	addr := config.Addr + ":" + config.Port
	emailPool, err := email.NewPool(addr, config.PoolCount, auth, tlsConfig...)
	if err != nil {
		log.Error("new email pool err", "err", err)
		return nil, err
	}
	return &mail{
		log:  log,
		pool: emailPool,
	}, nil
}

func (m *mail) Send(e *email.Email, timeout time.Duration) (err error) {
	return m.pool.Send(e, timeout)
}

func (m *mail) Close() {
	m.pool.Close()
}
