package startpage

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

	sedefaults.PagesColRequest(Info.Name, pagesCol, ctx, &retError)
	sedefaults.PagesColError(Info.Name, pagesCol)
	sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	sedefaults.ColRequest(Info.Name, col, &ctx, &retError)
	sedefaults.ColError(Info.Name, col, &retError)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, _ := dom.Find(dompaths.Link).Attr("href")
		linkText := parse.ParseURL(linkHref)
		titleText := strings.TrimSpace(dom.Find(dompaths.Title).Text())
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
			bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
			pageRankCounter[page]++
		} else {
			log.Trace().Msgf("%v: Matched Result, but couldn't retrieve data.\nURL:%v\nTitle:%v\nDescription:%v", Info.Name, linkText, titleText, descText)
		}
	})

	col.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "to prevent possible abuse of our service") {
			log.Error().Msgf("%v: Request blocked by engine due to scraping.", Info.Name)
		} else if strings.Contains(string(r.Body), "This page cannot function without javascript") {
			log.Error().Msgf("%v: Engine couldn't load requests, needs javascript.", Info.Name)
		}
	})

	safeSearch := getSafeSeach(&options)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	err := col.Request("GET", Info.URL+query+safeSearch, nil, colCtx, nil)
	if engines.IsTimeoutError(err) {
		log.Trace().Err(err).Msgf("%v: failed requesting with GET method", Info.Name)
	} else if err != nil {
		log.Error().Err(err).Msgf("%v: failed requesting with GET method", Info.Name)
	}

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		err := col.Request("GET", Info.URL+query+"&page="+strconv.Itoa(i+1)+safeSearch, nil, colCtx, nil)
		if engines.IsTimeoutError(err) {
			log.Trace().Err(err).Msgf("%v: failed requesting with GET method on page", Info.Name)
		} else if err != nil {
			log.Error().Err(err).Msgf("%v: failed requesting with GET method on page", Info.Name)
		}
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
