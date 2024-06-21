package search

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/config"
	"github.com/hearchco/agent/src/search/category"
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/utils/anonymize"
)

func Search(query string, category category.Name, opts options.Options, catConf config.Category) ([]result.Result, error) {
	// Capture start time.
	startTime := time.Now()

	if err := validateParams(query, opts); err != nil {
		return nil, err
	}

	log.Debug().
		Str("category", category.String()).
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
	searchers := initializeSearchers(searchCtx, catConf.Engines, catConf.Timings)

	// Create a channel of channels to receive the results from each engine.
	engChan := make(chan chan result.ResultScraped, len(catConf.Engines))

	// Create a map for the results with RWMutex.
	resMap := result.Map()

	// Start a goroutine to receive the results from each engine and add them to results map.
	go createReceiver(engChan, &resMap, len(catConf.Engines))

	// Create a sync.Once wrapper for each searcher.Search() to ensure that the engine is only run once.
	searchOnce := initOnceWrapper(catConf.Engines)

	// Run all required engines. WaitGroup should be awaited unless the hard timeout is reached.
	var wgRequiredEngines sync.WaitGroup
	runRequiredEngines(searchers, &wgRequiredEngines, query, opts, catConf.RequiredEngines, engChan, searchOnce)

	// Run all required by origin engines. Cond should be awaited unless the hard timeout is reached.
	var wgRequiredByOriginEngines sync.WaitGroup
	runRequiredByOriginEngines(searchers, &wgRequiredByOriginEngines, query, opts, catConf.RequiredByOriginEngines, catConf.Engines, engChan, searchOnce)

	// Run all preferred engines. WaitGroup should be awaited unless the preferred timeout is reached.
	var wgPreferredEngines sync.WaitGroup
	runPreferredEngines(searchers, &wgPreferredEngines, query, opts, catConf.PreferredEngines, engChan, searchOnce)

	// Run all preferred by origin engines. Cond should be awaited unless the preferred timeout is reached.
	var wgPreferredByOriginEngines sync.WaitGroup
	runPreferredByOriginEngines(searchers, &wgPreferredByOriginEngines, query, opts, catConf.PreferredByOriginEngines, catConf.Engines, engChan, searchOnce)

	// Close the channel of channels (it's safe because each sending already happened sequentially).
	close(engChan)

	// Cancel the hard timeout after all required engines have finished and all required by origin engines have finished.
	go cancelHardTimeout(startTime, cancelHardTimeoutFunc, query, &wgRequiredEngines, catConf.RequiredEngines, &wgRequiredByOriginEngines, catConf.RequiredByOriginEngines)

	// Cancel the preferred timeout after all preferred engines have finished and all preferred by origin engines have finished.
	go cancelPreferredTimeout(startTime, cancelPreferredTimeoutFunc, query, &wgPreferredEngines, catConf.PreferredEngines, &wgPreferredByOriginEngines, catConf.PreferredByOriginEngines)

	// Wait for both hard timeout and preferred timeout to finish.
	<-searchCtx.Done()

	// Extract the results and responders from the map.
	// TODO: Make title and desc length configurable.
	results, responders := resMap.ExtractResultsAndResponders(len(catConf.Engines), 100, 1000)

	log.Debug().
		Int("results", len(results)).
		Str("query", anonymize.String(query)).
		Str("responders", fmt.Sprintf("%v", responders)).
		Dur("duration", time.Since(startTime)).
		Msg("Scraping finished")

	// Return the results.
	return results, nil
}
