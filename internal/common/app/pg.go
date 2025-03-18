package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func getPgxpool(ctx context.Context, cfg PgConfig) (*pgxpool.Pool, error) {
	cfgpool, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, err
	}
	cfgpool.MaxConns = cfg.MaxOpenConns
	cfgpool.MinConns = cfg.MinConns
	cfgpool.MaxConnLifetime = cfg.MaxConnLifetime
	dbpool, err := pgxpool.NewWithConfig(ctx, cfgpool)
	if err != nil {
		return nil, err
	}

	if err = dbpool.Ping(ctx); err != nil {
		return nil, err
	}
	return dbpool, nil
}
