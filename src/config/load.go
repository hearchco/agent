package config

import (
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/exchange/currency"
	exchengines "github.com/hearchco/agent/src/exchange/engines"
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
					Currencies: moretime.ConvertToFancyTime(c.Server.Cache.TTL.Currencies),
				},
				Redis:    c.Server.Cache.Redis,
				DynamoDB: c.Server.Cache.DynamoDB,
			},
			ImageProxy: ReaderImageProxy{
				SecretKey: c.Server.ImageProxy.SecretKey,
				Timeout:   moretime.ConvertToFancyTime(c.Server.ImageProxy.Timeout),
			},
		},
		// Initialize the engines config map.
		REngines: map[string]ReaderEngineConfig{},
		// Exchange config.
		RExchange: ReaderExchange{
			BaseCurrency: c.Exchange.BaseCurrency.String(),
			REngines:     map[string]ReaderExchangeEngine{},
			RTimings: ReaderExchangeTimings{
				HardTimeout: moretime.ConvertToFancyTime(c.Exchange.Timings.HardTimeout),
			},
		},
	}

	// Set the engines config map.
	for _, engName := range c.Engines.NoWeb {
		eng := rc.REngines[engName.String()]
		eng.NoWeb = true
		rc.REngines[engName.String()] = eng
	}
	for _, engName := range c.Engines.NoImages {
		eng := rc.REngines[engName.String()]
		eng.NoImages = true
		rc.REngines[engName.String()] = eng
	}
	for _, engName := range c.Engines.NoSuggestions {
		eng := rc.REngines[engName.String()]
		eng.NoSuggestions = true
		rc.REngines[engName.String()] = eng
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
	if rc.Server.ImageProxy.SecretKey == "" {
		log.Fatal().
			Caller().
			Msg("Image proxy secret key is empty")
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
					Currencies: moretime.ConvertFromFancyTime(rc.Server.Cache.TTL.Currencies),
				},
				Redis:    rc.Server.Cache.Redis,
				DynamoDB: rc.Server.Cache.DynamoDB,
			},
			ImageProxy: ImageProxy{
				SecretKey: rc.Server.ImageProxy.SecretKey,
				Timeout:   moretime.ConvertFromFancyTime(rc.Server.ImageProxy.Timeout),
			},
		},
		// Initialize the disabled engines slices.
		Engines: EngineConfig{
			NoWeb:         make([]engines.Name, 0),
			NoImages:      make([]engines.Name, 0),
			NoSuggestions: make([]engines.Name, 0),
		},
		// Exchange config.
		Exchange: Exchange{
			BaseCurrency: currency.ConvertBase(rc.RExchange.BaseCurrency),
			Engines:      []exchengines.Name{},
			Timings: ExchangeTimings{
				HardTimeout: moretime.ConvertFromFancyTime(rc.RExchange.RTimings.HardTimeout),
			},
		},
	}

	// Set the disabled engines slices.
	for engNameS, engConf := range rc.REngines {
		engName, err := engines.NameString(engNameS)
		if err != nil {
			log.Panic().
				Err(err).
				Str("name", engNameS).
				Msg("Couldn't convert engine name string to type")
		}

		if engConf.NoWeb {
			nc.Engines.NoWeb = append(nc.Engines.NoWeb, engName)
		}
		if engConf.NoImages {
			nc.Engines.NoImages = append(nc.Engines.NoImages, engName)
		}
		if engConf.NoSuggestions {
			nc.Engines.NoSuggestions = append(nc.Engines.NoSuggestions, engName)
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
