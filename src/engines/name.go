package engines

import "strings"

type Name int

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
		return "Google"
	case Mojeek:
		return "Mojeek"
	case DuckDuckGo:
		return "DuckDuckGo"
	case Qwant:
		return "Qwant"
	case Etools:
		return "Etools"
	case Swisscows:
		return "Swisscows"
	case Brave:
		return "Brave"
	case Bing:
		return "Bing"
	case Startpage:
		return "Startpage"
	case Yandex:
		return "Yandex"
	case Yep:
		return "Yep"
	default:
		return "Undefined"
	}
}

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
	case Yandex.ToLower():
		return Yandex
	case Yep.ToLower():
		return Yep
	default:
		return Undefined
	}
}
