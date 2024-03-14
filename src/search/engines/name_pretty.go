package engines

import (
	"strings"

	"github.com/rs/zerolog/log"
)

var prettyNames = [...]string{
	"Undefined",
	"Bing",
	"Bing Images",
	"Brave",
	"DuckDuckGo",
	"Etools",
	"Google",
	"Google Images",
	"Google Scholar",
	"Mojeek",
	"Presearch",
	"Qwant",
	"Startpage",
	"Swisscows",
	"Yahoo",
	"Yep",
}

func (n Name) Pretty() string {
	return prettyNames[n]
}

func init() {
	if len(prettyNames) != len(_NameValues) {
		log.Panic().
			Msg("PrettyNames and _NameValues have different lengths")
	}
	for i, pn := range prettyNames {
		name := _NameValues[i]
		prettyNameLowered := strings.ReplaceAll(strings.ToLower(pn), " ", "")
		if name.ToLower() != prettyNameLowered {
			log.Panic().
				Str("name", name.String()).
				Str("prettyname", pn).
				Msg("PrettyNames and _NameValues are not in sync")
		}
	}
}

// return engines' pretty names without Undefined
func PrettyNames() []string {
	if len(prettyNames) > 1 {
		return prettyNames[1:]
	}
	return prettyNames[:]
}
