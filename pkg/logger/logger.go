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
	cfg.Level = zap.NewAtomicLevelAt(level) // Устанавливаем уровень логирования

	if logLevel == "dev" {
		cfg = zap.NewDevelopmentConfig() // Читаемый формат для Dev
	}

	zapLogger, err = cfg.Build()
	if err != nil {
		log.Fatalf("Cannot initialize zap logger: %v", err)
	}

	Log = zapLogger.Sugar()
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}
