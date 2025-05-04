package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
)

func MakeLoginEndpoint(s *service.LoginService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*authv1.LoginRequest)
		if !ok {
			return nil, service.ErrInvalidCredentials
		}

		tokenPair, err := s.Login(ctx, req.Username, req.Password)
		if err != nil {
			return nil, err
		}

		return &authv1.LoginResponse{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresAt:    tokenPair.ExpiresAt.Unix(),
		}, nil
	}
}
