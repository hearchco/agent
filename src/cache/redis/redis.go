package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/cache"
	"github.com/tminaorg/brzaguza/src/config"
)

type DB struct {
	rdb *redis.Client
	ctx context.Context
}

func New(ctx context.Context, config config.Redis) *DB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.Host, config.Port),
		Password: config.Password,
		DB:       int(config.Database),
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal().Msgf("Error connecting to redis (addr: %v:%v/%v): %v", config.Host, config.Port, config.Database, err)
	} else {
		log.Info().Msgf("Successful connection to redis (addr: %v:%v/%v)", config.Host, config.Port, config.Database)
	}

	return &DB{rdb: rdb, ctx: ctx}
}

func (db *DB) Close() {
	if err := db.rdb.Close(); err != nil {
		log.Fatal().Msgf("Error disconnecting from redis: %v", err)
	} else {
		log.Debug().Msg("Successfully disconnected from redis")
	}
}

func (db *DB) Set(k string, v cache.Value) {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()

	if val, err := cbor.Marshal(v); err != nil {
		log.Error().Msgf("Error marshaling value: %v", err)
	} else if err := db.rdb.Set(db.ctx, k, val, 0).Err(); err != nil {
		log.Fatal().Msgf("Error setting KV to redis: %v", err)
	} else {
		cacheTimeSince := time.Since(cacheTimer)
		log.Debug().Msgf("Cached results in %vms (%vns)", cacheTimeSince.Milliseconds(), cacheTimeSince.Nanoseconds())
	}
}

func (db *DB) Get(k string, o cache.Value) {
	v, err := db.rdb.Get(db.ctx, k).Result()
	val := []byte(v) // copy data before closing, casting needed for unmarshal

	if err == redis.Nil {
		log.Trace().Msgf("Found no value in redis for key (%v): %v", k, err)
	} else if err != nil {
		log.Fatal().Msgf("Error getting value from redis for key (%v): %v", k, err)
	} else if err := cbor.Unmarshal(val, o); err != nil {
		log.Error().Msgf("Failed unmarshaling value from redis for key (%v): %v", k, err)
	}
}
