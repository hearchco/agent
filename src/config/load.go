package config

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/hearchco/hearchco/src/moretime"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

// passed as pointer since config is modified
func (c *Config) fromReader(rc ReaderConfig) {
	if rc.Server.Proxy.Salt == "" {
		log.Fatal().Msg("config.fromReader(): proxy salt is empty")
	}

	nc := Config{
		Server: Server{
			Environment:  rc.Server.Environment,
			Port:         rc.Server.Port,
			FrontendUrls: strings.Split(rc.Server.FrontendUrls, ","),
			Cache: Cache{
				Type: rc.Server.Cache.Type,
				TTL: TTL{
					Time:        moretime.ConvertFromFancyTime(rc.Server.Cache.TTL.Time),
					RefreshTime: moretime.ConvertFromFancyTime(rc.Server.Cache.TTL.RefreshTime),
				},
				Badger: rc.Server.Cache.Badger,
				Redis:  rc.Server.Cache.Redis,
			},
			Proxy: Proxy{
				Salt: rc.Server.Proxy.Salt,
				Timeouts: ProxyTimeouts{
					Dial:         moretime.ConvertFromFancyTime(rc.Server.Proxy.Timeouts.Dial),
					KeepAlive:    moretime.ConvertFromFancyTime(rc.Server.Proxy.Timeouts.KeepAlive),
					TLSHandshake: moretime.ConvertFromFancyTime(rc.Server.Proxy.Timeouts.TLSHandshake),
				},
			},
		},
		Settings:   map[engines.Name]Settings{},
		Categories: map[category.Name]Category{},
	}

	for key, val := range rc.Settings {
		keyName, err := engines.NameString(key)
		if err != nil {
			log.Panic().
				Err(err).
				Str("engine", key).
				Msg("config.fromReader(): invalid engine name")
			// ^PANIC
		}
		nc.Settings[keyName] = val
	}

	for key, val := range rc.RCategories {
		engArr := []engines.Name{}
		for name, eng := range val.REngines {
			if eng.Enabled {
				engineName, nameErr := engines.NameString(name)
				if nameErr != nil {
					log.Panic().Err(nameErr).Msg("failed converting string to engine name")
					// ^PANIC
				}

				engArr = append(engArr, engineName)
			}
		}
		tim := Timings{
			// HardTimeout: time.Duration(val.RTimings.HardTimeout) * time.Millisecond,
			Timeout:     time.Duration(val.RTimings.Timeout) * time.Millisecond,
			PageTimeout: time.Duration(val.RTimings.PageTimeout) * time.Millisecond,
			Delay:       time.Duration(val.RTimings.Delay) * time.Millisecond,
			RandomDelay: time.Duration(val.RTimings.RandomDelay) * time.Millisecond,
			Parallelism: val.RTimings.Parallelism,
		}
		nc.Categories[key] = Category{
			Ranking: val.Ranking,
			Engines: engArr,
			Timings: tim,
		}
	}

	*c = nc
}

// called when loading default config, before merging with yaml and env
func (c Config) getReader() ReaderConfig {
	rc := ReaderConfig{
		Server: ReaderServer{
			Environment:  c.Server.Environment,
			Port:         c.Server.Port,
			FrontendUrls: strings.Join(c.Server.FrontendUrls, ","),
			Cache: ReaderCache{
				Type: c.Server.Cache.Type,
				TTL: ReaderTTL{
					Time:        moretime.ConvertToFancyTime(c.Server.Cache.TTL.Time),
					RefreshTime: moretime.ConvertToFancyTime(c.Server.Cache.TTL.RefreshTime),
				},
				Badger: c.Server.Cache.Badger,
				Redis:  c.Server.Cache.Redis,
			},
			Proxy: ReaderProxy{
				Salt: c.Server.Proxy.Salt,
				Timeouts: ReaderProxyTimeouts{
					Dial:         moretime.ConvertToFancyTime(c.Server.Proxy.Timeouts.Dial),
					KeepAlive:    moretime.ConvertToFancyTime(c.Server.Proxy.Timeouts.KeepAlive),
					TLSHandshake: moretime.ConvertToFancyTime(c.Server.Proxy.Timeouts.TLSHandshake),
				},
			},
		},
		RCategories: map[category.Name]ReaderCategory{},
		Settings:    map[string]Settings{},
	}

	for key, val := range c.Categories {
		tim := ReaderTimings{
			// HardTimeout: uint(val.Timings.HardTimeout.Milliseconds()),
			Timeout:     uint(val.Timings.Timeout.Milliseconds()),
			PageTimeout: uint(val.Timings.PageTimeout.Milliseconds()),
			Delay:       uint(val.Timings.Delay.Milliseconds()),
			RandomDelay: uint(val.Timings.RandomDelay.Milliseconds()),
			Parallelism: val.Timings.Parallelism,
		}
		rc.RCategories[key] = ReaderCategory{
			Ranking:  val.Ranking,
			REngines: map[string]ReaderEngine{},
			RTimings: tim,
		}
		for _, eng := range val.Engines {
			rc.RCategories[key].REngines[eng.ToLower()] = ReaderEngine{Enabled: true}
		}
	}

	for key, val := range c.Settings {
		rc.Settings[key.ToLower()] = val
	}

	return rc
}

// passed as pointer since config is modified
func (c *Config) Load(dataDirPath string) {
	rc := c.getReader()

	// Use "." as the key path delimiter. This can be "/" or any character.
	k := koanf.New(".")

	// Load default values using the structs provider.
	// We provide a struct along with the struct tag `koanf` to the
	// provider.
	if err := k.Load(structs.Provider(&rc, "koanf"), nil); err != nil {
		log.Panic().Err(err).Msg("config.Load(): failed loading default values")
		// ^PANIC
	}

	// Load YAML config
	yamlPath := path.Join(dataDirPath, "hearchco.yaml")
	if _, err := os.Stat(yamlPath); err != nil {
		log.Trace().
			Str("path", yamlPath).
			Msg("config.Load(): no yaml config found, looking for .yml")
		yamlPath = path.Join(dataDirPath, "hearchco.yml")
		if _, errr := os.Stat(yamlPath); errr != nil {
			log.Trace().
				Str("path", yamlPath).
				Msg("config.Load(): no yaml config found")
		} else if errr := k.Load(file.Provider(yamlPath), yaml.Parser()); errr != nil {
			log.Panic().Err(err).Msg("config.Load(): error loading yaml config")
			// ^PANIC
		}
	} else if err := k.Load(file.Provider(yamlPath), yaml.Parser()); err != nil {
		log.Panic().Err(err).Msg("config.Load(): error loading yaml config")
		// ^PANIC
	}

	// Load ENV config
	if err := k.Load(env.Provider("HEARCHCO_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "HEARCHCO_")), "_", ".", -1)
	}), nil); err != nil {
		log.Panic().Err(err).Msg("config.Load(): error loading env config")
		// ^PANIC
	}

	// Unmarshal config into struct
	if err := k.Unmarshal("", &rc); err != nil {
		log.Panic().Err(err).Msg("config.Load(): failed unmarshaling koanf config")
		// ^PANIC
	}

	c.fromReader(rc)
}
