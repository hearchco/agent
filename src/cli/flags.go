package cli

type Flags struct {
	Version    versionFlag `name:"version" help:"Print version information and quit"`
	Pretty     bool        `type:"bool" default:"false" env:"HEARCHCO_PRETTY" help:"Make logs pretty"`
	Verbosity  int8        `type:"counter" default:"0" short:"v" env:"HEARCHCO_VERBOSITY" help:"Log level verbosity"`
	ConfigPath string      `type:"path" default:"hearchco.yaml" env:"HEARCHCO_CONFIG_PATH" help:"Config file path"`

	Profiler
}

type Profiler struct {
	ProfilerServe  bool `type:"bool" default:"false" help:"Run the profiler and serve at /debug/pprof/ http endpoint"`
	ProfilerCPU    bool `type:"bool" default:"false" help:"Use cpu profiling"`
	ProfilerHeap   bool `type:"bool" default:"false" help:"Use heap profiling"`
	ProfilerGOR    bool `type:"bool" default:"false" help:"Use goroutine profiling"`
	ProfilerThread bool `type:"bool" default:"false" help:"Use threadcreate profiling"`
	ProfilerAlloc  bool `type:"bool" default:"false" help:"Use alloc profiling"`
	ProfilerBlock  bool `type:"bool" default:"false" help:"Use block profiling"`
	ProfilerMutex  bool `type:"bool" default:"false" help:"Use mutex profiling"`
	ProfilerClock  bool `type:"bool" default:"false" help:"Use clock profiling"`
	ProfilerTrace  bool `type:"bool" default:"false" help:"Use trace profiling"`
}
