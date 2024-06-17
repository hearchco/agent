package profiler

import (
	"github.com/pkg/profile"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/cli"
)

type profiler struct {
	enabled bool
	profile func(p *profile.Profile)
}

func Run(cliFlags cli.Flags) (bool, func()) {
	/*
		goroutine — stack traces of all current goroutines
		heap — a sampling of memory allocations of live objects
		allocs — a sampling of all past memory allocations
		threadcreate — stack traces that led to the creation of new OS threads
		block — stack traces that led to blocking on synchronization primitives
		mutex — stack traces of holders of contended mutexes
	*/

	profilers := [...]profiler{{
		enabled: cliFlags.ProfilerCPU,
		profile: profile.CPUProfile,
	}, {
		enabled: cliFlags.ProfilerHeap,
		profile: profile.MemProfileHeap,
	}, {
		enabled: cliFlags.ProfilerGOR,
		profile: profile.GoroutineProfile,
	}, {
		enabled: cliFlags.ProfilerThread,
		profile: profile.ThreadcreationProfile,
	}, {
		enabled: cliFlags.ProfilerBlock,
		profile: profile.BlockProfile,
	}, {
		enabled: cliFlags.ProfilerAlloc,
		profile: profile.MemProfileAllocs,
	}, {
		enabled: cliFlags.ProfilerMutex,
		profile: profile.MutexProfile,
	}, {
		enabled: cliFlags.ProfilerClock,
		profile: profile.ClockProfile,
	}, {
		enabled: cliFlags.ProfilerTrace,
		profile: profile.TraceProfile,
	}}

	profilerToRun := profiler{enabled: false}
	for _, p := range profilers {
		if profilerToRun.enabled && p.enabled {
			log.Fatal().
				Caller().
				Msg("Only one profiler can be run at a time")
			// ^FATAL
		} else if p.enabled {
			profilerToRun = p
		}
	}
	if !profilerToRun.enabled {
		return false, func() {}
	}

	p := profile.Start(profilerToRun.profile, profile.ProfilePath("./profiling/"), profile.NoShutdownHook)
	return true, func() {
		p.Stop()
	}
}
