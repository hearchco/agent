package engines

import "strings"

type Name int

//go:generate stringer -type=Name
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

func ConvertToName(s string) Name {
	switch strings.ToLower(s) {
	case Google.ToLower():
		return Google
	case Mojeek.ToLower():
		return Mojeek
	case DuckDuckGo.ToLower():
		return DuckDuckGo
	case Qwant.ToLower():
		return Qwant
	case Etools.ToLower():
		return Etools
	case Swisscows.ToLower():
		return Swisscows
	case Brave.ToLower():
		return Brave
	case Bing.ToLower():
		return Bing
	case Startpage.ToLower():
		return Startpage
	case Yahoo.ToLower():
		return Yahoo
	case Yandex.ToLower():
		return Yandex
	case Yep.ToLower():
		return Yep
	default:
		return Undefined
	}
}
