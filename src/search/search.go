package search

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/bucket/result"
	"github.com/tminaorg/brzaguza/src/category"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/rank"
)

func PerformSearch(query string, options engines.Options, conf *config.Config) []result.Result {
	searchTimer := time.Now()

	relay := bucket.Relay{
		ResultMap: make(map[string]*result.Result),
	}

	deadline, toRun := procBang(&query, &options, conf)

	query = url.QueryEscape(query)
	log.Debug().Msg(query)

	resTimer := time.Now()
	log.Debug().Msg("Waiting for results from engines...")
	var worker conc.WaitGroup
	runEngines(deadline, toRun, conf.Settings, query, &worker, &relay, options)
	worker.Wait()
	log.Debug().Msgf("Got results in %vms", time.Since(resTimer).Milliseconds())

	rankTimer := time.Now()
	log.Debug().Msg("Ranking...")
	results := rank.Rank(relay.ResultMap, conf.Categories[options.Category].Ranking) // have to make copy, since its a map value
	rankTimeSince := time.Since(rankTimer)
	log.Debug().Msgf("Finished ranking in %vms (%vns)", rankTimeSince.Milliseconds(), rankTimeSince.Nanoseconds())

	log.Debug().Msgf("Found results in %vms", time.Since(searchTimer).Milliseconds())

	return results
}

// engine_searcher, NewEngineStarter()  use this.
type EngineSearch func(context.Context, string, *bucket.Relay, engines.Options, config.Settings) error

func runEngines(deadline time.Duration, engs []engines.Name, settings map[engines.Name]config.Settings, query string, worker *conc.WaitGroup, relay *bucket.Relay, options engines.Options) {
	config.EnabledEngines = engs
	log.Info().Msgf("Enabled engines (%v): %v", len(config.EnabledEngines), config.EnabledEngines)

	engineStarter := NewEngineStarter()
	for i := range engs {
		eng := engs[i] // dont change for to `for _, eng := range engs {`, eng retains the same address throughout the whole loop
		worker.Go(func() {
			ctxTimer := time.Now()

			deadline := 1000 * time.Millisecond
			ctx, cancelCtx := context.WithTimeout(context.Background(), deadline)

			worker.Go(func() {
				err := engineStarter[eng](ctx, query, relay, options, settings[eng])
				if err != nil {
					log.Error().Err(err).Msgf("failed searching %v", eng)
				}
				cancelCtx()
			})

			<-ctx.Done()
			log.Trace().Msgf("%v: context done in %vms (deadline: %vms)", eng, time.Since(ctxTimer).Milliseconds(), deadline.Milliseconds())
		})
	}
}

func procBang(query *string, options *engines.Options, conf *config.Config) (time.Duration, []engines.Name) {
	useSpec, specEng := procSpecificEngine(*query, options, conf)
	goodCat := procCategory(*query, options)
	if !goodCat && !useSpec && (*query)[0] == '!' {
		log.Error().Msgf("invalid bang (not category or engine shortcut). query: %v", *query)
	}

	trimBang(query)

	if useSpec {
		return 5000 * time.Millisecond, []engines.Name{specEng}
	} else {
		return time.Duration(conf.Categories[options.Category].Deadline) * time.Millisecond, conf.Categories[options.Category].Engines
	}
}

func trimBang(query *string) {
	if (*query)[0] == '!' {
		*query = strings.SplitN(*query, " ", 2)[1]
	}
}

func procSpecificEngine(query string, options *engines.Options, conf *config.Config) (bool, engines.Name) {
	if query[0] != '!' {
		return false, engines.UNDEFINED
	}
	sp := strings.SplitN(query, " ", 2)
	specE := sp[0][1:]
	for key, val := range conf.Settings {
		if val.Shortcut == specE {
			return true, key
		}
	}

	return false, engines.UNDEFINED
}

func procCategory(query string, options *engines.Options) bool {
	cat := category.FromQuery(query)
	if cat != "" {
		options.Category = cat
	}
	if options.Category == "" {
		options.Category = category.GENERAL
	}
	return cat != ""
}
