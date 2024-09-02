package config

import (
	"time"

	"github.com/hearchco/agent/src/exchange/engines"
)

var exchangeEngines = []engines.Name{
	engines.CURRENCYAPI,
	engines.EXCHANGERATEAPI,
	engines.FRANKFURTER,
}

var exchangeTimings = ExchangeTimings{
	HardTimeout: 1 * time.Second,
}
