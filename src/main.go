package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
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

	var cpup interface{}
	var gorp interface{}
	var blockp interface{}
	var threadp interface{}
	var heapp interface{}
	var allocp interface{}
	var mutexp interface{}

	cpup = profile.Start(profile.CPUProfile)
	gorp = profile.Start(profile.GoroutineProfile)
	blockp = profile.Start(profile.BlockProfile)
	threadp = profile.Start(profile.ThreadcreationProfile)
	heapp = profile.Start(profile.MemProfileHeap)
	allocp = profile.Start(profile.MemProfileAllocs)
	mutexp = profile.Start(profile.MutexProfile)

	var cpuFile *os.File
	if cli.CPUProfile != "" {
		cpuFile, err := os.Create("profiling/" + cli.CPUProfile)
		if err != nil {
			log.Fatal().Err(err).Msgf("couldn't create cpu profile. couldn't create file.")
		}
		if err := pprof.StartCPUProfile(cpuFile); err != nil {
			log.Fatal().Err(err).Msgf("couldn't create cpu profile. couldn't run StartCPUProfile")
		}
	}

	profile.Start(profile.CPUProfile)

	return func() {
		if cli.CPUProfile != "" {
			pprof.StopCPUProfile()
			if err := cpuFile.Close(); err != nil {
				log.Fatal().Err(err).Msgf("couldn't create cpu profile. couldn't close file.")
			}
		}

		if cli.MEMProfile != "" {
			f, err := os.Create("profiling/" + cli.MEMProfile)
			if err != nil {
				log.Fatal().Err(err).Msgf("couldn't create memory profile. couldn't create file.")
			}
			runtime.GC()
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal().Err(err).Msgf("couldn't create memory profile. couldn't WriteHeapProfile.")
			}
			if err := f.Close(); err != nil {
				log.Fatal().Err(err).Msgf("couldn't create memory profile. couldn't close file.")
			}
		}

		if cli.GORProfile != "" {
			f, err := os.Create("profiling/" + cli.MEMProfile)
			if err != nil {
				log.Fatal().Err(err).Msgf("couldn't create goroutine profile. couldn't create file.")
			}
			if err := pprof.Lookup("goroutine").WriteTo(f, 0); err != nil {
				log.Fatal().Err(err).Msgf("couldn't create goroutine profile. failed profile write.")
			}
			if err := f.Close(); err != nil {
				log.Fatal().Err(err).Msgf("couldn't create goroutine profile. couldn't close file.")
			}
		}

		if cli.ThreadProfile != "" {

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
