package category

// CategoryJSON is format in which the config is passed from the user.
type CategoryJSON struct {
	Engines map[string]EngineJSON `koanf:"engines"`
	Ranking Ranking               `koanf:"ranking"`
	Timings TimingsJSON           `koanf:"timings"`
}

// EngineJSON is format in which the config is passed from the user.
type EngineJSON struct {
	// If false, the engine will not be used and other options will be ignored.
	// This adds the engine to engines slice during conversion.
	Enabled bool `koanf:"enabled"`
	// If true, the engine will be awaited unless the hard timeout is reached.
	// This adds the engine to required engines slice during conversion.
	Required bool `koanf:"required"`
	// If true, the fastest engine that has this engine in "Origins" will be awaited unless the hard timeout is reached.
	// This means that we want to get results from this engine or any engine that has this engine in "Origins", whichever responds the fastest.
	// This adds the engine to required engines by origin slice during conversion.
	RequiredByOrigin bool `koanf:"requiredbyorigin"`
	// If true, the engine will be awaited unless the preferred timeout is reached.
	// This adds the engine to preferred engines slice during conversion.
	Preferred bool `koanf:"preferred"`
	// If true, the fastest engine that has this engine in "Origins" will be awaited unless the preferred timeout is reached.
	// This means that we want to get results from this engine or any engine that has this engine in "Origins", whichever responds the fastest.
	// This adds the engine to preferred by origin slice during conversion.
	PreferredByOrigin bool `koanf:"preferredbyorigin"`
}

// TimingsJSON is format in which the config is passed from the user.
// In <number><unit> format.
// Example: 1s, 1m, 1h, 1d, 1w, 1M, 1y.
// If unit is not specified, it is assumed to be milliseconds.
type TimingsJSON struct {
	// Maximum amount of time to wait for the PreferredEngines (or ByOrigin) to respond.
	// If the search is still waiting for the RequiredEngines (or ByOrigin) after this time, the search will continue.
	PreferredTimeout string `koanf:"preferredtimeout"`
	// Hard timeout after which the search is forcefully stopped (even if the engines didn't respond).
	HardTimeout string `koanf:"hardtimeout"`
}
