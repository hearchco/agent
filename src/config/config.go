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
	"github.com/tminaorg/brzaguza/src/structures"
)

type SESettings struct {
	Timeout  int                 `koanf:"timeout"`
	Shortcut string              `koanf:"shortcut"`
	Crawlers []structures.Engine `koanf:"crawlers"`
}

type Engine struct {
	Enabled  bool       `koanf:"enabled"`
	Settings SESettings `koanf:"settings"`
}

// Config struct for Koanf
type Config struct {
	Engines map[structures.Engine]Engine `koanf:"engines"`
}

var EnabledEngines []structures.Engine = make([]structures.Engine, 0)

func SetupConfig(path string, name string) *Config {
	// Use "." as the key path delimiter. This can be "/" or any character.
	k := koanf.New(".")

	// Load default values using the structs provider.
	// We provide a struct along with the struct tag `koanf` to the
	// provider.
	k.Load(structs.Provider(Config{
		Engines: map[structures.Engine]Engine{
			"bing": {
				Enabled: true,
				Settings: SESettings{
					Shortcut: "bi",
					Crawlers: []structures.Engine{structures.Bing},
				},
			},
			"brave": {
				Enabled: true,
				Settings: SESettings{
					Shortcut: "br",
					Crawlers: []structures.Engine{structures.Brave, structures.Google},
				},
			},
			"duckduckgo": {
				Enabled: true,
				Settings: SESettings{
					Shortcut: "ddg",
					Crawlers: []structures.Engine{structures.Bing},
				},
			},
			"etools": {
				Enabled: true,
				Settings: SESettings{
					Shortcut: "ets",
					Crawlers: []structures.Engine{structures.Bing, structures.Google, structures.Mojeek, structures.Yandex},
				},
			},
			"google": {
				Enabled: true,
				Settings: SESettings{
					Shortcut: "go",
					Crawlers: []structures.Engine{structures.Google},
				},
			},
			"mojeek": {
				Enabled: true,
				Settings: SESettings{
					Shortcut: "mjk",
					Crawlers: []structures.Engine{structures.Mojeek},
				},
			},
			"qwant": {
				Enabled: true,
				Settings: SESettings{
					Shortcut: "qw",
					Crawlers: []structures.Engine{structures.Qwant, structures.Bing},
				},
			},
			"startpage": {
				Enabled: true,
				Settings: SESettings{
					Shortcut: "sp",
					Crawlers: []structures.Engine{structures.Google},
				},
			},
			"swisscows": {
				Enabled: true,
				Settings: SESettings{
					Shortcut: "sc",
					Crawlers: []structures.Engine{structures.Bing},
				},
			},
		},
	}, "koanf"), nil)

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

	// Count number of engines
	for name, engine := range config.Engines {
		if engine.Enabled {
			EnabledEngines = append(EnabledEngines, name)
		}
	}

	return &config
}

const InsertDefaultRank bool = true       // this should be moved to config
const LogDumpLocation string = "logdump/" // this should be moved to config
