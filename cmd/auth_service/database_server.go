package main

import (
	"context"
	"fmt"
	"time"

	"github.com/vwency/microservices_golang/pkg/config"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"google.golang.org/grpc"
)

func newDatabaseConnection(cfg config.ServiceConfig) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.UserDatabase.URL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к user_database: %w", err)
	}
	return conn, nil
}

func newDatabaseClient(conn *grpc.ClientConn) databasev1.DatabaseInitServiceClient {
	return databasev1.NewDatabaseInitServiceClient(conn)
}
