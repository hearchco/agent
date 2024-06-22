package google

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/hearchco/agent/src/utils/morestrings"
)

func (se Engine) Suggest(query string, locale options.Locale, sugChan chan result.SuggestionScraped) (error, bool) {
	foundResults := false

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
				Msg("Failed to convert response to suggestions")
		} else {
			for i, sug := range suggs {
				log.Trace().
					Caller().
					Int("index", i).
					Str("suggestion", sug).
					Msg("Sending suggestion to channel")
				sugChan <- result.NewSuggestionScraped(sug, se.Name, i)
				if !foundResults {
					foundResults = true
				}
			}
		}
	})

	ctx := colly.NewContext()
	combinedParams := morestrings.JoinNonEmpty("?", "&", sugParamClient)

	urll := fmt.Sprintf("%v%v&q=%v", suggestURL, combinedParams, query)
	anonUrll := fmt.Sprintf("%v?q=%v%v", suggestURL, combinedParams, anonymize.String(query))
	err := se.Get(ctx, urll, anonUrll)

	se.Wait()
	close(sugChan)
	return err, foundResults
}
