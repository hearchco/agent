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
	ttl     uint64 // in minutes
	db      *sql.DB
	queries *Queries
}

func Connect(ctx context.Context, ttl time.Duration, conf config.SQLite) (DB, error) {
	connString := path.Join(conf.Path, "hearchco.db")
	if !conf.Persist {
		connString = ":memory:"
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

	return DB{ctx, uint64(ttl.Minutes()), db, New(db)}, nil
}

func (db DB) Close() error {
	return db.db.Close()
}
