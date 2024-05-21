package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type DB struct {
	ctx       context.Context
	keyPrefix string
	rdb       *redis.Client
}

func New(ctx context.Context, keyPrefix string, config config.Redis) (DB, error) {
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

	return DB{ctx: ctx, keyPrefix: keyPrefix, rdb: rdb}, nil
}

func (db DB) Close() {
	if err := db.rdb.Close(); err != nil {
		log.Error().Err(err).Msg("redis.Close(): error disconnecting from redis")
	} else {
		log.Debug().Msg("Successfully disconnected from redis")
	}
}

func (db DB) Set(k string, v any, ttl ...time.Duration) error {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()

	var setTTL time.Duration = 0
	if len(ttl) > 0 {
		setTTL = ttl[0]
	}

	key := anonymize.HashToSHA256B64(fmt.Sprintf("%v%v", db.keyPrefix, k))
	if val, err := json.Marshal(v); err != nil {
		return fmt.Errorf("redis.Set(): error marshaling value: %w", err)
	} else if err := db.rdb.Set(db.ctx, key, val, setTTL).Err(); err != nil {
		return fmt.Errorf("redis.Set(): error setting KV to redis: %w", err)
	} else {
		log.Trace().
			Dur("duration", time.Since(cacheTimer)).
			Msg("Cached results")
	}

	return nil
}

func (db DB) SetResults(query string, category string, results []result.Result, ttl ...time.Duration) error {
	return db.Set(fmt.Sprintf("%v_%v", query, category), results, ttl...)
}

func (db DB) Get(k string, o any) error {
	key := anonymize.HashToSHA256B64(fmt.Sprintf("%v%v", db.keyPrefix, k))
	val, err := db.rdb.Get(db.ctx, key).Result()
	if err == redis.Nil {
		log.Trace().
			Str("key", key).
			Msg("Found no value in redis")
	} else if err != nil {
		return fmt.Errorf("redis.Get(): error getting value from redis for key %v: %w", key, err)
	} else if err := json.Unmarshal([]byte(val), o); err != nil {
		return fmt.Errorf("redis.Get(): failed unmarshaling value from redis for key %v: %w", key, err)
	}

	return nil
}

func (db DB) GetResults(query string, category string) ([]result.Result, error) {
	var results []result.Result
	err := db.Get(fmt.Sprintf("%v_%v", query, category), &results)
	return results, err
}

// returns time until the key expires, not the time it will be considered expired
func (db DB) GetTTL(k string) (time.Duration, error) {
	kInput := anonymize.HashToSHA256B64(k)

	// returns time with time.Second precision
	expiresIn, err := db.rdb.TTL(db.ctx, kInput).Result()
	if err == redis.Nil {
		log.Trace().
			Str("key", kInput).
			Msg("Found no value in redis")
	} else if err != nil {
		return expiresIn, fmt.Errorf("redis.Get(): error getting value from redis for key %v: %w", kInput, err)
	}

	/*
		In Redis 2.6 or older the command returns -1 if the key does not exist or if the key exist but has no associated expire.
		Starting with Redis 2.8 the return value in case of error changed:
		The command returns -2 if the key does not exist.
		The command returns -1 if the key exists but has no associated expire.
	*/
	if expiresIn < 0 {
		expiresIn = 0
	}

	return expiresIn, nil
}

func (db DB) GetResultsTTL(query string, category string) (time.Duration, error) {
	return db.GetTTL(fmt.Sprintf("%v_%v", query, category))
}
