package sedefaults

import (
	"context"
	"os"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/search/limit"
	"github.com/tminaorg/brzaguza/src/search/useragent"
	"github.com/tminaorg/brzaguza/src/structures"
)

func PagesColRequest(seName string, pagesCol *colly.Collector, ctx *context.Context, retError *error) {
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

func PagesColError(seName string, pagesCol *colly.Collector) {
	pagesCol.OnError(func(r *colly.Response, err error) {
		log.Debug().Msgf("%v: Pages Collector - OnError.\nURL: %v\nError: %v", seName, r.Ctx.Get("originalURL"), err)
	})
}

func PagesColResponse(seName string, pagesCol *colly.Collector, relay *structures.Relay) {
	pagesCol.OnResponse(func(r *colly.Response) {
		urll := r.Ctx.Get("originalURL")
		bucket.SetResultResponse(urll, r, relay, seName)
	})
}

func ColRequest(seName string, col *colly.Collector, ctx *context.Context, retError *error) {
	col.OnRequest(func(r *colly.Request) {
		if err := (*ctx).Err(); err != nil { // dont fully understand this
			log.Error().Msgf("%v: SE Collector; Error OnRequest %v", seName, r)
			r.Abort()
			*retError = err
			return
		}
	})
}

func ColError(seName string, col *colly.Collector, retError *error) {
	col.OnError(func(r *colly.Response, err error) {
		log.Error().Msgf("%v: SE Collector - OnError.\nURL: %v\nError: %v", seName, r.Request.URL.String(), err)
		log.Error().Msgf("%v: HTML Response written to %v%v_col.log.html", seName, config.LogDumpLocation, seName)
		writeErr := os.WriteFile(config.LogDumpLocation+seName+"_col.log.html", r.Body, 0644)
		if writeErr != nil {
			log.Error().Err(writeErr)
		}
		*retError = err
	})
}

func FunctionPrepare(seName string, options *structures.Options, ctx *context.Context) error {
	if ctx == nil {
		*ctx = context.Background()
	} //^ not necessary as ctx is always passed in search.go, branch predictor will skip this if

	if err := limit.RateLimit.Wait(*ctx); err != nil {
		return err
	}

	if options.UserAgent == "" {
		options.UserAgent = useragent.RandomUserAgent()
	}
	log.Trace().Msgf("%v: UserAgent: %v", seName, options.UserAgent)

	return nil
}

func InitializeCollectors(colPtr **colly.Collector, pagesColPtr **colly.Collector, options *structures.Options, limitRule *colly.LimitRule) {
	if options.MaxPages == 1 {
		*colPtr = colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent)) // so there is no thread creation overhead
	} else {
		*colPtr = colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent), colly.Async())
	}
	*pagesColPtr = colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent), colly.Async())

	if limitRule != nil {
		(*colPtr).Limit(limitRule)
	}
}
