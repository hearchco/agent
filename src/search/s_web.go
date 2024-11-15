package search

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/category"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/utils/anonymize"
)

// Searches for web using the provided category config.
func Web(query string, opts options.Options, catConf category.Category) ([]result.Result, error) {
	// Capture start time.
	startTime := time.Now()

	if err := validateParams(query, opts); err != nil {
		return nil, err
	}

	log.Debug().
		Str("query", anonymize.String(query)).
		Int("pages_start", opts.Pages.Start).
		Int("pages_max", opts.Pages.Max).
		Str("locale", opts.Locale.String()).
		Bool("safesearch", opts.SafeSearch).
		Str("engines", fmt.Sprintf("%v", catConf.Engines)).
		Str("required_engines", fmt.Sprintf("%v", catConf.RequiredEngines)).
		Str("required_by_origin_engines", fmt.Sprintf("%v", catConf.RequiredByOriginEngines)).
		Str("preferred_engines", fmt.Sprintf("%v", catConf.PreferredEngines)).
		Str("preferred_by_origin_engines", fmt.Sprintf("%v", catConf.PreferredByOriginEngines)).
		Dur("preferred_timeout", catConf.Timings.PreferredTimeout).
		Dur("hard_timeout", catConf.Timings.HardTimeout).
		Msg("Searching")

	// Create contexts with timeout for HardTimeout and PreferredTimeout.
	ctxHardTimeout, cancelHardTimeoutFunc := context.WithTimeout(context.Background(), catConf.Timings.HardTimeout)
	defer cancelHardTimeoutFunc()
	ctxPreferredTimeout, cancelPreferredTimeoutFunc := context.WithTimeout(context.Background(), catConf.Timings.PreferredTimeout)
	defer cancelPreferredTimeoutFunc()

	// Create a context that cancels when both HardTimeout and PreferredTimeout are done.
	searchCtx, cancelSearch := context.WithCancel(context.Background())
	defer cancelSearch()
	go func() {
		<-ctxHardTimeout.Done()
		<-ctxPreferredTimeout.Done()
		cancelSearch()
	}()

	// Initialize each engine.
	searchers := initializeSearchers(searchCtx, catConf.Engines)

	// Create a map for the results with RWMutex.
	// TODO: Make title and desc length configurable.
	concMap := result.NewResultMap(len(catConf.Engines), 100, 1000)

	// Create a sync.Once wrapper for each searcher.Search() to ensure that the engine is only run once.
	onceWrapMap := initOnceWrapper(catConf.Engines)

	// Run all required engines. WaitGroup should be awaited unless the hard timeout is reached.
	var wgRequiredEngines sync.WaitGroup
	runRequiredSearchers(catConf.RequiredEngines, searchers, &wgRequiredEngines, &concMap, query, opts, onceWrapMap)

	// Run all required by origin engines. Cond should be awaited unless the hard timeout is reached.
	var wgRequiredByOriginEngines sync.WaitGroup
	runRequiredByOriginSearchers(catConf.RequiredByOriginEngines, searchers, &wgRequiredByOriginEngines, &concMap, catConf.Engines, query, opts, onceWrapMap)

	// Run all preferred engines. WaitGroup should be awaited unless the preferred timeout is reached.
	var wgPreferredEngines sync.WaitGroup
	runPreferredSearchers(catConf.PreferredEngines, searchers, &wgPreferredEngines, &concMap, query, opts, onceWrapMap)

	// Run all preferred by origin engines. Cond should be awaited unless the preferred timeout is reached.
	var wgPreferredByOriginEngines sync.WaitGroup
	runPreferredByOriginSearchers(catConf.PreferredByOriginEngines, searchers, &wgPreferredByOriginEngines, &concMap, catConf.Engines, query, opts, onceWrapMap)

	// Cancel the hard timeout after all required engines have finished and all required by origin engines have finished.
	go cancelHardTimeout(startTime, cancelHardTimeoutFunc, query, &wgRequiredEngines, catConf.RequiredEngines, &wgRequiredByOriginEngines, catConf.RequiredByOriginEngines)

	// Cancel the preferred timeout after all preferred engines have finished and all preferred by origin engines have finished.
	go cancelPreferredTimeout(startTime, cancelPreferredTimeoutFunc, query, &wgPreferredEngines, catConf.PreferredEngines, &wgPreferredByOriginEngines, catConf.PreferredByOriginEngines)

	// Wait for both hard timeout and preferred timeout to finish.
	<-searchCtx.Done()

	// Extract the results and responders from the map.
	results, responders := concMap.ExtractWithResponders()

	log.Debug().
		Int("results", len(results)).
		Str("query", anonymize.String(query)).
		Str("responders", fmt.Sprintf("%v", responders)).
		Dur("duration", time.Since(startTime)).
		Msg("Scraping finished")

	// Return the results.
	return results, nil
}
