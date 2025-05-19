package service

import (
	"context"

	"github.com/vwency/microservices_golang/internal/auth_service/service/refresh"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
)

func (s *service) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	return refresh.Refresh(
		s.dbClient,
		s.jwtManager,
		s.logger,
		s.tokenPepper,
		ctx,
		req,
	)
}
