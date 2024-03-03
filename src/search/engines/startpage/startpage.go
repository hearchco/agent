package startpage

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
	"github.com/rs/zerolog/log"
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

		linkHref, hrefExists := dom.Find(dompaths.Link).Attr("href")
		linkText := parse.ParseURL(linkHref)
		titleText := strings.TrimSpace(dom.Find(dompaths.Title).Text())
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())

		if hrefExists && linkText != "" && linkText != "#" && titleText != "" {
			page := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
			bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol)
			pageRankCounter[page]++
		} else {
			log.Trace().
				Str("engine", Info.Name.String()).
				Str("url", linkText).
				Str("title", titleText).
				Str("description", descText).
				Msg("Matched result, but couldn't retrieve data")
		}
	})

	col.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "to prevent possible abuse of our service") {
			log.Error().
				Str("engine", Info.Name.String()).
				Msg("Request blocked by engine due to scraping")
		} else if strings.Contains(string(r.Body), "This page cannot function without javascript") {
			log.Error().
				Str("engine", Info.Name.String()).
				Msg("Engine couldn't load requests, needs javascript")
		}
	})

	safeSearch := getSafeSearch(options)

	retErrors := make([]error, options.MaxPages)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	urll := Info.URL + query + safeSearch
	anonUrll := Info.URL + anonymize.String(query) + safeSearch
	err = _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
	retErrors[0] = err

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		urll := Info.URL + query + "&page=" + strconv.Itoa(i+1) + safeSearch
		anonUrll := Info.URL + anonymize.String(query) + "&page=" + strconv.Itoa(i+1) + safeSearch
		err = _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
		retErrors[i] = err
	}

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

func getSafeSearch(options engines.Options) string {
	if options.SafeSearch {
		return "" // for startpage, Safe Search is the default
	}
	return "&qadf=none"
}
