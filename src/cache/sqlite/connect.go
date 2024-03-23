package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hearchco/hearchco/src/config"
	_ "modernc.org/sqlite"
)

type DB struct {
	ctx     context.Context
	ttl     time.Duration
	db      *sql.DB
	queries *Queries
}

func Connect(ctx context.Context, ttl time.Duration, conf config.SQLite) (DB, error) {
	connString := path.Join(conf.Path, "hearchco.db")
	if !conf.Persist {
		connString = ":memory:" // TODO: doesn't work since no migrations run in app
		return DB{}, fmt.Errorf("in-memory sqlite not supported yet")
	} else {
		_, err := os.Stat(conf.Path)
		if os.IsNotExist(err) {
			if err := os.MkdirAll(conf.Path, 0755); err != nil {
				return DB{}, fmt.Errorf("failed creating sqlite directory (%v): %w", conf.Path, err)
			}
		} else if err != nil {
			return DB{}, fmt.Errorf("failed checking sqlite directory (%v): %w", conf.Path, err)
		}
	}

	db, err := sql.Open("sqlite", connString)
	if err != nil {
		return DB{}, err
	}

	if err := db.Ping(); err != nil {
		return DB{}, err
	}

	return DB{ctx, ttl, db, New(db)}, nil
}

func (db DB) Close() error {
	return db.db.Close()
}

func (db DB) Timestamp() time.Time {
	return time.Now().Add(-db.ttl)
}
