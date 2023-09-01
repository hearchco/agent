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
	Server  Server            `koanf:"server"`
	Engines map[string]Engine `koanf:"engines"`
}

var EnabledEngines []string = make([]string, 0)

func SetupConfig(path string) *Config {
	// Use "." as the key path delimiter. This can be "/" or any character.
	k := koanf.New(".")

	// Load default values using the structs provider.
	// We provide a struct along with the struct tag `koanf` to the
	// provider.
	k.Load(structs.Provider(DefaultConfig, "koanf"), nil)

	// Load YAML config
	yamlPath := path + ".yaml"
	if _, err := os.Stat(yamlPath); err != nil {
		log.Trace().Msgf("no yaml config present at path: %v, looking for .yml", yamlPath)
		yamlPath = fullPath + ".yml"
		if _, errr := os.Stat(yamlPath); errr != nil {
			log.Trace().Msgf("no yaml config present at path: %v", yamlPath)
		} else if errr := k.Load(file.Provider(yamlPath), yaml.Parser()); errr != nil {
			log.Fatal().Msgf("error loading yaml config: %v", err)
		}
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
	config := DefaultConfig
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
