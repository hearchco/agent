package brave

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
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, &options, &settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

	pageRankCounter := make([]int, options.Pages.Max)

	localeCookie := getLocale(options)
	safeSearchCookie := getSafeSearch(options)

	col.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", localeCookie)
		r.Headers.Add("Cookie", safeSearchCookie)
	})

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		linkText, titleText, descText := _sedefaults.FieldsFromDOM(e.DOM, &dompaths, Info.Name)

		if descText == "" {
			descText = e.DOM.Find("div.product > div.flex-hcenter > div > div[class=\"text-sm text-gray\"]").Text()
		}
		if descText == "" {
			descText = e.DOM.Find("p.snippet-description").Text()
		}
		descText = _sedefaults.SanitizeDescription(descText)

		pageIndex := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)
		page := pageIndex + options.Pages.Start

		res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[pageIndex]+1)
		valid := bucket.AddSEResult(&res, Info.Name, relay, &options, pagesCol)
		if valid {
			pageRankCounter[pageIndex]++
		}
	})

	retErrors := make([]error, 0, options.Pages.Max)

	// starts from at least 0
	for i := options.Pages.Start; i < options.Pages.Start+options.Pages.Max; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i-options.Pages.Start))

		// dynamic params
		pageParam := "&source=web"
		// i == 0 is the first page
		if i > 0 {
			pageParam = "&spellcheck=0&offset=" + strconv.Itoa(i)
		}

		urll := Info.URL + query + pageParam
		anonUrll := Info.URL + anonymize.String(query) + pageParam

		err := _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
		if err != nil {
			retErrors = append(retErrors, err)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retErrors[:len(retErrors):len(retErrors)]
}

func getLocale(options engines.Options) string {
	region := strings.SplitN(strings.ToLower(options.Locale), "_", 2)[1]
	return "country=" + region
}

func getSafeSearch(options engines.Options) string {
	if options.SafeSearch {
		return "safesearch=strict"
	}
	return "safesearch=off"
}
