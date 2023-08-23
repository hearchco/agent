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
		Query      string `type:"string" default:"${query_string}" env:"BRZAGUZA_QUERY" help:"Query string used for search"`
		MaxPages   int    `type:"counter" default:"1" env:"BRZAGUZA_MAX_PAGES" help:"Number of pages to search"`
		Visit      bool   `type:"bool" default:"false" env:"BRZAGUZA_VISIT" help:"Should results be visited"`
		Silent     bool   `type:"bool" default:"false" short:"s" env:"BRZAGUZA_SILENT" help:"Should results be printed"`
		ConfigPath string `type:"path" default:"${config_folder}" env:"BRZAGUZA_CONFIG_FOLDER" help:"Config folder path"`
		Config     string `type:"string" default:"${config_file}" env:"BRZAGUZA_CONFIG_FILE" help:"Config file name"`
		Log        string `type:"path" default:"${log_file}" env:"BRZAGUZA_LOG_FILE" help:"Log file name"`
		Verbosity  int    `type:"counter" default:"0" short:"v" env:"BRZAGUZA_VERBOSITY" help:"Log level verbosity"`
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
	// this is used to determine if the container image is being built
	_, dockerEnv := os.LookupEnv("BRZAGUZA_DOCKER")
	var configPath string
	if dockerEnv {
		configPath = "/config"
	} else {
		configPath = "./"
	}

	ctx := kong.Parse(&cli,
		kong.Name("brzaguza"),
		kong.Description("Fastasst metasearch engine"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Summary: true,
			Compact: true,
		}),
		kong.Vars{
			"version":       fmt.Sprintf("%v (%v@%v)", Version, GitCommit, Timestamp),
			"config_folder": configPath,
			"config_file":   "brzaguza", // can be .yaml or any other type defined by Koanf
			"log_file":      "brzaguza.log",
			"query_string":  "banana death",
		},
	)

	if err := ctx.Validate(); err != nil {
		fmt.Println("Failed parsing cli:", err)
		os.Exit(1)
	}
}
