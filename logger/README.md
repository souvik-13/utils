# Logger Package

A flexible, high-performance logging package for Go, built on top of [zap](https://github.com/uber-go/zap), with enhancements for modular, structured, and configurable logging.

## Features

- **Module-based loggers**: Create and manage loggers for different modules or components using `GetLogger(name)`.
- **Multiple log levels**: Supports `debug`, `info`, `warn`, `error`, and `fatal` levels.
- **Structured logging**: Add fields and context to logs for better traceability.
- **Flexible output**: Log to console and/or files, with file rotation (via [lumberjack](https://github.com/natefinch/lumberjack)).
- **Dynamic configuration**: Change log levels and output settings at runtime.
- **Context integration**: Extract and log fields from Go `context.Context`.
- **Writer integration**: Implements `io.Writer` for redirecting output from other libraries.
- **Sampling and development mode**: Control verbosity and log format for production or development.

## Design Overview

### Logger Structure

- Each logger is an instance of the `Logger` struct, wrapping a zap logger and providing thread-safe operations.
- Loggers are registered and retrieved by name, allowing per-module configuration.
- Loggers can be enriched with fields and context for structured logging.

### Configuration

- Centralized via the `LoggingConfig` struct (see `config.go`).
- Supports default and per-module log levels, output destinations, log format (JSON/console), caller/stacktrace inclusion, and sampling.
- File output is managed with rotation, backup, and compression options.
- Configuration can be updated at runtime, and all loggers will refresh accordingly.

### Output and Formatting

- Console and file outputs are supported, with customizable format and colorization in development mode.
- Log format and encoder are selected based on configuration (see `formats.go`).

### Writer Support

- The package provides `LogWriter` and `MultiLogWriter` to redirect output from other sources (e.g., standard library) into the logger.

## Usage

### Basic Logging

```go
import "github.com/souvik-13/my-utils/logger"

log := logger.GetLogger("myapp")
log.Info("Application started")
log.Debugf("User %s logged in", username)
```

### Structured Logging

```go
log.WithFields(map[string]any{
    "user_id": "123",
    "action": "login",
}).Info("User logged in")

log.WithField("request_id", reqID).Warn("Slow request")
```

### Contextual Logging

```go
log.WithContext(ctx).Info("Request received")
```

### Configuration Example

```go
cfg := &logger.LoggingConfig{
    DefaultLevel: logger.InfoLevel,
    ModuleLevels: map[string]string{"myapp": logger.DebugLevel},
    Format:       "json",
    Development:  true,
    Output: logger.OutputConfig{
        Console: true,
        File: logger.FileOutputConfig{
            Enabled:    true,
            Path:       "./logs",
            MaxSize:    100, // MB
            MaxBackups: 5,
            MaxAge:     30, // days
            Compress:   true,
        },
    },
    IncludeCaller:     true,
    IncludeStacktrace: true,
    Sampling: logger.SamplingConfig{
        Enabled:    false,
        Initial:    100,
        Thereafter: 100,
    },
}
logger.Initialize(cfg)
```

### Changing Log Level at Runtime

```go
logger.SetLogLevel("myapp", logger.DebugLevel)
```

### Using as an io.Writer

```go
w := logger.NewLogWriter(log, logger.InfoLevel)
fmt.Fprintln(w, "This will be logged at info level")
```

## File Structure

- `logger.go`: Core logger logic and registry
- `config.go`: Configuration structures and management
- `formats.go`: Log formatting and encoder selection
- `module_logger.go`: Module-specific and structured logging methods
- `writer.go`: Writer implementations for integration with io.Writer

## Dependencies

- [zap](https://github.com/uber-go/zap) for fast, structured logging
- [lumberjack](https://github.com/natefinch/lumberjack) for file rotation

## License

[MIT License](LICENSE)
