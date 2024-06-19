package startpage

import (
	"github.com/hearchco/agent/src/search/engines"
)

const (
	seName    = engines.STARTPAGE
	searchURL = "https://www.startpage.com/sp/search"
)

var origins = [...]engines.Name{engines.STARTPAGE, engines.GOOGLE}
