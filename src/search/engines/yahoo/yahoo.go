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

	safeSearchCookieParam := getSafeSearch(options)

	col.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", "sB=v=1&pn=10&rw=new&userset=0"+safeSearchCookieParam)
	})

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		titleEl := dom.Find(dompaths.Title)
		titleAria, labelExists := titleEl.Attr("aria-label")
		titleText := strings.TrimSpace(titleAria)

		linkHref, hrefExists := titleEl.Attr("href")
		linkText := removeTelemetry(linkHref)

		descText := dom.Find(dompaths.Description).Text()

		linkText, titleText, descText = _sedefaults.SanitizeFields(linkText, titleText, descText)

		if !hrefExists {
			log.Error().
				Str("engine", Info.Name.String()).
				Str("url", linkText).
				Str("title", titleText).
				Str("description", descText).
				Str("link selector", dompaths.Link).
				Msg("yahoo.Search(): href attribute doesn't exist on matched URL element")

			return
		}

		if !labelExists {
			log.Error().
				Str("engine", Info.Name.String()).
				Str("url", linkText).
				Str("title", titleText).
				Str("description", descText).
				Str("title selector", dompaths.Title).
				Msg("yahoo.Search(): aria attribute doesn't exist on matched title element")

			return
		}

		pageIndex := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)
		page := pageIndex + options.Pages.Start + 1

		res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[pageIndex]+1)
		valid := bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol, enabledEngines)
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
		pageParam := ""
		// i == 0 is the first page
		if i > 0 {
			pageParam = "&b=" + strconv.Itoa((i+1)*10)
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
