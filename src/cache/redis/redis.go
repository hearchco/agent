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

var ctx = context.Background()

type DB struct {
	rdb *redis.Client
}

func New(config config.Redis) *DB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.Host, config.Port),
		Password: config.Password,
		DB:       int(config.Database),
	})

	if err := rdb.Ping(ctx); err != nil {
		log.Fatal().Msgf("Error connecting to redis (addr: %v:%v/%v): %v", config.Host, config.Port, config.Database, err)
	} else {
		log.Info().Msgf("Successful connection to redis (addr: %v:%v/%v)", config.Host, config.Port, config.Database)
	}

	return &DB{rdb: rdb}
}

// needed to comply with interface
func (db *DB) Close() {
	log.Debug().Msg("Successfully disconnected from redis")
}

func (db *DB) Set(k string, v cache.Value) {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()
	if val, err := cbor.Marshal(v); err != nil {
		log.Error().Msgf("Error marshaling value: %v", err)
	} else if err := db.rdb.Set(ctx, k, val, 0).Err(); err != nil {
		log.Fatal().Msgf("Error setting KV to redis: %v", err)
	} else {
		log.Debug().Msgf("Cached results in %vns", time.Since(cacheTimer).Nanoseconds())
	}
}

func (db *DB) Get(k string, o cache.Value) {
	v, err := db.rdb.Get(ctx, k).Result()
	val := []byte(v) // copy data before closing, casting needed for unmarshal
	if err == redis.Nil {
		log.Trace().Msgf("Found no value in redis for key (%v): %v", k, err)
	} else if err != nil {
		log.Fatal().Msgf("Error getting value from redis for key (%v): %v", k, err)
	} else if err := cbor.Unmarshal(val, o); err != nil {
		log.Error().Msgf("Failed unmarshaling value from redis for key (%v): %v", k, err)
	}
}
