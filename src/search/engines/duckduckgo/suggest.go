package duckduckgo

import (
	"sync/atomic"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/hearchco/agent/src/utils/moreurls"
	"github.com/hearchco/agent/src/utils/moreurls/parameters"
)

func (se Engine) Suggest(query string, options options.Options, sugChan chan result.SuggestionScraped) ([]error, bool) {
	foundResults := atomic.Bool{}
	retErrors := make([]error, 0, 1)

	se.OnResponse(func(e *colly.Response) {
		log.Trace().
			Caller().
			Bytes("body", e.Body).
			Msg("Got response")

		suggs, err := scraper.SuggestRespToSuggestions(e.Body)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Bytes("body", e.Body).
				Msg("Failed to convert response to suggestions")
		} else {
			log.Trace().
				Caller().
				Str("engine", se.Name.String()).
				Strs("suggestions", suggs).
				Msg("Sending suggestions to channel")
			for i, sug := range suggs {
				sugChan <- result.NewSuggestionScraped(sug, se.Name, i+1)
			}
			if !foundResults.Load() {
				foundResults.Store(true)
			}
		}
	})

	ctx := colly.NewContext()

	// Build the parameters.
	params := parameters.NewParams(
		paramQueryK, query,
		sugParamTypeK, sugParamTypeV,
	)

	// Build the url.
	urll := moreurls.Build(suggestURL, params)

	// Build anonymous url, by anonymizing the query.
	params.Set(paramQueryK, anonymize.String(query))
	anonUrll := moreurls.Build(suggestURL, params)

	// Send the request.
	if err := se.Get(ctx, urll, anonUrll); err != nil {
		retErrors = append(retErrors, err)
	}

	se.Wait()
	close(sugChan)
	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
}
