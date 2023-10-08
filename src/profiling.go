package main

import (
	"github.com/pkg/profile"
	"github.com/rs/zerolog/log"
)

type profiler struct {
	enabled bool
	profile func(p *profile.Profile)
}

func runProfiler(amProfiling *bool) func() {
	/*
		goroutine — stack traces of all current goroutines
		heap — a sampling of memory allocations of live objects
		allocs — a sampling of all past memory allocations
		threadcreate — stack traces that led to the creation of new OS threads
		block — stack traces that led to blocking on synchronization primitives
		mutex — stack traces of holders of contended mutexes
	*/

	profilers := [...]profiler{{
		enabled: cli.CPUProfile,
		profile: profile.CPUProfile,
	}, {
		enabled: cli.HeapProfile,
		profile: profile.MemProfileHeap,
	}, {
		enabled: cli.GORProfile,
		profile: profile.GoroutineProfile,
	}, {
		enabled: cli.ThreadProfile,
		profile: profile.ThreadcreationProfile,
	}, {
		enabled: cli.BlockProfile,
		profile: profile.BlockProfile,
	}, {
		enabled: cli.AllocProfile,
		profile: profile.MemProfileAllocs,
	}, {
		enabled: cli.MutexProfile,
		profile: profile.MutexProfile,
	}, {
		enabled: cli.ClockProfile,
		profile: profile.ClockProfile,
	}, {
		enabled: cli.TraceProfile,
		profile: profile.TraceProfile,
	}}

	profilerToRun := profiler{
		enabled: false,
	}
	for _, p := range profilers {
		if profilerToRun.enabled && p.enabled {
			log.Fatal().Msg("Only one profiler can be run at a time.")
			return func() {}
		} else if p.enabled {
			profilerToRun = p
		}
	}

	if profilerToRun.enabled {
		*amProfiling = true
	} else {
		return func() {}
	}

	return func() {
		profile.Start(profilerToRun.profile, profile.ProfilePath("./profiling/")).Stop()
	}
}
