package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// getEncoder returns a zapcore.Encoder based on the config
func getEncoder(config *LoggingConfig) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Enhanced formatting for development mode
	if config.Development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
		encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}

	// Choose encoder based on format
	switch config.Format {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

// createZapOptions converts our config to zap options
func createZapOptions(config *LoggingConfig) []zap.Option {
	options := []zap.Option{}

	// Add caller skip to show the correct file location
	options = append(options, zap.AddCallerSkip(1))

	if config.IncludeCaller {
		options = append(options, zap.AddCaller())
	}

	if config.IncludeStacktrace {
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	}

	if config.Development {
		options = append(options, zap.Development())
	}

	// Configure sampling if enabled
	if config.Sampling.Enabled {
		options = append(options, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSamplerWithOptions(
				core,
				time.Second,
				config.Sampling.Initial,
				config.Sampling.Thereafter,
			)
		}))
	}

	return options
}
