package service

import (
	"context"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type HelloService interface {
	SayHello(ctx context.Context, text string) (string, error)
}

type helloService struct {
	logger log.Logger
}

func NewHelloService(logger log.Logger) HelloService {
	return &helloService{
		logger: logger,
	}
}

func (s *helloService) SayHello(ctx context.Context, text string) (string, error) {
	level.Info(s.logger).Log("msg", "processing greeting", "text", text)

	if strings.Contains(strings.ToLower(text), "hello") {
		return "hello", nil
	}
	return "None", nil
}
