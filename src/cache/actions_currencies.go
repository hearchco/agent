package cache

import (
	"fmt"
	"time"

	"github.com/hearchco/agent/src/exchange/currency"
	"github.com/hearchco/agent/src/exchange/engines"
)

func (db DB) SetCurrencies(base currency.Currency, engs []engines.Name, currencies currency.Currencies, ttl ...time.Duration) error {
	key := combineBaseWithExchangeEnginesNames(base, engs)
	return db.driver.Set(key, currencies, ttl...)
}

func (db DB) GetCurrencies(base currency.Currency, engs []engines.Name) (currency.Currencies, error) {
	key := combineBaseWithExchangeEnginesNames(base, engs)
	var currencies currency.Currencies
	err := db.driver.Get(key, &currencies)
	return currencies, err
}

func (db DB) GetCurrenciesTTL(base currency.Currency, engs []engines.Name) (time.Duration, error) {
	key := combineBaseWithExchangeEnginesNames(base, engs)
	return db.driver.GetTTL(key)
}

func combineBaseWithExchangeEnginesNames(base currency.Currency, engs []engines.Name) string {
	return fmt.Sprintf("%v_%v", base.String(), combineExchangeEnginesNames(engs))
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
