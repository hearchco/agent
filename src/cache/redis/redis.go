package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type DB struct {
	ctx context.Context
	rdb *redis.Client
}

func New(ctx context.Context, config config.Redis) (DB, error) {
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

	return DB{rdb: rdb, ctx: ctx}, nil
}

func (db DB) Close() {
	if err := db.rdb.Close(); err != nil {
		log.Error().Err(err).Msg("redis.Close(): error disconnecting from redis")
	} else {
		log.Debug().Msg("Successfully disconnected from redis")
	}
}

func (db DB) Set(k string, v interface{}, ttl ...time.Duration) error {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()

	var setTtl time.Duration = 0
	if len(ttl) > 0 {
		setTtl = ttl[0]
	}

	if val, err := cbor.Marshal(v); err != nil {
		return fmt.Errorf("redis.Set(): error marshaling value: %w", err)
	} else if err := db.rdb.Set(db.ctx, anonymize.HashToSHA256B64(k), val, setTtl).Err(); err != nil {
		return fmt.Errorf("redis.Set(): error setting KV to redis: %w", err)
	} else {
		log.Trace().
			Dur("duration", time.Since(cacheTimer)).
			Msg("Cached results")
	}

	return nil
}

func (db DB) Get(k string, o interface{}, hashed ...bool) error {
	var kInput string
	if len(hashed) > 0 && hashed[0] {
		kInput = k
	} else {
		kInput = anonymize.HashToSHA256B64(k)
	}

	val, err := db.rdb.Get(db.ctx, kInput).Result()
	if err == redis.Nil {
		log.Trace().
			Str("key", kInput).
			Msg("Found no value in redis")
	} else if err != nil {
		return fmt.Errorf("redis.Get(): error getting value from redis for key %v: %w", kInput, err)
	} else if err := cbor.Unmarshal([]byte(val), o); err != nil {
		return fmt.Errorf("redis.Get(): failed unmarshaling value from redis for key %v: %w", kInput, err)
	}

	return nil
}

// returns time until the key expires, not the time it will be considered expired
func (db DB) GetTTL(k string, hashed ...bool) (time.Duration, error) {
	var kInput string
	if len(hashed) > 0 && hashed[0] {
		kInput = k
	} else {
		kInput = anonymize.HashToSHA256B64(k)
	}

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
