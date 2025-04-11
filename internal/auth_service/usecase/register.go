package auth_service_usecase

import (
	"context"
	"encoding/base64"
	"fmt"

	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
)

func (uc *AuthUsecase) Register(ctx context.Context, username, password, email string) (*authv1.RegisterResponse, error) {
	salt := []byte(username)
	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	if len(hashedPassword) == 0 {
		uc.logger.Error("Failed to hash password with argon2")
		return nil, fmt.Errorf("failed to hash password: %v", "hashing error")
	}

	encodedPassword := base64.StdEncoding.EncodeToString(hashedPassword)

	addUserReq := &databasev1.AddUserRequest{
		Username:       username,
		HashedPassword: encodedPassword,
		HashedRt:       "qweqwe",       // Placeholder for refresh token (you'll need to handle it properly)
		AccessRt:       "access-token", // Placeholder for access token
		Email:          email,
	}

	addUserResp, err := uc.dbClient.AddUser(ctx, addUserReq)
	if err != nil {
		uc.logger.Error("Failed to add user to database", zap.Error(err))
		return nil, err
	}

	if !addUserResp.Success {
		uc.logger.Error("Failed to add user to database", zap.String("message", addUserResp.Message))
		return nil, fmt.Errorf("failed to add user: %v", addUserResp.Message)
	}

	accessToken, expiresAt, err := uc.jwtManager.GenerateAccessToken(username, []string{"user"})
	if err != nil {
		uc.logger.Error("Failed to generate access token", zap.Error(err))
		return nil, err
	}

	refreshToken, _, err := uc.jwtManager.GenerateRefreshToken(username, []string{"user"})
	if err != nil {
		uc.logger.Error("Failed to generate refresh token", zap.Error(err))
		return nil, err
	}

	return &authv1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}
