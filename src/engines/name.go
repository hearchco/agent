package engines

import "strings"

type Name int

//go:generate go run github.com/dmarkham/enumer -type=Name
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
