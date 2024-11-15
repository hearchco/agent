package _engines_test

import (
	"github.com/hearchco/agent/src/search/engines/options"
)

type TestCaseHasAnyResults struct {
	Query   string
	Options options.Options
}

type TestCaseContainsResults struct {
	Query      string
	ResultURLs []string
	Options    options.Options
}

type TestCaseRankedResults struct {
	Query      string
	ResultURLs []string
	Options    options.Options
}

func NewOpts() options.Options {
	return options.Options{
		Pages:      options.Pages{Start: 0, Max: 1},
		Locale:     options.LocaleDefault,
		SafeSearch: false,
	}
}
