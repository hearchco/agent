package sedefaults

import (
	"context"
	"os"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/search/useragent"
)

func PagesColRequest(seName engines.Name, pagesCol *colly.Collector, ctx *context.Context, retError *error) {
	pagesCol.OnRequest(func(r *colly.Request) {
		if err := (*ctx).Err(); err != nil { // dont fully understand this
			log.Error().Msgf("%v: Pages Collector; Error OnRequest %v", seName, r)
			r.Abort()
			*retError = err
			return
		}
		r.Ctx.Put("originalURL", r.URL.String())
	})
}

func PagesColError(seName engines.Name, pagesCol *colly.Collector) {
	pagesCol.OnError(func(r *colly.Response, err error) {
		log.Debug().Msgf("%v: Pages Collector - OnError.\nURL: %v\nError: %v", seName, r.Ctx.Get("originalURL"), err)
	})
}

func PagesColResponse(seName engines.Name, pagesCol *colly.Collector, relay *bucket.Relay) {
	pagesCol.OnResponse(func(r *colly.Response) {
		urll := r.Ctx.Get("originalURL")
		bucket.SetResultResponse(urll, r, relay, seName)
	})
}

func ColRequest(seName engines.Name, col *colly.Collector, ctx *context.Context, retError *error) {
	col.OnRequest(func(r *colly.Request) {
		if err := (*ctx).Err(); err != nil {
			log.Error().Msgf("%v: SE Collector; Error OnRequest %v", seName, r)
			r.Abort()
			*retError = err
			return
		}
	})
}

func ColError(seName engines.Name, col *colly.Collector, retError *error) {
	col.OnError(func(r *colly.Response, err error) {
		log.Error().Err(err).Msgf("%v: SE Collector - OnError.\nURL: %v", seName, r.Request.URL.String())
		log.Debug().Msgf("%v: HTML Response written to %v%v_col.log.html", seName, config.LogDumpLocation, seName)
		writeErr := os.WriteFile(config.LogDumpLocation+string(seName)+"_col.log.html", r.Body, 0644)
		if writeErr != nil {
			log.Error().Err(writeErr)
		}
		*retError = err
	})
}

func Prepare(seName engines.Name, options *engines.Options, settings *config.Settings, support *engines.SupportedSettings, info *engines.Info, ctx *context.Context) error {
	if ctx == nil {
		*ctx = context.Background()
	} //^ not necessary as ctx is always passed in search.go, branch predictor will skip this if

	if options.UserAgent == "" {
		options.UserAgent = useragent.RandomUserAgent()
	}
	log.Trace().Msgf("%v: UserAgent: %v", seName, options.UserAgent)

	// These two ifs, could be moved to config.SetupConfig
	if settings.RequestedResultsPerPage != 0 && !support.RequestedResultsPerPage {
		log.Error().Msgf("%v: Variable settings.RequestedResultsPerPage is set, but not supported in this search engine. Its value is: %v", seName, settings.RequestedResultsPerPage)
		panic("sedefaults.Prepare(): Setting not supported.")
	}
	if settings.RequestedResultsPerPage == 0 && support.RequestedResultsPerPage {
		// If its used in the code but not set, give it the default value.
		settings.RequestedResultsPerPage = info.ResultsPerPage
	}

	if options.Mobile && !support.Mobile {
		options.Mobile = false // this line shouldn't matter [1]
		log.Debug().Msgf("%v: Mobile set but not supported. Value: %v", seName, options.Mobile)
	}
	if options.Locale != "" && !support.Locale {
		options.Locale = config.DefaultLocale // [1]
		log.Debug().Msgf("%v: Locale set but not supported. Value: %v", seName, options.Mobile)
	}
	if options.Locale == "" && support.Locale {
		options.Locale = config.DefaultLocale
	}
	if options.SafeSearch && !support.SafeSearch {
		options.SafeSearch = false // [1]
		log.Debug().Msgf("%v: SafeSearch set but not supported.", seName)
	}

	return nil
}

func InitializeCollectors(colPtr **colly.Collector, pagesColPtr **colly.Collector, options *engines.Options, timings *config.Timings) {
	if options.MaxPages == 1 {
		*colPtr = colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent)) // so there is no thread creation overhead
	} else {
		*colPtr = colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent), colly.Async())
	}
	*pagesColPtr = colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent), colly.Async())

	if timings != nil {
		var limitRule *colly.LimitRule = &colly.LimitRule{
			DomainGlob:  "*",
			Delay:       timings.Delay,
			RandomDelay: timings.RandomDelay,
			Parallelism: timings.Parallelism,
		}

		if err := (*colPtr).Limit(limitRule); err != nil {
			log.Error().Err(err).Msg("sedefaults: failed adding a new limit rule")
		}
		if timings.Timeout != 0 {
			(*colPtr).SetRequestTimeout(timings.Timeout)
		}
		if timings.PageTimeout != 0 {
			(*pagesColPtr).SetRequestTimeout(timings.PageTimeout)
		}
	}
}
