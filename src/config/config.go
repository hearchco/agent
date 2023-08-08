package config

import (
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/structures"
)

// Config struct for Koanf
type Config struct {
	Engines []structures.Engine `koanf:"engines"`
	/* example
	   Type       string              `koanf:"type"`
	   Empty      map[string]string   `koanf:"empty"`

	   	GrandChild struct {
	   		Ids []int `koanf:"ids"`
	   		On  bool  `koanf:"on"`
	   	} `koanf:"grandchild1"`
	*/
}

func SetupConfig(path string, name string) *Config {
	// Use "." as the key path delimiter. This can be "/" or any character.
	k := koanf.New(".")

	// Check if path ends with "/" and add it otherwise
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	fullPath := path + name

	// Load JSON config
	jsonPath := fullPath + ".json"
	if _, err := os.Stat(jsonPath); err != nil {
		log.Trace().Msgf("no json config present at path: %v", jsonPath)
	} else if err := k.Load(file.Provider(jsonPath), json.Parser()); err != nil {
		log.Fatal().Msgf("error loading json config: %v", err)
	}

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

	// Lowercase all engine names
	for index, _ := range config.Engines {
		config.Engines[index] = structures.Engine(strings.ToLower(string(config.Engines[index])))
	}

	return &config
}

const InsertDefaultRank bool = true       // this should be moved to config
const LogDumpLocation string = "logdump/" // this should be moved to config

// update this when loading config
var NumberOfEngines int = 10
