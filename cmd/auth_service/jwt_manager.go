package main

import (
	"fmt"
	"time"

	"github.com/vwency/microservices_golang/pkg/config"
	"github.com/vwency/microservices_golang/pkg/jwt"
)

func newJWTManager(cfg config.ServiceConfig) (*jwt.JWTManager, error) {
	accessTokenTTL, err := time.ParseDuration(cfg.Jwt.AccessTokenTtl)
	if err != nil {
		return nil, fmt.Errorf("неверный access_token_ttl: %w", err)
	}

	refreshTokenTTL, err := time.ParseDuration(cfg.Jwt.RefreshTokenTtl)
	if err != nil {
		return nil, fmt.Errorf("неверный refresh_token_ttl: %w", err)
	}

	return jwt.NewJWTManager(cfg.Jwt.Secret, accessTokenTTL, refreshTokenTTL)
}
