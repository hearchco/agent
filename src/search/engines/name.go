package engines

import "strings"

type Name int

//go:generate enumer -type=Name -json -text
//go:generate go run github.com/hearchco/agent/generate/enginer -type=Name -packagename search -output ../engine_enginer.go
const (
	UNDEFINED     Name = iota
	BING               // enginer,websearcher,imagesearcher
	BRAVE              // enginer,websearcher
	DUCKDUCKGO         // enginer,websearcher,suggester
	ETOOLS             // enginer,websearcher
	GOOGLE             // enginer,websearcher,imagesearcher,suggester
	GOOGLESCHOLAR      // enginer,websearcher
	MOJEEK             // enginer,websearcher
	PRESEARCH          // enginer,websearcher
	QWANT              // enginer,websearcher
	STARTPAGE          // enginer,websearcher
	SWISSCOWS          // enginer,websearcher
	YAHOO              // enginer,websearcher
	YEP                // disabled
)

// Returns engine names without UNDEFINED.
func Names() []Name {
	return _NameValues[1:]
}

func (n Name) ToLower() string {
	return strings.ToLower(n.String())
}
