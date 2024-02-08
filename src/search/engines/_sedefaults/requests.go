package _sedefaults

import (
	"fmt"
	"io"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

func DoGetRequest(urll string, anonurll string, colCtx *colly.Context, collector *colly.Collector, packageName engines.Name, retError *error) {
	log.Trace().
		Str("engine", packageName.String()).
		Str("url", anonurll).
		Msg("GET")
	err := collector.Request("GET", urll, nil, colCtx, nil)
	if err != nil {
		*retError = fmt.Errorf("%v.Search(): failed GET request to %v with %w", packageName.ToLower(), urll, err)
	}
}

func DoPostRequest(urll string, requestData io.Reader, colCtx *colly.Context, collector *colly.Collector, packageName engines.Name, retError *error) {
	log.Trace().
		Str("engine", packageName.String()).
		Str("url", urll).
		Msg("POST")
	err := collector.Request("POST", urll, requestData, colCtx, nil)
	if err != nil {
		*retError = fmt.Errorf("%v.Search(): failed POST request to %v and body %v. error %w", packageName.ToLower(), requestData, urll, err)
	}
}
