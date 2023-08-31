package engines

import "strings"

type Name int64

const (
	Undefined Name = iota
	Google
	Mojeek
	DuckDuckGo
	Qwant
	Etools
	Swisscows
	Brave
	Bing
	Startpage
	Yandex
	Yep
)

func (n Name) String() string {
	switch n {
	case Google:
		return "google"
	case Mojeek:
		return "mojeek"
	case DuckDuckGo:
		return "duckduckgo"
	case Qwant:
		return "qwant"
	case Etools:
		return "etools"
	case Swisscows:
		return "swisscows"
	case Brave:
		return "brave"
	case Bing:
		return "bing"
	case Startpage:
		return "startpage"
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
