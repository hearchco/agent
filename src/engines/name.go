package engines

import "strings"

type Name uint8

//go:generate enumer -type=Name -json -text -yaml -sql
//go:generate go run github.com/tminaorg/brzaguza/generate/searcher -type=Name -packagename search -output ../search/engine_searcher.go
const (
	Undefined Name = iota
	Bing
	Brave
	DuckDuckGo
	Etools
	Google
	Mojeek
	Presearch
	Qwant
	Startpage
	Swisscows
	Yahoo
	Yandex
	Yep
)

func (n Name) ToLower() string {
	return strings.ToLower(n.String())
}
