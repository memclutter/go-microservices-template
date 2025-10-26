package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		cleanup func()
		wantErr bool
		check   func(*testing.T, *Config)
	}{
		{
			name: "load from defaults",
			setup: func() {
				_ = os.Setenv("APP_ENV", "test")
			},
			cleanup: func() {
				_ = os.Unsetenv("APP_ENV")
			},
			wantErr: false,
			check: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "test", cfg.App.Env)
				assert.Equal(t, 8080, cfg.HTTP.Port)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			cfg, err := Load(".")
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}

func TestGetDatabaseDSN(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "pass",
		Name:     "db",
		SSLMode:  "disable",
	}

	dsn := cfg.GetDatabaseDSN()
	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "user=user")
	assert.Contains(t, dsn, "dbname=db")
}
