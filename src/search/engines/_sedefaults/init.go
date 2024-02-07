package _sedefaults

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

func InitializeCollectors(colPtr **colly.Collector, pagesColPtr **colly.Collector, options *engines.Options, timings *config.Timings) {
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
}
