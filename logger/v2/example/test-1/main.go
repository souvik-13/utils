package main

import (
	"context"

	"github.com/souvik-13/my-utils/logger/v2/example/logger"
)

func main() {
	lgr := logger.InitLogger(logger.Config{})

	lgr.Info("Logger initialized successfully")

	ctx := context.Background()
	logger.Logger(ctx).Info("This is an info message")
	logger.Logger(ctx).Info("This is an info message with a key-value pair")
	logger.Logger(ctx).DebugContext(ctx, "This is a debug message", "key", "value")
}
