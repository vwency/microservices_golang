package auth_service_usecase

import (
	"context"
	"errors"
	"fmt"

	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"go.uber.org/zap"
)

// ValidateResult содержит результаты валидации токена.
type ValidateResult struct {
	Valid     bool
	UserID    string
	Roles     []string
	ExpiresAt int64
}

// ValidateAccessToken проверяет валидность access токена.
func (uc *AuthUsecase) ValidateAccessToken(ctx context.Context, token string) (*ValidateResult, error) {
	// Проверка, был ли передан токен
	if token == "" {
		return nil, errors.New("access token is required")
	}

	// Валидация токена и извлечение claims
	claims, err := uc.jwtManager.ValidateToken(token)
	if err != nil {
		uc.logger.Error("Failed to validate access token", zap.Error(err))
		return nil, fmt.Errorf("failed to validate access token: %w", err)
	}

	// Извлечение userID из claims
	userID, ok := claims["UserID"].(string)
	if !ok {
		uc.logger.Error("userID is missing or not a string", zap.Any("claims", claims))
		return nil, errors.New("invalid token: userID missing or not a string")
	}

	// Извлечение ролей из claims и проверка их типа
	rolesInterface, ok := claims["Roles"].([]interface{})
	if !ok {
		uc.logger.Error("roles are missing or not an array", zap.Any("claims", claims))
		return nil, errors.New("invalid token: roles missing or not an array")
	}

	// Преобразование ролей в []string
	var rolesStr []string
	for _, role := range rolesInterface {
		roleStr, ok := role.(string)
		if !ok {
			uc.logger.Error("invalid role type", zap.Any("role", role))
			return nil, fmt.Errorf("invalid role type: %v", role)
		}
		rolesStr = append(rolesStr, roleStr)
	}

	// Получение пользователя из базы данных
	getUserResp, err := uc.dbClient.GetUser(ctx, &databasev1.GetUserRequest{Username: &userID})
	if err != nil {
		uc.logger.Error("Failed to fetch user during token validation", zap.Error(err), zap.String("userID", userID))
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Проверка наличия пользователя в базе данных
	if !getUserResp.Found {
		uc.logger.Error("User not found during token validation", zap.String("userID", userID))
		return nil, errors.New("user not found")
	}

	// Сравнение токена с хэшированным значением из базы данных
	match, err := authutils.ComparePasswordAndHash(uc.tokenPepper, token, getUserResp.HashedAccessToken)
	if err != nil {
		uc.logger.Error("Failed to compare token hashes", zap.Error(err), zap.String("userID", userID))
		return nil, fmt.Errorf("failed to compare token hashes: %w", err)
	}

	// Проверка на несовпадение хэшей
	if !match {
		uc.logger.Warn("Access token hash mismatch", zap.String("userID", userID))
		return nil, errors.New("invalid access token")
	}

	// Проверка срока действия токена (проверяем, что "exp" в claims — это int64)
	expiry, ok := claims["exp"].(float64) // JWT "exp" обычно в формате float64
	if !ok {
		uc.logger.Error("Token expiry missing or invalid", zap.Any("claims", claims))
		return nil, errors.New("invalid token: expiry missing or not a valid int64")
	}

	// Возвращаем результаты валидации
	return &ValidateResult{
		Valid:     true,
		UserID:    userID,
		Roles:     rolesStr,
		ExpiresAt: int64(expiry), // Преобразуем float64 в int64
	}, nil
}
