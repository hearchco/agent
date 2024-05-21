package cache

import (
	"context"
	"fmt"

	"github.com/hearchco/hearchco/src/cache/nocache"
	"github.com/hearchco/hearchco/src/cache/redis"
	"github.com/hearchco/hearchco/src/config"
	"github.com/rs/zerolog/log"
)

func New(ctx context.Context, fileDbPath string, cacheConf config.Cache) (DB, error) {
	var db DB
	var err error

	switch cacheConf.Type {
	case "redis":
		db, err = redis.New(ctx, cacheConf.KeyPrefix, cacheConf.Redis)
		if err != nil {
			err = fmt.Errorf("failed creating a redis cache: %w", err)
		}
	default:
		db, err = nocache.New()
		if err != nil {
			err = fmt.Errorf("failed creating a nocache: %w", err)
		}
		log.Warn().Msg("Running without caching!")
	}

	return db, err
}
