package search

import (
	"context"

	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
)

type EngineSearch func(context.Context, string, *bucket.Relay, engines.Options, config.Settings) error
