package yahoo

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/engines/_sedefaults"
	"github.com/hearchco/hearchco/src/search/parse"
	"github.com/rs/zerolog/log"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) error {
	if err := _sedefaults.Prepare(Info.Name, &options, &settings, &Support, &Info, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	_sedefaults.InitializeCollectors(&col, &pagesCol, &options, &timings)

	_sedefaults.PagesColRequest(Info.Name, pagesCol, ctx)
	_sedefaults.PagesColError(Info.Name, pagesCol)
	_sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	_sedefaults.ColRequest(Info.Name, col, ctx)
	_sedefaults.ColError(Info.Name, col)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	safeSearchCookieParam := getSafeSearch(&options)

	col.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", "sB=v=1&pn=10&rw=new&userset=0"+safeSearchCookieParam)
	})

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		titleEl := dom.Find(dompaths.Title)
		linkHref, hrefExists := titleEl.Attr("href")
		linkText := parse.ParseURL(removeTelemetry(linkHref))
		titleAria, labelExists := titleEl.Attr("aria-label")
		titleText := strings.TrimSpace(titleAria)
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())

		if labelExists && hrefExists && linkText != "" && linkText != "#" && titleText != "" {
			page := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
			bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
			pageRankCounter[page]++
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	_sedefaults.DoGetRequest(Info.URL+query, colCtx, col, Info.Name, &retError)

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		_sedefaults.DoGetRequest(Info.URL+query+"&b="+strconv.Itoa((i+1)*10), colCtx, col, Info.Name, &retError)
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
	link = "http" + strings.SplitN(suff, "/RK=", 2)[0]
	newLink, err := url.QueryUnescape(link)
	if err != nil {
		log.Error().
			Err(err).
			Str("url", link).
			Msg("yahoo.removeTelemetry(): couldn't parse url, url.QueryUnescape() failed")
		return ""
	}
	return newLink
}

func getSafeSearch(options *engines.Options) string {
	if options.SafeSearch {
		return "&vm=r"
	}
	return "&vm=p"
}
