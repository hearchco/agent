package search

import (
	"context"

	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/engines/bing"
	"github.com/tminaorg/brzaguza/src/engines/brave"
	"github.com/tminaorg/brzaguza/src/engines/duckduckgo"
	"github.com/tminaorg/brzaguza/src/engines/etools"
	"github.com/tminaorg/brzaguza/src/engines/google"
	"github.com/tminaorg/brzaguza/src/engines/mojeek"
	"github.com/tminaorg/brzaguza/src/engines/presearch"
	"github.com/tminaorg/brzaguza/src/engines/qwant"
	"github.com/tminaorg/brzaguza/src/engines/startpage"
	"github.com/tminaorg/brzaguza/src/engines/swisscows"
	"github.com/tminaorg/brzaguza/src/engines/yahoo"
	"github.com/tminaorg/brzaguza/src/engines/yep"
)

type EngineSearch func(context.Context, string, *bucket.Relay, engines.Options, config.Settings) error

func NewEngineStarter() []EngineSearch {
	// I would be a happy man if we could do this function at compile time
	mm := make([]EngineSearch, 100)

	//alphabetically sorted
	mm[engines.Bing] = bing.Search
	mm[engines.Brave] = brave.Search
	mm[engines.DuckDuckGo] = duckduckgo.Search
	mm[engines.Etools] = etools.Search
	mm[engines.Google] = google.Search
	mm[engines.Mojeek] = mojeek.Search
	mm[engines.Presearch] = presearch.Search
	mm[engines.Qwant] = qwant.Search
	mm[engines.Startpage] = startpage.Search
	mm[engines.Swisscows] = swisscows.Search
	mm[engines.Yahoo] = yahoo.Search
	//mm[engines.Yandex] = yandex.Search
	mm[engines.Yep] = yep.Search

	return mm
}
