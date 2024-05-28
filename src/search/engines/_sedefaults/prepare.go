package _sedefaults

import (
	"context"

	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

func Prepare(ctx context.Context, info engines.Info, support engines.SupportedSettings, options engines.Options, settings config.Settings) (context.Context, error) {
	// TODO: move to config initialization
	if settings.RequestedResultsPerPage != 0 && !support.RequestedResultsPerPage {
		log.Panic().
			Caller().
			Str("engine", info.Name.String()).
			Int("requestedResultsPerPage", settings.RequestedResultsPerPage).
			Msg("Setting not supported by engine")
		// ^PANIC
	}

	if options.Locale != "" && !support.Locale {
		log.Debug().
			Str("engine", info.Name.String()).
			Str("locale", options.Locale).
			Msg("Setting not supported by engine")
	}

	if options.SafeSearch && !support.SafeSearch {
		log.Debug().
			Str("engine", info.Name.String()).
			Bool("safeSearch", options.SafeSearch).
			Msg("Setting not supported by engine")
	}

	if ctx == nil {
		return context.Background(), nil
	} else {
		return ctx, nil
	}
}
