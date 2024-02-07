package engines

import "strings"

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

// Returns Engine Names without UNDEFINED
func Names() []Name {
	return _NameValues[1:]
}

func (n Name) ToLower() string {
	return strings.ToLower(n.String())
}
