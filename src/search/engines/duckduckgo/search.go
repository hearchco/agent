package duckduckgo

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/PuerkitoBio/goquery"
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
		r.Headers.Add("Cookie", localeCookieString(opts.Locale))
	})

	se.OnHTML(dompaths.ResultsContainer, func(e *colly.HTMLElement) {
		log.Trace().
			Caller().
			Msg("Matched results container")

		var urlText, linkScheme, titleText, descText string
		var hrefExists bool

		pageIndex := se.PageFromContext(e.Request.Ctx)
		page := pageIndex + opts.Pages.Start + 1

		e.DOM.Children().Each(func(i int, row *goquery.Selection) {
			switch i % 4 {
			case 0:
				var urlHref string
				urlHref, hrefExists = row.Find(dompaths.URL).Attr("href")
				if strings.Contains(urlHref, "https") {
					linkScheme = "https://"
				} else {
					linkScheme = "http://"
				}
				titleText = parse.SanitizeTitle(row.Find(dompaths.Title).Text())
			case 1:
				descText = parse.SanitizeDescription(row.Find(dompaths.Description).Text())
			case 2:
				rawURL := linkScheme + row.Find("td > span.link-text").Text()
				urlText = parse.SanitizeURL(rawURL)
			case 3:
				if !hrefExists {
					log.Error().
						Caller().
						Str("engine", se.Name.String()).
						Str("url", urlText).
						Str("title", titleText).
						Str("description", descText).
						Str("link selector", dompaths.URL).
						Msg("Href attribute doesn't exist on matched URL element")
					return
				}

				onPageRank := (i/4 + 1)

				r, err := result.ConstructResult(se.Name, urlText, titleText, descText, page, onPageRank)
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
						Int("rank", onPageRank).
						Str("result", fmt.Sprintf("%v", r)).
						Msg("Sending result to channel")
					resChan <- r
				}
			}
		})
	})

	for i := range opts.Pages.Max {
		pageNum0 := i + opts.Pages.Start
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(i))

		if pageNum0 == 0 {
			// Build the parameters.
			params := parameters.NewParams(
				paramQueryK, query,
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
		} else {
			// Build the parameters.
			params := parameters.NewParams(
				paramQueryK, query,
				// This value changes depending on how many results were returned on the first page, so it's set to the lowest seen value.
				paramPageK, strconv.Itoa(pageNum0*20),
			)

			// Build the POST body query params Reader.
			body := strings.NewReader(params.QueryEscape())

			// Build anonymous POST body query params.
			params.Set(paramQueryK, anonymize.String(query))
			anonBody := params.QueryEscape()

			// Send the POST request.
			if err := se.Post(ctx, searchURL, body, anonBody); err != nil {
				retErrors = append(retErrors, err)
			}
		}
	}

	se.Wait()
	close(resChan)
	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
}
