package engines

import "strings"

type Name int

//go:generate enumer -type=Name -json -text -sql
//go:generate go run github.com/hearchco/agent/generate/searcher -type=Name -packagename search -output ../engine_searcher.go
const (
	UNDEFINED Name = iota
	BING
	BINGIMAGES
	BRAVE
	DUCKDUCKGO
	ETOOLS
	GOOGLE
	GOOGLEIMAGES
	GOOGLESCHOLAR
	MOJEEK
	PRESEARCH
	QWANT
	STARTPAGE
	SWISSCOWS
	YAHOO
	// YEP
)

// Returns engine names without UNDEFINED.
func Names() []Name {
	return _NameValues[1:]
}

func (n Name) ToLower() string {
	return strings.ToLower(n.String())
}
