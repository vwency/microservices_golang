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

		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "failed to get user by credentials: %v", err)
	}

	return user, nil
}
