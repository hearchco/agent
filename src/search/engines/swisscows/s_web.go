package swisscows

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper/parse"
	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/hearchco/agent/src/utils/moreurls"
)

func (se Engine) WebSearch(query string, opts options.Options, resChan chan result.ResultScraped) ([]error, bool) {
	foundResults := atomic.Bool{}
	retErrors := make([]error, 0, opts.Pages.Max)

	se.OnRequest(func(r *colly.Request) {
		if r.Method == "OPTIONS" {
			return
		}

		var qry string = "?" + r.URL.RawQuery
		nonce, sig, err := generateAuth(qry)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Msg("Failed building request, couldn't generate auth")
			return
		}

		r.Headers.Set("X-Request-Nonce", nonce)
		r.Headers.Set("X-Request-Signature", sig)
		r.Headers.Set("Pragma", "no-cache")
	})

	se.OnResponse(func(r *colly.Response) {
		query := r.Request.URL.Query().Get("query")
		urll := r.Request.URL.String()
		anonUrll := anonymize.Substring(urll, query)
		log.Trace().
			Str("engine", se.Name.String()).
			Str("url", anonUrll).
			Str("nonce", r.Request.Headers.Get("X-Request-Nonce")).
			Str("signature", r.Request.Headers.Get("X-Request-Signature")).
			Msg("Got response")

		pageIndex := se.PageFromContext(r.Request.Ctx)
		page := pageIndex + opts.Pages.Start + 1

		var parsedResponse jsonResponse
		if err := json.Unmarshal(r.Body, &parsedResponse); err != nil {
			log.Error().
				Caller().
				Err(err).
				Bytes("body", r.Body).
				Msg("Failed to parse response, couldn't unmarshal JSON")
			return
		}

		counter := 1
		for _, jsonResult := range parsedResponse.Items {
			goodURL, goodTitle, goodDesc := parse.SanitizeFields(jsonResult.URL, jsonResult.Title, jsonResult.Desc)

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
	})

	// Constant params.
	paramLocaleV := localeParamValue(opts.Locale)

	for i := range opts.Pages.Max {
		pageNum0 := i + opts.Pages.Start
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(i))

		// Build the parameters.
		params := moreurls.NewParams(
			paramFreshnessK, paramFreshnessV,
			paramItemsK, paramItemsV,
			paramPageK, strconv.Itoa(pageNum0*10),
			paramQueryK, query,
			paramLocaleK, paramLocaleV,
		)

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
