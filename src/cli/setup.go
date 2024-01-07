package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
)

func validateLocale(locale string) {
	if len(locale) != 5 {
		log.Fatal().Msg("cli.validateLocale(): failed parsing cli locale argument: isn't 5 characters long")
		// ^FATAL
	}
	if !(('a' <= locale[0] && locale[0] <= 'z') && ('a' <= locale[1] && locale[1] <= 'z')) {
		log.Fatal().Msg("cli.validateLocale(): failed parsing cli locale argument: first two characters must be lowercase ASCII letters")
		// ^FATAL
	}
	if !(('A' <= locale[3] && locale[3] <= 'Z') && ('A' <= locale[4] && locale[4] <= 'Z')) {
		log.Fatal().Msg("cli.validateLocale(): failed parsing cli locale argument: last two characters must be uppercase ASCII letters")
		// ^FATAL
	}
	if locale[2] != '_' {
		log.Fatal().Msg("cli.validateLocale(): failed parsing cli locale argument: third character must be underscore (_)")
		// ^FATAL
	}
}

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
			"log_folder":   "./log",
			"query_string": "banana death",
		},
	)

	if err := ctx.Validate(); err != nil {
		log.Panic().Err(err).Msg("cli.Setup(): failed parsing cli") // panic is also run inside the library. when does this happen?
		// ^PANIC
	}

	validateLocale(cli.Locale)

	return cli
}
