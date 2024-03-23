package cache

import (
	"context"
	"fmt"

	"github.com/hearchco/hearchco/src/cache/nocache"
	"github.com/hearchco/hearchco/src/cache/postgres"
	"github.com/hearchco/hearchco/src/cache/sqlite"
	"github.com/hearchco/hearchco/src/config"
	"github.com/rs/zerolog/log"
)

func New(ctx context.Context, cacheConf config.Cache) (DB, error) {
	var db DB
	var err error

	switch cacheConf.Type {
	case "sqlite", "sqlite3":
		db, err = sqlite.Connect(ctx, cacheConf.TTL.Time, cacheConf.SQLite)
		if err != nil {
			err = fmt.Errorf("failed creating sqlite cache: %w", err)
		} else {
			log.Info().
				Str("path", cacheConf.SQLite.Path).
				Bool("persist", cacheConf.SQLite.Persist).
				Msg("Running with sqlite cache")
		}
	case "postgres", "postgresql":
		db, err = postgres.Connect(ctx, cacheConf.TTL.Time, cacheConf.Postgres)
		if err != nil {
			err = fmt.Errorf("failed creating postgres cache: %w", err)
		} else {
			log.Info().
				Str("uri", cacheConf.Postgres.URI).
				Msg("Running with postgres cache")
		}
	default:
		db, err = nocache.Connect()
		if err != nil {
			err = fmt.Errorf("failed creating a nocache: %w", err)
		} else {
			log.Warn().Msg("Running without caching!")
		}
	}

	return db, err
}
