package google

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/hearchco/agent/src/utils/moreurls"
)

func (se Engine) ImageSearch(query string, opts options.Options, resChan chan result.ResultScraped) ([]error, bool) {
	foundResults := atomic.Bool{}
	retErrors := make([]error, 0, opts.Pages.Max)
	pageRankCounter := scraper.NewPageRankCounter(opts.Pages.Max)

	se.OnResponse(func(e *colly.Response) {
		body := string(e.Body)
		index := strings.Index(body, `{"ischj":`)

		if index == -1 {
			log.Error().
				Caller().
				Str("engine", se.Name.String()).
				Str("body", body).
				Msg("Failed parsing response, couldn't find the start of JSON")
			return
		}

		body = body[index:]
		var jsonResponse imgJsonResponse
		if err := json.Unmarshal([]byte(body), &jsonResponse); err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("engine", se.Name.String()).
				Str("body", body).
				Msg("Failed parsing response, couldn't unmarshal JSON")
			return
		}

		pageIndex := se.PageFromContext(e.Request.Ctx)
		page := pageIndex + opts.Pages.Start + 1

		for _, metadata := range jsonResponse.ISCHJ.Metadata {
			origImg := metadata.OriginalImage
			thmbImg := metadata.Thumbnail
			resultJson := metadata.Result
			textInGridJson := metadata.TextInGrid

			// Google Images sometimes inverts original height and width.
			if (thmbImg.Height > thmbImg.Width) != (origImg.Height > origImg.Width) {
				origImg.Height, origImg.Width = origImg.Width, origImg.Height
			}

			if resultJson.ReferrerUrl != "" && origImg.Url != "" && thmbImg.Url != "" {
				r, err := result.ConstructImagesResult(
					se.Name, origImg.Url, resultJson.PageTitle, textInGridJson.Snippet, page, pageRankCounter.GetPlusOne(pageIndex),
					origImg.Height, origImg.Width, thmbImg.Height, thmbImg.Width, thmbImg.Url, resultJson.SiteTitle, resultJson.ReferrerUrl,
				)
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
				}
			} else {
				log.Error().
					Caller().
					Str("engine", se.Name.String()).
					Str("jsonMetadata", fmt.Sprintf("%v", metadata)).
					Str("url", resultJson.ReferrerUrl).
					Str("original", origImg.Url).
					Str("thumbnail", thmbImg.Url).
					Msg("Couldn't find image URL")
			}
		}
	})

	// Constant params.
	paramLocaleV, paramLocaleSecV := localeParamValues(opts.Locale)
	paramSafeSearchV := safeSearchParamValue(opts.SafeSearch)

	for i := range opts.Pages.Max {
		pageNum0 := i + opts.Pages.Start
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(i))

		// Build the parameters.
		params := moreurls.NewParams(
			paramQueryK, query,
			imgParamTbmK, imgParamTbmV,
			imgParamAsearchK, imgParamAsearchV,
			paramFilterK, paramFilterV,
			imgParamPageK, imgParamPageVPrefix+"1",
			paramLocaleK, paramLocaleV,
			paramLocaleSecK, paramLocaleSecV,
			paramSafeSearchK, paramSafeSearchV,
		)
		if pageNum0 > 0 {
			params.Set(imgParamPageK, fmt.Sprintf(imgParamPageVPrefix+"%v", pageNum0*10))
		}

		// Build the url.
		urll := moreurls.Build(imageSearchURL, params)

		// Build anonymous url, by anonymizing the query.
		params.Set(paramQueryK, anonymize.String(query))
		anonUrll := moreurls.Build(imageSearchURL, params)

		// Send the request.
		if err := se.Get(ctx, urll, anonUrll); err != nil {
			retErrors = append(retErrors, err)
		}
	}

	se.Wait()
	close(resChan)
	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
}
