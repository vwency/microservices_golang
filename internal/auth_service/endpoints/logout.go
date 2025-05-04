package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
)

func MakeLogoutEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*authv1.LogoutRequest)
		return s.Logout(ctx, req)
	}
}
