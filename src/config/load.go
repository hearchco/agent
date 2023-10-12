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
	"github.com/tminaorg/brzaguza/src/category"
	"github.com/tminaorg/brzaguza/src/engines"
)

var EnabledEngines []engines.Name = make([]engines.Name, 0)
var LogDumpLocation string = "dump/"

func (c *Config) FromReader(rc *ReaderConfig) {
	nc := Config{
		Server:   rc.Server,
		Settings: rc.Settings,
	}
	nc.Categories = map[category.Name]Category{}
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
		nc.Categories[key] = Category{
			Ranking: val.Ranking,
			Engines: engArr,
		}
	}

	*c = nc
}

func ReaderFromConfig(c *Config) ReaderConfig {
	rc := ReaderConfig{
		Server:   c.Server,
		Settings: c.Settings,
	}
	rc.RCategories = map[category.Name]ReaderCategory{}
	for key, val := range c.Categories {
		rc.RCategories[key] = ReaderCategory{
			Ranking:  val.Ranking,
			REngines: map[string]ReaderEngine{},
		}
		for _, eng := range val.Engines {
			rc.RCategories[key].REngines[eng.ToLower()] = ReaderEngine{Enabled: true}
		}
	}

	return rc
}

func (c *Config) Load(path string, logPath string) {
	rc := ReaderFromConfig(c)

	// Load vars
	loadVars(logPath)

	// Use "." as the key path delimiter. This can be "/" or any character.
	k := koanf.New(".")

	// Load default values using the structs provider.
	// We provide a struct along with the struct tag `koanf` to the
	// provider.
	k.Load(structs.Provider(&rc, "koanf"), nil)

	// Load YAML config
	yamlPath := path + "/brzaguza.yaml"
	if _, err := os.Stat(yamlPath); err != nil {
		log.Trace().Msgf("no yaml config present at path: %v, looking for .yml", yamlPath)
		yamlPath = path + "/brzaguza.yml"
		if _, errr := os.Stat(yamlPath); errr != nil {
			log.Trace().Msgf("no yaml config present at path: %v", yamlPath)
		} else if errr := k.Load(file.Provider(yamlPath), yaml.Parser()); errr != nil {
			log.Panic().Msgf("error loading yaml config: %v", err)
		}
	} else if err := k.Load(file.Provider(yamlPath), yaml.Parser()); err != nil {
		log.Panic().Msgf("error loading yaml config: %v", err)
	}

	// Load ENV config
	if err := k.Load(env.Provider("BRZAGUZA_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "BRZAGUZA_")), "_", ".", -1)
	}), nil); err != nil {
		log.Panic().Msgf("error loading env config: %v", err)
	}

	// Unmarshal config into struct
	k.Unmarshal("", &rc)

	c.FromReader(&rc)
}

func loadVars(logPath string) {
	LogDumpLocation = logPath + "/" + LogDumpLocation
}
