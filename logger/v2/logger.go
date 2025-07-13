package logger

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/m-mizutani/masq"
	"github.com/souvik-13/utils/logger/v2/constants"
)

type contextKey struct{}

var LoggerCtxKey contextKey = contextKey{}

type Logger struct {
	*slog.Logger

	cfg *Config
}

type Option func(*Config)

type Config struct {
	Kind            Kind
	Level           slog.Level
	AddSource       bool
	AddStackTraceAt slog.Level
	CallerSkipCount int
	MaskOptions     []masq.Option
}

func New(opts ...Option) *Logger {
	cfg := &Config{
		Kind:            KindZap,
		Level:           slog.LevelInfo,
		AddSource:       true,
		AddStackTraceAt: slog.LevelError,
		CallerSkipCount: 1,
		MaskOptions:     make([]masq.Option, 0, len(constants.PIIFieldKeys)),
	}

	for _, key := range constants.PIIFieldKeys {
		cfg.MaskOptions = append(cfg.MaskOptions, masq.WithFieldPrefix(key))
	}

	// cfg.MaskOptions = append(cfg.MaskOptions, masq.WithCustomTagKey(), masq.WithTag(), masq.WithRedactMessage())

	for _, opt := range opts {
		opt(cfg)
	}

	masqr := masq.New(
		cfg.MaskOptions...,
	)

	var handler slog.Handler

	switch cfg.Kind {
	default:
		handler = newZapSlogHandler(cfg, masqr)
	}

	slogger := slog.New(handler)

	return &Logger{
		slogger,
		cfg,
	}
}

func WithKind(kind Kind) Option {
	return func(cfg *Config) {
		cfg.Kind = kind
	}
}

func WithLevel(level slog.Level) Option {
	return func(cfg *Config) {
		cfg.Level = level
	}
}

func WithAddSource(addSource bool) Option {
	return func(cfg *Config) {
		cfg.AddSource = addSource
	}
}

func WithAddStackTraceAt(level slog.Level) Option {
	return func(cfg *Config) {
		cfg.AddStackTraceAt = level
	}
}

func WithCallerSkipCount(skipCount int) Option {
	return func(cfg *Config) {
		cfg.CallerSkipCount = skipCount
	}
}

func WithSecureFields(fields ...string) Option {
	return func(cfg *Config) {
		for _, field := range fields {
			cfg.MaskOptions = append(cfg.MaskOptions, masq.WithFieldPrefix(field))
		}
	}
}

func WithMaskOptions(opts ...masq.Option) Option {
	return func(cfg *Config) {
		cfg.MaskOptions = append(cfg.MaskOptions, opts...)
	}
}

func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, LoggerCtxKey, logger)
}

func FromContext(ctx context.Context) (*Logger, error) {
	logger, ok := ctx.Value(LoggerCtxKey).(*Logger)
	if !ok {
		return nil, fmt.Errorf("logger not found in context")
	}
	return logger, nil
}

// func (l *Logger) Error(msg string, err error) {
// 	l.Logger.Error(msg, slog.Any("error", err))
// }

func (l *Logger) With(args ...any) *Logger {
	if len(args) == 0 {
		return l
	}

	lc := l.Clone()
	lc.Logger = l.Logger.With(args...)

	return lc
}

func (l *Logger) Clone() *Logger {
	clone := *l
	return &clone
}
