package redis_test

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hearchco/hearchco/src/cache/redis"
	"github.com/hearchco/hearchco/src/config"
)

func newRedisConf() config.Redis {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	var redisPort uint16
	redisPortStr := os.Getenv("REDIS_PORT")

	if redisPortStr == "" {
		redisPort = 6379
	} else {
		redisPortInt, err := strconv.Atoi(redisPortStr)
		if err != nil {
			panic(err)
		}
		redisPort = uint16(redisPortInt)
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	var redisDatabase uint8
	redisDatabaseStr := os.Getenv("REDIS_DATABASE")
	if redisDatabaseStr == "" {
		redisDatabase = 0
	} else {
		redisDatabaseInt, err := strconv.Atoi(redisDatabaseStr)
		if err != nil {
			panic(err)
		}
		redisDatabase = uint8(redisDatabaseInt)
	}

	return config.Redis{
		Host:     redisHost,
		Port:     redisPort,
		Password: redisPassword,
		Database: redisDatabase,
	}
}

var redisConf = newRedisConf()

func TestNew(t *testing.T) {
	ctx := context.Background()
	_, err := redis.New(ctx, redisConf)
	if err != nil {
		t.Errorf("error creating redis: %v", err)
	}
}

func TestClose(t *testing.T) {
	ctx := context.Background()
	db, err := redis.New(ctx, redisConf)
	if err != nil {
		t.Errorf("error creating redis: %v", err)
	}

	db.Close()
}

func TestSet(t *testing.T) {
	ctx := context.Background()
	db, err := redis.New(ctx, redisConf)
	if err != nil {
		t.Errorf("error creating redis: %v", err)
	}

	defer db.Close()

	err = db.Set("testkeyset", "testvalue")
	if err != nil {
		t.Errorf("error setting key-value pair: %v", err)
	}
}

func TestSetTTL(t *testing.T) {
	ctx := context.Background()
	db, err := redis.New(ctx, redisConf)
	if err != nil {
		t.Errorf("error creating redis: %v", err)
	}

	defer db.Close()

	err = db.Set("testkeysetttl", "testvalue", 100*time.Second)
	if err != nil {
		t.Errorf("error setting key-value pair with TTL: %v", err)
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	db, err := redis.New(ctx, redisConf)
	if err != nil {
		t.Errorf("error creating redis: %v", err)
	}

	defer db.Close()

	err = db.Set("testkeyget", "testvalue")
	if err != nil {
		t.Errorf("error setting key-value pair: %v", err)
	}

	var value string
	err = db.Get("testkeyget", &value)
	if err != nil {
		t.Errorf("error getting value: %v", err)
	}

	if value != "testvalue" {
		t.Errorf("expected value: testvalue, got: %v", value)
	}
}

func TestGetTTL(t *testing.T) {
	ctx := context.Background()
	db, err := redis.New(ctx, redisConf)
	if err != nil {
		t.Errorf("error creating redis: %v", err)
	}

	defer db.Close()

	err = db.Set("testkeygetttl", "testvalue", 100*time.Second)
	if err != nil {
		t.Errorf("error setting key-value pair with TTL: %v", err)
	}

	ttl, err := db.GetTTL("testkeygetttl")
	if err != nil {
		t.Errorf("error getting TTL: %v", err)
	}

	// TTL is not exact, so we check for a range
	if ttl > 100*time.Second || ttl < 99*time.Second {
		t.Errorf("expected 100s >= ttl >= 99s, got: %v", ttl)
	}
}
