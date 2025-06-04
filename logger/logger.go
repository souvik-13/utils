package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

// Logger wraps the zap logger with additional functionality
type Logger struct {
	name          string
	zapLogger     *zap.Logger
	sugaredLogger *zap.SugaredLogger
	mw            sync.RWMutex
	level         zap.AtomicLevel
	fields        map[string]interface{}
}

var (
	// Global registry of loggers
	loggers     = make(map[string]*Logger)
	loggerMutex sync.RWMutex

	// Channel for configuration change notifications
	configRefresh = make(chan struct{}, 1)
)

// GetLogger returns a logger for the specified module
func GetLogger(name string) *Logger {
	loggerMutex.RLock()
	logger, exists := loggers[name]
	loggerMutex.RUnlock()

	if exists {
		return logger
	}

	// Create a new logger if it doesn't exist
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	// Double check in case another goroutine created it
	if logger, exists = loggers[name]; exists {
		return logger
	}

	logger = newLogger(name)
	loggers[name] = logger
	return logger
}

// newLogger creates a new logger instance
func newLogger(name string) *Logger {
	config := GetConfig()

	// Determine log level for this module
	levelStr := config.DefaultLevel
	if moduleLevel, exists := config.ModuleLevels[name]; exists {
		levelStr = moduleLevel
	}

	level := getZapLevel(levelStr)

	// Create output writers
	var cores []zapcore.Core
	encoder := getEncoder(config)

	// Console output
	if config.Output.Console {
		consoleCore := zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// File output
	if config.Output.File.Enabled {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.Output.File.Path + "/" + name + ".log",
			MaxSize:    config.Output.File.MaxSize,
			MaxBackups: config.Output.File.MaxBackups,
			MaxAge:     config.Output.File.MaxAge,
			Compress:   config.Output.File.Compress,
		})

		fileCore := zapcore.NewCore(
			encoder,
			fileWriter,
			level,
		)
		cores = append(cores, fileCore)
	}

	// Combine cores
	core := zapcore.NewTee(cores...)

	// Create logger with options
	zapLogger := zap.New(core, createZapOptions(config)...)

	// Name the logger
	zapLogger = zapLogger.Named(name)

	return &Logger{
		name:          name,
		zapLogger:     zapLogger,
		sugaredLogger: zapLogger.Sugar(),
		level:         level,
		fields:        make(map[string]interface{}),
	}
}

// getZapLevel converts string level to zap.AtomicLevel
func getZapLevel(level string) zap.AtomicLevel {
	switch level {
	case DebugLevel:
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case InfoLevel:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case WarnLevel:
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case ErrorLevel:
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case FatalLevel:
		return zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

// triggerConfigRefresh notifies all loggers to refresh their config
func triggerConfigRefresh() {
	select {
	case configRefresh <- struct{}{}:
		// Signal sent
	default:
		// Channel already has signal pending
	}
}

// Reload recreates the logger with updated configuration
func (l *Logger) Reload() {
	l.mw.Lock()
	defer l.mw.Unlock()

	newL := newLogger(l.name)
	l.zapLogger = newL.zapLogger
	l.sugaredLogger = newL.sugaredLogger
	l.level = newL.level

	// Re-add any fields
	for k, v := range l.fields {
		l.sugaredLogger = l.sugaredLogger.With(k, v)
	}
}
