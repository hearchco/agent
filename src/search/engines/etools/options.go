package etools

import (
	"github.com/hearchco/hearchco/src/search/engines"
)

const pageURL string = "https://www.etools.ch/search.do?page="

var Info = engines.Info{
	Domain:         "www.etools.ch",
	Name:           engines.ETOOLS,
	URL:            "https://www.etools.ch/searchSubmit.do",
	ResultsPerPage: 10,
}

var dompaths = engines.DOMPaths{
	Result:      "table.result > tbody > tr",
	Link:        "td.record > a",
	Title:       "td.record > a",
	Description: "td.record > div.text",
}

var Support = engines.SupportedSettings{
	SafeSearch: true,
}
