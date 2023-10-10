package cli

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket/result"
	"github.com/tminaorg/brzaguza/src/cache"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
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

func Run(flags Flags, db cache.DB, conf *config.Config) {
	log.Info().
		Str("query", flags.Query).
		Str("max-pages", fmt.Sprintf("%v", flags.MaxPages)).
		Str("visit", fmt.Sprintf("%v", flags.Visit)).
		Msg("Started searching")

	options := engines.Options{
		MaxPages:   flags.MaxPages,
		VisitPages: flags.Visit,
		Category:   flags.Category,
	}

	start := time.Now()

	// todo: this should be refactor to cliMode package with ctx cancelling as well
	var results []result.Result
	db.Get(flags.Query, &results)
	if results != nil {
		log.Debug().Msgf("Found results for query (%v) in cache", flags.Query)
	} else {
		log.Debug().Msg("Nothing found in cache, doing a clean search")
		results = search.PerformSearch(flags.Query, options, conf)
		db.Set(flags.Query, results)
	}

	duration := time.Since(start)

	if !flags.Silent {
		printResults(results)
	}
	log.Info().Msgf("Found %v results in %vms", len(results), duration.Milliseconds())
}
