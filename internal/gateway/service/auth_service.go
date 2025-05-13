package service

import (
	"context"
	"fmt"

	"github.com/vwency/microservices_golang/proto/auth_service"
	"google.golang.org/grpc"
)

type AuthServiceClient struct {
	Client auth_service.AuthServiceClient
}

func NewAuthServiceClient(authAddress string) (*AuthServiceClient, error) {
	authConn, err := grpc.Dial(authAddress, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to auth service: %v", err)
	}

	authClient := auth_service.NewAuthServiceClient(authConn)

	return &AuthServiceClient{
		Client: authClient,
	}, nil
}

func (c *AuthServiceClient) Register(ctx context.Context, username, password, email string) (*auth_service.RegisterResponse, error) {
	req := &auth_service.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
	}
	return c.Client.Register(ctx, req)
}

func (c *AuthServiceClient) Login(ctx context.Context, username, password string) (*auth_service.LoginResponse, error) {
	req := &auth_service.LoginRequest{
		Username: username,
		Password: password,
	}
	return c.Client.Login(ctx, req)
}

func (c *AuthServiceClient) Logout(ctx context.Context, username, accessToken string) (*auth_service.LogoutResponse, error) {
	req := &auth_service.LogoutRequest{
		Username:    username,
		AccessToken: accessToken,
	}
	return c.Client.Logout(ctx, req)
}

func (c *AuthServiceClient) Validate(ctx context.Context, token string) (*auth_service.ValidateResponse, error) {
	req := &auth_service.ValidateRequest{
		AccessToken: token,
	}
	return c.Client.ValidateAccessToken(ctx, req)
}

func (c *AuthServiceClient) Refresh(ctx context.Context, refreshToken string) (*auth_service.RefreshResponse, error) {
	req := &auth_service.RefreshRequest{
		RefreshToken: refreshToken,
	}
	return c.Client.Refresh(ctx, req)
}
