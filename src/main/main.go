package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/structures"
)

func printResults(results []structures.Result) {
	for _, r := range results {
		fmt.Printf("%v -----\n\t\"%s\"\n\t\"%s\"\n\t\"%s\"\n", r.Rank, r.Title, r.URL, r.Description)
	}
}

func main() {
	setupCli()
	setupLog()

	log.Info().
		Str("query", cli.Query).
		Msg("Started searching")
	results := performSearch(cli.Query)
	log.Info().
		Msg(fmt.Sprintf("Found %d results", len(results)))

	printResults(results)

}
