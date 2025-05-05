package main

import (
	"fmt"
	stdlog "log"
	"net"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/vwency/microservices_golang/internal/auth_service/endpoints"
	"github.com/vwency/microservices_golang/internal/auth_service/service"
	"github.com/vwency/microservices_golang/internal/auth_service/transport"
	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/jwt"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var Cfg config.ServiceConfig

func main() {
	// Load config
	env := config.DetectEnv()
	config.Init(env, "auth_service", &Cfg)

	// Initialize zap logger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		stdlog.Fatalf("failed to initialize zap logger: %v", err)
	}
	defer zapLogger.Sync()

	// Wrap zap.Logger as a go-kit log.Logger
	kitLogger := log.NewLogfmtLogger(os.Stdout) // simple logfmt fallback
	kitLogger = level.NewFilter(kitLogger, level.AllowDebug())

	// Parse JWT durations
	accessTokenTTL, err := time.ParseDuration(Cfg.Jwt.AccessTokenTtl)
	if err != nil {
		level.Error(kitLogger).Log("msg", "invalid access_token_ttl", "err", err)
		stdlog.Fatalf("invalid access_token_ttl: %v", err)
	}
	refreshTokenTTL, err := time.ParseDuration(Cfg.Jwt.RefreshTokenTtl)
	if err != nil {
		level.Error(kitLogger).Log("msg", "invalid refresh_token_ttl", "err", err)
		stdlog.Fatalf("invalid refresh_token_ttl: %v", err)
	}

	// Init JWT manager
	jwtManager, err := jwt.NewJWTManager(Cfg.Jwt.Secret, accessTokenTTL, refreshTokenTTL)
	if err != nil {
		level.Error(kitLogger).Log("msg", "failed to create JWT manager", "err", err)
		stdlog.Fatalf("failed to create JWT manager: %v", err)
	}

	// Connect to user_database
	dbConn, err := grpc.Dial(Cfg.UserDatabase.URL, grpc.WithInsecure())
	if err != nil {
		level.Error(kitLogger).Log("msg", "failed to connect to user_database", "err", err)
		stdlog.Fatalf("failed to connect to user_database: %v", err)
	}
	defer dbConn.Close()

	dbClient := databasev1.NewDatabaseInitServiceClient(dbConn)

	// Initialize service layer (business logic)
	authService := service.NewService(dbClient, jwtManager, kitLogger, Cfg.Jwt.HashPepper)

	// Create Go-Kit endpoints
	authEndpoints := endpoints.MakeEndpoints(authService)

	// Setup and start gRPC server
	addr := fmt.Sprintf(":%s", Cfg.App.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		level.Error(kitLogger).Log("msg", "failed to listen", "addr", addr, "err", err)
		stdlog.Fatalf("failed to listen on %s: %v", addr, err)
	}

	grpcServer := grpc.NewServer()
	transport.RegisterGRPCServer(grpcServer, authEndpoints)

	level.Info(kitLogger).Log("msg", "auth_service gRPC server started", "env", env, "addr", addr)

	if err := grpcServer.Serve(listener); err != nil {
		level.Error(kitLogger).Log("msg", "failed to serve gRPC", "err", err)
		stdlog.Fatalf("failed to serve gRPC: %v", err)
	}
}
