package user_usecase

import (
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

		if status.Code(err) != codes.Unknown {
			return err
		}
		return status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	uc.logger.Info("user created successfully",
		zap.String("username", params.Username))
	return nil
}
