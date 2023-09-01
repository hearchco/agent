package engines

import "strings"

type Name int64

const (
	Undefined Name = iota
<<<<<<< HEAD
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
=======
	Google
	Mojeek
	DuckDuckGo
	Qwant
	Etools
	Swisscows
	Brave
	Bing
	Startpage
>>>>>>> a0694f64bea02e5a0f681448bf43971a1d74148e
	Yandex
	Yep
)

func (n Name) String() string {
	switch n {
<<<<<<< HEAD
	case Bing:
		return "Bing"
	case Brave:
		return "Brave"
	case DuckDuckGo:
		return "DuckDuckGo"
	case Etools:
		return "Etools"
=======
>>>>>>> a0694f64bea02e5a0f681448bf43971a1d74148e
	case Google:
		return "Google"
	case Mojeek:
		return "Mojeek"
<<<<<<< HEAD
	case Presearch:
		return "Presearch"
	case Qwant:
		return "Qwant"
	case Startpage:
		return "Startpage"
	case Swisscows:
		return "Swisscows"
=======
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
>>>>>>> a0694f64bea02e5a0f681448bf43971a1d74148e
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

func (n Name) Equals(s string) bool {
	return n.ToLower() == strings.ToLower(s)
}

func ConvertToName(s string) Name {
	switch {
	case Google.Equals(s):
		return Google
	case Mojeek.Equals(s):
		return Mojeek
	case DuckDuckGo.Equals(s):
		return DuckDuckGo
	case Qwant.Equals(s):
		return Qwant
	case Etools.Equals(s):
		return Etools
	case Swisscows.Equals(s):
		return Swisscows
	case Brave.Equals(s):
		return Brave
	case Bing.Equals(s):
		return Bing
	case Startpage.Equals(s):
		return Startpage
	case Yandex.Equals(s):
		return Yandex
	case Yep.Equals(s):
		return Yep
	default:
		return Undefined
	}
}
