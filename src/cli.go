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
		Query         string `type:"string" default:"${query_string}" env:"BRZAGUZA_QUERY" help:"Query string used for search"`
		MaxPages      int    `type:"counter" default:"1" env:"BRZAGUZA_MAX_PAGES" help:"Number of pages to search"`
		Cli           bool   `type:"bool" default:"false" env:"BRZAGUZA_CLI" help:"Use CLI mode"`
		Visit         bool   `type:"bool" default:"false" env:"BRZAGUZA_VISIT" help:"Should results be visited"`
		Silent        bool   `type:"bool" default:"false" short:"s" env:"BRZAGUZA_SILENT" help:"Should results be printed"`
		Config        string `type:"path" default:"${config_path}" env:"BRZAGUZA_CONFIG" help:"Config folder path"`
		Log           string `type:"path" default:"${log_path}" env:"BRZAGUZA_LOG" help:"Log file path"`
		Verbosity     int    `type:"counter" default:"0" short:"v" env:"BRZAGUZA_VERBOSITY" help:"Log level verbosity"`
		CPUProfile    bool   `type:"bool" default:"false" env:"BRZAGUZA_CPUPROFILE" help:"Use cpu profiling"`
		HeapProfile   bool   `type:"bool" default:"false" env:"BRZAGUZA_HEAPPROFILE" help:"Use heap profiling"`
		GORProfile    bool   `type:"bool" default:"false" env:"BRZAGUZA_GORPROFILE" help:"Use goroutine profiling"`
		ThreadProfile bool   `type:"bool" default:"false" env:"BRZAGUZA_THREADPROFILE" help:"Use threadcreate profiling"`
		AllocProfile  bool   `type:"bool" default:"false" env:"BRZAGUZA_MEMALLOCPROFILE" help:"Use alloc profiling"`
		BlockProfile  bool   `type:"bool" default:"false" env:"BRZAGUZA_BLOCKPROFILE" help:"Use block profiling"`
		MutexProfile  bool   `type:"bool" default:"false" env:"BRZAGUZA_MUTEXPROFILE" help:"Use mutex profiling"`
		ClockProfile  bool   `type:"bool" default:"false" env:"BRZAGUZA_CLOCKPROFILE" help:"Use clock profiling"`
		TraceProfile  bool   `type:"bool" default:"false" env:"BRZAGUZA_TRACEPROFILE" help:"Use trace profiling"`
		ServeProfiler bool   `type:"bool" default:"false" env:"BRZAGUZA_SERVEPROFILER" help:"Run the profiler and serve at /debug/pprof/ http endpoint"`
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
			"config_path":  ".",
			"log_path":     "./log",
			"query_string": "banana death",
		},
	)

	if err := ctx.Validate(); err != nil {
		fmt.Println("Failed parsing cli:", err)
		os.Exit(1)
	}
}
