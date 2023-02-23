// Package pprof
//
// @author: xwc1125
// @date: 2021/7/22
package pprof

import (
	"net/http"

	"github.com/chain5j/logger"
)

type Pprof struct {
	log    logger.Logger
	server *http.Server
}

func NewPprof(addr string) *Pprof {
	return &Pprof{
		log: logger.Log("pprof"),
		server: &http.Server{
			Addr:    addr,
			Handler: nil,
		},
	}
}
func (p *Pprof) Start() error {
	return p.server.ListenAndServe()
}

func (p *Pprof) Stop() {
	p.log.Info("pprof stopping...")
	p.server.Close()
	p.log.Info("pprof stopped")
}
