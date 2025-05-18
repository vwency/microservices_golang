package gokit

import (
	"github.com/vwency/microservices_golang/internal/auth_service/endpoints"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
)

func NewAuthEndpoints(svc service.AuthService) endpoints.Endpoints {
	return endpoints.MakeEndpoints(svc)
}
