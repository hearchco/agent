package yep

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/engines/_sedefaults"
	"github.com/hearchco/hearchco/src/search/parse"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, &options, &settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

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
			bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol)
			counter += 1
		}
	})

	localeParam := getLocale(options)
	nRequested := settings.RequestedResultsPerPage
	safeSearchParam := getSafeSearch(options)

	var urll string
	if nRequested == Info.ResultsPerPage {
		urll = Info.URL + "client=web" + localeParam + "&no_correct=false&q=" + query + safeSearchParam + "&type=web"
	} else {
		urll = Info.URL + "client=web" + localeParam + "&limit=" + strconv.Itoa(nRequested) + "&no_correct=false&q=" + query + safeSearchParam + "&type=web"
	}
	var anonUrll string
	if nRequested == Info.ResultsPerPage {
		anonUrll = Info.URL + "client=web" + localeParam + "&no_correct=false&q=" + anonymize.String(query) + safeSearchParam + "&type=web"
	} else {
		anonUrll = Info.URL + "client=web" + localeParam + "&limit=" + strconv.Itoa(nRequested) + "&no_correct=false&q=" + anonymize.String(query) + safeSearchParam + "&type=web"
	}

	retErrors := make([]error, options.MaxPages)

	err = _sedefaults.DoGetRequest(urll, anonUrll, nil, col, Info.Name)
	retErrors[0] = err

	// TODO: missing pages request loop?

	col.Wait()
	pagesCol.Wait()

	realRetErrors := make([]error, 0)
	for _, err := range retErrors {
		if err != nil {
			realRetErrors = append(realRetErrors, err)
		}
	}
	return realRetErrors
}

func getLocale(options engines.Options) string {
	locale := strings.Split(options.Locale, "_")[1]
	return "&gl=" + locale
}

func getSafeSearch(options engines.Options) string {
	if options.SafeSearch {
		return "&safeSearch=strict"
	}
	return "&safeSearch=off"
}
