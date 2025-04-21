package user_usecase

import (
	"github.com/vwency/microservices_golang/internal/user_database/models"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (uc *UserUsecase) GetUserByID(userID string) (*models.User, error) {
	if userID == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id must be provided")
	}
	user, err := uc.repo.GetUserByID(userID)
	if err != nil {
		uc.logger.Error("failed to get user by ID",
			zap.String("user_id", userID),
			zap.Error(err))

		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "failed to get user by ID: %v", err)
	}

	return user, nil
}
