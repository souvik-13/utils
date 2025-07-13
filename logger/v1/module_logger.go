package logger

import (
	"context"

	"go.uber.org/zap"
)

// WithFields returns a new Logger with the given fields added
func (l *Logger) WithFields(fields map[string]any) *Logger {
	l.mw.Lock()
	defer l.mw.Unlock()

	newLogger := &Logger{
		name:          l.name,
		zapLogger:     l.zapLogger,
		sugaredLogger: l.sugaredLogger,
		level:         l.level,
		fields:        make(map[string]any),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
		newLogger.sugaredLogger = newLogger.sugaredLogger.With(k, v)
	}

	return newLogger
}

// WithField returns a new Logger with the given field added
func (l *Logger) WithField(key string, value any) *Logger {
	return l.WithFields(map[string]any{key: value})
}

// WithContext adds context-specific fields to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract relevant information from the context
	// This can be customized based on your context structure
	fields := extractContextFields(ctx)
	if len(fields) == 0 {
		return l
	}
	return l.WithFields(fields)
}

// SetLevel changes the log level for this specific logger
func (l *Logger) SetLevel(level string) {
	l.mw.Lock()
	defer l.mw.Unlock()

	newLevel := getZapLevel(level)
	l.level.SetLevel(newLevel.Level())
}

// Logging methods

// Debug logs a message at debug level
func (l *Logger) Debug(args ...any) {
	l.sugaredLogger.Debug(args...)
}

// Debugf logs a formatted message at debug level
func (l *Logger) Debugf(format string, args ...any) {
	l.sugaredLogger.Debugf(format, args...)
}

// Info logs a message at info level
func (l *Logger) Info(args ...any) {
	l.sugaredLogger.Info(args...)
}

// Infof logs a formatted message at info level
func (l *Logger) Infof(format string, args ...any) {
	l.sugaredLogger.Infof(format, args...)
}

// Warn logs a message at warning level
func (l *Logger) Warn(args ...any) {
	l.sugaredLogger.Warn(args...)
}

// Warnf logs a formatted message at warning level
func (l *Logger) Warnf(format string, args ...any) {
	l.sugaredLogger.Warnf(format, args...)
}

// Error logs a message at error level
func (l *Logger) Error(args ...any) {
	l.sugaredLogger.Error(args...)
}

// Errorf logs a formatted message at error level
func (l *Logger) Errorf(format string, args ...any) {
	l.sugaredLogger.Errorf(format, args...)
}

// Fatal logs a message at fatal level then calls os.Exit(1)
func (l *Logger) Fatal(args ...any) {
	l.sugaredLogger.Fatal(args...)
}

// Fatalf logs a formatted message at fatal level then calls os.Exit(1)
func (l *Logger) Fatalf(format string, args ...any) {
	l.sugaredLogger.Fatalf(format, args...)
}

// Structured logging methods with explicit fields

// DebugWithFields logs a debug message with additional fields
func (l *Logger) DebugWithFields(msg string, fields map[string]any) {
	if ce := l.zapLogger.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(fieldsToZapFields(fields)...)
	}
}

// InfoWithFields logs an info message with additional fields
func (l *Logger) InfoWithFields(msg string, fields map[string]any) {
	if ce := l.zapLogger.Check(zap.InfoLevel, msg); ce != nil {
		ce.Write(fieldsToZapFields(fields)...)
	}
}

// WarnWithFields logs a warning message with additional fields
func (l *Logger) WarnWithFields(msg string, fields map[string]any) {
	if ce := l.zapLogger.Check(zap.WarnLevel, msg); ce != nil {
		ce.Write(fieldsToZapFields(fields)...)
	}
}

// ErrorWithFields logs an error message with additional fields
func (l *Logger) ErrorWithFields(msg string, fields map[string]any) {
	if ce := l.zapLogger.Check(zap.ErrorLevel, msg); ce != nil {
		ce.Write(fieldsToZapFields(fields)...)
	}
}

// Helper functions

// extractContextFields extracts relevant fields from a context
func extractContextFields(ctx context.Context) map[string]any {
	fields := make(map[string]any)

	// Standard fields
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		fields["trace_id"] = traceID
	}
	if requestID, ok := ctx.Value("request_id").(string); ok {
		fields["request_id"] = requestID
	}
	if userID, ok := ctx.Value("user_id").(string); ok {
		fields["user_id"] = userID
	}

	// Additional fields that might be useful
	if sessionID, ok := ctx.Value("session_id").(string); ok {
		fields["session_id"] = sessionID
	}
	if clientIP, ok := ctx.Value("client_ip").(string); ok {
		fields["client_ip"] = clientIP
	}
	if operation, ok := ctx.Value("operation").(string); ok {
		fields["operation"] = operation
	}

	return fields
}

// fieldsToZapFields converts a map to zap.Field slice
func fieldsToZapFields(fields map[string]any) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}
