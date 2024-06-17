package cache

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/cache/nocache"
	"github.com/hearchco/agent/src/cache/redis"
	"github.com/hearchco/agent/src/config"
)

type DB struct {
	driver Driver
}

func New(ctx context.Context, cacheConf config.Cache) (DB, error) {
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

	return DB{drv}, err
}

func (db DB) Close() {
	db.driver.Close()
}
