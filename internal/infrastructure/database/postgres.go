package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/memclutter/go-microservices-template/pkg/config"
	"github.com/memclutter/go-microservices-template/pkg/logger"
)

// NewPostgresPool creates a new PostgreSQL connection pool
func NewPostgresPool(ctx context.Context, cfg *config.DatabaseConfig, log *logger.Logger) (*pgxpool.Pool, error) {
	dsn := cfg.GetDatabaseDSN()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure connection pool
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("PostgreSQL connection pool created successfully")

	return pool, nil
}

// Close closes the database connection pool
func ClosePostgresPool(pool *pgxpool.Pool, log *logger.Logger) {
	if pool != nil {
		pool.Close()
		log.Info("PostgreSQL connection pool closed")
	}
}
