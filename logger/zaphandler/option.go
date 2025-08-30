package zaphandler

import "log/slog"

type HandlerOption interface {
	apply(*Handler)
}

type handlerOptionFunc func(*Handler)

func (f handlerOptionFunc) apply(h *Handler) {
	f(h)
}

// WithName sets the name of the handler.
func WithName(name string) HandlerOption {
	return handlerOptionFunc(func(h *Handler) {
		h.name = name
	})
}

// WithAddCaller enables or disables the addition of caller information.
func WithAddCaller(add bool) HandlerOption {
	return handlerOptionFunc(func(h *Handler) {
		h.addCaller = add
	})
}

// WithAddStackTraceAt sets the level at which stack traces are added.
func WithAddStackTraceAt(level slog.Level) HandlerOption {
	return handlerOptionFunc(func(h *Handler) {
		h.addStackAt = level
	})
}

// WithCallerSkipCount sets the number of stack frames to skip when adding caller information.
func WithCallerSkipCount(skip int) HandlerOption {
	return handlerOptionFunc(func(h *Handler) {
		h.callerSkipCount = skip
	})
}

// WithReplaceAttr sets a function to replace attributes before logging.
func WithReplaceAttr(replace func(groups []string, a slog.Attr) slog.Attr) HandlerOption {
	return handlerOptionFunc(func(h *Handler) {
		h.replaceAttr = replace
	})
}
