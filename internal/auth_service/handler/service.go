package auth_service_handler

import (
	"context"

	auth_service_usecase "github.com/vwency/microservices_golang/internal/auth_service/usecase"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	authv1.UnimplementedAuthServiceServer
	loginHandler    *LoginHandler
	registerHandler *RegisterHandler
	validateHandler *ValidateHandler
}

func NewServer(uc *auth_service_usecase.AuthUsecase, logger *zap.Logger) *Server {
	return &Server{
		loginHandler:    NewLoginHandler(uc, logger),
		registerHandler: NewRegisterHandler(uc, logger),
		validateHandler: NewValidateHandler(uc, logger),
	}
}

// RegisterService переименован, чтобы избежать конфликта с методом Register
func (s *Server) RegisterService(grpcServer *grpc.Server) {
	authv1.RegisterAuthServiceServer(grpcServer, s)
}

func (s *Server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	return s.loginHandler.Login(ctx, req)
}

func (s *Server) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	return s.registerHandler.Register(ctx, req)
}

func (s *Server) Validate(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error) {
	return s.validateHandler.Validate(ctx, req)
}
