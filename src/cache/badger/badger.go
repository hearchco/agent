package badger

import (
	"fmt"
	"path"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/fxamacker/cbor/v2"
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/rs/zerolog/log"
)

type DB struct {
	bdb *badger.DB
}

func New(dataDirPath string, config config.Badger) (*DB, error) {
	badgerPath := path.Join(dataDirPath, "database")

	var opt badger.Options
	if config.Persist {
		opt = badger.DefaultOptions(badgerPath).WithLoggingLevel(badger.WARNING)
	} else {
		opt = badger.DefaultOptions("").WithInMemory(true).WithLoggingLevel(badger.WARNING)
	}

	bdb, err := badger.Open(opt)

	if err != nil {
		log.Error().
			Err(err).
			Bool("persistence", config.Persist).
			Str("path", badgerPath).
			Msg("badger.New(): error opening badger")
	} else if config.Persist {
		log.Info().
			Bool("persistence", config.Persist).
			Str("path", badgerPath).
			Msg("Successfully opened badger")
	} else {
		log.Info().
			Bool("persistence", config.Persist).
			Msg("Successfully opened in-memory badger")
	}

	return &DB{bdb: bdb}, err
}

func (db *DB) Close() {
	if err := db.bdb.Close(); err != nil {
		log.Error().Err(err).Msg("badger.Close(): error closing badger")
	} else {
		log.Debug().Msg("Successfully closed badger")
	}
}

func (db *DB) Set(k string, v cache.Value, ttl ...time.Duration) error {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()

	var setTtl time.Duration = 0
	if len(ttl) > 0 {
		setTtl = ttl[0]
	}

	if val, err := cbor.Marshal(v); err != nil {
		return fmt.Errorf("badger.Set(): error marshaling value: %w", err)
	} else if err := db.bdb.Update(func(txn *badger.Txn) error {
		var e *badger.Entry
		if setTtl != 0 {
			e = badger.NewEntry([]byte(anonymize.HashToSHA256B64(k)), val).WithTTL(ttl[0])
		} else {
			e = badger.NewEntry([]byte(anonymize.HashToSHA256B64(k)), val)
		}
		return txn.SetEntry(e)
		// ^returns error into else if
	}); err != nil {
		return fmt.Errorf("badger.Set(): error setting KV to badger: %w", err)
	} else {
		cacheTimeSince := time.Since(cacheTimer)
		log.Trace().
			Int64("ms", cacheTimeSince.Milliseconds()).
			Int64("ns", cacheTimeSince.Nanoseconds()).
			Msg("Cached results")
	}

	return nil
}

func (db *DB) Get(k string, o cache.Value, hashed ...bool) error {
	var kInput string
	if len(hashed) > 0 && hashed[0] {
		kInput = k
	} else {
		kInput = anonymize.HashToSHA256B64(k)
	}

	var val []byte
	err := db.bdb.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(kInput))
		if err != nil {
			return err
		}

		v, err := item.ValueCopy(nil)
		val = v

		return err
	})

	if err == badger.ErrKeyNotFound {
		log.Trace().
			Str("key", kInput).
			Msg("Found no value in badger")
	} else if err != nil {
		return fmt.Errorf("badger.Get(): error getting value from badger for key %v: %w", kInput, err)
	} else if err := cbor.Unmarshal(val, o); err != nil {
		return fmt.Errorf("badger.Get(): failed unmarshaling value from badger for key %v: %w", kInput, err)
	}

	return nil
}

// returns time until the key expires, not the time it will be considered expired
func (db *DB) GetTTL(k string, hashed ...bool) (time.Duration, error) {
	var kInput string
	if len(hashed) > 0 && hashed[0] {
		kInput = k
	} else {
		kInput = anonymize.HashToSHA256B64(k)
	}

	var expiresIn time.Duration
	err := db.bdb.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(kInput))
		if err != nil {
			return err
		}

		expiresAtUnix := time.Unix(int64(item.ExpiresAt()), 0)
		expiresIn = time.Until(expiresAtUnix)

		// returns negative time.Since() if expiresAtUnix is in the past
		if expiresIn < 0 {
			expiresIn = 0
		}

		return err
	})

	if err == badger.ErrKeyNotFound {
		log.Trace().
			Str("key", kInput).
			Msg("Found no value in badger")
	} else if err != nil {
		return expiresIn, fmt.Errorf("badger.Get(): error getting value from badger for key %v: %w", kInput, err)
	}

	return expiresIn, nil
}
