package yahoo

import (
	"context"
	"net/url"
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

	safeSearchCookieParam := getSafeSearch(options)

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
			bucket.AddSEResult(res, Info.Name, relay, options, pagesCol)
			pageRankCounter[page]++
		}
	})

	errChannel := make(chan error, options.MaxPages)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	urll := Info.URL + query
	anonUrll := Info.URL + anonymize.String(query)
	_sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, errChannel)

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		urll := Info.URL + query + "&b=" + strconv.Itoa((i+1)*10)
		anonUrll := Info.URL + anonymize.String(query) + "&b=" + strconv.Itoa((i+1)*10)
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

func getSafeSearch(options engines.Options) string {
	if options.SafeSearch {
		return "&vm=r"
	}
	return "&vm=p"
}
