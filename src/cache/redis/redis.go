package redis

import (
	"context"
	"fmt"

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
	log.Info().Msgf("Successful connection to redis (addr: %v:%v/%v)", config.Host, config.Port, config.Database)
	return &DB{rdb: rdb}
}

// needed to comply with interface
func (db *DB) Close() {
	log.Debug().Msg("Successfully disconnected from redis")
}

func (db *DB) Set(k string, v cache.Value) {
	if val, err := cbor.Marshal(v); err != nil {
		log.Error().Msgf("Error marshaling value: %v", err)
	} else if err := db.rdb.Set(ctx, k, val, 0).Err(); err != nil {
		log.Panic().Msgf("Error setting KV to redis: %v", err)
	}
}

func (db *DB) Get(k string, o cache.Value) {
	v, err := db.rdb.Get(ctx, k).Result()
	val := []byte(v) // copy data before closing, casting needed for unmarshal
	if err == redis.Nil {
		log.Trace().Msgf("Found no value in redis for key (%v): %v", k, err)
	} else if err != nil {
		log.Panic().Msgf("Error getting value from redis for key (%v): %v", k, err)
	} else if err := cbor.Unmarshal(val, o); err != nil {
		log.Error().Msgf("Failed unmarshaling value from redis for key (%v): %v", k, err)
	}
}
