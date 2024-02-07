package etools

import "github.com/hearchco/hearchco/src/search/engines"

const pageURL string = "https://www.etools.ch/search.do?page="

var Info engines.Info = engines.Info{
	Domain:         "www.etools.ch",
	Name:           engines.ETOOLS,
	URL:            "https://www.etools.ch/searchSubmit.do",
	ResultsPerPage: 10,
	Crawlers:       []engines.Name{engines.BING, engines.GOOGLE, engines.MOJEEK, engines.YANDEX},
}

var dompaths engines.DOMPaths = engines.DOMPaths{
	Result:      "table.result > tbody > tr",
	Link:        "td.record > a",
	Description: "td.record > div.text",
}

var Support engines.SupportedSettings = engines.SupportedSettings{
	SafeSearch: true,
}
