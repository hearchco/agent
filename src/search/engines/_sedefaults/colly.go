package _sedefaults

import (
	"context"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

func colRequest(col *colly.Collector, ctx context.Context, seName engines.Name, saveOrigUrl bool) {
	col.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			if engines.IsTimeoutError(err) {
				log.Trace().
					Caller().
					Err(err).
					Str("engine", seName.String()).
					Msg("Context timeout error")
			} else {
				log.Error().
					Caller().
					Err(err).
					Str("engine", seName.String()).
					Msg("Context error")
			}
			r.Abort()
			return
		}
		if saveOrigUrl {
			r.Ctx.Put("originalURL", r.URL.String())
		}
	})
}

func colError(col *colly.Collector, seName engines.Name, visiting bool) {
	col.OnError(func(r *colly.Response, err error) {
		if engines.IsTimeoutError(err) {
			log.Trace().
				// Err(err). // timeout error produces Get "url" error with the query
				Str("engine", seName.String()).
				// Str("url", urll). // can't reliably anonymize it (because it's engine dependent)
				Msg("_sedefaults.colError(): request timeout error for url")
		} else {
			event := log.Error()
			if visiting {
				event = log.Trace()
			}
			event.
				Caller().
				Err(err).
				Str("engine", seName.String()).
				// Str("url", urll). // can't reliably anonymize it (because it's engine dependent)
				Int("statusCode", r.StatusCode).
				Bytes("response", r.Body). // WARN: query can be present, depending on the response from the engine (example: google has the query in 3 places)
				Msg("Request error for url")
		}
	})
}

func pagesColResponse(pagesCol *colly.Collector, seName engines.Name, relay *bucket.Relay) {
	pagesCol.OnResponse(func(r *colly.Response) {
		urll := r.Ctx.Get("originalURL")
		if urll == "" {
			log.Error().
				Caller().
				Msg("Error getting original url")
			return
		}

		err := bucket.SetResultResponse(urll, r, relay, seName)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Msg("Error setting result")
		}
	})
}
