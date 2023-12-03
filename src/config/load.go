package config

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
	"github.com/hearchco/hearchco/src/category"
	"github.com/hearchco/hearchco/src/engines"
)

var EnabledEngines []engines.Name = make([]engines.Name, 0)
var LogDumpLocation string = "dump/"

func (c *Config) fromReader(rc *ReaderConfig) {
	nc := Config{
		Server:     rc.Server,
		Settings:   map[engines.Name]Settings{},
		Categories: map[category.Name]Category{},
	}

	for key, val := range rc.Settings {
		keyName, err := engines.NameString(key)
		if err != nil {
			log.Panic().Err(err).Msgf("failed reading config. invalid engine name: %v", key)
			return
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
					return
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

func (c *Config) getReader() ReaderConfig {
	rc := ReaderConfig{
		Server:      c.Server,
		Settings:    map[string]Settings{},
		RCategories: map[category.Name]ReaderCategory{},
	}

	for key, val := range c.Settings {
		rc.Settings[key.ToLower()] = val
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

	return rc
}

func (c *Config) Load(dataDirPath string, logDirPath string) {
	rc := c.getReader()

	// Load vars
	loadVars(logDirPath)

	// Use "." as the key path delimiter. This can be "/" or any character.
	k := koanf.New(".")

	// Load default values using the structs provider.
	// We provide a struct along with the struct tag `koanf` to the
	// provider.
	if err := k.Load(structs.Provider(&rc, "koanf"), nil); err != nil {
		log.Panic().Err(err).Msg("failed loading default values")
	}

	// Load YAML config
	yamlPath := path.Join(dataDirPath, "hearchco.yaml")
	if _, err := os.Stat(yamlPath); err != nil {
		log.Trace().Msgf("no yaml config present at path: %v, looking for .yml", yamlPath)
		yamlPath = path.Join(dataDirPath, "hearchco.yml")
		if _, errr := os.Stat(yamlPath); errr != nil {
			log.Trace().Msgf("no yaml config present at path: %v", yamlPath)
		} else if errr := k.Load(file.Provider(yamlPath), yaml.Parser()); errr != nil {
			log.Panic().Err(err).Msg("error loading yaml config")
		}
	} else if err := k.Load(file.Provider(yamlPath), yaml.Parser()); err != nil {
		log.Panic().Err(err).Msg("error loading yaml config")
	}

	// Load ENV config
	if err := k.Load(env.Provider("HEARCHCO_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "HEARCHCO_")), "_", ".", -1)
	}), nil); err != nil {
		log.Panic().Err(err).Msg("error loading env config")
	}

	// Unmarshal config into struct
	if err := k.Unmarshal("", &rc); err != nil {
		log.Panic().Err(err).Msg("failed unmarshaling koanf config")
	}

	c.fromReader(&rc)
}

func loadVars(logDirPath string) {
	LogDumpLocation = path.Join(logDirPath, LogDumpLocation)
}
