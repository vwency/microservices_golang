package user_usecase

import (
	"github.com/vwency/microservices_golang/internal/database/models"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (uc *UserUsecase) GetUser(params UserParams) (*models.User, error) {
	if params.UserID == "" && params.Username == "" && params.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "at least one search parameter must be provided")
	}

	if params.UserID != "" {
		user, err := uc.repo.GetUserByID(params.UserID)
		if err != nil {
			uc.logger.Error("failed to get user by ID",
				zap.String("user_id", params.UserID),
				zap.Error(err))

			// Пробрасываем все gRPC-ошибки как есть
			if status.Code(err) != codes.Unknown {
				return nil, err
			}
			return nil, status.Errorf(codes.Internal, "failed to get user by ID: %v", err)
		}
		return user, nil
	}

	user, err := uc.repo.GetUserByUsernameOrEmail(params.Username, params.Email)
	if err != nil {
		uc.logger.Error("failed to get user by username/email",
			zap.String("username", params.Username),
			zap.String("email", params.Email),
			zap.Error(err))

		// Пробрасываем все gRPC-ошибки как есть
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "failed to get user by credentials: %v", err)
	}

	return user, nil
}

func (uc *UserUsecase) CreateUser(params CreateUserParams) error {
	if err := params.Validate(); err != nil {
		uc.logger.Warn("validation failed",
			zap.String("username", params.Username),
			zap.Error(err))
		return status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	// Проверка существования пользователя
	existingUser, err := uc.repo.GetUserByUsernameOrEmail(params.Username, params.Email)
	if err != nil {
		// Игнорируем только NotFound ошибки
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

	user := &models.User{
		Username:           params.Username,
		HashedPassword:     params.HashedPassword,
		HashedRefreshToken: params.HashedRt,
		HashedAccessToken:  params.HashedAt,
		Email:              &params.Email,
	}

	if err := uc.repo.AddUser(user); err != nil {
		uc.logger.Error("failed to create user",
			zap.String("username", params.Username),
			zap.Error(err))

		// Сохраняем оригинальный статус ошибки
		if status.Code(err) != codes.Unknown {
			return err
		}
		return status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	uc.logger.Info("user created successfully",
		zap.String("username", params.Username))
	return nil
}

func (uc *UserUsecase) UpdateTokens(params UpdateTokensParams) error {
	if err := params.Validate(); err != nil {
		uc.logger.Warn("validation failed",
			zap.String("user_id", params.UserID),
			zap.Error(err))
		return status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	// Проверка существования пользователя
	if _, err := uc.repo.GetUserByID(params.UserID); err != nil {
		uc.logger.Error("failed to get user for token update",
			zap.String("user_id", params.UserID),
			zap.Error(err))

		// Пробрасываем оригинальную ошибку
		if status.Code(err) != codes.Unknown {
			return err
		}
		return status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	// Обновление токенов
	if err := uc.repo.UpdateUserTokens(params.UserID, params.HashedRefreshToken, params.HashedAccessToken); err != nil {
		uc.logger.Error("failed to update tokens",
			zap.String("user_id", params.UserID),
			zap.Error(err))

		// Пробрасываем оригинальную ошибку
		if status.Code(err) != codes.Unknown {
			return err
		}
		return status.Errorf(codes.Internal, "failed to update tokens: %v", err)
	}

	uc.logger.Info("tokens updated successfully",
		zap.String("user_id", params.UserID))
	return nil
}

func (uc *UserUsecase) GetUserByID(userID string) (*models.User, error) {
	// Ensure that the user_id is provided
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id must be provided")
	}
	user, err := uc.repo.GetUserByID(userID)
	if err != nil {
		uc.logger.Error("failed to get user by ID",
			zap.String("user_id", userID),
			zap.Error(err))

		// Propagate the error, if it's a known gRPC error, pass it as is
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "failed to get user by ID: %v", err)
	}

	// Successfully retrieved user
	return user, nil
}
