package zaphandler

import (
	"context"
	"log/slog"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ slog.Handler = (*Handler)(nil)

type Handler struct {
	core zapcore.Core
	name string
	addCaller bool
	addStackAt slog.Level
	callerSkipCount int
	replaceAttr func(groups []string, a slog.Attr) slog.Attr
	groups []string
}

func New(core zapcore.Core, opts ...HandlerOption) *Handler {
	h := &Handler{
		core:       core,
		addStackAt: slog.LevelError, // Default to Error level for stack traces
	}

	for _, opt := range opts {
		opt.apply(h)
	}

	return h
}

type groupObject []slog.Attr

func (g groupObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, attr := range g {
		convertAttrToField(attr).AddTo(enc)
	}
	return nil
}

func convertAttrToField(attr slog.Attr) zapcore.Field {
	if attr.Equal(slog.Attr{}) {
		return zap.Skip()
	}

	switch attr.Value.Kind() {
	case slog.KindBool:
		return zap.Bool(attr.Key, attr.Value.Bool())
	case slog.KindInt64:
		return zap.Int64(attr.Key, attr.Value.Int64())
	case slog.KindUint64:
		return zap.Uint64(attr.Key, attr.Value.Uint64())
	case slog.KindFloat64:
		return zap.Float64(attr.Key, attr.Value.Float64())
	case slog.KindString:
		return zap.String(attr.Key, attr.Value.String())
	case slog.KindDuration:
		return zap.Duration(attr.Key, attr.Value.Duration())
	case slog.KindTime:
		return zap.Time(attr.Key, attr.Value.Time())
	case slog.KindAny:
		return zap.Any(attr.Key, attr.Value.Any())
	case slog.KindGroup:
		if attr.Key != "" {
			return zap.Inline(groupObject(attr.Value.Group()))
		}
		return zap.Object(attr.Key, groupObject(attr.Value.Group()))
	case slog.KindLogValuer:
		return convertAttrToField(slog.Attr{
			Key:   attr.Key,
			Value: attr.Value.Resolve(),
		})
	default:
		return zap.Any(attr.Key, attr.Value.Any())
	}
}

var SlogToZapLevel = map[slog.Level]zapcore.Level{
	slog.LevelDebug: zapcore.DebugLevel,
	slog.LevelInfo:  zapcore.InfoLevel,
	slog.LevelWarn:  zapcore.WarnLevel,
	slog.LevelError: zapcore.ErrorLevel,
}

func (h *Handler) appendGroups(fields []zapcore.Field) []zapcore.Field {
	for _, group := range h.groups {
		fields = append(fields, zap.Namespace(group))
	}
	return fields
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	zapLevel := SlogToZapLevel[level]
	return h.core.Enabled(zapLevel)
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	ent := zapcore.Entry{
		Level:      SlogToZapLevel[record.Level],
		Time:       record.Time,
		Message:    record.Message,
		LoggerName: h.name,
	}
	ce := h.core.Check(ent, nil)
	if ce == nil {
		return nil
	}

	if h.addCaller && record.PC != 0 {
		frame, _ := runtime.CallersFrames([]uintptr{record.PC}).Next()
		if frame.PC != 0 {
			ce.Caller = zapcore.EntryCaller{
				Defined:  true,
				PC:       frame.PC,
				File:     frame.File,
				Line:     frame.Line,
				Function: frame.Function,
			}
		}
	}

	fields := make([]zapcore.Field, 0, record.NumAttrs()+len(h.groups))

	var addedNamespace bool
	record.Attrs(func(attr slog.Attr) bool {
		if h.replaceAttr != nil {
			attr = h.replaceAttr(h.groups, attr)
		}

		attr.Value = attr.Value.Resolve()

		f := convertAttrToField(attr)
		if !addedNamespace && len(h.groups) > 0 && f != zap.Skip() {
			fields = h.appendGroups(fields)
			addedNamespace = true
		}
		fields = append(fields, f)
		return true
	})

	ce.Write(fields...)
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := make([]zapcore.Field, 0, len(attrs)+len(h.groups))
	var addedNamespace bool
	for _, attr := range attrs {
		if h.replaceAttr != nil {
			attr = h.replaceAttr(h.groups, attr)
		}

		attr.Value = attr.Value.Resolve()

		f := convertAttrToField(attr)
		if !addedNamespace && len(h.groups) > 0 && f != zap.Skip() {
			fields = h.appendGroups(fields)
			addedNamespace = true
		}
		fields = append(fields, f)
	}

	cloned := *h
	cloned.core = cloned.core.With(fields)
	if addedNamespace {
		cloned.groups = nil // Reset groups if namespace was added
	}
	return &cloned
}

func (h *Handler) WithGroup(group string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = group
	cloned := *h
	cloned.groups = newGroups
	return &cloned
}
