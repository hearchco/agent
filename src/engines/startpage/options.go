package startpage

import "github.com/tminaorg/brzaguza/src/structures"

var Info structures.SEInfo = structures.SEInfo{
	Domain:         "www.startpage.com",
	Name:           "Startpage",
	URL:            "https://www.startpage.com/sp/search?q=",
	ResultsPerPage: 10,
	Crawlers:       []structures.EngineName{structures.Google},
}

var dompaths structures.SEDOMPaths = structures.SEDOMPaths{
	Result:      "section > div.w-gl__result",
	Link:        "a.result-link",
	Title:       "a.w-gl__result-title",
	Description: "p.w-gl__description",
}

var Support structures.SupportedSettings = structures.SupportedSettings{
	SafeSearch: true,
}
