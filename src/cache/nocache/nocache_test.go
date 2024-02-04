package nocache_test

import (
	"testing"

	"github.com/hearchco/hearchco/src/cache/nocache"
)

func TestNew(t *testing.T) {
	_, err := nocache.New()
	if err != nil {
		t.Errorf("error creating nocache: %v", err)
	}
}

func TestClose(t *testing.T) {
	db, err := nocache.New()
	if err != nil {
		t.Errorf("error creating nocache: %v", err)
	}
	db.Close()
}

func TestSet(t *testing.T) {
	db, err := nocache.New()
	if err != nil {
		t.Errorf("error creating nocache: %v", err)
	}
	defer db.Close()

	err = db.Set("testkey", "testvalue")
	if err != nil {
		t.Errorf("error setting key-value pair: %v", err)
	}
}

func TestGet(t *testing.T) {
	db, err := nocache.New()
	if err != nil {
		t.Errorf("error creating nocache: %v", err)
	}
	defer db.Close()

	err = db.Set("testkey", "testvalue")
	if err != nil {
		t.Errorf("error setting key-value pair: %v", err)
	}

	var value string = "testvalue"
	err = db.Get("testkey", &value)
	if err != nil {
		t.Errorf("error getting value: %v", err)
	}
	if value != "testvalue" {
		t.Errorf("expected value: testvalue, got: %v", value)
	}
}

func TestGetTTL(t *testing.T) {
	db, err := nocache.New()
	if err != nil {
		t.Errorf("error creating nocache: %v", err)
	}
	defer db.Close()

	err = db.Set("testkey", "testvalue", 1)
	if err != nil {
		t.Errorf("error setting key-value pair with TTL: %v", err)
	}

	ttl, err := db.GetTTL("testkey")
	if err != nil {
		t.Errorf("error getting TTL: %v", err)
	}
	if ttl != 0 {
		t.Errorf("expected TTL: 0, got: %v", ttl)
	}
}
