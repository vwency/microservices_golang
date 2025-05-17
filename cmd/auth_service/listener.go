package main

import (
	"fmt"
	"net"

	"github.com/vwency/microservices_golang/pkg/config"
)

func newListener(cfg config.ServiceConfig) (net.Listener, error) {
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("не удалось прослушивать %s: %w", addr, err)
	}
	return listener, nil
}
