package startpage

import "github.com/hearchco/hearchco/src/search/engines"

var Info = engines.Info{
	Domain:         "www.startpage.com",
	Name:           engines.STARTPAGE,
	URL:            "https://www.startpage.com/do/search?language=english&q=",
	ResultsPerPage: 10,
}

var dompaths = engines.DOMPaths{
	Result:      "section > div.w-gl__result",
	Link:        "a.result-link",
	Title:       "a.w-gl__result-title",
	Description: "p.w-gl__description",
}

var Support = engines.SupportedSettings{
	SafeSearch: true,
}
