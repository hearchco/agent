package _sedefaults

import (
	"context"

	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/useragent"
	"github.com/rs/zerolog/log"
)

func Prepare(seName engines.Name, options *engines.Options, settings *config.Settings, support *engines.SupportedSettings, info *engines.Info, ctx *context.Context) error {
	if ctx == nil {
		*ctx = context.Background()
	}

	if options.UserAgent == "" {
		options.UserAgent = useragent.RandomUserAgent()
	}
	log.Trace().
		Str("engine", seName.String()).
		Str("userAgent", options.UserAgent).
		Msg("Prepare")

	// TODO: move to config.SetupConfig
	if settings.RequestedResultsPerPage != 0 && !support.RequestedResultsPerPage {
		log.Panic().
			Str("engine", seName.String()).
			Int("requestedResultsPerPage", settings.RequestedResultsPerPage).
			Msg("_sedefaults.Prepare(): setting not supported by engine")
		// ^PANIC
	}
	if settings.RequestedResultsPerPage == 0 && support.RequestedResultsPerPage {
		// if its used in the code but not set, give it the default value
		settings.RequestedResultsPerPage = info.ResultsPerPage
	}

	if options.Mobile && !support.Mobile {
		options.Mobile = false // this line shouldn't matter [1]
		log.Debug().
			Str("engine", seName.String()).
			Bool("mobile", options.Mobile).
			Msg("Mobile set but not supported")
	}

	if options.Locale != "" && !support.Locale {
		options.Locale = config.DefaultLocale // [1]
		log.Debug().
			Str("engine", seName.String()).
			Str("locale", options.Locale).
			Msg("Locale set but not supported")
	}

	if options.Locale == "" && support.Locale {
		options.Locale = config.DefaultLocale
	}

	if options.SafeSearch && !support.SafeSearch {
		options.SafeSearch = false // [1]
		log.Debug().
			Str("engine", seName.String()).
			Bool("safeSearch", options.SafeSearch).
			Msg("SafeSearch set but not supported")
	}

	return nil
}
