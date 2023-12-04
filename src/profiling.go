package main

import (
	"log"

	"github.com/pkg/profile"
	"github.com/hearchco/hearchco/src/cli"
)

type profiler struct {
	enabled bool
	profile func(p *profile.Profile)
}

func runProfiler(cliFlags *cli.Flags) (bool, func()) {
	/*
		goroutine — stack traces of all current goroutines
		heap — a sampling of memory allocations of live objects
		allocs — a sampling of all past memory allocations
		threadcreate — stack traces that led to the creation of new OS threads
		block — stack traces that led to blocking on synchronization primitives
		mutex — stack traces of holders of contended mutexes
	*/

	profilers := [...]profiler{{
		enabled: cliFlags.CPUProfile,
		profile: profile.CPUProfile,
	}, {
		enabled: cliFlags.HeapProfile,
		profile: profile.MemProfileHeap,
	}, {
		enabled: cliFlags.GORProfile,
		profile: profile.GoroutineProfile,
	}, {
		enabled: cliFlags.ThreadProfile,
		profile: profile.ThreadcreationProfile,
	}, {
		enabled: cliFlags.BlockProfile,
		profile: profile.BlockProfile,
	}, {
		enabled: cliFlags.AllocProfile,
		profile: profile.MemProfileAllocs,
	}, {
		enabled: cliFlags.MutexProfile,
		profile: profile.MutexProfile,
	}, {
		enabled: cliFlags.ClockProfile,
		profile: profile.ClockProfile,
	}, {
		enabled: cliFlags.TraceProfile,
		profile: profile.TraceProfile,
	}}

	profilerToRun := profiler{enabled: false}
	for _, p := range profilers {
		if profilerToRun.enabled && p.enabled {
			log.Fatal("main.runProfiler(): only one profiler can be run at a time.")
			return false, func() {}
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
