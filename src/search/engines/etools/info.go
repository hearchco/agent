package etools

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName    = engines.ETOOLS
	searchURL = "https://www.etools.ch/searchSubmit.do"
	pageURL   = "https://www.etools.ch/search.do"
)

var origins = [...]engines.Name{engines.ETOOLS, engines.BING, engines.BRAVE, engines.DUCKDUCKGO, engines.GOOGLE, engines.MOJEEK, engines.QWANT, engines.YAHOO}
