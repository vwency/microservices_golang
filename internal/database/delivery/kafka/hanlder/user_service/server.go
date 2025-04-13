package handler_user_service_kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/vwency/microservices_golang/internal/database/usecase/user_usecase"
	"go.uber.org/zap"
)

type Server struct {
	addUserHandler    *AddUserHandler
	getUserHandler    *GetUserHandler
	updateUserHandler *UpdateUserHandler
	logger            *zap.Logger
	writer            *kafka.Writer
}

func NewServer(uc *user_usecase.UserUsecase, logger *zap.Logger, writer *kafka.Writer) *Server {
	return &Server{
		addUserHandler:    NewAddUserHandler(uc, logger, writer),
		getUserHandler:    NewGetUserHandler(uc, logger, writer),
		updateUserHandler: NewUpdateUserHandler(uc, logger, writer),
		logger:            logger,
		writer:            writer,
	}
}

func (s *Server) HandleMessages(ctx context.Context, reader *kafka.Reader) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("stopping kafka message handler")
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				s.logger.Error("failed to read message", zap.Error(err))
				continue
			}

			var handlerFunc func(context.Context, kafka.Message) error

			switch string(msg.Key) {
			case "add_user":
				handlerFunc = s.addUserHandler.Handle
			case "get_user":
				handlerFunc = s.getUserHandler.Handle
			case "update_user":
				handlerFunc = s.updateUserHandler.Handle
			default:
				s.logger.Warn("unknown message key",
					zap.String("key", string(msg.Key)),
					zap.ByteString("value", msg.Value))
				continue
			}

			if err := handlerFunc(ctx, msg); err != nil {
				s.logger.Error("failed to handle message",
					zap.String("key", string(msg.Key)),
					zap.Error(err))

				errorResponse := map[string]interface{}{
					"success": false,
					"error":   err.Error(),
				}
				errorBytes, _ := json.Marshal(errorResponse)
				s.writer.WriteMessages(ctx, kafka.Message{
					Key:   msg.Key,
					Value: errorBytes,
				})
			}
		}
	}
}
