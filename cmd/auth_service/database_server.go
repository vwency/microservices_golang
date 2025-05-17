package main

import (
	"context"
	"fmt"
	"time"

	"github.com/vwency/microservices_golang/pkg/config"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func newDatabaseConnection(cfg config.ServiceConfig) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(cfg.UserDatabase.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к user_database: %w", err)
	}

	conn.Connect()

	state := conn.GetState()
	if state != connectivity.Ready {
		if !conn.WaitForStateChange(ctx, state) {
			conn.Close()
			return nil, fmt.Errorf("connection timeout: context deadline exceeded")
		}
	}

	return conn, nil
}

func newDatabaseClient(conn *grpc.ClientConn) databasev1.DatabaseInitServiceClient {
	return databasev1.NewDatabaseInitServiceClient(conn)
}
