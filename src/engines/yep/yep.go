package yep

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/bucket"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/engines"
	"github.com/hearchco/hearchco/src/search/parse"
	"github.com/hearchco/hearchco/src/sedefaults"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) error {
	if err := sedefaults.Prepare(Info.Name, &options, &settings, &Support, &Info, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	sedefaults.InitializeCollectors(&col, &pagesCol, &options, &timings)

	sedefaults.PagesColRequest(Info.Name, pagesCol, ctx)
	sedefaults.PagesColError(Info.Name, pagesCol)
	sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	sedefaults.ColRequest(Info.Name, col, ctx)
	sedefaults.ColError(Info.Name, col)

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

	localeParam := getLocale(&options)
	nRequested := settings.RequestedResultsPerPage
	safeSearchParam := getSafeSearch(&options)

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

	sedefaults.DoGetRequest(urll, anonUrll, nil, col, Info.Name, &retError)

	col.Wait()
	pagesCol.Wait()

	return retError
}

func getLocale(options *engines.Options) string {
	locale := strings.Split(options.Locale, "_")[1]
	return "&gl=" + locale
}

func getSafeSearch(options *engines.Options) string {
	if options.SafeSearch {
		return "&safeSearch=strict"
	}
	return "&safeSearch=off"
}
