package service

import (
	"context"
	"fmt"

	"github.com/vwency/microservices_golang/proto/auth_service"
	database_init "github.com/vwency/microservices_golang/proto/database"
	"google.golang.org/grpc"
)

type AuthServiceClient struct {
	Client   auth_service.AuthServiceClient
	DBClient database_init.DatabaseInitServiceClient
}

func NewAuthServiceClient(authAddress, dbAddress string) (*AuthServiceClient, error) {
	authConn, err := grpc.Dial(authAddress, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to auth service: %v", err)
	}

	dbConn, err := grpc.Dial(dbAddress, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to database service: %v", err)
	}

	authClient := auth_service.NewAuthServiceClient(authConn)
	dbClient := database_init.NewDatabaseInitServiceClient(dbConn)

	return &AuthServiceClient{
		Client:   authClient,
		DBClient: dbClient,
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

func (c *AuthServiceClient) GetUser(ctx context.Context, username, email string) (*database_init.GetUserResponse, error) {
	req := &database_init.GetUserRequest{
		Username: username,
		Email:    email,
	}
	return c.DBClient.GetUser(ctx, req)
}
