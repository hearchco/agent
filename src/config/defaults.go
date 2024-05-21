package config

import (
	"time"

	"github.com/hearchco/hearchco/src/moretime"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

const DefaultLocale string = "en_US"

func EmptyRanking() CategoryRanking {
	rnk := CategoryRanking{
		REXP:    0.5,
		A:       1,
		B:       0,
		C:       1,
		D:       0,
		TRA:     1,
		TRB:     0,
		TRC:     1,
		TRD:     0,
		Engines: map[string]CategoryEngineRanking{},
	}

	for _, eng := range engines.Names() {
		rnk.Engines[eng.ToLower()] = CategoryEngineRanking{
			Mul:   1,
			Const: 0,
		}
	}

	return rnk
}

func NewRanking() CategoryRanking {
	return EmptyRanking()
}

func NewSettings() map[engines.Name]Settings {
	mp := map[engines.Name]Settings{
		engines.BING: {
			Shortcut: "bi",
		},
		engines.BINGIMAGES: {
			Shortcut: "biimg",
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
		engines.GOOGLEIMAGES: {
			Shortcut: "gimg",
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
			log.Fatal().
				Str("engine", eng.String()).
				Msg("config.NewSettings(): no shortcut set")
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

func NewImage() []engines.Name {
	return []engines.Name{
		engines.BINGIMAGES,
		engines.GOOGLEIMAGES,
	}
}

func NewQuick() []engines.Name {
	return []engines.Name{
		engines.BING,
		engines.BRAVE,
		engines.DUCKDUCKGO,
		engines.GOOGLE,
		engines.MOJEEK,
	}
}

func NewScience() []engines.Name {
	return []engines.Name{
		engines.GOOGLESCHOLAR,
	}
}

func New() Config {
	return Config{
		Server: Server{
			Environment:  "normal",
			Port:         3030,
			FrontendUrls: []string{"http://localhost:5173"},
			Cache: Cache{
				Type: "badger",
				TTL: TTL{
					Time:        moretime.Week,
					RefreshTime: 3 * moretime.Day,
				},
				Redis: Redis{
					Host: "localhost",
					Port: 6379,
				},
			},
			Proxy: ImageProxy{
				Timeouts: ImageProxyTimeouts{
					Dial:         3 * time.Second,
					KeepAlive:    3 * time.Second,
					TLSHandshake: 2 * time.Second,
				},
			},
		},
		Settings: NewSettings(),
		Categories: map[category.Name]Category{
			category.GENERAL: {
				Engines: NewGeneral(),
				Ranking: NewRanking(),
				Timings: CategoryTimings{
					PreferredTimeoutMin:    1 * time.Second,
					PreferredTimeoutMax:    2 * time.Second,
					PreferredResultsNumber: 20,
					StepTime:               50 * time.Millisecond,
					MinimumResultsNumber:   10,
					HardTimeout:            3 * time.Second,
				},
			},
			category.IMAGES: {
				Engines: NewImage(),
				Ranking: NewRanking(),
				Timings: CategoryTimings{
					PreferredTimeoutMin:    1 * time.Second,
					PreferredTimeoutMax:    2 * time.Second,
					PreferredResultsNumber: 40,
					StepTime:               100 * time.Millisecond,
					MinimumResultsNumber:   20,
					HardTimeout:            3 * time.Second,
				},
			},
			category.SCIENCE: {
				Engines: NewScience(),
				Ranking: NewRanking(),
				Timings: CategoryTimings{
					PreferredTimeoutMin:    1 * time.Second,
					PreferredTimeoutMax:    2 * time.Second,
					PreferredResultsNumber: 10,
					StepTime:               100 * time.Millisecond,
					MinimumResultsNumber:   5,
					HardTimeout:            3 * time.Second,
				},
			},
			category.QUICK: {
				Engines: NewQuick(),
				Ranking: NewRanking(),
				Timings: CategoryTimings{
					PreferredTimeoutMin:    500 * time.Millisecond,
					PreferredTimeoutMax:    1500 * time.Millisecond,
					PreferredResultsNumber: 10,
					StepTime:               25 * time.Millisecond,
					MinimumResultsNumber:   5,
					HardTimeout:            3 * time.Second,
				},
			},
			category.THOROUGH: {
				Engines: NewGeneral(),
				Ranking: NewRanking(),
				Timings: CategoryTimings{
					PreferredTimeoutMin:    1 * time.Second,
					PreferredTimeoutMax:    4 * time.Second,
					PreferredResultsNumber: 70,
					StepTime:               100 * time.Millisecond,
					MinimumResultsNumber:   50,
					HardTimeout:            5 * time.Second,
				},
			},
		},
	}
}
