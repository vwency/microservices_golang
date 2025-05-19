package service

import (
	"context"

	"github.com/vwency/microservices_golang/internal/auth_service/service/login"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
)

func (s *service) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	return login.Login(
		s.dbClient,
		s.jwtManager,
		s.logger,
		s.tokenPepper,
		ctx,
		req,
	)
}
