package log

import (
	"log/slog"
	"os"
)

type Logger struct {
	stderr *slog.Logger
	stdout *slog.Logger
}

func New() *Logger {
	return &Logger{
		stderr: slog.New(slog.NewJSONHandler(os.Stderr, nil)),
		stdout: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func NewWithSetLevel(logLevel string) *Logger {
	var level slog.Level

	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	}

	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}

	return &Logger{
		stderr: slog.New(slog.NewJSONHandler(os.Stderr, nil)),
		stdout: slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts)),
	}
}

// Debug logs at [LevelDebug].
func (l Logger) Debug(msg string, args ...any) {
	l.stdout.Debug(msg, args...)
}

// Info logs at [LevelInfo].
func (l Logger) Info(msg string, args ...any) {
	l.stdout.Info(msg, args...)
}

// Error logs at [LevelError].
func (l Logger) Error(msg string, args ...any) {
	l.stderr.Error(msg, args...)
}
