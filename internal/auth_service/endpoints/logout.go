package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MakeLogoutEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*authv1.LogoutRequest)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "invalid request type")
		}
		return s.Logout(ctx, req)
	}
}
