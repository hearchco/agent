package search

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/utils/anonymize"
)

func Suggest(query string, locale options.Locale) ([]string, error) {
	// Capture start time.
	startTime := time.Now()

	if err := validateSuggestParams(query, locale); err != nil {
		return nil, err
	}

	log.Debug().
		Str("query", anonymize.String(query)).
		Str("locale", locale.String()).
		Msg("Suggesting")

	// TODO: Implement the suggest function.
	suggestions := []string{}
	responders := []engines.Name{}

	log.Debug().
		Int("suggestions", len(suggestions)).
		Str("query", anonymize.String(query)).
		Str("responders", fmt.Sprintf("%v", responders)).
		Dur("duration", time.Since(startTime)).
		Msg("Scraping finished")

	// Return the suggestions.
	return suggestions, nil
}
