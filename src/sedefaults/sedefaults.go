package sedefaults

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/hearchco/hearchco/src/bucket"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/engines"
	"github.com/hearchco/hearchco/src/search/useragent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func PagesColRequest(seName engines.Name, pagesCol *colly.Collector, ctx context.Context) {
	pagesCol.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			if engines.IsTimeoutError(err) {
				log.Trace().
					Err(err).
					Str("engine", seName.String()).
					Msg("sedefaults.PagesColRequest() -> pagesCol.OnRequest(): context timeout error")
			} else {
				log.Error().
					Err(err).
					Str("engine", seName.String()).
					Msg("sedefaults.PagesColRequest() -> pagesCol.OnRequest(): context error")
			}
			r.Abort()
			return
		}
		r.Ctx.Put("originalURL", r.URL.String())
	})
}

func PagesColError(seName engines.Name, pagesCol *colly.Collector) {
	pagesCol.OnError(func(r *colly.Response, err error) {
		urll := r.Ctx.Get("originalURL")
		if engines.IsTimeoutError(err) {
			log.Trace().
				Err(err).
				Str("engine", seName.String()).
				Str("url", urll).
				Msg("sedefaults.PagesColError() -> pagesCol.OnError(): request timeout error for url")
		} else {
			log.Trace().
				Err(err).
				Str("engine", seName.String()).
				Str("url", urll).
				Str("response", string(r.Body)).
				Msg("sedefaults.PagesColError() -> pagesCol.OnError(): request error for url")
		}
	})
}

func PagesColResponse(seName engines.Name, pagesCol *colly.Collector, relay *bucket.Relay) {
	pagesCol.OnResponse(func(r *colly.Response) {
		urll := r.Ctx.Get("originalURL")
		err := bucket.SetResultResponse(urll, r, relay, seName)
		if err != nil {
			log.Error().Err(err).Msg("sedefaults.PagesColResponse(): error setting result")
		}
	})
}

func ColRequest(seName engines.Name, col *colly.Collector, ctx context.Context) {
	col.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			if engines.IsTimeoutError(err) {
				log.Trace().
					Err(err).
					Str("engine", seName.String()).
					Msg("sedefaults.ColRequest() -> col.OnRequest(): context timeout error")
			} else {
				log.Error().
					Err(err).
					Str("engine", seName.String()).
					Msg("sedefaults.ColRequest() -> col.OnRequest(): context error")
			}
			r.Abort()
			return
		}
	})
}

func ColError(seName engines.Name, col *colly.Collector) {
	col.OnError(func(r *colly.Response, err error) {
		urll := r.Request.URL.String()
		if engines.IsTimeoutError(err) {
			log.Trace().
				Err(err).
				Str("engine", seName.String()).
				Str("url", urll).
				Msg("sedefaults.ColError() -> col.OnError(): request timeout error for url")
		} else {
			log.Error().
				Err(err).
				Str("engine", seName.String()).
				Str("url", urll).
				Int("statusCode", r.StatusCode).
				Str("response", string(r.Body)).
				Msg("sedefaults.ColError() -> col.OnError(): request error for url")

			dumpPath := fmt.Sprintf("%v%v_col.log.html", config.LogDumpLocation, seName.String())
			log.Debug().
				Str("engine", seName.String()).
				Str("responsePath", dumpPath).
				Func(func(e *zerolog.Event) {
					bodyWriteErr := os.WriteFile(dumpPath, r.Body, 0644)
					if bodyWriteErr != nil {
						log.Error().
							Err(bodyWriteErr).
							Str("engine", seName.String()).
							Msg("sedefaults.ColError() -> col.OnError(): error writing html response body to file")
					}
				}).
				Msg("sedefaults.ColError() -> col.OnError(): html response written")
		}
	})
}

func Prepare(seName engines.Name, options *engines.Options, settings *config.Settings, support *engines.SupportedSettings, info *engines.Info, ctx *context.Context) error {
	if ctx == nil {
		*ctx = context.Background()
	}

	if options.UserAgent == "" {
		options.UserAgent = useragent.RandomUserAgent()
	}
	log.Trace().
		Str("engine", seName.String()).
		Str("userAgent", options.UserAgent).
		Msg("Prepare")

	// TODO: move to config.SetupConfig
	if settings.RequestedResultsPerPage != 0 && !support.RequestedResultsPerPage {
		log.Panic().
			Str("engine", seName.String()).
			Int("requestedResultsPerPage", settings.RequestedResultsPerPage).
			Msg("sedefaults.Prepare(): setting not supported by engine")
		// ^PANIC
	}
	if settings.RequestedResultsPerPage == 0 && support.RequestedResultsPerPage {
		// if its used in the code but not set, give it the default value
		settings.RequestedResultsPerPage = info.ResultsPerPage
	}

	if options.Mobile && !support.Mobile {
		options.Mobile = false // this line shouldn't matter [1]
		log.Debug().
			Str("engine", seName.String()).
			Bool("mobile", options.Mobile).
			Msg("Mobile set but not supported")
	}

	if options.Locale != "" && !support.Locale {
		options.Locale = config.DefaultLocale // [1]
		log.Debug().
			Str("engine", seName.String()).
			Str("locale", options.Locale).
			Msg("Locale set but not supported")
	}

	if options.Locale == "" && support.Locale {
		options.Locale = config.DefaultLocale
	}

	if options.SafeSearch && !support.SafeSearch {
		options.SafeSearch = false // [1]
		log.Debug().
			Str("engine", seName.String()).
			Bool("safeSearch", options.SafeSearch).
			Msg("SafeSearch set but not supported")
	}

	return nil
}

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
				Msg("sedefaults.InitializeCollectors(): failed adding new limit rule")
		}

		if timings.Timeout != 0 {
			(*colPtr).SetRequestTimeout(timings.Timeout)
		}

		if timings.PageTimeout != 0 {
			(*pagesColPtr).SetRequestTimeout(timings.PageTimeout)
		}
	}

	if options.Proxies != nil {
		log.Debug().
			Strs("proxies", options.Proxies).
			Msg("Using proxies")

		// Rotate proxies
		rp, err := proxy.RoundRobinProxySwitcher(options.Proxies...)
		if err != nil {
			log.Fatal().
				Err(err).
				Strs("proxies", options.Proxies).
				Msg("sedefaults.InitializeCollectors(): failed creating proxy switcher")
		}

		(*colPtr).SetProxyFunc(rp)
		(*pagesColPtr).SetProxyFunc(rp)
	}
}

func DoGetRequest(urll string, colCtx *colly.Context, collector *colly.Collector, packageName engines.Name, retError *error) {
	log.Trace().
		Str("engine", packageName.String()).
		Str("url", urll).
		Msg("GET")
	err := collector.Request("GET", urll, nil, colCtx, nil)
	if err != nil {
		*retError = fmt.Errorf("%v.Search(): failed GET request to %v with %w", packageName.ToLower(), urll, err)
	}
}

func DoPostRequest(urll string, requestData io.Reader, colCtx *colly.Context, collector *colly.Collector, packageName engines.Name, retError *error) {
	log.Trace().
		Str("engine", packageName.String()).
		Str("url", urll).
		Msg("POST")
	err := collector.Request("POST", urll, requestData, colCtx, nil)
	if err != nil {
		*retError = fmt.Errorf("%v.Search(): failed POST request to %v and body %v. error %w", packageName.ToLower(), requestData, urll, err)
	}
}

func PageFromContext(ctx *colly.Context, seName engines.Name) int {
	var pageStr string = ctx.Get("page")
	page, converr := strconv.Atoi(pageStr)
	if converr != nil {
		log.Panic().
			Err(converr).
			Str("engine", seName.String()).
			Str("page", pageStr).
			Msg("sedefaults.PageFromContext(): failed to convert page number to int")
		// ^PANIC
	}
	return page
}
