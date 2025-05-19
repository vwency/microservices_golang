package service

import (
	"context"

	"github.com/vwency/microservices_golang/internal/auth_service/service/logout"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
)

func (s *service) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	return logout.Logout(
		s.dbClient,
		s.logger,
		s.tokenPepper,
		ctx,
		req,
	)
}
