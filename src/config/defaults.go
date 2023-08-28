package config

import "github.com/tminaorg/brzaguza/src/engines"

const DefaultLocale string = "en-US"

var DefaultConfig Config = Config{
	Engines: map[engines.Name]Engine{
		engines.Bing: {
			Enabled: true,
			Settings: Settings{
				Shortcut: "bi",
			},
		},
		engines.Brave: {
			Enabled: true,
			Settings: Settings{
				Shortcut: "br",
			},
		},
		engines.DuckDuckGo: {
			Enabled: true,
			Settings: Settings{
				Shortcut: "ddg",
			},
		},
		engines.Etools: {
			Enabled: true,
			Settings: Settings{
				Shortcut: "ets",
			},
		},
		engines.Google: {
			Enabled: true,
			Settings: Settings{
				Shortcut: "go",
			},
		},
		engines.Mojeek: {
			Enabled: true,
			Settings: Settings{
				Shortcut: "mjk",
			},
		},
		engines.Qwant: {
			Enabled: true,
			Settings: Settings{
				Shortcut: "qw",
			},
		},
		engines.Startpage: {
			Enabled: true,
			Settings: Settings{
				Shortcut: "sp",
			},
		},
		engines.Swisscows: {
			Enabled: true,
			Settings: Settings{
				Shortcut: "sc",
			},
		},
	},
}
