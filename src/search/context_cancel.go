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

// Hard timeout is associated with the required engines.
func cancelHardTimeout(start time.Time, cancel context.CancelFunc, query string, wgEngs *sync.WaitGroup, engs []engines.Name, wgByOriginEngs *sync.WaitGroup, byOriginEngs []engines.Name) {
	groupNames := [...]string{groupRequired, groupRequiredByOrigin}
	cancelTimeout(groupNames, start, cancel, query, wgEngs, engs, wgByOriginEngs, byOriginEngs)
}

// Preferred timeout is associated with the preferred engines.
func cancelPreferredTimeout(start time.Time, cancel context.CancelFunc, query string, wgEngs *sync.WaitGroup, engs []engines.Name, wgByOriginEngs *sync.WaitGroup, byOriginEngs []engines.Name) {
	groupNames := [...]string{groupPreferred, groupPreferredByOrigin}
	cancelTimeout(groupNames, start, cancel, query, wgEngs, engs, wgByOriginEngs, byOriginEngs)
}

// Cancel timeout for the provided engines.
func cancelTimeout(groupNames [2]string, start time.Time, cancel context.CancelFunc, query string, wgEngs *sync.WaitGroup, engs []engines.Name, wgByOriginEngs *sync.WaitGroup, byOriginEngs []engines.Name) {
	var wg sync.WaitGroup

	// Wait for all required engines to finish.
	wg.Add(1)
	go func() {
		defer wg.Done()
		wgEngs.Wait()
		log.Debug().
			Str("query", anonymize.String(query)).
			Str("group", groupNames[0]).
			Str("engines", fmt.Sprintf("%v", engs)).
			Dur("duration", time.Since(start)).
			Msg("Scraping group finished")
	}()

	// Wait for all required by origin engines to finish.
	wg.Add(1)
	go func() {
		defer wg.Done()
		wgByOriginEngs.Wait()
		log.Debug().
			Str("query", anonymize.String(query)).
			Str("group", groupNames[1]).
			Str("engines", fmt.Sprintf("%v", byOriginEngs)).
			Dur("duration", time.Since(start)).
			Msg("Scraping group finished")
	}()

	wg.Wait()
	cancel()
}
