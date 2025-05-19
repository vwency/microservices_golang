package service

import (
	"context"

	"github.com/vwency/microservices_golang/internal/auth_service/service/validate"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
)

func (s *service) ValidateAccessToken(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error) {
	return validate.ValidateAccessToken(validate.ValidateDeps{
		DbClient:    s.dbClient,
		JwtManager:  s.jwtManager,
		Logger:      s.logger,
		TokenPepper: s.tokenPepper,
	}, ctx, req)
}
