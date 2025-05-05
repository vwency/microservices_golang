package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MakeRegisterEndpoint(s service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*authv1.RegisterRequest)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "invalid request type")
		}

		res, err := s.Register(ctx, req)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "register failed: %v", err)
		}

		return res, nil
	}
}
