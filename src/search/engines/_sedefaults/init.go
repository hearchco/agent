package _sedefaults

import (
	"context"
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/useragent"
	"github.com/rs/zerolog/log"
)

func InitializeCollectors(ctx context.Context, engineName engines.Name, options engines.Options, settings config.Settings, timings config.CategoryTimings, relay *bucket.Relay) (*colly.Collector, *colly.Collector) {
	// get random user agent and corresponding Sec-Ch-Ua header
	userAgent, secChUa := useragent.RandomUserAgentWithHeader()

	// create collectors
	col := colly.NewCollector(
		colly.Async(),
		colly.MaxDepth(1),
		colly.UserAgent(userAgent),
		colly.IgnoreRobotsTxt(),
		colly.Headers(map[string]string{
			"Sec-Ch-Ua": secChUa,
		}),
	)
	pagesCol := colly.NewCollector(
		colly.Async(),
		colly.MaxDepth(1),
		colly.UserAgent(userAgent),
		colly.IgnoreRobotsTxt(),
		colly.Headers(map[string]string{
			"Sec-Ch-Ua": secChUa,
		}),
	)

	// set collector limit rules
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

	// set collector proxies
	if settings.Proxies != nil {
		log.Debug().
			Strs("proxies", settings.Proxies).
			Msg("Using proxies")

		// rotate proxies
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

	// set up collector
	colRequest(col, ctx, engineName, false)
	colError(col, engineName, false)

	// set up pages collector
	colRequest(pagesCol, ctx, engineName, true)
	colError(pagesCol, engineName, true)
	pagesColResponse(pagesCol, engineName, relay)

	return col, pagesCol
}
