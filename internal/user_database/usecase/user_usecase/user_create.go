package user_usecase

import (
	"github.com/google/uuid"
	"github.com/vwency/microservices_golang/internal/user_database/models"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (uc *UserUsecase) CreateUser(params CreateUserParams) error {
	if err := params.Validate(); err != nil {
		uc.logger.Warn("validation failed",
			zap.String("username", params.Username),
			zap.Error(err))
		return status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	// Проверяем что UserID передан (теперь он обязателен)
	if params.UserID == "" {
		uc.logger.Error("user_id is required")
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Парсим переданный UserID
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		uc.logger.Warn("invalid user_id format",
			zap.String("user_id", params.UserID),
			zap.Error(err))
		return status.Errorf(codes.InvalidArgument, "invalid user_id format")
	}

	// Проверяем существование пользователя по ID
	existingUserByID, err := uc.repo.GetUserByID(params.UserID)
	if err != nil && status.Code(err) != codes.NotFound {
		uc.logger.Error("failed to check user existence by ID",
			zap.String("user_id", params.UserID),
			zap.Error(err))
		return status.Errorf(codes.Internal, "check user existence by ID failed: %v", err)
	}
	if existingUserByID != nil {
		uc.logger.Warn("user with this ID already exists",
			zap.String("user_id", params.UserID))
		return ErrUserAlreadyExists
	}

	// Проверяем существование по username/email
	existingUser, err := uc.repo.GetUserByUsernameOrEmail(params.Username, params.Email)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			uc.logger.Error("failed to check user existence",
				zap.String("username", params.Username),
				zap.Error(err))
			return status.Errorf(codes.Internal, "check user existence failed: %v", err)
		}
	}

	if existingUser != nil {
		uc.logger.Warn("user already exists",
			zap.String("username", params.Username))
		return ErrUserAlreadyExists
	}

	// Создаем модель пользователя с переданным ID
	user := models.User{
		ID:                 userID,
		Username:           params.Username,
		HashedPassword:     params.HashedPassword,
		HashedRefreshToken: params.HashedRt,
		HashedAccessToken:  params.HashedAt,
	}

	if params.Email != "" {
		user.Email = &params.Email
	}

	if err := uc.repo.AddUser(&user); err != nil {
		uc.logger.Error("failed to create user",
			zap.String("username", params.Username),
			zap.String("user_id", params.UserID),
			zap.Error(err))

		if status.Code(err) != codes.Unknown {
			return err
		}
		return status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	uc.logger.Info("user created successfully",
		zap.String("username", params.Username),
		zap.String("user_id", params.UserID))
	return nil
}
