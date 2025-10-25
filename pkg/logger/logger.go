package logger

import (
	"context"
	"log/slog"
	"os"
)

// Logger is a wrapper around slog.Logger with additional methods
type Logger struct {
	*slog.Logger
}

// New creates a new logger instance based on environment
func New(env string) *Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     getLogLevel(env),
		AddSource: env == "development",
	}

	if env == "production" {
		// JSON format for production (easier to parse in log aggregators)
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// Text format for development (human-readable)
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

func getLogLevel(env string) slog.Level {
	switch env {
	case "production":
		return slog.LevelInfo
	case "development":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

// WithContext adds context fields to logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract common fields from context (request ID, user ID, etc.)
	// Example:
	// if reqID := ctx.Value("request_id"); reqID != nil {
	//     return &Logger{l.With("request_id", reqID)}
	// }
	return l
}

// WithError adds error field to logger
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger: l.With("error", err),
	}
}

// WithField adds a single field
func (l *Logger) WithField(key string, value any) *Logger {
	return &Logger{
		Logger: l.With(key, value),
	}
}

// WithFields adds multiple fields
func (l *Logger) WithFields(fields map[string]any) *Logger {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		Logger: l.With(args...),
	}
}
