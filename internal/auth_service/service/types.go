package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/vwency/microservices_golang/pkg/jwt"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type contextKey string

const ipContextKey = contextKey("ip")

type service struct {
	dbClient    databasev1.DatabaseInitServiceClient
	jwtManager  *jwt.JWTManager
	logger      log.Logger
	tokenPepper string
}

type Service struct {
	DBClient    databasev1.DatabaseInitServiceClient
	JWTManager  *jwt.JWTManager
	Logger      log.Logger
	TokenPepper string
}

type AuthService interface {
	Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error)
	Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error)
	Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error)
	Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error)
	ValidateAccessToken(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error)
}

type RegisterRequest struct {
	Username string
	Password string
	Email    string
}

type RegisterResponse struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
