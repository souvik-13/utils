package logger

import (
	"log/slog"

	"github.com/souvik-13/utils/logger/v3/zaphandler"
	"go.uber.org/zap/zapcore"
)

type Kind string

const (
	KindZap Kind = "ZapHandler"
)

func (k Kind) String() string {
	return string(k)
}

func newZapSlogHandler(cfg *Config, masqr func(groups []string, a slog.Attr) slog.Attr) slog.Handler {
	level := zaphandler.SlogToZapLevel[cfg.Level]
	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "time",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(cfg.output),
		level,
	)

	return zaphandler.New(
		zapCore,
		zaphandler.WithName(cfg.Kind.String()),
		zaphandler.WithAddCaller(cfg.AddSource),
		zaphandler.WithAddStackTraceAt(cfg.AddStackTraceAt),
		zaphandler.WithCallerSkipCount(cfg.CallerSkipCount),
		zaphandler.WithReplaceAttr(masqr),
	)
}
