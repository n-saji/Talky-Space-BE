package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func dbConfig() *pgxpool.Config {
	const defaultMaxConns = 20
	const defaultMinConns = 5
	const defaultMaxConnLifetime = time.Hour * 1
	const defaultMaxConnIdleTime = time.Minute * 10
	const defaultConnTimeout = time.Second * 30
	const defaultHealthCheckPeriod = time.Minute * 1

	DbConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("Unable to parse DATABASE_URL: " + err.Error())
	}

	DbConfig.MaxConns = defaultMaxConns
	DbConfig.MinConns = defaultMinConns
	DbConfig.MaxConnLifetime = defaultMaxConnLifetime
	DbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	DbConfig.ConnConfig.ConnectTimeout = defaultConnTimeout
	DbConfig.HealthCheckPeriod = defaultHealthCheckPeriod

	return DbConfig
}

func PingDb(ctx context.Context, pool *pgxpool.Pool) error {
	return pool.Ping(ctx)
}

func ConnectToDb(ctx context.Context) (*pgxpool.Pool, error) {

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	return pool, nil
}
