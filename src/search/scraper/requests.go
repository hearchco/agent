package scraper

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
)

func (e EngineBase) Get(ctx *colly.Context, urll string, anonurll string) error {
	log.Trace().
		Str("engine", e.Name.String()).
		Str("url", anonurll).
		Str("method", http.MethodGet).
		Msg("Making a new request")

	if err := e.collector.Request(http.MethodGet, urll, nil, ctx, nil); err != nil {
		return fmt.Errorf("%v: failed GET request to %v with %w", e.Name.String(), anonurll, err)
	}

	return nil
}

func (e EngineBase) Post(ctx *colly.Context, urll string, body io.Reader, anonBody string) error {
	log.Trace().
		Str("engine", e.Name.String()).
		Str("url", urll).
		Str("body", anonBody).
		Str("method", http.MethodPost).
		Msg("Making a new request")

	if err := e.collector.Request(http.MethodPost, urll, body, ctx, nil); err != nil {
		return fmt.Errorf("%v: failed POST request to %v with %w", e.Name.String(), urll, err)
	}

	return nil
}
