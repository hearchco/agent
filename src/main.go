package main

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/search"
	"github.com/tminaorg/brzaguza/src/structures"
)

func printResults(results []structures.Result) {
	for _, r := range results {
		fmt.Printf("%v -----\n\t\"%v\"\n\t\"%v\"\n\t\"%v\"\n", r.Rank, r.Title, r.URL, r.Description)
	}
}

func main() {
	setupCli()
	setupLog()

	log.Info().
		Str("query", cli.Query).
		Msg("Started searching")

	start := time.Now()
	results := search.PerformSearch(cli.Query, cli.MaxPages, cli.Visit)
	duration := time.Since(start)

	printResults(results)
	log.Info().
		Msg(fmt.Sprintf("Found %v results in %vms", len(results), duration.Milliseconds()))
}
