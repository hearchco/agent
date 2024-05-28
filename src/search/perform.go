package search

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/rank"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

func PerformSearch(query string, options engines.Options, categoryConf config.Category, settings map[engines.Name]config.Settings, salt string) []result.Result {
	// check for empty query
	if query == "" {
		log.Trace().
			Caller().
			Msg("Empty search query.")
		return []result.Result{}
	}

	// start searching
	searchTimer := time.Now()
	log.Debug().
		Str("queryAnon", anonymize.String(query)).
		Str("queryHash", anonymize.HashToSHA256B64(query)).
		Msg("Searching...")

	// getting results from engines
	resTimer := time.Now()
	log.Debug().Msg("Waiting for results from engines...")

	resultMap := runEngines(categoryConf.Engines, url.QueryEscape(query), options, settings, categoryConf.Timings, salt)

	log.Debug().
		Dur("duration", time.Since(resTimer)).
		Msg("Got results")

	// ranking results
	rankTimer := time.Now()
	log.Debug().Msg("Ranking...")

	results := rank.Rank(resultMap, categoryConf.Ranking)

	log.Debug().
		Dur("duration", time.Since(rankTimer)).
		Msg("Finished ranking")

	// finish searching
	log.Debug().
		Dur("duration", time.Since(searchTimer)).
		Msg("Found results")

	return results
}

func runEngines(engs []engines.Name, query string, options engines.Options, settings map[engines.Name]config.Settings, timings config.CategoryTimings, salt string) map[string]*result.Result {
	// create engine strings slice for logging
	engsStrs := make([]string, 0, len(engs))
	for _, eng := range engs {
		engsStrs = append(engsStrs, eng.String())
	}

	log.Info().
		Int("number", len(engs)).
		Strs("engines", engsStrs).
		Msg("Enabled engines")

	// create a relay to store results
	relay := bucket.Relay{
		ResultMap: make(map[string]*result.Result),
	}

	// create a wait group to wait for all engines to finish
	var wg sync.WaitGroup
	engineStarter := NewEngineStarter()

	start := time.Now()
	// initially set the preferred timeout minimum (will be reassigned to step time later)
	ctx, cancelCtx := context.WithTimeout(context.Background(), timings.PreferredTimeoutMin)
	ctxHard, cancelCtxHard := context.WithTimeout(context.Background(), timings.HardTimeout)

	// run all engines concurrently
	for _, eng := range engs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// if an error can be handled inside, it won't be returned
			// runs the Search function in the engine package
			errs := engineStarter[eng].Search(context.Background(), query, &relay, options, settings[eng], timings, salt, len(engs))
			if len(errs) > 0 {
				log.Error().
					Caller().
					Errs("errors", errs).
					Str("engine", eng.String()).
					Msg("Error(s) while searching")
			}
		}()
	}

	// wait for all engines to finish
	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		waitCh <- struct{}{}
	}()

	// break the loop if the preferred number of results is found before the preferred timeout is reached
	// otherwise break the loop when the minimum number of results if found
	// or if the hard timeout is reached
	// or if all engines finished
Outer:
	for {
		select {
		// preferred timeout (min/max) or step time reached
		case <-ctx.Done():
			currTimeout := time.Since(start)
			if currTimeout < timings.PreferredTimeoutMax {
				// if the preferred number of results isn't reached, continue additional step time
				if len(relay.ResultMap) < timings.PreferredResultsNumber {
					log.Debug().
						Dur("duration", currTimeout).
						Int("results", len(relay.ResultMap)).
						Msg("Timeout reached while waiting for engines, waiting additional step time")
					cancelCtx() // cancel the current context before creating a new one to prevent context leak
					ctx, cancelCtx = context.WithTimeout(context.Background(), timings.StepTime)
				} else {
					log.Debug().
						Dur("duration", currTimeout).
						Int("results", len(relay.ResultMap)).
						Msg("Timeout reached while waiting for engines")
					break Outer
				}
			} else {
				// if the minimum number of results isn't reached, continue additional step time
				if len(relay.ResultMap) < timings.MinimumResultsNumber {
					log.Debug().
						Dur("duration", currTimeout).
						Int("results", len(relay.ResultMap)).
						Msg("Preferred timeout maximum reached, waiting for minimum results required")
					cancelCtx() // cancel the current context before creating a new one to prevent context leak
					ctx, cancelCtx = context.WithTimeout(context.Background(), timings.StepTime)
				} else {
					log.Debug().
						Dur("duration", currTimeout).
						Int("results", len(relay.ResultMap)).
						Msg("Preferred timeout maximum reached")
					break Outer
				}
			}

		// hard timeout reached
		case <-ctxHard.Done():
			log.Debug().
				Dur("duration", time.Since(start)).
				Int("results", len(relay.ResultMap)).
				Msg("Hard timeout reached while waiting for engines")
			break Outer

		// all engines finished
		case <-waitCh:
			log.Debug().
				Dur("duration", time.Since(start)).
				Int("results", len(relay.ResultMap)).
				Msg("All engines finished")
			break Outer
		}
	}

	// cancel the current contexts to prevent context leak
	cancelCtx()
	cancelCtxHard()

	return relay.ResultMap
}
