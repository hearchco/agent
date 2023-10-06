package main

import (
	"github.com/pkg/profile"
	"github.com/rs/zerolog/log"
)

func runProfiler(amProfiling *bool) func() {
	/*
		goroutine — stack traces of all current goroutines
		heap — a sampling of memory allocations of live objects
		allocs — a sampling of all past memory allocations
		threadcreate — stack traces that led to the creation of new OS threads
		block — stack traces that led to blocking on synchronization primitives
		mutex — stack traces of holders of contended mutexes
	*/

	var cpup interface{ Stop() }
	var gorp interface{ Stop() }
	var blockp interface{ Stop() }
	var threadp interface{ Stop() }
	var heapp interface{ Stop() }
	var allocp interface{ Stop() }
	var mutexp interface{ Stop() }
	var tracep interface{ Stop() }
	var clockp interface{ Stop() }

	profileCnt := uint(0)
	if cli.CPUProfile {
		profileCnt += 1
	}
	if cli.HeapProfile {
		profileCnt += 1
	}
	if cli.GORProfile {
		profileCnt += 1
	}
	if cli.ThreadProfile {
		profileCnt += 1
	}
	if cli.BlockProfile {
		profileCnt += 1
	}
	if cli.AllocProfile {
		profileCnt += 1
	}
	if cli.MutexProfile {
		profileCnt += 1
	}
	if cli.ClockProfile {
		profileCnt += 1
	}
	if cli.TraceProfile {
		profileCnt += 1
	}

	if profileCnt > 1 {
		log.Fatal().Msg("only one profiler can be run at a time.")
		return func() {}
	}

	if cli.CPUProfile {
		cpup = profile.Start(profile.CPUProfile, profile.ProfilePath("./profiling/"))
	}
	if cli.HeapProfile {
		heapp = profile.Start(profile.MemProfileHeap, profile.ProfilePath("./profiling/"))
	}
	if cli.GORProfile {
		gorp = profile.Start(profile.GoroutineProfile, profile.ProfilePath("./profiling/"))
	}
	if cli.ThreadProfile {
		threadp = profile.Start(profile.ThreadcreationProfile, profile.ProfilePath("./profiling/"))
	}
	if cli.BlockProfile {
		blockp = profile.Start(profile.BlockProfile, profile.ProfilePath("./profiling/"))
	}
	if cli.AllocProfile {
		allocp = profile.Start(profile.MemProfileAllocs, profile.ProfilePath("./profiling/"))
	}
	if cli.MutexProfile {
		mutexp = profile.Start(profile.MutexProfile, profile.ProfilePath("./profiling/"))
	}
	if cli.ClockProfile {
		clockp = profile.Start(profile.ClockProfile, profile.ProfilePath("./profiling/"))
	}
	if cli.TraceProfile {
		tracep = profile.Start(profile.TraceProfile, profile.ProfilePath("./profiling/"))
	}

	if profileCnt == 1 {
		*amProfiling = true
	}

	return func() {
		if cli.CPUProfile {
			cpup.Stop()
		}
		if cli.HeapProfile {
			heapp.Stop()
		}
		if cli.GORProfile {
			gorp.Stop()
		}
		if cli.ThreadProfile {
			threadp.Stop()
		}
		if cli.BlockProfile {
			blockp.Stop()
		}
		if cli.AllocProfile {
			allocp.Stop()
		}
		if cli.MutexProfile {
			mutexp.Stop()
		}
		if cli.ClockProfile {
			clockp.Stop()
		}
		if cli.TraceProfile {
			tracep.Stop()
		}
	}
}
