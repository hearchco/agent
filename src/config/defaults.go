package config

import "github.com/tminaorg/brzaguza/src/engines"

const DefaultLocale string = "en-US"

var DefaultConfig Config = Config{
	Engines: map[string]Engine{
		engines.Bing.String(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "bi",
			},
		},
		engines.Brave.String(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "br",
			},
		},
		engines.DuckDuckGo.String(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "ddg",
			},
		},
		engines.Etools.String(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "ets",
			},
		},
		engines.Google.String(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "go",
			},
		},
		engines.Mojeek.String(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "mjk",
			},
		},
		engines.Qwant.String(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "qw",
			},
		},
		engines.Startpage.String(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "sp",
			},
		},
		engines.Swisscows.String(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "sc",
			},
		},
	},
}
