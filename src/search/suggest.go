package search

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/utils/anonymize"
)

func Suggest(query string, locale options.Locale, hardTimeout time.Duration) ([]result.Suggestion, error) {
	// Capture start time.
	startTime := time.Now()

	if err := validateSuggestParams(query, locale); err != nil {
		return nil, err
	}

	log.Debug().
		Str("query", anonymize.String(query)).
		Str("locale", locale.String()).
		Msg("Suggesting")

	// Create context with timeout for HardTimeout.
	suggestCtx, cancelSuggestFunc := context.WithTimeout(context.Background(), hardTimeout)
	defer cancelSuggestFunc()

	// Initialize each engine.
	// TODO: Make enabled engines and timing configurable.
	enabledEngines := []engines.Name{engines.DUCKDUCKGO, engines.GOOGLE}
	suggesters := initializeSuggesters(suggestCtx, enabledEngines, config.CategoryTimings{
		HardTimeout: hardTimeout,
	})

	// Create a channel of channels to receive the suggestions from each engine.
	engChan := make(chan chan result.SuggestionScraped, len(suggesters))

	// Create a map for the suggestions with RWMutex.
	sugMap := result.SuggestionMap(len(suggesters))

	// Start a goroutine to receive the suggestions from each engine and add them to suggestions map.
	go createSuggestionsReceiver(engChan, &sugMap, len(suggesters))

	// Run the suggesters, cancellin the context when one engine finishes successfully or all engines finish (successful or not).
	runSuggestionsEngines(suggesters, cancelSuggestFunc, query, locale, enabledEngines, engChan)

	// Close the channel of channels (it's safe because each sending already happened sequentially).
	close(engChan)

	// Wait for the suggesters to finish, either by success or by timeout.
	<-suggestCtx.Done()

	suggestions, responders := sugMap.ExtractSuggestionsAndResponders()

	log.Debug().
		Int("suggestions", len(suggestions)).
		Str("query", anonymize.String(query)).
		Str("responders", fmt.Sprintf("%v", responders)).
		Dur("duration", time.Since(startTime)).
		Msg("Scraping finished")

	// Return the suggestions.
	return suggestions, nil
}
