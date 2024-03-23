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
	conn    *pgx.Conn
	queries *Queries
}

func Connect(ctx context.Context, ttl time.Duration, conf config.Postgres) (DB, error) {
	conn, err := pgx.Connect(ctx, conf.URI)
	if err != nil {
		return DB{}, err
	}

	return DB{ctx, uint64(ttl.Minutes()), conn, New(conn)}, nil
}

func (db DB) Close() error {
	return db.conn.Close(db.ctx)
}
