package auth_service_handler

import (
	"context"

	auth_service_usecase "github.com/vwency/microservices_golang/internal/auth_service/usecase"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RegisterHandler struct {
	usecase *auth_service_usecase.AuthUsecase
	logger  *zap.Logger
}

func NewRegisterHandler(usecase *auth_service_usecase.AuthUsecase, logger *zap.Logger) *RegisterHandler {
	return &RegisterHandler{usecase: usecase, logger: logger}
}

func (h *RegisterHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "username, password, and email are required")
	}

	response, err := h.usecase.Register(ctx, req.Username, req.Password, req.Email)
	if err != nil {
		h.logger.Error("Registration failed",
			zap.Error(err),
			zap.String("username", req.Username),
			zap.String("email", req.Email))
		return nil, status.Errorf(codes.Internal, "registration failed: %v", err)
	}

	return response, nil
}
