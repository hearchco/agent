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
)

func (n Name) String() string {
	switch n {
	case Bing:
		return "bing"
	case Brave:
		return "brave"
	case DuckDuckGo:
		return "duckduckgo"
	case Etools:
		return "etools"
	case Google:
		return "google"
	case Mojeek:
		return "mojeek"
	case Presearch:
		return "presearch"
	case Qwant:
		return "qwant"
	case Startpage:
		return "startpage"
	case Swisscows:
		return "swisscows"
	case Yandex:
		return "yandex"
	case Yep:
		return "yep"
	default:
		return "undefined"
	}
}

func (n Name) ToUpper() string {
	return strings.ToUpper(n.String())
}

func (n Name) ToLower() string {
	return strings.ToLower(n.String())
}

func (n Name) Equals(s string) bool {
	return n.ToUpper() == strings.ToUpper(s)
}
