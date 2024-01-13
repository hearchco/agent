package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
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
		log.Error().
			Err(err).
			Str("address", fmt.Sprintf("%v:%v/%v", config.Host, config.Port, config.Database)).
			Msg("redis.New(): error connecting to redis")
	} else {
		log.Info().
			Str("address", fmt.Sprintf("%v:%v/%v", config.Host, config.Port, config.Database)).
			Msg("Successful connection to redis")
	}

	return &DB{rdb: rdb, ctx: ctx}
}

func (db *DB) Close() {
	if err := db.rdb.Close(); err != nil {
		log.Error().Err(err).Msg("redis.Close(): error disconnecting from redis")
	} else {
		log.Debug().Msg("Successfully disconnected from redis")
	}
}

func (db *DB) Set(k string, v cache.Value) error {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()

	if val, err := cbor.Marshal(v); err != nil {
		return fmt.Errorf("redis.Set(): error marshaling value: %w", err)
	} else if err := db.rdb.Set(db.ctx, cache.HashString(k), val, 0).Err(); err != nil {
		return fmt.Errorf("redis.Set(): error setting KV to redis: %w", err)
	} else {
		cacheTimeSince := time.Since(cacheTimer)
		log.Trace().
			Int64("ms", cacheTimeSince.Milliseconds()).
			Int64("ns", cacheTimeSince.Nanoseconds()).
			Msg("Cached results")
	}
	return nil
}

func (db *DB) Get(k string, o cache.Value) error {
	v, err := db.rdb.Get(db.ctx, cache.HashString(k)).Result()
	val := []byte(v) // copy data before closing, casting needed for unmarshal

	if err == redis.Nil {
		log.Trace().
			Str("key", k).
			Msg("Found no value in redis")
	} else if err != nil {
		return fmt.Errorf("redis.Get(): error getting value from redis for key %v: %w", k, err)
	} else if err := cbor.Unmarshal(val, o); err != nil {
		return fmt.Errorf("redis.Get(): failed unmarshaling value from redis for key %v: %w", k, err)
	}
	return nil
}
