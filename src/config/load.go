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
	"github.com/tminaorg/brzaguza/src/engines"
)

var EnabledEngines []engines.Name = make([]engines.Name, 0)

<<<<<<< HEAD
func (c *Config) Load(path string) {
=======
func (c *Config) Load(path string, name string) {
>>>>>>> 9a569e1 (Rebasing to logger conflicts)
	// Use "." as the key path delimiter. This can be "/" or any character.
	k := koanf.New(".")

	// Load default values using the structs provider.
	// We provide a struct along with the struct tag `koanf` to the
	// provider.
	k.Load(structs.Provider(c, "koanf"), nil)

	// Load YAML config
	yamlPath := path + ".yaml"
	if _, err := os.Stat(yamlPath); err != nil {
		log.Trace().Msgf("no yaml config present at path: %v, looking for .yml", yamlPath)
		yamlPath = path + ".yml"
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
	k.Unmarshal("", &c)

	// Add enabled engines names and remove disabled ones
	for name, engine := range c.Engines {
		if engine.Enabled {
			EnabledEngines = append(EnabledEngines, engines.ConvertToName(name))
		} else {
			delete(c.Engines, name)
		}
	}
}

const InsertDefaultRank bool = true       // this should be moved to config
const LogDumpLocation string = "logdump/" // this should be moved to config
