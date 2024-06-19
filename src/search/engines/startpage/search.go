package startpage

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

func New() scraper.Enginer {
	return &Engine{scraper.EngineBase{
		Name:    seName,
		Origins: origins[:],
	}}
}

func (se Engine) Search(query string, opts options.Options, resChan chan result.ResultScraped) ([]error, bool) {
	foundResults := atomic.Bool{}
	retErrors := make([]error, 0, opts.Pages.Max)
	pageRankCounter := scraper.NewPageRankCounter(opts.Pages.Max)

	se.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		urlText, titleText, descText := parse.FieldsFromDOM(e.DOM, dompaths, se.Name)

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

	se.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "to prevent possible abuse of our service") {
			log.Error().
				Str("engine", se.Name.String()).
				Msg("Request blocked due to scraping")
		} else if strings.Contains(string(r.Body), "This page cannot function without javascript") {
			log.Error().
				Str("engine", se.Name.String()).
				Msg("Couldn't load requests, needs javascript")
		}
	})

	// Static params.
	paramSafeSearch := safeSearchParamString(opts.SafeSearch)

	for i := range opts.Pages.Max {
		pageNum0 := i + opts.Pages.Start
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(i))

		// Dynamic params.
		paramPage := ""
		if pageNum0 > 0 {
			paramPage = fmt.Sprintf("%v=%v", paramKeyPage, pageNum0+1)
		}

		combinedParams := morestrings.JoinNonEmpty([]string{paramPage, paramSafeSearch}, "&", "&")

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
