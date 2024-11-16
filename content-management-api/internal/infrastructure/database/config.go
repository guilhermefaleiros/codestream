package database

import (
	"context"
	"fmt"
	appConfig "github.com/guilhermefaleiros/codestream/content-management-system/internal/infrastructure/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func NewConnection(ctx context.Context, config appConfig.AppConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	poolConfig.MaxConns = config.Database.MaxConnection
	poolConfig.MinConns = config.Database.MinConnection
	poolConfig.MaxConnIdleTime = time.Duration(config.Database.MaxIdleTime) * time.Second
	poolConfig.MaxConnLifetime = time.Duration(config.Database.MaxLifeTime) * time.Second

	pool, err := pgxpool.NewWithConfig(ctxWithTimeout, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	return pool, nil
}
