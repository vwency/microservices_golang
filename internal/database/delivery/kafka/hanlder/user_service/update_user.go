package handler_user_service_kafka

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase"
	"go.uber.org/zap"
)

type UpdateUserHandler struct {
	uc     *user_usecase.UserUsecase
	logger *zap.Logger
	writer *kafka.Writer
}

func NewUpdateUserHandler(uc *user_usecase.UserUsecase, logger *zap.Logger, writer *kafka.Writer) *UpdateUserHandler {
	return &UpdateUserHandler{
		uc:     uc,
		logger: logger.With(zap.String("handler", "update_user")),
		writer: writer,
	}
}

func (h *UpdateUserHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var req struct {
		Username string `json:"username"`
		HashedRt string `json:"hashed_rt"`
		AccessRt string `json:"access_rt"`
	}

	if err := json.Unmarshal(msg.Value, &req); err != nil {
		h.logger.Error("failed to unmarshal message", zap.Error(err))
		return err
	}

	if req.Username == "" {
		h.logger.Warn("empty username in request")
		return errors.New("username is required")
	}

	if req.HashedRt == "" || req.AccessRt == "" {
		h.logger.Warn("empty token fields in request",
			zap.String("username", req.Username))
		return errors.New("token fields cannot be empty")
	}

	updateParams := user_usecase.UpdateTokensParams{
		Username: req.Username,
		HashedRt: req.HashedRt,
		HashedAt: req.AccessRt,
	}

	err := h.uc.UpdateTokens(updateParams)
	if err != nil {
		h.logger.Error("failed to update user tokens",
			zap.String("username", req.Username),
			zap.Error(err))

		if errors.Is(err, user_usecase.ErrUserNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	h.logger.Info("tokens updated successfully",
		zap.String("username", req.Username))

	response := map[string]interface{}{
		"success": true,
		"message": "Tokens updated successfully",
	}
	responseBytes, _ := json.Marshal(response)
	return h.writer.WriteMessages(ctx, kafka.Message{
		Key:   msg.Key,
		Value: responseBytes,
	})
}
