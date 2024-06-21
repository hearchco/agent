package duckduckgo

import (
	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
)

func (se Engine) Suggest(query string, locale options.Locale, sugChan chan result.SuggestionScraped) ([]error, bool) {
	return nil, false
}
