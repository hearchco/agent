package main

import (
	"fmt"
	"time"

	"github.com/pkg/profile"
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

func runProfiler() func() {
	/*
		goroutine — stack traces of all current goroutines
		heap — a sampling of memory allocations of live objects
		allocs — a sampling of all past memory allocations
		threadcreate — stack traces that led to the creation of new OS threads
		block — stack traces that led to blocking on synchronization primitives
		mutex — stack traces of holders of contended mutexes
	*/

	var cpup interface{ Stop() }
	var gorp interface{ Stop() }
	var blockp interface{ Stop() }
	var threadp interface{ Stop() }
	var heapp interface{ Stop() }
	var allocp interface{ Stop() }
	var mutexp interface{ Stop() }

	if cli.CPUProfile != "" {
		cpup = profile.Start(profile.CPUProfile, profile.ProfilePath("./profiling/"+cli.CPUProfile))
	}
	if cli.HeapProfile != "" {
		heapp = profile.Start(profile.MemProfileHeap, profile.ProfilePath("./profiling/"+cli.HeapProfile))
	}
	if cli.GORProfile != "" {
		gorp = profile.Start(profile.GoroutineProfile, profile.ProfilePath("./profiling/"+cli.GORProfile))
	}
	if cli.ThreadProfile != "" {
		threadp = profile.Start(profile.ThreadcreationProfile, profile.ProfilePath("./profiling/"+cli.ThreadProfile))
	}
	if cli.BlockProfile != "" {
		blockp = profile.Start(profile.BlockProfile, profile.ProfilePath("./profiling/"+cli.BlockProfile))
	}
	if cli.AllocProfile != "" {
		allocp = profile.Start(profile.MemProfileAllocs, profile.ProfilePath("./profiling/"+cli.AllocProfile))
	}
	if cli.MutexProfile != "" {
		mutexp = profile.Start(profile.MutexProfile, profile.ProfilePath("./profiling/"+cli.MutexProfile))
	}

	return func() {
		if cli.CPUProfile != "" {
			cpup.Stop()
		}
		if cli.HeapProfile != "" {
			heapp.Stop()
		}
		if cli.GORProfile != "" {
			gorp.Stop()
		}
		if cli.ThreadProfile != "" {
			threadp.Stop()
		}
		if cli.BlockProfile != "" {
			blockp.Stop()
		}
		if cli.AllocProfile != "" {
			allocp.Stop()
		}
		if cli.MutexProfile != "" {
			mutexp.Stop()
		}
	}
}

func main() {
	mainTimer := time.Now()

	// parse cli arguments
	setupCli()

	defer runProfiler()() //runs the profiler, and defers the closing

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
		if db != nil {
			db.Get(cli.Query, &results)
		}
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

	if db != nil {
		db.Close()
	}

	//closeProfiler(cli.CPUProfile, cli.MEMProfile)
	log.Debug().Msgf("Program finished in %vms", time.Since(mainTimer).Milliseconds())
}
