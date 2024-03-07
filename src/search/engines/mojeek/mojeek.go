package mojeek

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

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		titleEl := dom.Find(dompaths.Title)
		linkHref, hrefExists := titleEl.Attr("href")
		linkText := parse.ParseURL(linkHref)
		titleText := strings.TrimSpace(titleEl.Text())
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())

		if hrefExists && linkText != "" && linkText != "#" && titleText != "" {
			page := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
			bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol)
			pageRankCounter[page]++
		}
	})

	localeParam := getLocale(options)
	safeSearchParam := getSafeSearch(options)

	retErrors := make([]error, options.MaxPages)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	urll := Info.URL + query + localeParam + safeSearchParam
	anonUrll := Info.URL + anonymize.String(query) + localeParam + safeSearchParam
	err = _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
	retErrors[0] = err

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		urll := Info.URL + query + "&s=" + strconv.Itoa(i*10+1) + localeParam + safeSearchParam
		anonUrll := Info.URL + anonymize.String(query) + "&s=" + strconv.Itoa(i*10+1) + localeParam + safeSearchParam
		err = _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
		retErrors[i] = err
	}

	col.Wait()
	pagesCol.Wait()

	return _sedefaults.NonNilErrorsFromSlice(retErrors)
}

func getLocale(options engines.Options) string {
	spl := strings.SplitN(strings.ToLower(options.Locale), "_", 2)
	return "&lb=" + spl[0] + "&arc=" + spl[1]
}

func getSafeSearch(options engines.Options) string {
	if options.SafeSearch {
		return "&safe=1"
	}
	return "&safe=0"
}
