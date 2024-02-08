package _sedefaults

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

func InitializeCollectors(colPtr **colly.Collector, pagesColPtr **colly.Collector, settings *config.Settings, options *engines.Options, timings *config.Timings) {
	*colPtr = colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent), colly.Async())
	*pagesColPtr = colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent), colly.Async())

	if timings != nil {
		var limitRule *colly.LimitRule = &colly.LimitRule{
			DomainGlob:  "*",
			Delay:       timings.Delay,
			RandomDelay: timings.RandomDelay,
			Parallelism: timings.Parallelism,
		}

		if err := (*colPtr).Limit(limitRule); err != nil {
			log.Error().
				Err(err).
				Str("limitRule", fmt.Sprintf("%v", limitRule)).
				Msg("_sedefaults.InitializeCollectors(): failed adding new limit rule")
		}

		if timings.Timeout != 0 {
			(*colPtr).SetRequestTimeout(timings.Timeout)
		}

		if timings.PageTimeout != 0 {
			(*pagesColPtr).SetRequestTimeout(timings.PageTimeout)
		}
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

		(*colPtr).SetProxyFunc(rp)
		(*pagesColPtr).SetProxyFunc(rp)
	}

	if settings.InsecureSkipVerify {
		tp := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		(*colPtr).WithTransport(tp)
		(*pagesColPtr).WithTransport(tp)
	}
}
