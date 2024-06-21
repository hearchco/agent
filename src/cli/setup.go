package cli

import (
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"
)

// Returns flags struct from parsed cli arguments.
func Setup(ver string, timestamp string, commit string) (Flags, string) {
	verStruct := version{
		ver:       ver,
		timestamp: timestamp,
		commit:    commit,
	}

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
			"version": verStruct.String(),
		},
	)

	if err := ctx.Validate(); err != nil {
		log.Panic().
			Caller().
			Err(err).
			Msg("Failed parsing cli")
		// ^PANIC
	}

	return cli, verStruct.String()
}
