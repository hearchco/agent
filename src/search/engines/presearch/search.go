package presearch

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
	"github.com/hearchco/agent/src/search/scraper/parse"
	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/hearchco/agent/src/utils/moreurls"
	"github.com/hearchco/agent/src/utils/moreurls/parameters"
)

func (se Engine) Search(query string, opts options.Options, resChan chan result.ResultScraped) ([]error, bool) {
	foundResults := atomic.Bool{}
	retErrors := make([]error, 0, opts.Pages.Max)

	se.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", "presearch_session=;")
		r.Headers.Add("Cookie", "use_local_search_results=false")
		r.Headers.Add("Cookie", "ai_results_disable=1")
		r.Headers.Add("Cookie", safeSearchCookieString(opts.SafeSearch))
	})

	se.OnResponse(func(r *colly.Response) {
		pageIndex := se.PageFromContext(r.Request.Ctx)
		page := pageIndex + opts.Pages.Start + 1

		apiStr := r.Request.Ctx.Get("isAPI")
		isApi, err := strconv.ParseBool(apiStr)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("engine", se.Name.String()).
				Msg("Failed to parse isAPI")
		}

		if isApi {
			// JSON response.
			var pr jsonResponse
			err := json.Unmarshal(r.Body, &pr)
			if err != nil {
				log.Error().
					Caller().
					Err(err).
					Str("engine", se.Name.String()).
					Bytes("body", r.Body).
					Msg("Failed to parse response, couldn't unmarshal JSON")
			}

			counter := 1
			for _, jsonR := range pr.Results.StandardResults {
				goodURL, goodTitle, goodDesc := parse.SanitizeFields(jsonR.Link, jsonR.Title, jsonR.Desc)

				r, err := result.ConstructResult(se.Name, goodURL, goodTitle, goodDesc, page, counter)
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
						Int("rank", counter).
						Str("result", fmt.Sprintf("%v", r)).
						Msg("Sending result to channel")
					resChan <- r
					counter++
				}
			}
		} else {
			// HTML response, forward call to API.
			suff := strings.SplitN(string(r.Body), "window.searchId = \"", 2)[1]
			searchId := strings.SplitN(suff, "\"", 2)[0]

			nextCtx := colly.NewContext()
			nextCtx.Put("page", strconv.Itoa(page))
			nextCtx.Put("isAPI", "true")

			urll := fmt.Sprintf("https://presearch.com/results?id=%v", searchId)
			anonUrll := fmt.Sprintf("https://presearch.com/results?id=%v", anonymize.CalculateHashBase64(searchId))

			if err := se.Get(nextCtx, urll, anonUrll); err != nil {
				retErrors = append(retErrors, err)
			}
		}
	})

	for i := range opts.Pages.Max {
		pageNum0 := i + opts.Pages.Start
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(i))

		// Build the parameters.
		params := parameters.NewParams(
			paramQueryK, query,
		)
		if pageNum0 > 0 {
			params = parameters.NewParams(
				paramQueryK, query,
				paramPageK, strconv.Itoa(pageNum0+1),
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
