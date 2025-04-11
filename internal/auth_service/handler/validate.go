package auth_service_handler

import (
	"context"

	auth_service_usecase "github.com/vwency/microservices_golang/internal/auth_service/usecase"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValidateHandler struct {
	usecase *auth_service_usecase.AuthUsecase
	logger  *zap.Logger
}

func NewValidateHandler(usecase *auth_service_usecase.AuthUsecase, logger *zap.Logger) *ValidateHandler {
	return &ValidateHandler{usecase: usecase, logger: logger}
}

func (h *ValidateHandler) Validate(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error) {
	result, err := h.usecase.ValidateAccessToken(req.AccessToken)
	if err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return &authv1.ValidateResponse{
		Valid:     result.Valid,
		UserId:    result.UserID,
		Roles:     result.Roles,
		ExpiresAt: result.ExpiresAt,
	}, nil
}
