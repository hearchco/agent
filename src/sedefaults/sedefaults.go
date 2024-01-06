package sedefaults

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/bucket"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/engines"
	"github.com/hearchco/hearchco/src/search/useragent"
	"github.com/rs/zerolog/log"
)

func PagesColRequest(seName engines.Name, pagesCol *colly.Collector, ctx context.Context) {
	pagesCol.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			if engines.IsTimeoutError(err) {
				log.Trace().Err(err).Msgf("sedefaults.PagesColRequest() from %v -> pagesCol.OnRequest(): context timeout error", seName)
			} else {
				log.Error().Err(err).Msgf("sedefaults.PagesColRequest() from %v -> pagesCol.OnRequest(): context error", seName)
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
			log.Trace().Err(err).Msgf("sedefaults.PagesColError() from %v -> pagesCol.OnError(): request timeout error for %v", seName, urll)
		} else {
			log.Trace().Err(err).Msgf("sedefaults.PagesColError() from %v -> pagesCol.OnError(): request error for %v\nresponse: %v", seName, urll, r)
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
				log.Trace().Err(err).Msgf("sedefaults.ColRequest() from %v -> col.OnRequest(): context timeout error", seName)
			} else {
				log.Error().Err(err).Msgf("sedefaults.ColRequest() from %v -> col.OnRequest(): context error", seName)
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
			log.Trace().Err(err).Msgf("sedefaults.ColError() from %v -> col.OnError(): request timeout error for %v", seName, urll)
		} else {
			log.Error().Err(err).Msgf("sedefaults.ColError() from %v -> col.OnError(): request error for %v\nresponse(%v): %v", seName, urll, r.StatusCode, string(r.Body))
			log.Debug().Msgf("sedefaults.ColError() from %v -> col.OnError(): html response written to %v%v_col.log.html", seName, config.LogDumpLocation, seName)

			bodyWriteErr := os.WriteFile(config.LogDumpLocation+seName.String()+"_col.log.html", r.Body, 0644)
			if bodyWriteErr != nil {
				log.Error().Err(bodyWriteErr).Msgf("sedefaults.ColError() from %v -> col.OnError(): error writing html response body to file", seName)
			}
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
	log.Trace().Msgf("%v: UserAgent: %v", seName, options.UserAgent)

	// TODO: move to config.SetupConfig
	if settings.RequestedResultsPerPage != 0 && !support.RequestedResultsPerPage {
		log.Panic().Msgf("sedefaults.Prepare() from %v: setting not supported. variable settings.RequestedResultsPerPage is set in the config for %v. that setting is not supported for this search engine. the settings value is: %v", seName, seName, settings.RequestedResultsPerPage)
		// ^PANIC
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
			log.Error().Err(err).Msgf("sedefaults.InitializeCollectors(): failed adding new limit rule: %v", limitRule)
		}
		if timings.Timeout != 0 {
			(*colPtr).SetRequestTimeout(timings.Timeout)
		}
		if timings.PageTimeout != 0 {
			(*pagesColPtr).SetRequestTimeout(timings.PageTimeout)
		}
	}
}

func DoGetRequest(urll string, colCtx *colly.Context, collector *colly.Collector, packageName engines.Name, retError *error) {
	log.Trace().Msgf("%v GET: %v", strings.ToUpper(packageName.String()), urll)
	err := collector.Request("GET", urll, nil, colCtx, nil)
	if err != nil {
		*retError = fmt.Errorf("%v.Search(): failed GET request to %v with %w", packageName.ToLower(), urll, err)
	}
}

func DoPostRequest(urll string, requestData io.Reader, colCtx *colly.Context, collector *colly.Collector, packageName engines.Name, retError *error) {
	log.Trace().Msgf("%v POST: %v", strings.ToUpper(packageName.String()), urll)
	err := collector.Request("POST", urll, requestData, colCtx, nil)
	if err != nil {
		*retError = fmt.Errorf("%v.Search(): failed POST request to %v and body %v. error %w", packageName.ToLower(), requestData, urll, err)
	}
}

func PageFromContext(ctx *colly.Context, seName engines.Name) int {
	var pageStr string = ctx.Get("page")
	page, converr := strconv.Atoi(pageStr)
	if converr != nil {
		log.Panic().Err(converr).Msgf("sedefaults.PageFromContext from %v: failed to convert page number to int. pageStr: %v", seName, pageStr)
		// ^PANIC
	}
	return page
}
