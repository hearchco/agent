package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/hearchco/hearchco/src/config"
	_ "modernc.org/sqlite"
)

type DB struct {
	ctx     context.Context
	ttl     uint64 // in minutes
	db      *sql.DB
	queries *Queries
}

func Connect(ctx context.Context, ttl time.Duration, conf config.SQLite) (DB, error) {
	connString := conf.Path
	if !conf.Persist {
		connString = ":memory:" // TODO: check if this is correct
	}

	db, err := sql.Open("sqlite", connString)
	if err != nil {
		return DB{}, err
	}

	return DB{ctx, uint64(ttl.Minutes()), db, New(db)}, nil
}

func (db DB) Close() error {
	return db.db.Close()
}
