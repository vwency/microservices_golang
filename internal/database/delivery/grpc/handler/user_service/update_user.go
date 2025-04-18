package handler_user_service_grpc

import (
	"context"
	"errors"

	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase"
	pb "github.com/vwency/microservices_golang/proto/database"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateUserHandler struct {
	pb.UnimplementedDatabaseInitServiceServer
	uc     *user_usecase.UserUsecase
	logger *zap.Logger
}

func NewUpdateUserHandler(uc *user_usecase.UserUsecase, logger *zap.Logger) *UpdateUserHandler {
	return &UpdateUserHandler{
		uc:     uc,
		logger: logger.With(zap.String("handler", "update_user")),
	}
}

func (h *UpdateUserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.GetUserId() == "" {
		h.logger.Warn("empty user_id in request")
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetHashedRefreshToken() == "" || req.GetHashedAccessToken() == "" {
		h.logger.Warn("empty token fields in request",
			zap.String("user_id", req.GetUserId()))
		return nil, status.Error(codes.InvalidArgument, "token fields cannot be empty")
	}

	updateParams := user_usecase.UpdateTokensParams{
		UserID:             req.GetUserId(),
		HashedRefreshToken: req.GetHashedRefreshToken(),
		HashedAccessToken:  req.GetHashedAccessToken(),
	}

	err := h.uc.UpdateTokens(updateParams)
	if err != nil {
		h.logger.Error("failed to update user tokens",
			zap.String("user_id", req.GetUserId()),
			zap.Error(err))

		switch {
		case errors.Is(err, user_usecase.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, "user not found")
		default:
			return nil, status.Error(codes.Internal, "failed to update tokens")
		}
	}

	h.logger.Info("tokens updated successfully",
		zap.String("user_id", req.GetUserId()))

	return &pb.UpdateUserResponse{
		Success: true,
		Message: "Tokens updated successfully",
	}, nil
}
