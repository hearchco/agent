package engines

import "strings"

type Name int

//go:generate enumer -type=Name -json -text -sql
//go:generate go run github.com/hearchco/agent/generate/enginer -type=Name -packagename search -output ../engine_enginer.go
const (
	UNDEFINED     Name = iota
	BING               // enginer,searcher
	BINGIMAGES         // enginer,searcher
	BRAVE              // enginer,searcher
	DUCKDUCKGO         // enginer,searcher,suggester
	ETOOLS             // enginer,searcher
	GOOGLE             // enginer,searcher,suggester
	GOOGLEIMAGES       // enginer,searcher
	GOOGLESCHOLAR      // enginer,searcher
	MOJEEK             // enginer,searcher
	PRESEARCH          // enginer,searcher
	QWANT              // enginer,searcher
	STARTPAGE          // enginer,searcher
	SWISSCOWS          // enginer,searcher
	YAHOO              // enginer,searcher
	YEP
)

// Returns engine names without UNDEFINED.
func Names() []Name {
	return _NameValues[1:]
}

func (n Name) ToLower() string {
	return strings.ToLower(n.String())
}
