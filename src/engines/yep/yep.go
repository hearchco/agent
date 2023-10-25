package yep

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/search/parse"
	"github.com/tminaorg/brzaguza/src/sedefaults"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) error {
	if err := sedefaults.Prepare(Info.Name, &options, &settings, &Support, &Info, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	sedefaults.InitializeCollectors(&col, &pagesCol, &options, &timings)

	sedefaults.PagesColRequest(Info.Name, pagesCol, &ctx, &retError)
	sedefaults.PagesColError(Info.Name, pagesCol)
	sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	sedefaults.ColRequest(Info.Name, col, &ctx, &retError)
	sedefaults.ColError(Info.Name, col, &retError)

	col.OnRequest(func(r *colly.Request) {
		r.Headers.Del("Accept")
	})

	col.OnResponse(func(r *colly.Response) {
		content := parseJSON(r.Body)

		counter := 1
		for _, result := range content.Results {
			if result.TType != "Organic" {
				continue
			}

			goodURL := parse.ParseURL(result.URL)
			goodTitle := parse.ParseTextWithHTML(result.Title)
			goodDescription := parse.ParseTextWithHTML(result.Snippet)

			res := bucket.MakeSEResult(goodURL, goodTitle, goodDescription, Info.Name, 1, counter)
			bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
			counter += 1
		}
	})

	locale := getLocale(&options)
	nRequested := settings.RequestedResultsPerPage
	safeSearch := getSafeSearch(&options)

	var apiURL string
	if nRequested == Info.ResultsPerPage {
		apiURL = Info.URL + "client=web&gl=" + locale + "&no_correct=false&q=" + query + "&safeSearch=" + safeSearch + "&type=web"
	} else {
		apiURL = Info.URL + "client=web&gl=" + locale + "&limit=" + strconv.Itoa(nRequested) + "&no_correct=false&q=" + query + "&safeSearch=" + safeSearch + "&type=web"
	}

	err := col.Request("GET", apiURL, nil, nil, nil)
	if err != nil && !engines.IsTimeoutError(err) {
		log.Error().Err(err).Msgf("%v: failed requesting with GET method", Info.Name)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}

func getLocale(options *engines.Options) string {
	locale := strings.Split(options.Locale, "-")[1]
	return locale
}

func getSafeSearch(options *engines.Options) string {
	if options.SafeSearch {
		return "strict"
	}
	return "off"
}
