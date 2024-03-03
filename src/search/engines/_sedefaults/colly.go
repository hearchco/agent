package _sedefaults

import (
	"context"
	"fmt"
	"os"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func colRequest(col *colly.Collector, ctx context.Context, seName engines.Name, saveOrigUrl bool) {
	col.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			if engines.IsTimeoutError(err) {
				log.Trace().
					Err(err).
					Str("engine", seName.String()).
					Msg("_sedefaults.colRequest(): context timeout error")
			} else {
				log.Error().
					Err(err).
					Str("engine", seName.String()).
					Msg("_sedefaults.colRequest(): context error")
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
				Err(err).
				Str("engine", seName.String()).
				// Str("url", urll). // can't reliably anonymize it (because it's engine dependent)
				Int("statusCode", r.StatusCode).
				Str("response", string(r.Body)). // WARN: query can be present, depending on the response from the engine (example: google has the query in 3 places)
				Msg("_sedefaults.colError(): request error for url")

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
							Msg("_sedefaults.colError(): error writing html response body to file")
					}
				}).
				Msg("_sedefaults.colError(): html response written")
		}
	})
}

func pagesColResponse(pagesCol *colly.Collector, seName engines.Name, relay *bucket.Relay) {
	pagesCol.OnResponse(func(r *colly.Response) {
		urll := r.Ctx.Get("originalURL")
		if urll == "" {
			log.Error().
				Msg("_sedefaults.pagesColResponse(): error getting original url")
		} else {
			err := bucket.SetResultResponse(urll, r, relay, seName)
			if err != nil {
				log.Error().
					Err(err).
					Msg("_sedefaults.pagesColResponse(): error setting result")
			}
		}
	})
}
