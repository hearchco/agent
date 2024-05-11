package _sedefaults

import (
	"context"
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

// it's okay to return pointers to collectors since colly.NewCollector() returns a pointer
func InitializeCollectors(ctx context.Context, engineName engines.Name, options engines.Options, settings config.Settings, timings config.CategoryTimings, relay *bucket.Relay) (*colly.Collector, *colly.Collector) {
	col := colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent), colly.Async())
	pagesCol := colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent), colly.Async())

	limitRule := colly.LimitRule{
		DomainGlob:  "*",
		Delay:       timings.Delay,
		RandomDelay: timings.RandomDelay,
		Parallelism: timings.Parallelism,
	}
	if err := col.Limit(&limitRule); err != nil {
		log.Error().
			Err(err).
			Str("limitRule", fmt.Sprintf("%v", limitRule)).
			Msg("_sedefaults.InitializeCollectors(): failed adding new limit rule")
	}

	if timings.Timeout != 0 {
		col.SetRequestTimeout(timings.Timeout)
	}

	if timings.PageTimeout != 0 {
		pagesCol.SetRequestTimeout(timings.PageTimeout)
	}

	if settings.Proxies != nil {
		log.Debug().
			Strs("proxies", settings.Proxies).
			Msg("Using proxies")

		// Rotate proxies
		rp, err := proxy.RoundRobinProxySwitcher(settings.Proxies...)
		if err != nil {
			log.Fatal().
				Err(err).
				Strs("proxies", settings.Proxies).
				Msg("_sedefaults.InitializeCollectors(): failed creating proxy switcher")
		}

		col.SetProxyFunc(rp)
		pagesCol.SetProxyFunc(rp)
	}

	// Set up collector
	colRequest(col, ctx, engineName, false)
	colError(col, engineName, false)

	// Set up pages collector
	colRequest(pagesCol, ctx, engineName, true)
	colError(pagesCol, engineName, true)
	pagesColResponse(pagesCol, engineName, relay)

	return col, pagesCol
}
