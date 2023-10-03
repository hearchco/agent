package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
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
	log.Info().Msgf("successful connection to redis (addr: %v:%v/%v)", config.Host, config.Port, config.Database)
	return &DB{rdb: rdb}
}

// needed to comply with interface
func (db *DB) Close() {
	log.Debug().Msg("successfully disconnected from redis")
}

func (db *DB) Set(k string, v interface{}) {
	if err := db.rdb.Set(ctx, k, v, 0).Err(); err != nil {
		log.Trace().Msgf("key = %v, value = %v", k, v)
		log.Panic().Msgf("error setting KV to redis: %v", err)
	} else {
		log.Trace().Msgf("success: key = %v, value = %v", k, v)
	}
}

func (db *DB) Get(k string) string {
	v, err := db.rdb.Get(ctx, k).Result()

	if err == redis.Nil {
		log.Trace().Msgf("error: key = %v, value = %v", k, v)
		log.Error().Msgf("found no value in redis for key (%v): %v", k, err)
	} else if err != nil {
		log.Trace().Msgf("error: key = %v, value = %v", k, v)
		log.Panic().Msgf("error getting value from redis for key (%v): %v", k, err)
	}

	log.Trace().Msgf("success: key = %v, value = %v", k, v)
	return v
}
