package main

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.uber.org/zap"
)

func newKitLogger(zapLogger *zap.Logger) log.Logger {
	kitLogger := log.NewLogfmtLogger(os.Stdout)
	kitLogger = level.NewFilter(kitLogger, level.AllowDebug())
	kitLogger = log.With(kitLogger, "ts", log.DefaultTimestampUTC)
	return kitLogger
}
