package main

import (
	kitlog "github.com/go-kit/kit/log"
	"go.uber.org/zap"
)

func NewKitLogger(logger *zap.Logger) kitlog.Logger {
	kitLogger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(zap.NewStdLog(logger).Writer()))
	return kitlog.With(kitLogger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)
}
