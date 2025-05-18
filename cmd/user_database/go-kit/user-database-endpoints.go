package gokit

import (
	"github.com/vwency/microservices_golang/internal/user_database/endpoints"
	"github.com/vwency/microservices_golang/internal/user_database/service"
)

func NewEndpoints(svc service.Service) endpoints.Endpoints {
	return endpoints.MakeEndpoints(svc)
}
