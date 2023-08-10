package etools

import "github.com/tminaorg/brzaguza/src/structures"

var Info structures.SEInfo = structures.SEInfo{
	Domain:     "www.etools.ch",
	Name:       "Etools",
	URL:        "https://www.etools.ch/searchSubmit.do",
	ResPerPage: 10,
	Crawlers:   []structures.EngineName{structures.Bing, structures.Google, structures.Mojeek, structures.Yandex},
}

var dompaths structures.SEDOMPaths = structures.SEDOMPaths{
	Result:      "table.result > tbody > tr",
	Link:        "td.record > a",
	Description: "td.record > div.text",
}
