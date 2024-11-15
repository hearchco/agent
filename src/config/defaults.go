package config

import (
	"time"

	exchengines "github.com/hearchco/agent/src/exchange/engines"
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/utils/moretime"
)

func New() Config {
	return Config{
		Server: Server{
			Environment:  "normal",
			Port:         8000,
			FrontendUrls: []string{"http://localhost:5173"},
			Cache: Cache{
				Type:      "none",
				KeyPrefix: "HEARCHCO_",
				TTL: TTL{
					Currencies: moretime.Day,
				},
				Redis: Redis{
					Host: "localhost",
					Port: 6379,
				},
				DynamoDB: DynamoDB{
					Table: "hearchco",
				},
			},
			ImageProxy: ImageProxy{
				Timeout: 3 * time.Second,
			},
		},
		Engines: EngineConfig{
			NoWeb:         []engines.Name{},
			NoImages:      []engines.Name{},
			NoSuggestions: []engines.Name{},
		},
		Exchange: Exchange{
			BaseCurrency: "EUR",
			Engines: []exchengines.Name{
				exchengines.CURRENCYAPI,
				exchengines.EXCHANGERATEAPI,
				exchengines.FRANKFURTER,
			},
			Timings: ExchangeTimings{
				HardTimeout: 1 * time.Second,
			},
		},
	}
}
