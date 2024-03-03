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
	for seInd := uint8(0); seInd < r.TimesReturned; seInd++ {
		fmt.Printf("%v", r.EngineRanks[seInd].SearchEngine.ToLower())
		if seInd != r.TimesReturned-1 {
			fmt.Print(", ")
		}
	}
	fmt.Printf("\n")
}

func printResult(r result.Result) {
	fmt.Printf("%v (%.2f) -----\n\t%q\n\t%q\n\t%q\n\t-", r.Rank, r.Score, r.Title, r.URL, r.Description)
	for seInd := uint8(0); seInd < r.TimesReturned; seInd++ {
		fmt.Printf("%v", r.EngineRanks[seInd].SearchEngine.ToLower())
		if seInd != r.TimesReturned-1 {
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
		Int("maxPages", flags.MaxPages).
		Bool("visit", flags.Visit).
		Msg("Started hearching")

	options := engines.Options{
		MaxPages:   flags.MaxPages,
		VisitPages: flags.Visit,
		Category:   category.FromString[flags.Category],
		UserAgent:  flags.UserAgent,
		Locale:     flags.Locale,
		SafeSearch: flags.SafeSearch,
		Mobile:     flags.Mobile,
	}

	start := time.Now()

	results, foundInDB := search.Search(flags.Query, options, db, conf.Settings, conf.Categories)

	duration := time.Since(start)
	if !flags.Silent {
		printResults(results)
	}
	log.Info().
		Int("number", len(results)).
		Int64("ms", duration.Milliseconds()).
		Msg("Found results")

	search.CacheAndUpdateResults(flags.Query, options, db, conf.Server.Cache.TTL, conf.Settings, conf.Categories, results, foundInDB)
}
