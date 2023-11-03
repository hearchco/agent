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
		log.Fatal().Err(err).Msgf("redis.New(): error connecting to redis with addr: %v:%v/%v", config.Host, config.Port, config.Database)
		return nil
	} else {
		log.Info().Msgf("Successful connection to redis (addr: %v:%v/%v)", config.Host, config.Port, config.Database)
	}

	return &DB{rdb: rdb, ctx: ctx}
}

func (db *DB) Close() {
	if err := db.rdb.Close(); err != nil {
		log.Fatal().Err(err).Msg("redis.Close(): error disconnecting from redis")
		return
	} else {
		log.Debug().Msg("Successfully disconnected from redis")
	}
}

func (db *DB) Set(k string, v cache.Value) error {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()

	if val, err := cbor.Marshal(v); err != nil {
		return fmt.Errorf("redis.Set(): error marshaling value: %w", err)
	} else if err := db.rdb.Set(db.ctx, k, val, 0).Err(); err != nil {
		log.Fatal().Err(err).Msg("redis.Set(): error setting KV to redis")
		return nil
	} else {
		cacheTimeSince := time.Since(cacheTimer)
		log.Debug().Msgf("Cached results in %vms (%vns)", cacheTimeSince.Milliseconds(), cacheTimeSince.Nanoseconds())
	}
	return nil
}

func (db *DB) Get(k string, o cache.Value) error {
	v, err := db.rdb.Get(db.ctx, k).Result()
	val := []byte(v) // copy data before closing, casting needed for unmarshal

	if err == redis.Nil {
		log.Trace().Msgf("found no value in redis for key %v", k)
	} else if err != nil {
		log.Fatal().Err(err).Msgf("redis.Get(): error getting value from redis for key %v", k)
		return nil
	} else if err := cbor.Unmarshal(val, o); err != nil {
		return fmt.Errorf("redis.Set(): failed unmarshaling value from redis for key %v: %w", k, err)
	}
	return nil
}
