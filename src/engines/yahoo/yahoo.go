package yahoo

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
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

	sedefaults.PagesColRequest(Info.Name, pagesCol, ctx, &retError)
	sedefaults.PagesColError(Info.Name, pagesCol)
	sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	sedefaults.ColRequest(Info.Name, col, &ctx, &retError)
	sedefaults.ColError(Info.Name, col, &retError)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		titleEl := dom.Find(dompaths.Title)
		linkHref, _ := titleEl.Attr("href")
		linkText := parse.ParseURL(linkHref)
		linkText = removeTelemetry(linkText)
		titleAria, _ := titleEl.Attr("aria-label")
		titleText := strings.TrimSpace(titleAria)
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
			bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
			pageRankCounter[page]++
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	err := col.Request("GET", Info.URL+query, nil, colCtx, nil)
	if engines.IsTimeoutError(err) {
		log.Trace().Err(err).Msgf("%v: failed requesting with GET method", Info.Name)
	} else if err != nil {
		log.Error().Err(err).Msgf("%v: failed requesting with GET method", Info.Name)
	}

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		err := col.Request("GET", Info.URL+query+"&b="+strconv.Itoa((i+1)*10), nil, colCtx, nil)
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

func removeTelemetry(link string) string {
	if !strings.Contains(link, "://r.search.yahoo.com/") {
		return link
	}
	suff := strings.SplitAfterN(link, "/RU=http", 2)[1]
	newLink := "http" + strings.SplitN(suff, "/RK=", 2)[0]
	return parse.ParseURL(newLink)
}
