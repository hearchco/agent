package config

import "github.com/tminaorg/brzaguza/src/engines"

const DefaultLocale string = "en-US"

var DefaultConfig Config = Config{
	Engines: map[engines.Name]Engine{
		"bing": {
			Enabled: true,
			Settings: Settings{
				Shortcut: "bi",
			},
		},
		"brave": {
			Enabled: true,
			Settings: Settings{
				Shortcut: "br",
			},
		},
		"duckduckgo": {
			Enabled: true,
			Settings: Settings{
				Shortcut: "ddg",
			},
		},
		"etools": {
			Enabled: true,
			Settings: Settings{
				Shortcut: "ets",
			},
		},
		"google": {
			Enabled: true,
			Settings: Settings{
				Shortcut: "go",
			},
		},
		"mojeek": {
			Enabled: true,
			Settings: Settings{
				Shortcut: "mjk",
			},
		},
		"qwant": {
			Enabled: true,
			Settings: Settings{
				Shortcut: "qw",
			},
		},
		"startpage": {
			Enabled: true,
			Settings: Settings{
				Shortcut: "sp",
			},
		},
		"swisscows": {
			Enabled: true,
			Settings: Settings{
				Shortcut: "sc",
			},
		},
	},
}
