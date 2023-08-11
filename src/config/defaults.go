package config

import "github.com/tminaorg/brzaguza/src/structures"

const DefaultLocale string = "en-US"

var DefaultConfig Config = Config{
	Engines: map[structures.EngineName]Engine{
		"bing": {
			Enabled: true,
			Settings: SESettings{
				Shortcut: "bi",
			},
		},
		"brave": {
			Enabled: true,
			Settings: SESettings{
				Shortcut: "br",
			},
		},
		"duckduckgo": {
			Enabled: true,
			Settings: SESettings{
				Shortcut: "ddg",
			},
		},
		"etools": {
			Enabled: true,
			Settings: SESettings{
				Shortcut: "ets",
			},
		},
		"google": {
			Enabled: true,
			Settings: SESettings{
				Shortcut: "go",
			},
		},
		"mojeek": {
			Enabled: true,
			Settings: SESettings{
				Shortcut: "mjk",
			},
		},
		"qwant": {
			Enabled: true,
			Settings: SESettings{
				Shortcut: "qw",
			},
		},
		"startpage": {
			Enabled: true,
			Settings: SESettings{
				Shortcut: "sp",
			},
		},
		"swisscows": {
			Enabled: true,
			Settings: SESettings{
				Shortcut: "sc",
			},
		},
	},
}
