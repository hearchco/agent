package cli

import (
	"fmt"
	"time"

	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/cache"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

func printImageResult(r result.Result) {
	fmt.Printf("%v (%.2f) -----\n\t%q\n\t%q\n\t%q\n\t%q\n\t%q\n\t%q\n\t-", r.Rank, r.Score, r.Title, r.URL, r.Description, r.ImageResult.Source, r.ImageResult.SourceURL, r.ImageResult.ThumbnailURL)
	for seInd := range len(r.EngineRanks) {
		fmt.Printf("%v", r.EngineRanks[seInd].SearchEngine.ToLower())
		if seInd != len(r.EngineRanks)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Printf("\n")
}

func printResult(r result.Result) {
	fmt.Printf("%v (%.2f) -----\n\t%q\n\t%q\n\t%q\n\t-", r.Rank, r.Score, r.Title, r.URL, r.Description)
	for seInd := range len(r.EngineRanks) {
		fmt.Printf("%v", r.EngineRanks[seInd].SearchEngine.ToLower())
		if seInd != len(r.EngineRanks)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Printf("\n")
}

func printResults(results []result.Result) {
	fmt.Print("\n\tThe Search Results:\n\n")

	images := false
	if len(results) > 0 && results[0].ImageResult.Source != "" {
		images = true
	}

	for _, r := range results {
		if images {
			printImageResult(r)
		} else {
			printResult(r)
		}
	}
}

func Run(flags Flags, db cache.DB, conf config.Config) {
	log.Info().
		Str("queryAnon", anonymize.String(flags.Query)).
		Str("queryHash", anonymize.HashToSHA256B64(flags.Query)).
		Int("maxPages", flags.PagesMax).
		Bool("visit", flags.Visit).
		Msg("Started hearching")

	categoryName, err := category.FromString(flags.Category)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid category")
	}

	// all of these have default values set and are validated beforehand
	options := engines.Options{
		VisitPages: flags.Visit,
		SafeSearch: flags.SafeSearch,
		Pages: engines.Pages{
			Start: flags.PagesStart,
			Max:   flags.PagesMax,
		},
		Locale:   flags.Locale,
		Category: categoryName,
	}

	start := time.Now()

	results, foundInDB := search.Search(flags.Query, options, db, conf.Categories[options.Category], conf.Settings, conf.Server.Proxy.Salt)

	if !flags.Silent {
		printResults(results)
	}

	log.Info().
		Int("number", len(results)).
		Dur("duration", time.Since(start)).
		Msg("Found results")

	search.CacheAndUpdateResults(flags.Query, options, db, conf.Server.Cache.TTL, conf.Categories[options.Category], conf.Settings, results, foundInDB, conf.Server.Proxy.Salt)
}
