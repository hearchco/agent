package badger

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/rs/zerolog/log"
)

type DRV struct {
	keyPrefix string
	client    *badger.DB
}

func New(dataDirPath string, keyPrefix string, config config.Badger) (DRV, error) {
	badgerPath := path.Join(dataDirPath, "database")

	var opt badger.Options
	if config.Persist {
		opt = badger.DefaultOptions(badgerPath).WithLoggingLevel(badger.WARNING)
	} else {
		opt = badger.DefaultOptions("").WithInMemory(true).WithLoggingLevel(badger.WARNING)
	}

	client, err := badger.Open(opt)
	if err != nil {
		log.Error().
			Err(err).
			Bool("persistence", config.Persist).
			Str("path", badgerPath).
			Msg("Error opening badger")
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

	return DRV{keyPrefix, client}, err
}

func (drv DRV) Close() {
	if err := drv.client.Close(); err != nil {
		log.Error().
			Err(err).
			Msg("Error closing badger")
	} else {
		log.Debug().Msg("Successfully closed badger")
	}
}

func (drv DRV) Set(k string, v any, ttl ...time.Duration) error {
	log.Debug().Msg("Caching...")
	cacheTimer := time.Now()

	var setTtl time.Duration = 0
	if len(ttl) > 0 {
		setTtl = ttl[0]
	}

	key := fmt.Sprintf("%v%v", drv.keyPrefix, k)
	if val, err := json.Marshal(v); err != nil {
		return fmt.Errorf("badger.Set(): error marshaling value: %w", err)
	} else if err := drv.client.Update(func(txn *badger.Txn) error {
		var e *badger.Entry
		if setTtl != 0 {
			e = badger.NewEntry([]byte(key), val).WithTTL(ttl[0])
		} else {
			e = badger.NewEntry([]byte(key), val)
		}
		return txn.SetEntry(e)
		// ^returns error into else if
	}); err != nil {
		return fmt.Errorf("badger.Set(): error setting KV to badger: %w", err)
	} else {
		log.Trace().
			Dur("duration", time.Since(cacheTimer)).
			Msg("Cached results")
	}

	return nil
}

func (drv DRV) Get(k string, o any) error {
	key := anonymize.HashToSHA256B64(fmt.Sprintf("%v%v", drv.keyPrefix, k))

	var val []byte
	err := drv.client.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		v, err := item.ValueCopy(nil)
		val = v

		return err
	})

	if err == badger.ErrKeyNotFound {
		log.Trace().
			Str("key", key).
			Msg("Found no value in badger")
	} else if err != nil {
		return fmt.Errorf("badger.Get(): error getting value from badger for key %v: %w", key, err)
	} else if err := json.Unmarshal(val, o); err != nil {
		return fmt.Errorf("badger.Get(): failed unmarshaling value from badger for key %v: %w", key, err)
	}

	return nil
}

// returns time until the key expires, not the time it will be considered expired
func (drv DRV) GetTTL(k string) (time.Duration, error) {
	key := anonymize.HashToSHA256B64(fmt.Sprintf("%v%v", drv.keyPrefix, k))

	var expiresIn time.Duration
	err := drv.client.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
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
			Str("key", key).
			Msg("Found no value in badger")
	} else if err != nil {
		return expiresIn, fmt.Errorf("badger.Get(): error getting value from badger for key %v: %w", key, err)
	}

	return expiresIn, nil
}
