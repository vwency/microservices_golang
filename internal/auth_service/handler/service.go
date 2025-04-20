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
	refreshHandler  *RefreshHandler
	logoutHandler   *LogoutHandler
}

func NewServer(
	authUsecase *auth_service_usecase.AuthUsecase,
	logger *zap.Logger,
) *Server {
	return &Server{
		loginHandler:    NewLoginHandler(authUsecase, logger),
		registerHandler: NewRegisterHandler(authUsecase, logger),
		validateHandler: NewValidateHandler(authUsecase, logger),
		refreshHandler:  NewRefreshHandler(authUsecase, logger),
		logoutHandler:   NewLogoutHandler(authUsecase, logger),
	}
}

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

func (s *Server) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	return s.refreshHandler.Refresh(ctx, req)
}

func (s *Server) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	return s.logoutHandler.Logout(ctx, req)
}
