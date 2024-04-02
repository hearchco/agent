package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
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
		log.Panic().Err(err).Msg("cli.Setup(): failed parsing cli") // panic is also run inside the library. when does this happen?
		// ^PANIC
	}

	if locErr := engines.ValidateLocale(cli.Locale); locErr != nil {
		log.Fatal().Err(locErr).Msg("cli.Setup(): invalid locale flag")
		// ^FATAL
	}

	if category.SafeFromString(cli.Category) == category.UNDEFINED {
		log.Fatal().Msg("cli.Setup(): invalid category flag")
		// ^FATAL
	}

	if cli.StartPage < 1 {
		log.Fatal().
			Int("startpage", cli.StartPage).
			Msg("cli.Setup(): invalid start page flag (must be >= 1)")
		// ^FATAL
	} else {
		// since it's >=1, we decrement it to match the 0-based index
		cli.StartPage -= 1
	}

	if cli.MaxPages < 1 {
		log.Fatal().
			Int("maxpages", cli.MaxPages).
			Msg("cli.Setup(): invalid max pages flag (must be >= 1)")
		// ^FATAL
	}

	return cli
}
