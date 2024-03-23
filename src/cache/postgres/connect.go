package postgres

import (
	"context"
	"time"

	"github.com/hearchco/hearchco/src/config"
	"github.com/jackc/pgx/v5"
)

type DB struct {
	ctx     context.Context
	ttl     uint64 // in minutes
	db      *pgx.Conn
	queries *Queries
}

func Connect(ctx context.Context, ttl time.Duration, conf config.Postgres) (DB, error) {
	db, err := pgx.Connect(ctx, conf.URI)
	if err != nil {
		return DB{}, err
	}

	if err := db.Ping(ctx); err != nil {
		return DB{}, err
	}

	return DB{ctx, uint64(ttl.Minutes()), db, New(db)}, nil
}

func (db DB) Close() error {
	return db.db.Close(db.ctx)
}
