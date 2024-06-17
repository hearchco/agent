package yahoo

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
	"github.com/hearchco/agent/src/search/scraper/parse"
	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/hearchco/agent/src/utils/morestrings"
)

type Engine struct {
	scraper.EngineBase
}

func New() *Engine {
	return &Engine{EngineBase: scraper.EngineBase{
		Name:    info.Name,
		Origins: info.Origins,
	}}
}

func (se Engine) Search(query string, opts options.Options, resChan chan result.ResultScraped) ([]error, bool) {
	foundResults := atomic.Bool{}
	retErrors := make([]error, 0, opts.Pages.Max)
	pageRankCounter := scraper.NewPageRankCounter(opts.Pages.Max)

	se.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", fmt.Sprintf("%v&%v", safeSearchCookiePrefix, safeSearchCookieString(opts.SafeSearch)))
	})

	se.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		titleEl := dom.Find(dompaths.Title)
		titleAria, labelExists := titleEl.Attr("aria-label")
		if !labelExists {
			log.Error().
				Caller().
				Str("engine", se.Name.String()).
				Str("title selector", dompaths.Title).
				Msg("Aria attribute doesn't exist on matched title element")
			return
		}
		titleText := strings.TrimSpace(titleAria)

		urlHref, hrefExists := titleEl.Attr("href")
		if !hrefExists {
			log.Error().
				Caller().
				Str("engine", se.Name.String()).
				Str("link selector", dompaths.URL).
				Msg("Href attribute doesn't exist on matched URL element")
			return
		}

		urlText, err := removeTelemetry(urlHref)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("engine", se.Name.String()).
				Str("url", urlText).
				Msg("Failed to remove telemetry")
			return
		}

		descText := dom.Find(dompaths.Description).Text()

		urlText, titleText, descText = parse.SanitizeFields(urlText, titleText, descText)

		pageIndex := se.PageFromContext(e.Request.Ctx)
		page := pageIndex + opts.Pages.Start + 1

		r, err := result.ConstructResult(se.Name, urlText, titleText, descText, page, pageRankCounter.GetPlusOne(pageIndex))
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("result", fmt.Sprintf("%v", r)).
				Msg("Failed to construct result")
		} else {
			log.Trace().
				Caller().
				Int("page", page).
				Int("rank", pageRankCounter.GetPlusOne(pageIndex)).
				Str("result", fmt.Sprintf("%v", r)).
				Msg("Sending result to channel")
			resChan <- r
			pageRankCounter.Increment(pageIndex)
			if !foundResults.Load() {
				foundResults.Store(true)
			}
		}
	})

	for i := range opts.Pages.Max {
		pageNum0 := i + opts.Pages.Start
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(i))

		// Dynamic params.
		pageParam := ""
		if pageNum0 > 0 {
			pageParam = fmt.Sprintf("%v=%v", params.Page, (pageNum0-1)*7+8)
		}

		combinedParams := morestrings.JoinNonEmpty([]string{pageParam}, "&", "&")

		urll := fmt.Sprintf("%v?p=%v%v", info.URL, query, combinedParams)
		anonUrll := fmt.Sprintf("%v?p=%v%v", info.URL, anonymize.String(query), combinedParams)

		if err := se.Get(ctx, urll, anonUrll); err != nil {
			retErrors = append(retErrors, err)
		}
	}

	se.Wait()
	close(resChan)
	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
}
