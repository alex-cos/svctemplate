package logx

import (
	"log/slog"
	"os"
)

var globalLogger *slog.Logger = slog.New(
	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelInfo,
		ReplaceAttr: nil,
	}))

func Set(l *slog.Logger) {
	globalLogger = l
}

func L() *slog.Logger {
	return globalLogger
}
