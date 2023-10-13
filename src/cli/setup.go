package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
)

func Setup() Flags {
	var cli Flags
	ctx := kong.Parse(&cli,
		kong.Name("brzaguza"),
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
		log.Fatal().Err(err).Msg("Failed parsing cli")
	}

	return cli
}
