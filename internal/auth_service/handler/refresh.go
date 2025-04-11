// auth_service_handler/refresh_handler.go
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
	refreshUsecase *auth_service_usecase.RefreshUsecase
	logger         *zap.Logger
}

func NewRefreshHandler(refreshUsecase *auth_service_usecase.RefreshUsecase, logger *zap.Logger) *RefreshHandler {
	return &RefreshHandler{refreshUsecase: refreshUsecase, logger: logger}
}

func (h *RefreshHandler) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	tokens, err := h.refreshUsecase.Refresh(ctx, req.RefreshToken)
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
