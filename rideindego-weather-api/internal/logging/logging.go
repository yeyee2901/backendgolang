package logging

import (
	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewFileLogger(filepath string, serviceName string, level slog.Leveler) *slog.Logger {
	lj := &lumberjack.Logger{
		Filename: filepath,
		Compress: true,
	}

	logger := slog.New(slog.NewJSONHandler(lj, &slog.HandlerOptions{Level: level}))
	logger = logger.With(
		slog.String("service", serviceName),
	)

	return logger
}
