package main

import (
	"context"

	"github.com/souvik-13/utils/logger/v3/example/logger"
)

var lgr = logger.InitLogger(logger.Config{})

func main() {

	lgr.Info("Logger initialized successfully")

	ctx := context.Background()
	logger.Logger(ctx).Info("This is an info message")
	logger.Logger(ctx).Info("This is an info message with a key-value pair")
	logger.Logger(ctx).DebugContext(ctx, "This is a debug message", "key", "value")

	foo()
}

func foo() {
	log := lgr.WithOutputFile("test.log")
	log.Info("This is an info message from foo")
	log.Debug("This is a debug message from foo")
}

func bar() {

}
