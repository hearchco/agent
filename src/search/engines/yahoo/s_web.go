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
	"github.com/hearchco/agent/src/utils/moreurls"
)

func (se Engine) WebSearch(query string, opts options.Options, resChan chan result.ResultScraped) ([]error, bool) {
	foundResults := atomic.Bool{}
	retErrors := make([]error, 0, opts.Pages.Max)
	pageRankCounter := scraper.NewPageRankCounter(opts.Pages.Max)

	se.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", safeSearchCookieString(opts.SafeSearch))
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

		// Build the parameters.
		params := moreurls.NewParams(
			paramQueryK, query,
		)
		if pageNum0 > 0 {
			params = moreurls.NewParams(
				paramQueryK, query,
				paramPageK, strconv.Itoa((pageNum0-1)*7+8),
			)
		}

		// Build the url.
		urll := moreurls.Build(searchURL, params)

		// Build anonymous url, by anonymizing the query.
		params.Set(paramQueryK, anonymize.String(query))
		anonUrll := moreurls.Build(searchURL, params)

		// Send the request.
		if err := se.Get(ctx, urll, anonUrll); err != nil {
			retErrors = append(retErrors, err)
		}
	}

	se.Wait()
	close(resChan)
	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
}
