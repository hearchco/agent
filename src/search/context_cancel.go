package search

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/utils/anonymize"
)

func cancelHardTimeout(start time.Time, cancel context.CancelFunc, query string, wgRequiredEngines *sync.WaitGroup, requiredEngines []engines.Name, wgRequiredByOriginEngines *sync.WaitGroup, requiredByOriginEngines []engines.Name) {
	var wg sync.WaitGroup
	wg.Add(2)

	// Wait for all required engines to finish.
	go func() {
		defer wg.Done()
		wgRequiredEngines.Wait()
		log.Debug().
			Str("query", anonymize.String(query)).
			Str("group", "required").
			Str("engines", fmt.Sprintf("%v", requiredEngines)).
			Dur("duration", time.Since(start)).
			Msg("Scraping group finished")
	}()

	// Wait for all required by origin engines to finish.
	go func() {
		defer wg.Done()
		wgRequiredByOriginEngines.Wait()
		log.Debug().
			Str("query", anonymize.String(query)).
			Str("group", "required by origin").
			Str("engines", fmt.Sprintf("%v", requiredByOriginEngines)).
			Dur("duration", time.Since(start)).
			Msg("Scraping group finished")
	}()

	wg.Wait()
	cancel()
}

func cancelPreferredTimeout(start time.Time, cancel context.CancelFunc, query string, wgPreferredEngines *sync.WaitGroup, preferredEngines []engines.Name, wgPreferredByOriginEngines *sync.WaitGroup, preferredByOriginEngines []engines.Name) {
	var wg sync.WaitGroup

	// Wait for all preferred engines to finish.
	wg.Add(1)
	go func() {
		defer wg.Done()
		wgPreferredEngines.Wait()
		log.Debug().
			Str("query", anonymize.String(query)).
			Str("group", "preferred").
			Str("engines", fmt.Sprintf("%v", preferredEngines)).
			Dur("duration", time.Since(start)).
			Msg("Scraping group finished")
	}()

	// Wait for all preferred by origin engines to finish.
	wg.Add(1)
	go func() {
		defer wg.Done()
		wgPreferredByOriginEngines.Wait()
		log.Debug().
			Str("query", anonymize.String(query)).
			Str("group", "preferred by origin").
			Str("engines", fmt.Sprintf("%v", preferredByOriginEngines)).
			Dur("duration", time.Since(start)).
			Msg("Scraping group finished")
	}()

	// Wait for both groups to finish.
	wg.Wait()
	cancel()
}
