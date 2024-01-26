package engines

import (
	"log"
	"strings"
)

type Name uint8

//go:generate enumer -type=Name -json -text -yaml -sql
//go:generate go run github.com/hearchco/hearchco/generate/searcher -type=Name -packagename search -output ../search/engine_searcher.go
const (
	UNDEFINED Name = iota
	BING
	BRAVE
	DUCKDUCKGO
	ETOOLS
	GOOGLE
	GOOGLESCHOLAR
	MOJEEK
	PRESEARCH
	QWANT
	STARTPAGE
	SWISSCOWS
	YAHOO
	YEP
)

func initPrettyName() {
	// laid out like this instead of ["Undefined", "Bing", "Brave", ...] so the index doesn't matter
	PrettyName = make([]string, len(_NameValues))
	PrettyName[UNDEFINED] = "UNDEFINED"
	PrettyName[BING] = "Bing"
	PrettyName[BRAVE] = "Brave"
	PrettyName[DUCKDUCKGO] = "DuckDuckGo"
	PrettyName[ETOOLS] = "Etools"
	PrettyName[GOOGLE] = "Google"
	PrettyName[GOOGLESCHOLAR] = "GoogleScholar"
	PrettyName[MOJEEK] = "Mojeek"
	PrettyName[PRESEARCH] = "Presearch"
	PrettyName[QWANT] = "Qwant"
	PrettyName[STARTPAGE] = "Startpage"
	PrettyName[SWISSCOWS] = "Swisscows"
	PrettyName[YAHOO] = "Yahoo"
	PrettyName[YEP] = "Yep"

	// Check if all search engines have a pretty name set
	for _, eng := range NameValues() {
		if PrettyName[eng] == "" {
			log.Fatalf("engines.init() (names.go): %v doesn't have a pretty name set.", eng)
			// ^FATAL
		}
	}
}

func init() {
	initPrettyName()
}

var PrettyName []string

// Returns Engine Names without UNDEFINED
func Names() []Name {
	return _NameValues[1:]
}

func (n Name) ToLower() string {
	return strings.ToLower(n.String())
}
