package gokit

import (
	kitlog "github.com/go-kit/kit/log"
	"github.com/vwency/microservices_golang/internal/user_database/repository"
	"github.com/vwency/microservices_golang/internal/user_database/service"
)

func NewService(repo repository.Repository, logger kitlog.Logger) service.Service {
	return service.NewService(repo, logger)
}
