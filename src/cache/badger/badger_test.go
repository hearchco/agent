package badger_test

import (
	"testing"
	"time"

	"github.com/hearchco/hearchco/src/cache/badger"
	"github.com/hearchco/hearchco/src/config"
)

func TestNewInMemory(t *testing.T) {
	_, err := badger.New("", config.Badger{Persist: false})
	if err != nil {
		t.Errorf("error opening in-memory badger: %v", err)
	}
}

func TestNewPersistence(t *testing.T) {
	path := "./testdump/new"
	_, err := badger.New(path, config.Badger{Persist: true})
	if err != nil {
		t.Errorf("error opening badger at %v: %v", path, err)
	}
}

func TestCloseInMemory(t *testing.T) {
	db, err := badger.New("", config.Badger{Persist: false})
	if err != nil {
		t.Errorf("error opening in-memory badger: %v", err)
	}

	db.Close()
}

func TestClosePersistence(t *testing.T) {
	path := "./testdump/close"
	db, err := badger.New(path, config.Badger{Persist: true})
	if err != nil {
		t.Errorf("error opening badger at %v: %v", path, err)
	}

	db.Close()
}

func TestSetInMemory(t *testing.T) {
	db, err := badger.New("", config.Badger{Persist: false})
	if err != nil {
		t.Errorf("error opening in-memory badger: %v", err)
	}

	defer db.Close()

	err = db.Set("testkey", "testvalue")
	if err != nil {
		t.Errorf("error setting key-value pair: %v", err)
	}
}

func TestSetPersistence(t *testing.T) {
	path := "./testdump/set"
	db, err := badger.New(path, config.Badger{Persist: true})
	if err != nil {
		t.Errorf("error opening badger at %v: %v", path, err)
	}

	defer db.Close()

	err = db.Set("testkey", "testvalue")
	if err != nil {
		t.Errorf("error setting key-value pair: %v", err)
	}
}

func TestSetTTLInMemory(t *testing.T) {
	db, err := badger.New("", config.Badger{Persist: false})
	if err != nil {
		t.Errorf("error opening in-memory badger: %v", err)
	}

	defer db.Close()

	err = db.Set("testkey", "testvalue", 100*time.Second)
	if err != nil {
		t.Errorf("error setting key-value pair with TTL: %v", err)
	}
}

func TestSetTTLPersistence(t *testing.T) {
	path := "./testdump/setttl"
	db, err := badger.New(path, config.Badger{Persist: true})
	if err != nil {
		t.Errorf("error opening badger at %v: %v", path, err)
	}

	defer db.Close()

	err = db.Set("testkey", "testvalue", 100*time.Second)
	if err != nil {
		t.Errorf("error setting key-value pair with TTL: %v", err)
	}
}

func TestGetInMemory(t *testing.T) {
	db, err := badger.New("", config.Badger{Persist: false})
	if err != nil {
		t.Errorf("error opening in-memory badger: %v", err)
	}

	defer db.Close()

	err = db.Set("testkey", "testvalue")
	if err != nil {
		t.Errorf("error setting key-value pair: %v", err)
	}

	var value string
	err = db.Get("testkey", &value)
	if err != nil {
		t.Errorf("error getting key-value pair: %v", err)
	}

	if value != "testvalue" {
		t.Errorf("expected value: testvalue, got: %v", value)
	}
}

func TestGetPersistence(t *testing.T) {
	path := "./testdump/get"
	db, err := badger.New(path, config.Badger{Persist: true})
	if err != nil {
		t.Errorf("error opening badger at %v: %v", path, err)
	}

	defer db.Close()

	err = db.Set("testkey", "testvalue")
	if err != nil {
		t.Errorf("error setting key-value pair: %v", err)
	}

	var value string
	err = db.Get("testkey", &value)
	if err != nil {
		t.Errorf("error getting key-value pair: %v", err)
	}

	if value != "testvalue" {
		t.Errorf("expected value: testvalue, got: %v", value)
	}
}

func TestGetTTLInMemory(t *testing.T) {
	db, err := badger.New("", config.Badger{Persist: false})
	if err != nil {
		t.Errorf("error opening in-memory badger: %v", err)
	}

	defer db.Close()

	err = db.Set("testkey", "testvalue", 100*time.Second)
	if err != nil {
		t.Errorf("error setting key-value pair with TTL: %v", err)
	}

	ttl, err := db.GetTTL("testkey")
	if err != nil {
		t.Errorf("error getting TTL: %v", err)
	}

	if ttl > 100*time.Second || ttl < 99*time.Second {
		t.Errorf("expected 100s >= ttl >= 99s, got: %v", ttl)
	}
}

func TestGetTTLPersistence(t *testing.T) {
	path := "./testdump/getttl"
	db, err := badger.New(path, config.Badger{Persist: true})
	if err != nil {
		t.Errorf("error opening badger at %v: %v", path, err)
	}

	defer db.Close()

	err = db.Set("testkey", "testvalue", 100*time.Second)
	if err != nil {
		t.Errorf("error setting key-value pair with TTL: %v", err)
	}

	ttl, err := db.GetTTL("testkey")
	if err != nil {
		t.Errorf("error getting TTL: %v", err)
	}

	if ttl > 100*time.Second || ttl < 99*time.Second {
		t.Errorf("expected 100s >= ttl >= 99s, got: %v", ttl)
	}
}
