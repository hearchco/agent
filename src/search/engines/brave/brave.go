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
	"github.com/hearchco/hearchco/src/search/parse"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, &options, &settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	localeCookie := getLocale(options)
	safeSearchCookie := getSafeSearch(options)

	col.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", localeCookie)
		r.Headers.Add("Cookie", safeSearchCookie)
	})

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, hrefExists := dom.Find(dompaths.Link).Attr("href")
		linkText := parse.ParseURL(linkHref)
		titleText := strings.TrimSpace(dom.Find(dompaths.Title).Text())
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())

		if hrefExists && linkText != "" && linkText != "#" && titleText != "" {
			if descText == "" {
				descText = strings.TrimSpace(dom.Find("div.product > div.flex-hcenter > div > div[class=\"text-sm text-gray\"]").Text())
			}
			if descText == "" {
				descText = strings.TrimSpace(dom.Find("p.snippet-description").Text())
			}

			page := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
			bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol)
			pageRankCounter[page]++
		}
	})

	errChannel := make(chan error, options.MaxPages)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	urll := Info.URL + query + "&source=web"
	anonUrll := Info.URL + anonymize.String(query) + "&source=web"
	_sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, errChannel)

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		urll := Info.URL + query + "&spellcheck=0&offset=" + strconv.Itoa(i)
		anonUrll := Info.URL + anonymize.String(query) + "&spellcheck=0&offset=" + strconv.Itoa(i)
		_sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, errChannel)
	}

	retErrors := make([]error, 0)
	for i := 0; i < options.MaxPages; i++ {
		err := <-errChannel
		if err != nil {
			retErrors = append(retErrors, err)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retErrors
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
