package logger

import (
	"sync"
)

// LoggingConfig holds all configuration for the logging system
type LoggingConfig struct {
	// Default logging level
	DefaultLevel string `mapstructure:"default_level"`
	// Module specific log levels
	ModuleLevels map[string]string `mapstructure:"module_levels"`
	// Output configuration
	Output OutputConfig `mapstructure:"output"`
	// Format of logs (json, text, console)
	Format string `mapstructure:"format"`
	// Include caller information
	IncludeCaller bool `mapstructure:"include_caller"`
	// Include stack traces for errors
	IncludeStacktrace bool `mapstructure:"include_stacktrace"`
	// Development mode (more verbose, human-readable)
	Development bool `mapstructure:"development"`
	// Sampling configuration
	Sampling SamplingConfig `mapstructure:"sampling"`
}

// OutputConfig defines where logs are sent
type OutputConfig struct {
	// Console output
	Console bool `mapstructure:"console"`
	// File output
	File FileOutputConfig `mapstructure:"file"`
}

// FileOutputConfig for file-based logging
type FileOutputConfig struct {
	// Enable file logging
	Enabled bool `mapstructure:"enabled"`
	// Path to log files
	Path string `mapstructure:"path"`
	// Maximum size in MB before rotating
	MaxSize int `mapstructure:"max_size"`
	// Maximum number of backups to keep
	MaxBackups int `mapstructure:"max_backups"`
	// Maximum age of log files in days
	MaxAge int `mapstructure:"max_age"`
	// Compress rotated logs
	Compress bool `mapstructure:"compress"`
}

// SamplingConfig controls log sampling for high-volume logs
type SamplingConfig struct {
	// Enable sampling
	Enabled bool `mapstructure:"enabled"`
	// Initial logs to process without sampling
	Initial int `mapstructure:"initial"`
	// Sample rate after initial logs (1 in N)
	Thereafter int `mapstructure:"thereafter"`
}

var (
	configOnce   sync.Once
	globalConfig *LoggingConfig
)

// GetConfig returns the global logging configuration
func GetConfig() *LoggingConfig {
	if globalConfig == nil {
		// Default configuration if Initialize hasn't been called
		configOnce.Do(initDefaultConfig)
	}
	return globalConfig
}

// initDefaultConfig initializes a default configuration
func initDefaultConfig() {
	globalConfig = &LoggingConfig{
		DefaultLevel:      InfoLevel,
		ModuleLevels:      make(map[string]string),
		Format:            "json",
		IncludeCaller:     true,
		IncludeStacktrace: true,
		Development:       false,
		Output: OutputConfig{
			Console: true,
			File: FileOutputConfig{
				Enabled:    false,
				Path:       "./logs",
				MaxSize:    100,
				MaxBackups: 5,
				MaxAge:     30,
				Compress:   true,
			},
		},
		Sampling: SamplingConfig{
			Enabled:    false,
			Initial:    100,
			Thereafter: 100,
		},
	}
}

// Initialize configures the logger from the main application config
func Initialize(cfg *LoggingConfig) {
	loggingConfig := &LoggingConfig{
		DefaultLevel:      cfg.DefaultLevel,
		ModuleLevels:      cfg.ModuleLevels,
		Format:            cfg.Format,
		IncludeCaller:     cfg.IncludeCaller,
		IncludeStacktrace: cfg.IncludeStacktrace,
		Development:       cfg.Development,
		Output: OutputConfig{
			Console: cfg.Output.Console,
			File: FileOutputConfig{
				Enabled:    cfg.Output.File.Enabled,
				Path:       cfg.Output.File.Path,
				MaxSize:    cfg.Output.File.MaxSize,
				MaxBackups: cfg.Output.File.MaxBackups,
				MaxAge:     cfg.Output.File.MaxAge,
				Compress:   cfg.Output.File.Compress,
			},
		},
		Sampling: SamplingConfig{
			Enabled:    cfg.Sampling.Enabled,
			Initial:    cfg.Sampling.Initial,
			Thereafter: cfg.Sampling.Thereafter,
		},
	}

	configOnce.Do(func() {
		globalConfig = loggingConfig
	})

	// If already initialized, update and trigger refresh
	if globalConfig != nil && globalConfig != loggingConfig {
		globalConfig = loggingConfig
		triggerConfigRefresh()
	}
}

// SetLogLevel changes the log level for a specific module
func SetLogLevel(module, level string) {
	config := GetConfig()
	config.ModuleLevels[module] = level
	// Signal all loggers to refresh their config
	triggerConfigRefresh()
}
