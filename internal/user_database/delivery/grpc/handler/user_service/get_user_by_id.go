package handler_user_service_grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/usecase/user_usecase"
	pb "github.com/vwency/microservices_golang/proto/user_database"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserByIDHandler struct {
	uc     *user_usecase.UserUsecase
	logger *zap.Logger
}

func NewGetUserByIDHandler(uc *user_usecase.UserUsecase, logger *zap.Logger) *GetUserByIDHandler {
	return &GetUserByIDHandler{
		uc:     uc,
		logger: logger.With(zap.String("handler", "get_user_by_id")),
	}
}

func (h *GetUserByIDHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	if req.GetUserId() == "" {
		h.logger.Warn("empty user_id")
		return nil, status.Error(codes.InvalidArgument, "user_id must be provided")
	}

	userID := req.GetUserId()
	user, err := h.uc.GetUserByID(userID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			h.logger.Info("user not found", zap.String("user_id", userID))
			return &pb.GetUserByIDResponse{
				Found:   false,
				Message: "User not found",
			}, nil
		}

		h.logger.Error("failed to get user", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	return &pb.GetUserByIDResponse{
		Found:    true,
		UserId:   user.ID.String(),
		Username: user.Username,
		Email:    email,
		Message:  "User found",
	}, nil
}
