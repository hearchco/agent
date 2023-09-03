package config

import "github.com/tminaorg/brzaguza/src/engines"

const DefaultLocale string = "en-US"

<<<<<<< HEAD
func New() *Config {
	return &Config{
		Engines: map[string]Engine{
			engines.Bing.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "bi",
				},
			},
			engines.Brave.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "br",
				},
			},
			engines.DuckDuckGo.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "ddg",
				},
			},
			engines.Etools.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "ets",
				},
			},
			engines.Google.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "go",
				},
			},
			engines.Mojeek.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "mjk",
				},
			},
			engines.Qwant.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "qw",
				},
			},
			engines.Startpage.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "sp",
				},
			},
			engines.Swisscows.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "sc",
				},
			},
			engines.Yep.ToLower(): {
				Enabled: true,
				Settings: Settings{
					Shortcut: "yep",
				},
			},
		},
	}
=======
var DefaultConfig Config = Config{
	Engines: map[string]Engine{
		engines.Bing.ToLower(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "bi",
			},
		},
		engines.Brave.ToLower(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "br",
			},
		},
		engines.DuckDuckGo.ToLower(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "ddg",
			},
		},
		engines.Etools.ToLower(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "ets",
			},
		},
		engines.Google.ToLower(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "go",
			},
		},
		engines.Mojeek.ToLower(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "mjk",
			},
		},
		engines.Qwant.ToLower(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "qw",
			},
		},
		engines.Startpage.ToLower(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "sp",
			},
		},
		engines.Swisscows.ToLower(): {
			Enabled: true,
			Settings: Settings{
				Shortcut: "sc",
			},
		},
	},
>>>>>>> 75bb49c (Using correct formating inside the code, while loading is lowercase)
}
