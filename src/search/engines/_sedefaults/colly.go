package _sedefaults

import (
	"context"
	"fmt"
	"os"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func colRequest(ctx context.Context, seName engines.Name, col *colly.Collector) {
	col.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			if engines.IsTimeoutError(err) {
				log.Trace().
					Err(err).
					Str("engine", seName.String()).
					Msg("_sedefaults.ColRequest() -> col.OnRequest(): context timeout error")
			} else {
				log.Error().
					Err(err).
					Str("engine", seName.String()).
					Msg("_sedefaults.ColRequest() -> col.OnRequest(): context error")
			}
			r.Abort()
			return
		}
	})
}

func colError(seName engines.Name, col *colly.Collector) {
	col.OnError(func(r *colly.Response, err error) {
		if engines.IsTimeoutError(err) {
			log.Trace().
				// Err(err). // timeout error produces Get "url" error with the query
				Str("engine", seName.String()).
				// Str("url", urll). // can't reliably anonymize it (because it's engine dependent and query isn't passed to this function)
				Msg("_sedefaults.ColError() -> col.OnError(): request timeout error for url")
		} else {
			log.Error().
				Err(err).
				Str("engine", seName.String()).
				// Str("url", urll). // can't reliably anonymize it (because it's engine dependent and query isn't passed to this function)
				Int("statusCode", r.StatusCode).
				Str("response", string(r.Body)). // query can be present, depending on the response from the engine (Google has the query in 3 places)
				Msg("_sedefaults.ColError() -> col.OnError(): request error for url")

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
							Msg("_sedefaults.ColError() -> col.OnError(): error writing html response body to file")
					}
				}).
				Msg("_sedefaults.ColError() -> col.OnError(): html response written")
		}
	})
}
