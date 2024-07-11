package config

import (
	"time"

	"github.com/hearchco/agent/src/exchange/engines"
)

// ReaderCategory is format in which the config is read from the config file and environment variables.
type ReaderExchange struct {
	REngines map[string]ReaderExchangeEngine `koanf:"engines"`
	RTimings ReaderExchangeTimings           `koanf:"timings"`
}
type Exchange struct {
	Engines []engines.Name
	Timings ExchangeTimings
}

// ReaderEngine is format in which the config is read from the config file and environment variables.
type ReaderExchangeEngine struct {
	// If false, the engine will not be used.
	Enabled bool `koanf:"enabled"`
}

// ReaderTimings is format in which the config is read from the config file and environment variables.
// In <number><unit> format.
// Example: 1s, 1m, 1h, 1d, 1w, 1M, 1y.
// If unit is not specified, it is assumed to be milliseconds.
type ReaderExchangeTimings struct {
	// Hard timeout after which the search is forcefully stopped (even if the engines didn't respond).
	HardTimeout string `koanf:"hardtimeout"`
}
type ExchangeTimings struct {
	// Hard timeout after which the search is forcefully stopped (even if the engines didn't respond).
	HardTimeout time.Duration
}
