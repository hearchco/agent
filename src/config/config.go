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
	Name     structures.Engine `koanf:"name"`
	Enabled  bool              `koanf:"enabled"`
	Settings SESettings        `koanf:"settings"`
}

// Config struct for Koanf
type Config struct {
	Engines []Engine `koanf:"engines"`
}

var EnabledEngines []structures.Engine = make([]structures.Engine, 0)
var engineDefaults []Engine = []Engine{
	{
		Name:    "bing",
		Enabled: true,
		Settings: SESettings{
			Shortcut: "bi",
			Crawlers: []structures.Engine{structures.Bing},
		},
	},
	{
		Name:    "brave",
		Enabled: true,
		Settings: SESettings{
			Shortcut: "br",
			Crawlers: []structures.Engine{structures.Brave, structures.Google},
		},
	},
	{
		Name:    "duckduckgo",
		Enabled: true,
		Settings: SESettings{
			Shortcut: "ddg",
			Crawlers: []structures.Engine{structures.Bing},
		},
	},
	{
		Name:    "etools",
		Enabled: true,
		Settings: SESettings{
			Shortcut: "ets",
			Crawlers: []structures.Engine{structures.Bing, structures.Google, structures.Mojeek, structures.Yandex},
		},
	},
	{
		Name:    "google",
		Enabled: true,
		Settings: SESettings{
			Shortcut: "go",
			Crawlers: []structures.Engine{structures.Google},
		},
	},
	{
		Name:    "mojeek",
		Enabled: true,
		Settings: SESettings{
			Shortcut: "mjk",
			Crawlers: []structures.Engine{structures.Mojeek},
		},
	},
	{
		Name:    "qwant",
		Enabled: true,
		Settings: SESettings{
			Shortcut: "qw",
			Crawlers: []structures.Engine{structures.Qwant, structures.Bing},
		},
	},
	{
		Name:    "startpage",
		Enabled: true,
		Settings: SESettings{
			Shortcut: "sp",
			Crawlers: []structures.Engine{structures.Google},
		},
	},
	{
		Name:    "swisscows",
		Enabled: true,
		Settings: SESettings{
			Shortcut: "sc",
			Crawlers: []structures.Engine{structures.Bing},
		},
	},
}

func addDefaults(config *Config) {
	// could be be algorithmicaly optimized by sorting both slices (defaultEngines could be presorted) and performing a merge (from merge sort)
	// questionable performance increase though
	for _, defEng := range engineDefaults {
		if defEng.Enabled {
			hasInside := false
			for _, confEng := range config.Engines {
				if defEng.Name == confEng.Name {
					hasInside = true
				}
			}
			if !hasInside {
				config.Engines = append(config.Engines, defEng)
			}
		}
	}
}

func SetupConfig(path string, name string) *Config {
	// Use "." as the key path delimiter. This can be "/" or any character.
	k := koanf.New(".")

	// Load default values using the structs provider.
	// We provide a struct along with the struct tag `koanf` to the
	// provider.
	k.Load(structs.Provider(Config{Engines: engineDefaults}, "koanf"), nil)

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

	// Add Defaults
	addDefaults(&config)

	// Add enabled engines names and remove disabled ones
	tmpEnabled := make([]Engine, 0, len(config.Engines))
	for _, engine := range config.Engines {
		if engine.Enabled {
			EnabledEngines = append(EnabledEngines, structures.Engine(engine.Name))
			tmpEnabled = append(tmpEnabled, engine)
		}
	}
	config.Engines = tmpEnabled

	return &config
}

const InsertDefaultRank bool = true       // this should be moved to config
const LogDumpLocation string = "logdump/" // this should be moved to config
