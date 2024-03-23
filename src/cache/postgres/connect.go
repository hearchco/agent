package postgres

import (
	"context"
	"time"

	"github.com/hearchco/hearchco/src/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type DB struct {
	ctx     context.Context
	ttl     time.Duration
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

	return DB{ctx, ttl, db, New(db)}, nil
}

func (db DB) Close() error {
	return db.db.Close(db.ctx)
}

func (db DB) Timestamp() pgtype.Timestamp {
	return pgtype.Timestamp{
		Time: time.Now().Add(-db.ttl),
	}
}
