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

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) error {
	if err := _sedefaults.Prepare(Info.Name, &options, &settings, &Support, &Info, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	_sedefaults.InitializeCollectors(&col, &pagesCol, &settings, &options, &timings)

	_sedefaults.PagesColRequest(Info.Name, pagesCol, ctx)
	_sedefaults.PagesColError(Info.Name, pagesCol)
	_sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	_sedefaults.ColRequest(Info.Name, col, ctx)
	_sedefaults.ColError(Info.Name, col)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		linkText, titleText, descText := _sedefaults.FieldsFromDOM(e.DOM, &dompaths, Info.Name)

		page := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

		res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
		bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
		pageRankCounter[page]++
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

	safeSearch := getSafeSeach(&options)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	urll := Info.URL + query + safeSearch
	anonUrll := Info.URL + anonymize.String(query) + safeSearch
	_sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, &retError)

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		urll := Info.URL + query + "&page=" + strconv.Itoa(i+1) + safeSearch
		anonUrll := Info.URL + anonymize.String(query) + "&page=" + strconv.Itoa(i+1) + safeSearch
		_sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, &retError)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}

func getSafeSeach(options *engines.Options) string {
	if options.SafeSearch {
		return "" // for startpage, Safe Search is the default
	}
	return "&qadf=none"
}
