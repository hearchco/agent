package _sedefaults

import (
	"context"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

func pagesColRequest(ctx context.Context, seName engines.Name, pagesCol *colly.Collector) {
	pagesCol.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			if engines.IsTimeoutError(err) {
				log.Trace().
					Err(err).
					Str("engine", seName.String()).
					Msg("_sedefaults.PagesColRequest() -> pagesCol.OnRequest(): context timeout error")
			} else {
				log.Error().
					Err(err).
					Str("engine", seName.String()).
					Msg("_sedefaults.PagesColRequest() -> pagesCol.OnRequest(): context error")
			}
			r.Abort()
			return
		}
		r.Ctx.Put("originalURL", r.URL.String())
	})
}

func pagesColError(seName engines.Name, pagesCol *colly.Collector) {
	pagesCol.OnError(func(r *colly.Response, err error) {
		urll := r.Ctx.Get("originalURL")
		if engines.IsTimeoutError(err) {
			log.Trace().
				Err(err).
				Str("engine", seName.String()).
				Str("url", urll).
				Msg("_sedefaults.PagesColError() -> pagesCol.OnError(): request timeout error for url")
		} else {
			log.Trace().
				Err(err).
				Str("engine", seName.String()).
				Str("url", urll).
				Str("response", string(r.Body)).
				Msg("_sedefaults.PagesColError() -> pagesCol.OnError(): request error for url")
		}
	})
}

func pagesColResponse(seName engines.Name, pagesCol *colly.Collector, relay *bucket.Relay) {
	pagesCol.OnResponse(func(r *colly.Response) {
		urll := r.Ctx.Get("originalURL")
		err := bucket.SetResultResponse(urll, r, relay, seName)
		if err != nil {
			log.Error().
				Err(err).
				Msg("_sedefaults.PagesColResponse(): error setting result")
		}
	})
}
