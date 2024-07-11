package cache

import (
	"fmt"
	"time"

	"github.com/hearchco/agent/src/exchange/currency"
	"github.com/hearchco/agent/src/exchange/engines"
)

func (db DB) SetCurrencies(engs []engines.Name, currencies map[currency.Currency]float64, ttl ...time.Duration) error {
	key := combineExchangeEnginesNames(engs)
	return db.driver.Set(key, currencies, ttl...)
}

func (db DB) GetCurrencies(engs []engines.Name) (map[currency.Currency]float64, error) {
	key := combineExchangeEnginesNames(engs)
	var currencies map[currency.Currency]float64
	err := db.driver.Get(key, &currencies)
	return currencies, err
}

func (db DB) GetCurrenciesTTL(engs []engines.Name) (time.Duration, error) {
	key := combineExchangeEnginesNames(engs)
	return db.driver.GetTTL(key)
}

func combineExchangeEnginesNames(engs []engines.Name) string {
	var key string
	for i, eng := range engs {
		if i == 0 {
			key = fmt.Sprintf("%v", eng.String())
		} else {
			key = fmt.Sprintf("%v_%v", key, eng.String())
		}
	}
	return key
}
