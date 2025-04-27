package handler_user_service_grpc

import (
	"context"

	"github.com/vwency/microservices_golang/internal/user_database/usecase/user_usecase"
	pb "github.com/vwency/microservices_golang/proto/user_database"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserHandler struct {
	uc     *user_usecase.UserUsecase
	logger *zap.Logger
}

func NewGetUserHandler(uc *user_usecase.UserUsecase, logger *zap.Logger) *GetUserHandler {
	return &GetUserHandler{
		uc:     uc,
		logger: logger.With(zap.String("handler", "get_user")),
	}
}

func (h *GetUserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if (req.UserId == nil || *req.UserId == "") &&
		(req.Username == nil || *req.Username == "") &&
		(req.Email == nil || *req.Email == "") {
		h.logger.Warn("empty request parameters")
		return nil, status.Error(codes.InvalidArgument, "user_id, username, or email must be provided")
	}

	params := user_usecase.UserParams{
		UserID:   req.GetUserId(),
		Username: req.GetUsername(), // req.GetUsername() вернет string
		Email:    req.GetEmail(),
	}

	user, err := h.uc.GetUser(params)
	if err != nil {
		userID := ""
		username := ""
		email := ""

		if req.UserId != nil {
			userID = *req.UserId
		}
		if req.Username != nil {
			username = *req.Username
		}
		if req.Email != nil {
			email = *req.Email
		}

		if status.Code(err) == codes.NotFound {
			h.logger.Info("user not found",
				zap.String("user_id", userID),
				zap.String("username", username),
				zap.String("email", email))
			return &pb.GetUserResponse{
				Found:   false,
				Message: "User not found",
			}, nil
		}

		h.logger.Error("failed to get user",
			zap.String("user_id", userID),
			zap.String("username", username),
			zap.String("email", email),
			zap.Error(err))
		return nil, err
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	userID := ""
	if user.ID != [16]byte{} {
		userID = user.ID.String()
	}

	return &pb.GetUserResponse{
		Found:              true,
		UserId:             userID,
		Username:           user.Username,
		Email:              email,
		HashedPassword:     user.HashedPassword,
		HashedRefreshToken: user.HashedRefreshToken,
		HashedAccessToken:  user.HashedAccessToken,
		Message:            "User found",
	}, nil
}
