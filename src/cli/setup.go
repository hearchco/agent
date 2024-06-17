package cli

import (
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
)

var (
	// Release variables.
	Version   string
	Timestamp string
	GitCommit string
)

// Returns flags struct from parsed cli arguments.
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
			"version": VersionString(),
		},
	)

	if err := ctx.Validate(); err != nil {
		log.Panic().
			Caller().
			Err(err).
			Msg("Failed parsing cli")
		// ^PANIC
	}

	return cli
}
