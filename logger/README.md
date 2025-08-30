# Logger

A flexible, feature-rich structured logging package for Go applications built on top of Go's standard library `log/slog` with Zap backend integration.

## Features

- **Structured Logging**: JSON-formatted logs with rich metadata
- **Context-Aware**: Seamless context propagation with logger
- **Configurable Levels**: Debug, Info, Warn, Error support
- **PII Protection**: Built-in PII field masking functionality
- **Stack Traces**: Configurable stack trace inclusion
- **Caller Information**: Source file and line number reporting
- **File Output**: Support for both stdout and file outputs
- **Zap Integration**: High-performance logging using Zap backend
- **Extensible**: Easily customizable through options pattern

## Installation

```bash
go get github.com/souvik-13/utils/logger/v3
```

## Quick Start

```go
package main

import (
    "context"
    "log/slog"

    "github.com/souvik-13/utils/logger/v3"
)

func main() {
    // Create a new logger
    log := logger.New(
        logger.WithLevel(slog.LevelDebug),
        logger.WithAddSource(true),
    )

    // Basic logging
    log.Info("Application started")

    // Structured logging with key-value pairs
    log.Info("User logged in", "userID", "123", "role", "admin")

    // Error logging with an error
    log.Error("Failed to connect to database", "error", err)

    // Context-aware logging
    ctx := context.Background()
    ctx = logger.ToContext(ctx, log)

    // Later in your application:
    logFromCtx, err := logger.FromContext(ctx)
    if err == nil {
        logFromCtx.Info("Retrieved logger from context")
    }
}
```

## Configuration Options

The logger can be customized with various options:

```go
logger.New(
    logger.WithKind(logger.KindZap),            // Logger implementation (currently only Zap supported)
    logger.WithLevel(slog.LevelDebug),          // Minimum log level (Debug, Info, Warn, Error)
    logger.WithAddSource(true),                 // Include source file and line information
    logger.WithAddStackTraceAt(slog.LevelError), // Include stack traces for logs at or above this level
    logger.WithCallerSkipCount(1),              // Number of stack frames to skip for caller info
    logger.WithSecureFields("password", "token"), // Additional fields to mask for PII protection
    logger.WithOutfile("app.log"),              // Output to file instead of stdout
)
```

## Secure Logging

The logger automatically masks PII fields to prevent sensitive information from appearing in logs:

```go
// Default PII fields masked include "email"
log.Info("User registered", "email", "user@example.com") // email will be masked

// Add additional fields to mask
logger.New(
    logger.WithSecureFields("password", "ssn", "credit_card"),
)
```

## Output to File

Direct logs to a file instead of stdout:

```go
// Set at initialization
log := logger.New(
    logger.WithOutfile("/var/log/app.log"),
)

// Or redirect an existing logger
fileLogger := log.WithOutputFile("/var/log/app.log")
fileLogger.Info("This goes to the file")
```

## Context Integration

```go
// Store logger in context
ctx := logger.ToContext(context.Background(), log)

// Retrieve logger from context
loggerFromCtx, err := logger.FromContext(ctx)
if err != nil {
    // Handle missing logger
}

// Common pattern for libraries/middlewares
func ProcessRequest(ctx context.Context, req Request) {
    log, err := logger.FromContext(ctx)
    if err != nil {
        // Use default logger or handle error
        log = defaultLogger
    }

    log.Info("Processing request", "requestID", req.ID)
}
```

## Creating a Logger Wrapper for Your Application

For larger applications, you might want to create a simplified logger wrapper:

```go
package myapp

import (
    "context"
    "log/slog"
    "sync"

    "github.com/souvik-13/utils/logger/v3"
)

var (
    defaultLogger *logger.Logger
    once sync.Once
)

// InitLogger initializes the global logger
func InitLogger() *logger.Logger {
    once.Do(func() {
        defaultLogger = logger.New(
            logger.WithLevel(slog.LevelDebug),
            logger.WithAddSource(true),
            logger.WithAddStackTraceAt(slog.LevelError),
        )
    })
    return defaultLogger
}

// Log retrieves the logger from context or falls back to default
func Log(ctx context.Context) *logger.Logger {
    l, err := logger.FromContext(ctx)
    if err == nil {
        return l
    }
    return defaultLogger
}
```

## Level-Specific Logging Methods

```go
log.Debug("Debug message")
log.Info("Info message")
log.Warn("Warning message")
log.Error("Error message")

// With structured data
log.Info("User logged in",
    "userID", "123",
    "ipAddress", "192.168.1.1",
    "loginTime", time.Now(),
)

// With context
log.InfoContext(ctx, "Request processed",
    "requestID", req.ID,
    "duration", duration,
)
```

## Advanced Usage

### With Additional Context

```go
// Create a derived logger with additional fields
userLogger := log.With("userID", "123", "sessionID", "abc")
userLogger.Info("User performed action", "action", "login")
```

### Clone Logger

```go
// Create a clone of the logger
clonedLogger := log.Clone()
```

## Dependencies

This logger package relies on:

- Go 1.24.1 or higher
- [github.com/m-mizutani/masq](https://github.com/m-mizutani/masq) - For masking PII fields
- [go.uber.org/zap](https://github.com/uber-go/zap) - For high-performance logging

## License

This project is licensed under the terms of the license included in the repository.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
