package cli

import (
	"fmt"
	"strings"
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
	query := flags.Query
	// convert query to only contain the actual query, w/o category or whitespaces
	query = strings.TrimSpace(query)
	catFromQuery := category.FromQuery(query)
	if catFromQuery != "" {
		// remove the category from the query
		query = strings.TrimSpace(strings.TrimPrefix(query, "!"+catFromQuery.String()))
	}

	if query == "" {
		log.Error().Msg("Empty query or only category found")
		return
	}

	// convert category string to category.Name, either from query (takes precedence) or from parameters
	var categoryName category.Name
	if catFromQuery != "" {
		categoryName = category.SafeFromString(catFromQuery.String())
	} else {
		categoryName = category.SafeFromString(flags.Category)
	}

	if categoryName == category.UNDEFINED {
		log.Error().Msg("Invalid category")
		return
	}

	options := engines.Options{
		VisitPages: flags.Visit,
		SafeSearch: flags.SafeSearch,
		Mobile:     flags.Mobile,
		Pages: engines.Pages{
			Start: flags.StartPage,
			Max:   flags.MaxPages,
		},
		UserAgent: flags.UserAgent,
		Locale:    flags.Locale,
		Category:  categoryName,
	}

	log.Info().
		Str("queryAnon", anonymize.String(query)).
		Str("queryHash", anonymize.HashToSHA256B64(query)).
		Str("category", options.Category.String()).
		Bool("visit", options.VisitPages).
		Bool("safeSearch", options.SafeSearch).
		Bool("mobile", options.Mobile).
		Int("startPage", options.Pages.Start).
		Int("maxPages", options.Pages.Max).
		Str("userAgent", options.UserAgent).
		Str("locale", options.Locale).
		Msg("Started hearching")

	start := time.Now()

	results, foundInDB := search.Search(query, options, db, conf.Settings, conf.Categories, conf.Server.Proxy.Salt)

	if !flags.Silent {
		printResults(results)
	}

	log.Info().
		Int("number", len(results)).
		Dur("duration", time.Since(start)).
		Msg("Found results")

	search.CacheAndUpdateResults(query, options, db, conf.Server.Cache.TTL, conf.Settings, conf.Categories, results, foundInDB, conf.Server.Proxy.Salt)
}
