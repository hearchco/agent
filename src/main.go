package main

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket/result"
	"github.com/tminaorg/brzaguza/src/cache"
	"github.com/tminaorg/brzaguza/src/cache/pebble"
	"github.com/tminaorg/brzaguza/src/cache/redis"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/logger"
	"github.com/tminaorg/brzaguza/src/router"
	"github.com/tminaorg/brzaguza/src/search"
)

func printResults(results []result.Result) {
	fmt.Print("\n\tThe Search Results:\n\n")
	for _, r := range results {
		fmt.Printf("%v (%.2f) -----\n\t\"%v\"\n\t\"%v\"\n\t\"%v\"\n\t-", r.Rank, r.Score, r.Title, r.URL, r.Description)
		for seInd := uint8(0); seInd < r.TimesReturned; seInd++ {
			fmt.Printf("%v", r.EngineRanks[seInd].SearchEngine)
			if seInd != r.TimesReturned-1 {
				fmt.Print(", ")
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	// parse cli arguments
	setupCli()

	// configure logging
	logger.Setup(cli.Log, cli.Verbosity)

	// load config file
	config := config.New()
	config.Load(cli.Config, cli.Log)

	// cache database
	var db cache.DB
	switch config.Server.Cache.Type {
	case "pebble":
		db = pebble.New(cli.Config)
	case "redis":
		db = redis.New(config.Server.Cache.Redis)
	default:
		log.Warn().Msg("Running without caching!")
	}

	if db != nil {
		defer db.Close()
	}

	// startup
	if cli.Cli {
		log.Info().
			Str("query", cli.Query).
			Str("max-pages", fmt.Sprintf("%v", cli.MaxPages)).
			Str("visit", fmt.Sprintf("%v", cli.Visit)).
			Msg("Started searching")

		options := engines.Options{
			MaxPages:   cli.MaxPages,
			VisitPages: cli.Visit,
		}

		start := time.Now()

		var results []result.Result
		db.Get(cli.Query, &results)
		if results != nil {
			log.Debug().Msgf("Found results for query (%v) in cache", cli.Query)
		} else {
			log.Debug().Msg("Nothing found in cache, doing a clean search")
			results = search.PerformSearch(cli.Query, options, config)
			cache.Save(db, cli.Query, results)
		}

		duration := time.Since(start)

		if !cli.Silent {
			printResults(results)
		}
		log.Info().Msgf("Found %v results in %vms", len(results), duration.Milliseconds())
	} else {
		router.Setup(config, db)
	}
}
