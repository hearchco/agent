package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
)

type Flags struct {
	Globals

	// flags
	Query       string `type:"string" default:"${query_string}" env:"HEARCHCO_QUERY" help:"Query string used for search"`
	Cli         bool   `type:"bool" default:"false" env:"HEARCHCO_CLI" help:"Use CLI mode"`
	Silent      bool   `type:"bool" default:"false" short:"s" env:"HEARCHCO_SILENT" help:"Should results be printed"`
	DataDirPath string `type:"path" default:"${data_folder}" env:"HEARCHCO_DATA_DIR" help:"Data folder path"`
	LogDirPath  string `type:"path" default:"${log_folder}" env:"HEARCHCO_LOG_DIR" help:"Log folder path"`
	LogToFile   bool   `type:"bool" default:"false" env:"HEARCHCO_LOG_TO_FILE" help:"Should logs be written to file"`
	Verbosity   int8   `type:"counter" default:"0" short:"v" env:"HEARCHCO_VERBOSITY" help:"Log level verbosity"`
	// options
	MaxPages   int    `type:"counter" default:"1" env:"HEARCHCO_MAX_PAGES" help:"Number of pages to search"`
	Visit      bool   `type:"bool" default:"false" env:"HEARCHCO_VISIT" help:"Should results be visited"`
	Category   string `type:"string" default:"" short:"c" env:"HEARCHCO_CATEGORY" help:"Search result category. Can also be supplied through the query (e.g. \"!info smartphone\"). Supported values: info[/wiki], science[/sci], news, blog, surf, newnews[/nnews]"`
	UserAgent  string `type:"string" default:"" env:"HEARCHCO_USER_AGENT" help:"The user agent"`
	Locale     string `type:"string" default:"" env:"HEARCHCO_LOCALE" help:"Locale string specifying result language and region preference. The format is en_US"`
	SafeSearch bool   `type:"bool" default:"false" env:"HEARCHCO_SAFE_SEARCH" help:"Whether to use safe search"`
	Mobile     bool   `type:"bool" default:"false" env:"HEARCHCO_MOBILE" help:"Whether to gear results towards mobile"`
	// profiler
	CPUProfile    bool `type:"bool" default:"false" env:"HEARCHCO_CPUPROFILE" help:"Use cpu profiling"`
	HeapProfile   bool `type:"bool" default:"false" env:"HEARCHCO_HEAPPROFILE" help:"Use heap profiling"`
	GORProfile    bool `type:"bool" default:"false" env:"HEARCHCO_GORPROFILE" help:"Use goroutine profiling"`
	ThreadProfile bool `type:"bool" default:"false" env:"HEARCHCO_THREADPROFILE" help:"Use threadcreate profiling"`
	AllocProfile  bool `type:"bool" default:"false" env:"HEARCHCO_MEMALLOCPROFILE" help:"Use alloc profiling"`
	BlockProfile  bool `type:"bool" default:"false" env:"HEARCHCO_BLOCKPROFILE" help:"Use block profiling"`
	MutexProfile  bool `type:"bool" default:"false" env:"HEARCHCO_MUTEXPROFILE" help:"Use mutex profiling"`
	ClockProfile  bool `type:"bool" default:"false" env:"HEARCHCO_CLOCKPROFILE" help:"Use clock profiling"`
	TraceProfile  bool `type:"bool" default:"false" env:"HEARCHCO_TRACEPROFILE" help:"Use trace profiling"`
	ServeProfiler bool `type:"bool" default:"false" env:"HEARCHCO_SERVEPROFILER" help:"Run the profiler and serve at /debug/pprof/ http endpoint"`
}

var (
	// release variables
	Version   string
	Timestamp string
	GitCommit string
)

type Globals struct {
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
