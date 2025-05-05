package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
)

// ... (остальные типы и константы из main.go остаются без изменений)

// ValidateAccessToken проверяет access token на валидность и возвращает результат.
func (s *service) ValidateAccessToken(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error) {
	token := req.GetAccessToken()
	if token == "" {
		return nil, errors.New("access token is required")
	}

	// 1. Проверка токена через JWTManager
	claims, err := s.jwtManager.ValidateToken(token)
	if err != nil {
		s.logger.Log("error", fmt.Sprintf("Failed to validate token: %v", err), "token", token)
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	// 2. Извлечение userID
	userID, ok := claims["UserID"].(string)
	if !ok || userID == "" {
		s.logger.Log("error", "Invalid or missing UserID in claims", "claims", claims)
		return nil, errors.New("invalid token: userID missing or not a string")
	}

	// 3. Извлечение ролей
	rolesRaw, ok := claims["Roles"].([]interface{})
	if !ok {
		s.logger.Log("error", "Invalid or missing Roles in claims", "claims", claims)
		return nil, errors.New("invalid token: roles missing or not an array")
	}

	var roles []string
	for _, r := range rolesRaw {
		roleStr, ok := r.(string)
		if !ok {
			s.logger.Log("error", fmt.Sprintf("Invalid role type: %v", r))
			return nil, fmt.Errorf("invalid role type: %v", r)
		}
		roles = append(roles, roleStr)
	}

	// 4. Получение пользователя из базы
	getUserResp, err := s.dbClient.GetUser(ctx, &databasev1.GetUserRequest{UserId: &userID})
	if err != nil {
		s.logger.Log("error", fmt.Sprintf("Failed to fetch user: %v", err), "userID", userID)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if !getUserResp.Found {
		s.logger.Log("warn", "User not found", "userID", userID)
		return nil, errors.New("user not found")
	}

	// 5. Проверка совпадения токена с хэшем из базы
	match, err := authutils.ComparePasswordAndHash(s.tokenPepper, token, getUserResp.HashedAccessToken)
	if err != nil {
		s.logger.Log("error", fmt.Sprintf("Token hash comparison failed: %v", err), "userID", userID)
		return nil, fmt.Errorf("token hash comparison failed: %w", err)
	}

	if !match {
		s.logger.Log("warn", "Token hash mismatch", "userID", userID)
		return nil, errors.New("invalid token: hash mismatch")
	}

	// 6. Проверка срока действия токена
	expiryFloat, ok := claims["exp"].(float64)
	if !ok {
		s.logger.Log("error", "Invalid or missing expiry in claims", "claims", claims)
		return nil, errors.New("invalid token: expiry missing or not a number")
	}
	expiry := int64(expiryFloat)

	if time.Now().Unix() > expiry {
		s.logger.Log("warn", "Token expired", "userID", userID, "expiry", expiry)
		return nil, errors.New("access token has expired")
	}

	// Успешная валидация
	return &authv1.ValidateResponse{
		Valid:     true,
		UserId:    userID,
		Roles:     roles,
		ExpiresAt: expiry,
	}, nil
}
