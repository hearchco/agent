package config

import (
	"time"

	"github.com/hearchco/agent/src/search/category"
	"github.com/hearchco/agent/src/utils/moretime"
)

func New() Config {
	return Config{
		Server: Server{
			Environment:  "normal",
			Port:         3030,
			FrontendUrls: []string{"http://localhost:5173"},
			Cache: Cache{
				Type:      "none",
				KeyPrefix: "HEARCHCO_",
				TTL: TTL{
					Time: moretime.Week,
				},
				Redis: Redis{
					Host: "localhost",
					Port: 6379,
				},
			},
			ImageProxy: ImageProxy{
				Timeouts: ImageProxyTimeouts{
					Dial:         3 * time.Second,
					KeepAlive:    3 * time.Second,
					TLSHandshake: 2 * time.Second,
				},
			},
		},
		Categories: map[category.Name]Category{
			category.GENERAL: {
				Engines:                  generalEngines,
				RequiredEngines:          generalRequiredEngines,
				RequiredByOriginEngines:  generalRequiredByOriginEngines,
				PreferredEngines:         generalPreferredEngines,
				PreferredByOriginEngines: generalPreferredByOriginEngines,
				Ranking:                  generalRanking(),
				Timings:                  generalTimings,
			},
			category.IMAGES: {
				Engines:                  imagesEngines,
				RequiredEngines:          imagesRequiredEngines,
				RequiredByOriginEngines:  imagesRequiredByOriginEngines,
				PreferredEngines:         imagesPreferredEngines,
				PreferredByOriginEngines: imagesPreferredByOriginEngines,
				Ranking:                  imagesRanking(),
				Timings:                  imagesTimings,
			},
			category.SCIENCE: {
				Engines:                  scienceEngines,
				RequiredEngines:          scienceRequiredEngines,
				RequiredByOriginEngines:  scienceRequiredByOriginEngines,
				PreferredEngines:         sciencePreferredEngines,
				PreferredByOriginEngines: sciencePreferredByOriginEngines,
				Ranking:                  scienceRanking(),
				Timings:                  scienceTimings,
			},
			category.THOROUGH: {
				Engines:                  thoroughEngines,
				RequiredEngines:          thoroughRequiredEngines,
				RequiredByOriginEngines:  thoroughRequiredByOriginEngines,
				PreferredEngines:         thoroughPreferredEngines,
				PreferredByOriginEngines: thoroughPreferredByOriginEngines,
				Ranking:                  thoroughRanking(),
				Timings:                  thoroughTimings,
			},
		},
	}
}
