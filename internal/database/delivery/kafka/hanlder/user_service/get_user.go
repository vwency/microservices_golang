package handler_user_service_kafka

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase"
	"go.uber.org/zap"
)

type GetUserHandler struct {
	uc     *user_usecase.UserUsecase
	logger *zap.Logger
	writer *kafka.Writer
}

func NewGetUserHandler(uc *user_usecase.UserUsecase, logger *zap.Logger, writer *kafka.Writer) *GetUserHandler {
	return &GetUserHandler{
		uc:     uc,
		logger: logger.With(zap.String("handler", "get_user")),
		writer: writer,
	}
}

func (h *GetUserHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := json.Unmarshal(msg.Value, &req); err != nil {
		h.logger.Error("failed to unmarshal message", zap.Error(err))
		return err
	}

	if req.Username == "" && req.Email == "" {
		h.logger.Warn("empty request parameters")
		return errors.New("username or email must be provided")
	}

	params := user_usecase.UserParams{
		Username: req.Username,
		Email:    req.Email,
	}

	user, err := h.uc.GetUser(params)
	if err != nil {
		h.logger.Error("failed to get user",
			zap.String("username", req.Username),
			zap.String("email", req.Email),
			zap.Error(err))
		return err
	}

	response := make(map[string]interface{})
	if user == nil {
		h.logger.Info("user not found",
			zap.String("username", req.Username),
			zap.String("email", req.Email))
		response["found"] = false
		response["message"] = "User not found"
	} else {
		email := ""
		if user.Email != nil {
			email = *user.Email
		}

		response["found"] = true
		response["username"] = user.Username
		response["email"] = email
		response["hashed_rt"] = user.HashedRefreshToken
		response["hashed_password"] = user.HashedPassword
		response["hashed_at"] = user.HashedAccessToken
		response["message"] = "User found"
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return h.writer.WriteMessages(ctx, kafka.Message{
		Key:   msg.Key,
		Value: responseBytes,
	})
}
