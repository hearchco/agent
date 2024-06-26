package etools

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
		// Ignore the first request if it's not the first page (see below).
		ignoreS := e.Request.Ctx.Get("ignore")
		if ignoreS == strconv.FormatBool(true) {
			return
		}

		urlText, titleText, descText := parse.FieldsFromDOM(e.DOM, dompaths, se.Name)

		// Need to perform this check here so the check below doesn't panic.
		if urlText == "" {
			log.Error().
				Caller().
				Str("title", titleText).
				Str("description", descText).
				Msg("Invalid result, url is empty")
			return
		}

		// Telemetry link, e.g. //web.search.ch/r/redirect?event=website&origin=result!u377d618861533351/https://de.wikipedia.org/wiki/Charles_Paul_Wilp.
		if urlText[0] != 'h' {
			urlText = "http" + strings.Split(urlText, "http")[1] // Works for https as well.
		}

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

	se.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "Sorry for the CAPTCHA") {
			log.Error().
				Caller().
				Msg("Captcha detected")
		}
	})

	firstRequest := true

	// Static params.
	paramSafeSearch := safeSearchParamString(opts.SafeSearch)

	for i := range opts.Pages.Max {
		pageNum0 := i + opts.Pages.Start
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(i))

		var err error
		// eTools requires a request for the first page.
		if pageNum0 == 0 || firstRequest {
			combinedParams := morestrings.JoinNonEmpty("&", "&", paramCountry, paramLanguage, paramSafeSearch)

			body := strings.NewReader(fmt.Sprintf("query=%v%v", query, combinedParams))
			anonBody := fmt.Sprintf("query=%v%v", anonymize.String(query), combinedParams)

			if firstRequest {
				firstCtx := colly.NewContext()
				firstCtx.Put("ignore", strconv.FormatBool(true))
				err = se.Post(firstCtx, searchURL, body, anonBody)
			} else {
				err = se.Post(ctx, searchURL, body, anonBody)
			}

			firstRequest = false
			se.Wait() // Needed to save the JSESSION cookie.
		}

		// Since the above can happen for the first request and then we need to request the wanted page.
		if pageNum0 > 0 {
			// Query isn't needed as it's saved in the JSESSION cookie.
			paramPage := fmt.Sprintf("%v=%v", paramKeyPage, pageNum0+1)
			urll := fmt.Sprintf("%v?%v", pageURL, paramPage)
			err = se.Get(ctx, urll, urll)
		}

		if err != nil {
			retErrors = append(retErrors, err)
		}
	}

	se.Wait()
	close(resChan)
	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
}
