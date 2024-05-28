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
			"Accept":             "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"Accept-Encoding":    "gzip, deflate, br, zstd", // Chromium-based browsers have "zstd" but that isn't supported by Firefox nor Safari
			"Accept-Language":    "en-US,en;q=0.9",
			"Sec-Ch-Ua":          secChUa, // "Google Chrome";v="119", "Chromium";v="119", "Not=A?Brand";v="24"
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": "\"Windows\"",
			"Sec-Fetch-Dest":     "document",
			"Sec-Fetch-Mode":     "navigate",
			"Sec-Fetch-Site":     "none",
		}),
	)
	pagesCol := colly.NewCollector(
		colly.Async(),
		colly.MaxDepth(1),
		colly.UserAgent(userAgent),
		colly.IgnoreRobotsTxt(),
		colly.Headers(map[string]string{
			"Accept":             "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"Accept-Encoding":    "gzip, deflate, br, zstd", // Chromium-based browsers have "zstd" but that isn't supported by Firefox nor Safari
			"Accept-Language":    "en-US,en;q=0.9",
			"Sec-Ch-Ua":          secChUa, // "Google Chrome";v="119", "Chromium";v="119", "Not=A?Brand";v="24"
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": "\"Windows\"",
			"Sec-Fetch-Dest":     "document",
			"Sec-Fetch-Mode":     "navigate",
			"Sec-Fetch-Site":     "none",
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
			Caller().
			Err(err).
			Str("limitRule", fmt.Sprintf("%v", limitRule)).
			Msg("Failed adding new limit rule")
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
				Caller().
				Err(err).
				Strs("proxies", settings.Proxies).
				Msg("Failed creating proxy switcher")
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
