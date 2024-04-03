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
	"github.com/rs/zerolog/log"
)

type Engine struct{}

func New() Engine {
	return Engine{}
}

func (e Engine) Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings, salt string, enabledEngines int) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, &options, &settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

	pageRankCounter := make([]int, options.Pages.Max)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		linkText, titleText, descText := _sedefaults.FieldsFromDOM(e.DOM, dompaths, Info.Name)

		pageIndex := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)
		page := pageIndex + options.Pages.Start + 1

		res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[pageIndex]+1)
		valid := bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol, enabledEngines)
		if valid {
			pageRankCounter[pageIndex]++
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

	retErrors := make([]error, 0, options.Pages.Max)

	// static params
	safeSearch := getSafeSearch(options)

	// starts from at least 0
	for i := options.Pages.Start; i < options.Pages.Start+options.Pages.Max; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i-options.Pages.Start))

		// dynamic params
		pageParam := ""
		// i == 0 is the first page
		if i > 0 {
			pageParam = "&page=" + strconv.Itoa(i+1)
		}

		urll := Info.URL + query + pageParam + safeSearch
		anonUrll := Info.URL + anonymize.String(query) + pageParam + safeSearch

		err := _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
		if err != nil {
			retErrors = append(retErrors, err)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retErrors[:len(retErrors):len(retErrors)]
}

func getSafeSearch(options engines.Options) string {
	if options.SafeSearch {
		return "" // for startpage, Safe Search is the default
	}
	return "&qadf=none"
}
