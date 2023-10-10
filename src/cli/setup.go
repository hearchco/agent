package cli

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
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
			"config_path":  ".",
			"log_path":     "./log",
			"query_string": "banana death",
		},
	)

	if err := ctx.Validate(); err != nil {
		fmt.Println("Failed parsing cli:", err)
		os.Exit(1)
	}

	return cli
}
