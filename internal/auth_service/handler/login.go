package auth_service_handler

import (
	"context"

	auth_service_usecase "github.com/vwency/microservices_golang/internal/auth_service/usecase"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginHandler struct {
	usecase *auth_service_usecase.AuthUsecase
	logger  *zap.Logger
}

func NewLoginHandler(usecase *auth_service_usecase.AuthUsecase, logger *zap.Logger) *LoginHandler {
	return &LoginHandler{usecase: usecase, logger: logger}
}

func (h *LoginHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	tokens, err := h.usecase.Login(ctx, req.Username, req.Password)
	if err != nil {
		h.logger.Error("Login failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}

	return &authv1.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt.Unix(),
	}, nil
}
