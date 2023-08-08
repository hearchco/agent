package main

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/search"
	"github.com/tminaorg/brzaguza/src/structures"
)

func printResults(results []structures.Result) {
	fmt.Print("\n\tThe Search Results:\n\n")
	for _, r := range results {
		fmt.Printf("%v -----\n\t\"%v\"\n\t\"%v\"\n\t\"%v\"\n\t-", r.Rank, r.Title, r.URL, r.Description)
		for seInd := 0; seInd < r.TimesReturned; seInd++ {
			fmt.Printf("%v", r.SearchEngines[seInd].SearchEngine)
			if seInd != r.TimesReturned-1 {
				fmt.Print(", ")
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	setupCli()
	setupLog()
	config := config.SetupConfig(cli.ConfigPath, cli.Config)

	log.Info().
		Str("query", cli.Query).
		Str("max-pages", fmt.Sprintf("%v", cli.MaxPages)).
		Str("visit", fmt.Sprintf("%v", cli.Visit)).
		Msg("Started searching")

	start := time.Now()
	results := search.PerformSearch(cli.Query, cli.MaxPages, cli.Visit, config)
	duration := time.Since(start)

	printResults(results)
	log.Info().Msgf("Found %v results in %vms", len(results), duration.Milliseconds())
}
