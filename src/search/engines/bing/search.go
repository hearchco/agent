package bing

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

func (se Engine) Search(query string, opts options.Options, resChan chan result.ResultScraped) ([]error, bool) {
	foundResults := atomic.Bool{}
	retErrors := make([]error, 0, opts.Pages.Max)
	pageRankCounter := scraper.NewPageRankCounter(opts.Pages.Max)

	se.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		log.Trace().
			Caller().
			Msg("Matched result")

		// The telemetry link is a valid link so it can be sanitized.
		urlText, titleText, descText := parse.FieldsFromDOM(e.DOM, dompaths, se.Name)

		urlWOTelemetry, err := removeTelemetry(urlText)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("url", urlText).
				Msg("Failed to remove telemetry")
			return
		}
		urlText = parse.SanitizeURL(urlWOTelemetry)

		if descText == "" {
			descText = e.DOM.Find("p.b_algoSlug").Text()
		}
		descText = strings.TrimPrefix(descText, "<span class=\"algoSlug_icon\" data-priority=\"2\">WEB</span>")
		descText = parse.SanitizeDescription(descText)

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

	// Static params.
	paramLocale := localeParamString(opts.Locale)

	for i := range opts.Pages.Max {
		pageNum0 := i + opts.Pages.Start
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(i))

		// Dynamic params.
		paramPage := ""
		if pageNum0 > 0 {
			paramPage = fmt.Sprintf("%v=%v", paramKeyPage, pageNum0*10+1)
		}

		combinedParams := morestrings.JoinNonEmpty("&", "&", paramPage, paramLocale)

		urll := fmt.Sprintf("%v?q=%v%v", searchURL, query, combinedParams)
		anonUrll := fmt.Sprintf("%v?q=%v%v", searchURL, anonymize.String(query), combinedParams)

		if err := se.Get(ctx, urll, anonUrll); err != nil {
			retErrors = append(retErrors, err)
		}
	}

	se.Wait()
	close(resChan)
	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
}
