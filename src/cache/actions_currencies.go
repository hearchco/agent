package cache

import (
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
	enginesNamesStrings := make([]string, 0, len(engs)+1)
	for _, eng := range engs {
		enginesNamesStrings = append(enginesNamesStrings, eng.String())
	}

	baseWithEnginesNamesStrings := append(enginesNamesStrings, base.String())
	return combineIntoKey(baseWithEnginesNamesStrings...)
}
