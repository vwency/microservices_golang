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

func (c *AuthServiceClient) GetUser(ctx context.Context, username, email string) (*database_init.GetUserResponse, error) {
	req := &database_init.GetUserRequest{
		Username: username,
		Email:    email,
	}

	resp, err := c.DBClient.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %v", err)
	}

	return resp, nil
}
