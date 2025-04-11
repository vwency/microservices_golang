package auth_service_handler

import (
	"context"

	auth_service_usecase "github.com/vwency/microservices_golang/internal/auth_service/usecase"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LogoutHandler struct {
	usecase *auth_service_usecase.LogoutUsecase
	logger  *zap.Logger
}

func NewLogoutHandler(usecase *auth_service_usecase.LogoutUsecase, logger *zap.Logger) *LogoutHandler {
	return &LogoutHandler{usecase: usecase, logger: logger}
}

func (h *LogoutHandler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	success, err := h.usecase.Logout(ctx, req.Username)
	if err != nil {
		h.logger.Error("Logout failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "logout failed: %v", err)
	}

	return &authv1.LogoutResponse{
		Success: success,
		Message: "Logged out successfully",
	}, nil
}
