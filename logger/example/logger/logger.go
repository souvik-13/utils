package logger

import (
	"context"
	"log/slog"
	"sync"

	"github.com/souvik-13/utils/logger/v3"
)

type Config struct {
}

var (
	lgr  *logger.Logger
	once sync.Once
)

func InitLogger(config Config) *logger.Logger {

	once.Do(func() {
		lgr = logger.New(
			logger.WithLevel(slog.LevelDebug),
			logger.WithAddSource(true),
			logger.WithAddStackTraceAt(slog.LevelDebug),
			logger.WithCallerSkipCount(0),
			logger.WithOutfile(""),
		)
	})

	return lgr
}

func Logger(ctx context.Context) *logger.Logger {
	l, err := logger.FromContext(ctx)
	if err == nil {
		return l
	}

	return lgr
}
