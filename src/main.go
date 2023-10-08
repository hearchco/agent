package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket/result"
	"github.com/tminaorg/brzaguza/src/cache"
	"github.com/tminaorg/brzaguza/src/cache/nocache"
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
	mainTimer := time.Now()

	// parse cli arguments
	setupCli()

	// configure logging
	logger.Setup(cli.Log, cli.Verbosity)

	// load config file
	config := config.New()
	config.Load(cli.Config, cli.Log)

	// signal interrupt (CTRL+C)
	ctx, stopCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// cache database
	var db cache.DB
	switch config.Server.Cache.Type {
	case "pebble":
		db = pebble.New(cli.Config)
	case "redis":
		db = redis.New(ctx, config.Server.Cache.Redis)
	default:
		db = nocache.New()
		log.Warn().Msg("Running without caching!")
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

		// todo: this should be refactor to cliStartup package with ctx cancelling as well
		var results []result.Result
		db.Get(cli.Query, &results)
		if results != nil {
			log.Debug().Msgf("Found results for query (%v) in cache", cli.Query)
		} else {
			log.Debug().Msg("Nothing found in cache, doing a clean search")
			results = search.PerformSearch(cli.Query, options, config)
			db.Set(cli.Query, results)
		}

		duration := time.Since(start)

		if !cli.Silent {
			printResults(results)
		}
		log.Info().Msgf("Found %v results in %vms", len(results), duration.Milliseconds())
	} else {
		if router, err := router.New(config); err != nil {
			log.Error().Msgf("Failed creating a router: %v", err)
		} else {
			router.Start(ctx, db)
		}
	}

	// program cleanup
	db.Close()
	stopCtx()

	log.Debug().Msgf("Program finished in %vms", time.Since(mainTimer).Milliseconds())
}
