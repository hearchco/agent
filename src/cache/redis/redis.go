package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/utils/anonymize"
)

type DRV struct {
	ctx       context.Context
	keyPrefix string
	client    *redis.Client
}

func New(ctx context.Context, keyPrefix string, config config.Redis) (DRV, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.Host, config.Port),
		Password: config.Password,
		DB:       int(config.Database),
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Error().
			Err(err).
			Str("address", fmt.Sprintf("%v:%v/%v", config.Host, config.Port, config.Database)).
			Msg("Error creating new connection to redis")
	} else {
		log.Info().
			Str("address", fmt.Sprintf("%v:%v/%v", config.Host, config.Port, config.Database)).
			Msg("Successfully connected to redis")
	}

	return DRV{ctx, keyPrefix, client}, err
}

func (drv DRV) Close() {
	if err := drv.client.Close(); err != nil {
		log.Error().
			Err(err).
			Msg("Error closing connection to redis")
	} else {
		log.Debug().Msg("Successfully disconnected from redis")
	}
}

func (drv DRV) Set(k string, v any, ttl ...time.Duration) error {
	var setTtl time.Duration
	if len(ttl) > 0 {
		setTtl = ttl[0]
	}

	// Create a hash of the key to store in redis.
	keyHash := anonymize.CalculateHashBase64(fmt.Sprintf("%v%v", drv.keyPrefix, k))

	// Serialize the value to JSON.
	valJSON, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("redis.Set(): error marshaling value: %w", err)
	}

	// Encrypt the value using the original key (not the hash).
	val, err := anonymize.Encrypt(valJSON, k)
	if err != nil {
		return fmt.Errorf("redis.Set(): error encrypting value: %w", err)
	}

	// Set the key-value pair in redis.
	if err := drv.client.Set(drv.ctx, keyHash, val, setTtl).Err(); err != nil {
		return fmt.Errorf("redis.Set(): error setting KV to redis: %w", err)
	}

	return nil
}

func (drv DRV) Get(k string, o any) error {
	// Create a hash of the key to retrieve from redis.
	keyHash := anonymize.CalculateHashBase64(fmt.Sprintf("%v%v", drv.keyPrefix, k))

	// Get the value from redis.
	val, err := drv.client.Get(drv.ctx, keyHash).Result()
	if err == redis.Nil {
		log.Trace().
			Str("key_hash", keyHash).
			Msg("Found no value in redis")
		return nil
	} else if err != nil {
		return fmt.Errorf("redis.Get(): error getting value from redis for key hash %v: %w", keyHash, err)
	}

	// Decrypt the value using the original key (not the hash).
	valJSON, err := anonymize.Decrypt(val, k)
	if err != nil {
		return fmt.Errorf("redis.Get(): error decrypting value: %w", err)
	}

	// Deserialize the value from JSON.
	if err := json.Unmarshal(valJSON, o); err != nil {
		return fmt.Errorf("redis.Get(): failed unmarshaling value from redis for key hash %v: %w", keyHash, err)
	}

	return nil
}

// Returns time until the key expires, not the time it will be considered expired.
func (drv DRV) GetTTL(k string) (time.Duration, error) {
	// Create a hash of the key to retrieve from redis.
	keyHash := anonymize.CalculateHashBase64(fmt.Sprintf("%v%v", drv.keyPrefix, k))

	// Returns time with time.Second precision.
	expiresIn, err := drv.client.TTL(drv.ctx, keyHash).Result()
	if err == redis.Nil {
		log.Trace().
			Str("key_hash", keyHash).
			Msg("Found no value in redis")
	} else if err != nil {
		return expiresIn, fmt.Errorf("redis.Get(): error getting value from redis for key hash %v: %w", keyHash, err)
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
