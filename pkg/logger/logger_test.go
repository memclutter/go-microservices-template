package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want slog.Level
	}{
		{
			name: "production logger",
			env:  "production",
			want: slog.LevelInfo,
		},
		{
			name: "development logger",
			env:  "development",
			want: slog.LevelDebug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.env)
			assert.NotNil(t, logger)
		})
	}
}

func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{})
	logger := &Logger{Logger: slog.New(handler)}

	logger.WithFields(map[string]any{
		"user_id": "123",
		"action":  "login",
	}).Info("User logged in")

	var logEntry map[string]any
	err := json.NewDecoder(&buf).Decode(&logEntry)
	require.NoError(t, err)

	assert.Equal(t, "User logged in", logEntry["msg"])
	assert.Equal(t, "123", logEntry["user_id"])
	assert.Equal(t, "login", logEntry["action"])
}
