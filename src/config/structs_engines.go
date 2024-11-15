package config

import (
	"github.com/hearchco/agent/src/search/engines"
)

// ReaderEngineConfig is format in which the config is read from the config file and environment variables.
// Used to disable certain search types for an engine. By default, all types are enabled.
type ReaderEngineConfig struct {
	NoWeb         bool // Whether this engine is disallowed to do web searches.
	NoImages      bool // Whether this engine is disallowed to do image searches.
	NoSuggestions bool // Whether this engine is disallowed to do suggestion searches.
}

// Slices of disabled engines for each search type, by default these are empty.
type EngineConfig struct {
	NoWeb         []engines.Name
	NoImages      []engines.Name
	NoSuggestions []engines.Name
}
