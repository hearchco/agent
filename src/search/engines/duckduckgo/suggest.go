package duckduckgo

import (
	"fmt"
	"sync/atomic"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/hearchco/agent/src/utils/morestrings"
)

func (se Engine) Suggest(query string, locale options.Locale, sugChan chan result.SuggestionScraped) ([]error, bool) {
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
	combinedParams := morestrings.JoinNonEmpty("&", "&", sugParamType)

	urll := fmt.Sprintf("%v?q=%v%v", suggestURL, query, combinedParams)
	anonUrll := fmt.Sprintf("%v?q=%v%v", suggestURL, anonymize.String(query), combinedParams)
	err := se.Get(ctx, urll, anonUrll)
	if err != nil {
		retErrors = append(retErrors, err)
	}

	se.Wait()
	close(sugChan)
	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
}
