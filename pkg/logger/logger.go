package logger

import (
	"log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func Init(logLevel string) {
	var zapLogger *zap.Logger
	var err error

	var level zapcore.Level
	err = level.Set(logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %s", logLevel)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)

	if logLevel == "dev" {
		cfg = zap.NewDevelopmentConfig()
	}

	zapLogger, err = cfg.Build()
	if err != nil {
		log.Fatalf("Cannot initialize zap logger: %v", err)
	}

	Log = zapLogger.Sugar()
}
