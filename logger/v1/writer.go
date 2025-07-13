package logger

import (
	"io"
	"sync"
)

// LogWriter implements io.Writer to redirect output to the logger
type LogWriter struct {
	logger *Logger
	level  string
	buffer []byte
	mutex  sync.Mutex
}

// NewLogWriter creates a new writer that sends output to the given logger at the specified level
func NewLogWriter(logger *Logger, level string) io.Writer {
	return &LogWriter{
		logger: logger,
		level:  level,
		buffer: make([]byte, 0, 1024),
	}
}

// Write implements the io.Writer interface
func (w *LogWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Append to buffer
	w.buffer = append(w.buffer, p...)

	// Process complete lines
	processed := 0
	for i, b := range w.buffer {
		if b == '\n' {
			line := string(w.buffer[:i])
			w.logLine(line)
			processed = i + 1
		}
	}

	// Keep unprocessed data
	if processed > 0 {
		w.buffer = w.buffer[processed:]
	}

	return len(p), nil
}

// Flush logs any buffered content
func (w *LogWriter) Flush() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if len(w.buffer) > 0 {
		w.logLine(string(w.buffer))
		w.buffer = w.buffer[:0]
	}
}

// logLine logs a single line at the configured level
func (w *LogWriter) logLine(line string) {
	switch w.level {
	case DebugLevel:
		w.logger.Debug(line)
	case InfoLevel:
		w.logger.Info(line)
	case WarnLevel:
		w.logger.Warn(line)
	case ErrorLevel:
		w.logger.Error(line)
	case FatalLevel:
		w.logger.Error(line) // Using Error instead of Fatal to avoid process termination
	default:
		w.logger.Info(line)
	}
}

// MultiLogWriter sends output to multiple writers
type MultiLogWriter struct {
	writers []io.Writer
}

// NewMultiLogWriter creates a writer that duplicates output to all provided writers
func NewMultiLogWriter(writers ...io.Writer) io.Writer {
	return &MultiLogWriter{writers: writers}
}

// Write implements io.Writer
func (w *MultiLogWriter) Write(p []byte) (n int, err error) {
	for _, writer := range w.writers {
		n, err = writer.Write(p)
		if err != nil {
			return n, err
		}
	}
	return len(p), nil
}
