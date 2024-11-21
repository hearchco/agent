package brave

import (
	"fmt"
	"strconv"
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
		r.Headers.Set("Accept-Encoding", "gzip, deflate") // Brotli lib used has some issues with Brave's brotli compression.
		r.Headers.Add("Cookie", localeCookieString(opts.Locale))
		r.Headers.Add("Cookie", safeSearchCookieString(opts.SafeSearch))
	})

	se.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		urlText, titleText, descText := parse.FieldsFromDOM(e.DOM, dompaths, se.Name)

		if descText == "" {
			descText = e.DOM.Find("div.product > div.flex-hcenter > div > div[class=\"text-sm text-gray\"]").Text()
		}
		if descText == "" {
			descText = e.DOM.Find("p.snippet-description").Text()
		}
		descText = parse.SanitizeDescription(descText)

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
			paramSourceK, paramSourceV,
		)
		if pageNum0 > 0 {
			params = moreurls.NewParams(
				paramQueryK, query,
				paramPageK, strconv.Itoa(pageNum0),
				paramSpellcheckK, paramSpellcheckV,
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
