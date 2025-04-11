package auth_service_usecase

import (
	"errors"

	"go.uber.org/zap"
)

type ValidateResult struct {
	Valid     bool
	UserID    string
	Roles     []string
	ExpiresAt int64
}

func (uc *AuthUsecase) ValidateAccessToken(token string) (*ValidateResult, error) {
	if token == "" {
		return nil, errors.New("access token is required")
	}

	claims, err := uc.jwtManager.ValidateToken(token)
	if err != nil {
		uc.logger.Error("failed to validate access token", zap.Error(err))
		return nil, err
	}

	return &ValidateResult{
		Valid:     true,
		UserID:    claims.UserID,
		Roles:     claims.Roles,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}
