package auth_service_handler

import (
	"context"

	auth_service_usecase "github.com/vwency/microservices_golang/internal/auth_service/usecase"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RefreshHandler struct {
	authUsecase *auth_service_usecase.AuthUsecase
	logger      *zap.Logger
}

func NewRefreshHandler(authUsecase *auth_service_usecase.AuthUsecase, logger *zap.Logger) *RefreshHandler {
	return &RefreshHandler{
		authUsecase: authUsecase,
		logger:      logger,
	}
}

func (h *RefreshHandler) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	tokens, err := h.authUsecase.Refresh(ctx, req.RefreshToken)
	if err != nil {
		h.logger.Error("Refresh failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "refresh failed: %v", err)
	}

	return &authv1.RefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt.Unix(),
	}, nil
}
