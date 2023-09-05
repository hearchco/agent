package engines

import "strings"

type Name uint8

//go:generate enumer -type=Name
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
