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
<<<<<<< HEAD
	Yahoo
=======
>>>>>>> f7daa83 (Fixed naming order, fixed name for Presearch)
)

func (n Name) String() string {
	switch n {
	case Bing:
<<<<<<< HEAD
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

=======
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

>>>>>>> f7daa83 (Fixed naming order, fixed name for Presearch)
func (n Name) ToLower() string {
	return strings.ToLower(n.String())
}

func (n Name) Equals(s string) bool {
<<<<<<< HEAD
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
=======
	return n.ToUpper() == strings.ToUpper(s)
>>>>>>> f7daa83 (Fixed naming order, fixed name for Presearch)
}
