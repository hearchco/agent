package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/hearchco/hearchco/src/gotypelimits"
	"github.com/hearchco/hearchco/src/search/category"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

func Setup() Flags {
	var cli Flags
	ctx := kong.Parse(&cli,
		kong.Name("hearchco"),
		kong.Description("Fastasst metasearch engine"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Summary: true,
			Compact: true,
		}),
		kong.Vars{
			"version":      fmt.Sprintf("%v (%v@%v)", Version, GitCommit, Timestamp),
			"data_folder":  ".",
			"query_string": "banana death",
		},
	)

	if err := ctx.Validate(); err != nil {
		log.Panic().
			Caller().
			Err(err).
			Msg("Failed parsing cli")
		// ^PANIC
	}

	if cli.Query == "" {
		log.Fatal().
			Caller().
			Msg("Query cannot be empty or whitespace")
		// ^FATAL
	}

	// TODO: make upper limit configurable
	pagesMaxUpperLimit := 10
	if cli.PagesMax < 1 || cli.PagesMax > pagesMaxUpperLimit {
		log.Fatal().
			Caller().
			Int("pages", cli.PagesMax).
			Int("min", 1).
			Int("max", pagesMaxUpperLimit).
			Msg("Pages value out of range")
		// ^FATAL
	}

	if cli.PagesStart < 1 || cli.PagesStart > gotypelimits.MaxInt-pagesMaxUpperLimit {
		log.Fatal().
			Caller().
			Int("start", cli.PagesStart).
			Int("min", 1).
			Int("max", gotypelimits.MaxInt-pagesMaxUpperLimit).
			Msg("Start value out of range")
		// ^FATAL
	} else {
		// since it's >=1, we decrement it to match the 0-based index
		cli.PagesStart -= 1
	}

	if err := engines.ValidateLocale(cli.Locale); err != nil {
		log.Fatal().
			Caller().
			Err(err).
			Msg("Invalid locale flag")
		// ^FATAL
	}

	if _, err := category.FromString(cli.Category); err != nil {
		log.Fatal().
			Caller().
			Msg("Invalid category flag")
		// ^FATAL
	}

	return cli
}
