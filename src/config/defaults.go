package config

import (
	"log"
	"time"

	"github.com/hearchco/hearchco/src/moretime"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
)

const DefaultLocale string = "en_US"

func EmptyRanking() Ranking {
	rnk := Ranking{
		REXP:    0.5,
		A:       1,
		B:       0,
		C:       1,
		D:       0,
		TRA:     1,
		TRB:     0,
		TRC:     1,
		TRD:     0,
		Engines: map[string]EngineRanking{},
	}

	for _, eng := range engines.Names() {
		rnk.Engines[eng.ToLower()] = EngineRanking{
			Mul:   1,
			Const: 0,
		}
	}

	return rnk
}

func NewRanking() Ranking {
	return EmptyRanking()
}

func NewSettings() map[engines.Name]Settings {
	mp := map[engines.Name]Settings{
		engines.BING: {
			Shortcut: "bi",
		},
		engines.BRAVE: {
			Shortcut: "br",
		},
		engines.DUCKDUCKGO: {
			Shortcut: "ddg",
		},
		engines.ETOOLS: {
			Shortcut: "ets",
		},
		engines.GOOGLE: {
			Shortcut: "g",
		},
		engines.GOOGLESCHOLAR: {
			Shortcut: "gs",
		},
		engines.MOJEEK: {
			Shortcut: "mjk",
		},
		engines.PRESEARCH: {
			Shortcut: "ps",
		},
		engines.QWANT: {
			Shortcut: "qw",
		},
		engines.STARTPAGE: {
			Shortcut: "sp",
		},
		engines.SWISSCOWS: {
			Shortcut: "sc",
		},
		engines.YAHOO: {
			Shortcut: "yh",
		},
		engines.YEP: {
			Shortcut: "yep",
		},
	}

	// Check if all search engines have a shortcut set
	for _, eng := range engines.Names() {
		if _, ok := mp[eng]; !ok {
			log.Fatalf("config.NewSettings(): %v doesn't have a shortcut set.", eng)
			// ^FATAL
		}
	}

	return mp
}

func NewAllEnabled() []engines.Name {
	return engines.Names()
}

func NewGeneral() []engines.Name {
	return []engines.Name{
		engines.BING,
		engines.BRAVE,
		engines.DUCKDUCKGO,
		engines.ETOOLS,
		engines.GOOGLE,
		engines.MOJEEK,
		engines.PRESEARCH,
		engines.QWANT,
		engines.STARTPAGE,
		engines.SWISSCOWS,
		engines.YAHOO,
		engines.YEP,
	}
}

func NewInfo() []engines.Name {
	return []engines.Name{
		engines.BING,
		engines.GOOGLE,
		engines.MOJEEK,
	}
}

func NewScience() []engines.Name {
	return []engines.Name{
		engines.GOOGLESCHOLAR,
	}
}

func New() *Config {
	return &Config{
		Server: Server{
			Port:        3030,
			FrontendUrl: "http://localhost:8000",
			Cache: Cache{
				Type: "badger",
				TTL: TTL{
					Time:        moretime.Week,
					RefreshTime: 3 * moretime.Day,
				},
				Badger: Badger{
					Persist: true,
				},
				Redis: Redis{
					Host: "localhost",
					Port: 6379,
				},
			},
		},
		Settings: NewSettings(),
		Categories: map[category.Name]Category{
			category.GENERAL: {
				Engines: NewGeneral(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.INFO: {
				Engines: NewInfo(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.SCIENCE: {
				Engines: NewScience(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     3000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.NEWS: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.BLOG: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     2500 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.SURF: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     2000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
			category.NEWNEWS: {
				Engines: NewAllEnabled(),
				Ranking: NewRanking(),
				Timings: Timings{
					Timeout:     1000 * time.Millisecond,
					PageTimeout: 1000 * time.Millisecond,
				},
			},
		},
	}
}
