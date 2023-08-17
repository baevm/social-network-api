package db

import (
	"context"
	"social-network-api/cfg"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New() (*pgxpool.Pool, error) {
	dbConf, err := pgxpool.ParseConfig(cfg.Get().DB_DSN)

	if err != nil {
		return nil, err
	}

	dbConf.MaxConns = cfg.Get().DB.MaxOpenConn
	dbConf.MinConns = cfg.Get().DB.MaxIdleConn
	dbConf.MaxConnIdleTime = cfg.Get().DB.ConnMaxLifeTime * time.Minute

	db, err := pgxpool.New(context.Background(), cfg.Get().DB_DSN)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Ping(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
