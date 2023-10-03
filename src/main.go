package main

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket/result"
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
		results := search.PerformSearch(cli.Query, options, config)
		duration := time.Since(start)
		if !cli.Silent {
			printResults(results)
		}
		log.Info().Msgf("Found %v results in %vms", len(results), duration.Milliseconds())
	} else {
		router.Setup(config)
	}
}
