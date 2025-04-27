package handler_jwt

import (
	"context"
	"fmt"
	"time"

	usecase_jwt "github.com/vwency/microservices_golang/internal/jwt_service/usecase"
	"github.com/vwency/microservices_golang/proto/jwt_service"
)

type JwtHandler struct {
	jwt_service.UnimplementedJwtServiceServer
	usecase *usecase_jwt.JwtUsecase
}

func NewJwtHandler(usecase *usecase_jwt.JwtUsecase) *JwtHandler {
	return &JwtHandler{usecase: usecase}
}

func (h *JwtHandler) GenerateToken(ctx context.Context, req *jwt_service.GenerateTokenRequest) (*jwt_service.GenerateTokenResponse, error) {
	expiresIn := time.Duration(req.GetExpiresIn()) * time.Second
	if req.GetExpiresIn() == 0 {
		expiresIn = 72 * time.Hour
	}

	token, expiresAt, err := h.usecase.GenerateToken(req.GetSecret(), req.GetPayload(), expiresIn)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	return &jwt_service.GenerateTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (h *JwtHandler) ValidateToken(ctx context.Context, req *jwt_service.ValidateTokenRequest) (*jwt_service.ValidateTokenResponse, error) {
	userID, payload, expiresAt, err := h.usecase.ValidateToken(req.GetToken(), req.GetSecret())
	if err != nil {
		return &jwt_service.ValidateTokenResponse{
			Valid:     false,
			Error:     err.Error(),
			ExpiresAt: expiresAt,
		}, nil
	}

	return &jwt_service.ValidateTokenResponse{
		UserId:    userID,
		Valid:     true,
		Payload:   payload,
		ExpiresAt: expiresAt,
	}, nil
}
