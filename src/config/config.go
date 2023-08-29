package config

import (
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/engines"
)

// Delegates Timeout, PageTimeout to colly.Collector.SetRequestTimeout(); Note: See https://github.com/gocolly/colly/issues/644
// Delegates Delay, RandomDelay, Parallelism to colly.Collector.Limit()
type Timings struct {
	Timeout     time.Duration `koanf:"timeout"`
	PageTimeout time.Duration `koanf:"pageTimeout"`
	Delay       time.Duration `koanf:"delay"`
	RandomDelay time.Duration `koanf:"randomDelay"`
	Parallelism int           `koanf:"parallelism"`
}

type Settings struct {
	RequestedResultsPerPage int     `koanf:"requestedResults"`
	Shortcut                string  `koanf:"shortcut"`
	Timings                 Timings `koanf:"timings"`
}

type Engine struct {
	Enabled  bool     `koanf:"enabled"`
	Settings Settings `koanf:"settings"`
}

// Server config
type Server struct {
	Port        int    `koanf:"port"`
	FrontendUrl string `koanf:"frontendUrl"`
	RedisUrl    string `koanf:"redisUrl"`
}

// Config struct for Koanf
type Config struct {
	Server  Server                  `koanf:"server"`
	Engines map[engines.Name]Engine `koanf:"engines"`
}

var EnabledEngines []engines.Name = make([]engines.Name, 0)

func SetupConfig(path string, name string) *Config {
	// Use "." as the key path delimiter. This can be "/" or any character.
	k := koanf.New(".")

	// Load default values using the structs provider.
	// We provide a struct along with the struct tag `koanf` to the
	// provider.
	k.Load(structs.Provider(DefaultConfig, "koanf"), nil)

	// Check if path ends with "/" and add it otherwise
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	fullPath := path + name // in case we add other config formats

	// Load YAML config
	yamlPath := fullPath + ".yaml"
	if _, err := os.Stat(yamlPath); err != nil {
		log.Trace().Msgf("no yaml config present at path: %v", yamlPath)
	} else if err := k.Load(file.Provider(yamlPath), yaml.Parser()); err != nil {
		log.Fatal().Msgf("error loading yaml config: %v", err)
	}

	// Load ENV config
	if err := k.Load(env.Provider("BRZAGUZA_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "BRZAGUZA_")), "_", ".", -1)
	}), nil); err != nil {
		log.Fatal().Msgf("error loading env config: %v", err)
	}

	// Unmarshal config into struct
	var config Config
	k.Unmarshal("", &config)

	// Add enabled engines names and remove disabled ones
	for name, engine := range config.Engines {
		if engine.Enabled {
			EnabledEngines = append(EnabledEngines, name)
		} else {
			delete(config.Engines, name)
		}
	}

	return &config
}

const InsertDefaultRank bool = true       // this should be moved to config
const LogDumpLocation string = "logdump/" // this should be moved to config
