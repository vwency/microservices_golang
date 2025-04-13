package handler_user_service_kafka

import (
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaConfig struct {
	Brokers       []string
	RequestTopic  string
	ResponseTopic string
}

func SetupKafka(cfg KafkaConfig, logger *zap.Logger) (*kafka.Writer, *kafka.Reader) {
	kafkaLogger := kafka.LoggerFunc(func(msg string, args ...interface{}) {
		fields := make([]zap.Field, 0, len(args)/2)
		for i := 0; i < len(args); i += 2 {
			if key, ok := args[i].(string); ok {
				fields = append(fields, zap.Any(key, args[i+1]))
			}
		}
		logger.Debug(msg, fields...)
	})

	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.ResponseTopic,
		Balancer: &kafka.LeastBytes{},
		Logger:   kafkaLogger,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.RequestTopic,
		GroupID:  "user-service-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
		Logger:   kafkaLogger,
	})

	return writer, reader
}
