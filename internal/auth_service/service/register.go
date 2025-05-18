package service

import (
	"context"
	"sync"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	authv1 "github.com/vwency/microservices_golang/proto/auth_service"
	databasev1 "github.com/vwency/microservices_golang/proto/user_database"
	"github.com/vwency/microservices_golang/utils/authutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	// Логирование в одну операцию
	level.Info(s.logger).Log(
		"msg", "Attempting registration",
		"username", req.Username,
		"ip", getIPFromContext(ctx),
	)

	// Валидация в начале
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "username, password and email are required")
	}

	// Параллельное выполнение операций, которые могут выполняться одновременно
	var (
		userID          = uuid.New().String()
		hashedPassword  string
		accessToken     string
		accessExpiresAt time.Time
		refreshToken    string
		err             error
	)

	// Группа для параллельного выполнения
	var wg sync.WaitGroup
	wg.Add(3)

	// Хеширование пароля
	go func() {
		defer wg.Done()
		hashedPassword, err = authutils.GenHash(req.Username, req.Password, nil)
	}()

	// Генерация access token
	go func() {
		defer wg.Done()
		payload := map[string]interface{}{
			"UserID": userID,
			"Roles":  []interface{}{"user"},
		}
		accessToken, accessExpiresAt, err = s.jwtManager.GenerateAccessToken(payload)
	}()

	// Генерация refresh token
	go func() {
		defer wg.Done()
		payload := map[string]interface{}{
			"UserID": userID,
			"Roles":  []interface{}{"user"},
		}
		refreshToken, _, err = s.jwtManager.GenerateRefreshToken(payload)
	}()

	wg.Wait()

	// Проверка ошибок после параллельного выполнения
	if err != nil {
		level.Error(s.logger).Log("msg", "Error during parallel operations", "err", err)
		return nil, status.Errorf(codes.Internal, "operation failed: %v", err)
	}

	// Хеширование токенов (последовательно, так как зависит от предыдущих результатов)
	hashedAccessToken, err := authutils.GenHash(s.tokenPepper, accessToken, nil)
	if err != nil {
		level.Error(s.logger).Log("msg", "Failed to hash access token", "err", err)
		return nil, status.Errorf(codes.Internal, "failed to hash access token: %v", err)
	}

	hashedRefreshToken, err := authutils.GenHash(s.tokenPepper, refreshToken, nil)
	if err != nil {
		level.Error(s.logger).Log("msg", "Failed to hash refresh token", "err", err)
		return nil, status.Errorf(codes.Internal, "failed to hash refresh token: %v", err)
	}

	// Добавление пользователя
	if _, err := s.dbClient.AddUser(ctx, &databasev1.AddUserRequest{
		Username:           req.Username,
		HashedPassword:     hashedPassword,
		Email:              req.Email,
		HashedAccessToken:  hashedAccessToken,
		HashedRefreshToken: hashedRefreshToken,
		UserId:             &userID,
	}); err != nil {
		level.Error(s.logger).Log("msg", "Failed to add user", "err", err)
		return nil, status.Errorf(codes.Internal, "failed to add user: %v", err)
	}

	level.Info(s.logger).Log(
		"msg", "User registered successfully",
		"user_id", userID,
		"username", req.Username,
	)

	return &authv1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}
