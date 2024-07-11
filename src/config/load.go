package config

import (
	"os"
	"slices"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"

	exchengines "github.com/hearchco/agent/src/exchange/engines"
	"github.com/hearchco/agent/src/search/category"
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/utils/moretime"
)

func (c *Config) Load(configPath string) {
	rc := c.getReader()

	// Use "." as the key path delimiter. This can be "/" or any character.
	k := koanf.New(".")

	// Load default values using the structs provider.
	if err := k.Load(structs.Provider(&rc, "koanf"), nil); err != nil {
		log.Panic().
			Caller().
			Err(err).
			Msg("Failed loading default values")
		// ^PANIC
	}

	// Load YAML config.
	if _, err := os.Stat(configPath); err != nil {
		log.Warn().
			Caller().
			Str("path", configPath).
			Msg("No config found on path")
	} else if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		log.Panic().
			Caller().
			Err(err).
			Msg("Error loading config")
		// ^PANIC
	}

	// Load ENV config.
	if err := k.Load(env.Provider("HEARCHCO_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "HEARCHCO_")), "_", ".", -1)
	}), nil); err != nil {
		log.Panic().
			Caller().
			Err(err).
			Msg("Error loading env config")
		// ^PANIC
	}

	// Unmarshal config into struct.
	if err := k.Unmarshal("", &rc); err != nil {
		log.Panic().
			Caller().
			Err(err).
			Msg("Failed unmarshaling koanf config")
		// ^PANIC
	}

	c.fromReader(rc)
}

// Called when loading default config, before merging with YAML and ENV.
func (c Config) getReader() ReaderConfig {
	rc := ReaderConfig{
		// Server config.
		Server: ReaderServer{
			Environment:  c.Server.Environment,
			Port:         c.Server.Port,
			FrontendUrls: strings.Join(c.Server.FrontendUrls, ","),
			Cache: ReaderCache{
				Type: c.Server.Cache.Type,
				TTL: ReaderTTL{
					Results:    moretime.ConvertToFancyTime(c.Server.Cache.TTL.Results),
					Currencies: moretime.ConvertToFancyTime(c.Server.Cache.TTL.Currencies),
				},
				Redis: c.Server.Cache.Redis,
			},
			ImageProxy: ReaderImageProxy{
				Salt:    c.Server.ImageProxy.Salt,
				Timeout: moretime.ConvertToFancyTime(c.Server.ImageProxy.Timeout),
			},
		},
		// Initialize the categories map.
		RCategories: map[category.Name]ReaderCategory{},
		// Exchange config.
		RExchange: ReaderExchange{
			REngines: map[string]ReaderExchangeEngine{},
			RTimings: ReaderExchangeTimings{
				HardTimeout: moretime.ConvertToFancyTime(c.Exchange.Timings.HardTimeout),
			},
		},
	}

	// Set the categories map config.
	for catName, catConf := range c.Categories {
		// Timings config.
		timingsConf := ReaderCategoryTimings{
			PreferredTimeout: moretime.ConvertToFancyTime(catConf.Timings.PreferredTimeout),
			HardTimeout:      moretime.ConvertToFancyTime(catConf.Timings.HardTimeout),
			Delay:            moretime.ConvertToFancyTime(catConf.Timings.Delay),
			RandomDelay:      moretime.ConvertToFancyTime(catConf.Timings.RandomDelay),
			Parallelism:      catConf.Timings.Parallelism,
		}

		// Set the category config.
		rc.RCategories[catName] = ReaderCategory{
			// Initialize the engines map.
			REngines: map[string]ReaderCategoryEngine{},
			Ranking:  catConf.Ranking,
			RTimings: timingsConf,
		}

		// Set the engines map config.
		for _, eng := range catConf.Engines {
			rc.RCategories[catName].REngines[eng.ToLower()] = ReaderCategoryEngine{
				Enabled:           true,
				Required:          slices.Contains(catConf.RequiredEngines, eng),
				RequiredByOrigin:  slices.Contains(catConf.RequiredByOriginEngines, eng),
				Preferred:         slices.Contains(catConf.PreferredEngines, eng),
				PreferredByOrigin: slices.Contains(catConf.PreferredByOriginEngines, eng),
			}
		}
	}

	// Set the exchange engines.
	for _, eng := range c.Exchange.Engines {
		rc.RExchange.REngines[eng.ToLower()] = ReaderExchangeEngine{
			Enabled: true,
		}
	}

	return rc
}

// Passed as pointer since config is modified.
func (c *Config) fromReader(rc ReaderConfig) {
	if rc.Server.ImageProxy.Salt == "" {
		log.Fatal().
			Caller().
			Msg("Image proxy salt is empty")
	}

	nc := Config{
		// Server config.
		Server: Server{
			Environment:  rc.Server.Environment,
			Port:         rc.Server.Port,
			FrontendUrls: strings.Split(rc.Server.FrontendUrls, ","),
			Cache: Cache{
				Type: rc.Server.Cache.Type,
				TTL: TTL{
					Results:    moretime.ConvertFromFancyTime(rc.Server.Cache.TTL.Results),
					Currencies: moretime.ConvertFromFancyTime(rc.Server.Cache.TTL.Currencies),
				},
				Redis: rc.Server.Cache.Redis,
			},
			ImageProxy: ImageProxy{
				Salt:    rc.Server.ImageProxy.Salt,
				Timeout: moretime.ConvertFromFancyTime(rc.Server.ImageProxy.Timeout),
			},
		},
		// Initialize the categories map.
		Categories: map[category.Name]Category{},
		// Exchange config.
		Exchange: Exchange{
			Engines: []exchengines.Name{},
			Timings: ExchangeTimings{
				HardTimeout: moretime.ConvertFromFancyTime(rc.RExchange.RTimings.HardTimeout),
			},
		},
	}

	// Set the categories map config.
	for catName, catRConf := range rc.RCategories {
		// Initialize the engines slices.
		engEnabled := make([]engines.Name, 0)
		engRequired := make([]engines.Name, 0)
		engRequiredByOrigin := make([]engines.Name, 0)
		engPreferred := make([]engines.Name, 0)
		engPreferredByOrigin := make([]engines.Name, 0)

		// Set the engines slices according to the reader config.
		for engS, engRConf := range catRConf.REngines {
			engName, err := engines.NameString(engS)
			if err != nil {
				log.Panic().
					Caller().
					Err(err).
					Msg("Failed converting string to engine name")
				// ^PANIC
			}

			if engRConf.Enabled {
				engEnabled = append(engEnabled, engName)

				if engRConf.Required {
					engRequired = append(engRequired, engName)
				} else if engRConf.RequiredByOrigin {
					engRequiredByOrigin = append(engRequiredByOrigin, engName)
				} else if engRConf.Preferred {
					engPreferred = append(engPreferred, engName)
				} else if engRConf.PreferredByOrigin {
					engPreferredByOrigin = append(engPreferredByOrigin, engName)
				}
			}
		}

		// Timings config.
		timingsConf := CategoryTimings{
			PreferredTimeout: moretime.ConvertFromFancyTime(catRConf.RTimings.PreferredTimeout),
			HardTimeout:      moretime.ConvertFromFancyTime(catRConf.RTimings.HardTimeout),
			Delay:            moretime.ConvertFromFancyTime(catRConf.RTimings.Delay),
			RandomDelay:      moretime.ConvertFromFancyTime(catRConf.RTimings.RandomDelay),
			Parallelism:      catRConf.RTimings.Parallelism,
		}

		// Set the category config.
		nc.Categories[catName] = Category{
			Engines:                  engEnabled,
			RequiredEngines:          engRequired,
			RequiredByOriginEngines:  engRequiredByOrigin,
			PreferredEngines:         engPreferred,
			PreferredByOriginEngines: engPreferredByOrigin,
			Ranking:                  catRConf.Ranking,
			Timings:                  timingsConf,
		}
	}

	// Set the exchange engines.
	for engS, engRConf := range rc.RExchange.REngines {
		engName, err := exchengines.NameString(engS)
		if err != nil {
			log.Panic().
				Caller().
				Err(err).
				Str("engine", engS).
				Msg("Failed converting string to engine name")
			// ^PANIC
		}
		if engRConf.Enabled {
			nc.Exchange.Engines = append(nc.Exchange.Engines, engName)
		}
	}

	// Set the new config.
	*c = nc
}
