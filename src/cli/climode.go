package cli

import (
	"fmt"
	"time"

	"github.com/hearchco/hearchco/src/bucket/result"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/category"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/engines"
	"github.com/hearchco/hearchco/src/search"
	"github.com/rs/zerolog/log"
)

func printResults(results []result.Result) {
	fmt.Print("\n\tThe Search Results:\n\n")
	for _, r := range results {
		fmt.Printf("%v (%.2f) -----\n\t\"%v\"\n\t\"%v\"\n\t\"%v\"\n\t-", r.Rank, r.Score, r.Title, r.URL, r.Description)
		for seInd := uint8(0); seInd < r.TimesReturned; seInd++ {
			fmt.Printf("%v", r.EngineRanks[seInd].SearchEngine.ToLower())
			if seInd != r.TimesReturned-1 {
				fmt.Print(", ")
			}
		}
		fmt.Printf("\n")
	}
}

func Run(flags Flags, db cache.DB, conf *config.Config) {
	log.Info().
		Str("query", flags.Query).
		Str("max-pages", fmt.Sprintf("%v", flags.MaxPages)).
		Str("visit", fmt.Sprintf("%v", flags.Visit)).
		Msg("Started searching")

	options := engines.Options{
		MaxPages:   flags.MaxPages,
		VisitPages: flags.Visit,
		Category:   category.FromString[flags.Category],
	}

	start := time.Now()

	// todo: ctx cancelling (important since pebble is NoSync)
	var results []result.Result
	gerr := db.Get(flags.Query, &results)
	if gerr != nil {
		log.Fatal().Err(gerr).Msgf("cli.Run(): failed accessing cache for query %v", flags.Query)
		// ^FATAL
	}

	if results != nil {
		log.Debug().Msgf("Found results for query (%v) in cache", flags.Query)
	} else {
		log.Debug().Msg("Nothing found in cache, doing a clean search")

		results = search.PerformSearch(flags.Query, options, conf)

		serr := db.Set(flags.Query, results)
		if serr != nil {
			log.Error().Err(serr).Msgf("cli.Run(): error updating database with search results")
		}
	}

	duration := time.Since(start)

	if !flags.Silent {
		printResults(results)
	}
	log.Info().Msgf("Found %v results in %vms", len(results), duration.Milliseconds())
}
