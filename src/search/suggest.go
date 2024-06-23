package search

import (
	"context"
	"fmt"
	"sync"
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
	ctxHardTimeout, cancelHardTimeoutFunc := context.WithTimeout(context.Background(), hardTimeout)
	defer cancelHardTimeoutFunc()

	// Create a condition signal that indicates when the suggestions are received from any engine.
	received := sync.Cond{L: &sync.Mutex{}}

	// Create a context that cancels when both HardTimeout and received signal are done.
	suggestCtx, cancelSuggest := context.WithCancel(context.Background())
	defer cancelSuggest()
	go func() {
		<-ctxHardTimeout.Done()
		received.L.Lock()
		received.Wait()
		received.L.Unlock()
		cancelSuggest()
	}()

	// Initialize each engine.
	// TODO: Make enabled engines and timing configurable.
	enabledEngines := []engines.Name{engines.DUCKDUCKGO, engines.GOOGLE}
	suggesters := initializeSuggesters(suggestCtx, enabledEngines, config.CategoryTimings{
		HardTimeout: hardTimeout,
	})

	// Create a channel of channels to receive the suggestions from each engine.
	engChan := make(chan chan []result.SuggestionScraped, len(suggesters))

	// Create a map for the suggestions with RWMutex.
	sugMap := result.SuggestionMap(len(suggesters))

	// Start a goroutine to receive the suggestions from each engine and add them to suggestions map.
	go createSuggestionsReceiver(&received, engChan, &sugMap, len(suggesters))

	// Run the suggesters, cancellin the context when one engine finishes successfully or all engines finish (successful or not).
	runSuggestionsEngines(suggesters, cancelHardTimeoutFunc, query, locale, enabledEngines, engChan)

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
