package googleimages

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

	se.OnResponse(func(e *colly.Response) {
		body := string(e.Body)
		index := strings.Index(body, "{\"ischj\":")

		if index == -1 {
			log.Error().
				Caller().
				Str("engine", se.Name.String()).
				Str("body", body).
				Msg("Failed parsing response, couldn't find the start of JSON")
			return
		}

		body = body[index:]
		var jsonResponse jsonResponse
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

	// Static params.
	localeParam := localeParamString(opts.Locale)
	safeSearchParam := safeSearchParamString(opts.SafeSearch)

	for i := range opts.Pages.Max {
		pageNum0 := i + opts.Pages.Start
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(i))

		// Dynamic params.
		pageParam := fmt.Sprintf("%v:1", params.Page)
		if pageNum0 > 0 {
			pageParam = fmt.Sprintf("%v:%v", params.Page, pageNum0*10)
		}

		combinedParams := morestrings.JoinNonEmpty([]string{tbmParam, asearchParam, filterParam, pageParam, localeParam, safeSearchParam}, "&", "&")

		urll := fmt.Sprintf("%v?q=%v%v", info.URL, query, combinedParams)
		anonUrll := fmt.Sprintf("%v?q=%v%v", info.URL, anonymize.String(query), combinedParams)

		if err := se.Get(ctx, urll, anonUrll); err != nil {
			retErrors = append(retErrors, err)
		}
	}

	se.Wait()
	close(resChan)
	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
}
