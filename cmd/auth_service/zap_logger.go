package main

import (
	"fmt"

	"go.uber.org/zap"
)

func newZapLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("не удалось инициализировать zap logger: %w", err)
	}
	return logger, nil
}
