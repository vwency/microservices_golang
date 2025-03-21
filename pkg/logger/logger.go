package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log — глобальный логгер
var Log *zap.SugaredLogger

// Init инициализирует логгер в зависимости от переданного logLevel
func Init(logLevel string) {
	var zapLogger *zap.Logger
	var err error

	// Устанавливаем уровень логирования
	var level zapcore.Level
	err = level.Set(logLevel)
	if err != nil {
		log.Fatalf("Invalid log level: %s", logLevel)
	}

	// Настройка логгера в зависимости от logLevel
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level) // Устанавливаем уровень логирования

	// Различие между Production и Development логированием
	if logLevel == "dev" {
		cfg = zap.NewDevelopmentConfig() // Читаемый формат для Dev
	}

	// Создание логгера
	zapLogger, err = cfg.Build()
	if err != nil {
		log.Fatalf("Cannot initialize zap logger: %v", err)
	}

	Log = zapLogger.Sugar()
}

// Info логирует сообщение уровня INFO
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Error логирует сообщение уровня ERROR
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Debug логирует сообщение уровня DEBUG
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

// Fatal логирует сообщение уровня FATAL и завершает процесс
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}
