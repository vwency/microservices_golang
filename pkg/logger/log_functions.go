package logger

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
