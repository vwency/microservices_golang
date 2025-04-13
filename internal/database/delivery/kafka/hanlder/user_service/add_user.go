package handler_user_service_kafka

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase"
	"go.uber.org/zap"
)

type AddUserHandler struct {
	uc     *user_usecase.UserUsecase
	logger *zap.Logger
	writer *kafka.Writer
}

func NewAddUserHandler(uc *user_usecase.UserUsecase, logger *zap.Logger, writer *kafka.Writer) *AddUserHandler {
	return &AddUserHandler{
		uc:     uc,
		logger: logger.With(zap.String("handler", "add_user")),
		writer: writer,
	}
}

func (h *AddUserHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var req struct {
		Username       string `json:"username"`
		HashedPassword string `json:"hashed_password"`
		Email          string `json:"email"`
		HashedRt       string `json:"hashed_rt"`
	}

	if err := json.Unmarshal(msg.Value, &req); err != nil {
		h.logger.Error("failed to unmarshal message", zap.Error(err))
		return err
	}

	if req.Username == "" || req.HashedPassword == "" || req.Email == "" {
		h.logger.Warn("missing required fields",
			zap.String("username", req.Username),
			zap.String("email", req.Email))
		return errors.New("username, hashed password, and email are required")
	}

	params := user_usecase.CreateUserParams{
		UserParams: user_usecase.UserParams{
			Username: req.Username,
			Email:    req.Email,
		},
		HashedPassword: req.HashedPassword,
		HashedRt:       req.HashedRt,
		HashedAt:       time.Now().Format(time.RFC3339),
	}

	if err := h.uc.CreateUser(params); err != nil {
		h.logger.Error("failed to create user",
			zap.String("username", req.Username),
			zap.Error(err))

		if errors.Is(err, user_usecase.ErrUserAlreadyExists) {
			return errors.New("user already exists")
		}

		return err
	}

	h.logger.Info("user created successfully",
		zap.String("username", req.Username))

	// Optionally produce a success message
	response := map[string]interface{}{
		"success": true,
		"message": "User created successfully",
	}
	responseBytes, _ := json.Marshal(response)
	return h.writer.WriteMessages(ctx, kafka.Message{
		Value: responseBytes,
	})
}
