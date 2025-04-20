package user_usecase

import (
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (uc *UserUsecase) UpdateTokens(params UpdateTokensParams) error {
	if err := params.Validate(); err != nil {
		uc.logger.Warn("validation failed",
			zap.String("user_id", params.UserID),
			zap.Error(err))
		return status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	if _, err := uc.repo.GetUserByID(params.UserID); err != nil {
		uc.logger.Error("failed to get user for token update",
			zap.String("user_id", params.UserID),
			zap.Error(err))

		if status.Code(err) != codes.Unknown {
			return err
		}
		return status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	if err := uc.repo.UpdateUserTokens(params.UserID, params.HashedRefreshToken, params.HashedAccessToken); err != nil {
		uc.logger.Error("failed to update tokens",
			zap.String("user_id", params.UserID),
			zap.Error(err))

		if status.Code(err) != codes.Unknown {
			return err
		}
		return status.Errorf(codes.Internal, "failed to update tokens: %v", err)
	}

	uc.logger.Info("tokens updated successfully",
		zap.String("user_id", params.UserID))
	return nil
}
