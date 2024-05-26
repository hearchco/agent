package cache

import (
	"context"
	"fmt"

	"github.com/hearchco/hearchco/src/cache/nocache"
	"github.com/hearchco/hearchco/src/cache/redis"
	"github.com/hearchco/hearchco/src/config"
	"github.com/rs/zerolog/log"
)

type DB struct {
	driver Driver
}

func New(ctx context.Context, fileDbPath string, cacheConf config.Cache) (DB, error) {
	var drv Driver
	var err error

	switch cacheConf.Type {
	case "redis":
		drv, err = redis.New(ctx, cacheConf.KeyPrefix, cacheConf.Redis)
		if err != nil {
			err = fmt.Errorf("failed creating a redis cache: %w", err)
		}
	default:
		drv, err = nocache.New()
		if err != nil {
			err = fmt.Errorf("failed creating a nocache: %w", err)
		}
		log.Warn().Msg("Running without caching!")
	}

	return DB{driver: drv}, err
}

func (db DB) Close() {
	db.driver.Close()
}
