package engines

import "strings"

type Name int64

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
	Yandex
	Yep
	Yahoo
)

func (n Name) String() string {
	switch n {
	case Bing:
		return "Bing"
	case Brave:
		return "Brave"
	case DuckDuckGo:
		return "DuckDuckGo"
	case Etools:
		return "Etools"
	case Google:
		return "Google"
	case Mojeek:
		return "Mojeek"
	case Presearch:
		return "Presearch"
	case Qwant:
		return "Qwant"
	case Startpage:
		return "Startpage"
	case Swisscows:
		return "Swisscows"
	case Yandex:
		return "Yandex"
	case Yep:
		return "Yep"
	case Yahoo:
		return "Yahoo"
	default:
		return "Undefined"
	}
}

func (n Name) ToLower() string {
	return strings.ToLower(n.String())
}
