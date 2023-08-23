package etools

import "github.com/tminaorg/brzaguza/src/engines"

const pageURL string = "https://www.etools.ch/search.do?page="

var Info engines.Info = engines.Info{
	Domain:         "www.etools.ch",
	Name:           "Etools",
	URL:            "https://www.etools.ch/searchSubmit.do",
	ResultsPerPage: 10,
	Crawlers:       []engines.Name{engines.Bing, engines.Google, engines.Mojeek, engines.Yandex},
}

var dompaths engines.DOMPaths = engines.DOMPaths{
	Result:      "table.result > tbody > tr",
	Link:        "td.record > a",
	Description: "td.record > div.text",
}

var Support engines.SupportedSettings = engines.SupportedSettings{}
