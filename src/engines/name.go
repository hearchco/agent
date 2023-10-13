package engines

import "strings"

type Name uint8

//go:generate enumer -type=Name -json -text -yaml -sql
//go:generate go run github.com/tminaorg/brzaguza/generate/searcher -type=Name -packagename search -output ../search/engine_searcher.go
const (
	UNDEFINED Name = iota
	BING
	BRAVE
	DUCKDUCKGO
	ETOOLS
	GOOGLE
	MOJEEK
	PRESEARCH
	QWANT
	STARTPAGE
	SWISSCOWS
	YAHOO
	YANDEX
	YEP
)

func (n Name) ToLower() string {
	return strings.ToLower(n.String())
}
