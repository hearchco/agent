package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

var (
	// release variables
	Version   string
	Timestamp string
	GitCommit string

	// CLI
	cli struct {
		globals

		// flags
		Query     string `type:"string" default:"${query_string}" env:"BRZAGUZA_QUERY" help:"Query string used for search"`
		MaxPages  int    `type:"counter" default:"1" env:"BRZAGUZA_MAX_PAGES" help:"Number of pages to search"`
		Visit     bool   `type:"bool" default:"false" env:"BRZAGUZA_VISIT" help:"Should results be visited"`
		Log       string `type:"path" default:"${log_file}" env:"BRZAGUZA_LOG_FILE" help:"Log file path"`
		Verbosity int    `type:"counter" default:"0" short:"v" env:"BRZAGUZA_VERBOSITY" help:"Log level verbosity"`
	}
)

type globals struct {
	Version versionFlag `name:"version" help:"Print version information and quit"`
}

type versionFlag string

func (v versionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v versionFlag) IsBool() bool                         { return true }
func (v versionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

func setupCli() {
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
			"log_file":     "brzaguza.log",
			"query_string": "cars for sale in Toronto, Canada",
		},
	)

	if err := ctx.Validate(); err != nil {
		fmt.Println("Failed parsing cli:", err)
		os.Exit(1)
	}
}
