package dynamodb

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hearchco/agent/src/config"
)

func newDynamoDBConf() config.DynamoDB {
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		panic("AWS_REGION environment variable not set")
	}

	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		panic("DYNAMODB_TABLE environment variable not set")
	}

	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		panic("DYNAMODB_ENDPOINT environment variable not set")
	}

	return config.DynamoDB{Region: awsRegion, Table: tableName, Endpoint: endpoint}
}

var (
	conf      = newDynamoDBConf()
	keyPrefix = "TEST_"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	_, err := New(ctx, keyPrefix, conf)
	if err != nil {
		t.Errorf("error creating dynamodb client: %v", err)
	}
}

func TestClose(t *testing.T) {
	ctx := context.Background()
	db, err := New(ctx, keyPrefix, conf)
	if err != nil {
		t.Errorf("error creating dynamodb client: %v", err)
	}

	db.Close()
}

func TestSet(t *testing.T) {
	ctx := context.Background()
	db, err := New(ctx, keyPrefix, conf)
	if err != nil {
		t.Errorf("error creating dynamodb client: %v", err)
	}

	err = db.Set("testkeyset", "testvalue")
	if err != nil {
		t.Errorf("error setting key-value pair: %v", err)
	}
}

func TestSetTTL(t *testing.T) {
	ctx := context.Background()
	db, err := New(ctx, keyPrefix, conf)
	if err != nil {
		t.Errorf("error creating dynamodb client: %v", err)
	}

	err = db.Set("testkeysetttl", "testvalue", 100*time.Second)
	if err != nil {
		t.Errorf("error setting key-value pair with TTL: %v", err)
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	db, err := New(ctx, keyPrefix, conf)
	if err != nil {
		t.Errorf("error creating dynamodb client: %v", err)
	}

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
	db, err := New(ctx, keyPrefix, conf)
	if err != nil {
		t.Errorf("error creating dynamodb client: %v", err)
	}

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

func TestGetExpired(t *testing.T) {
	ctx := context.Background()
	db, err := New(ctx, keyPrefix, conf)
	if err != nil {
		t.Errorf("error creating dynamodb client: %v", err)
	}

	err = db.Set("testkeygetexpired", "testvalue", 1*time.Second)
	if err != nil {
		t.Errorf("error setting key-value pair with TTL: %v", err)
	}

	time.Sleep(1 * time.Second)

	var value string
	err = db.Get("testkeygetexpired", &value)
	if err != nil {
		t.Errorf("error getting value: %v", err)
	}

	if value != "" {
		t.Errorf("expected no value, got: %v, err: %v", value, err)
	}
}
