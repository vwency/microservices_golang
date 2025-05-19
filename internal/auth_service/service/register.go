package service

import (
	"context"

	"github.com/vwency/microservices_golang/internal/auth_service/service/register"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
)

func (s *service) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	return register.Register(
		s.dbClient,
		s.logger,
		s.jwtManager,
		s.tokenPepper,
		ctx,
		req,
	)
}
