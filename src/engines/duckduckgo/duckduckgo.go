package duckduckgo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/search/limit"
	"github.com/tminaorg/brzaguza/src/search/useragent"
	"github.com/tminaorg/brzaguza/src/structures"
)

/*
	This module uses POST requests. GET requests could be used like
	https://lite.duckduckgo.com/lite/?q=<query>&dc=<rank of first result on page>
*/

const seName string = "DuckDuckGo"
const seURL string = "https://lite.duckduckgo.com/lite/"

func Search(ctx context.Context, query string, relay *structures.Relay, options *structures.Options) error {
	if ctx == nil {
		ctx = context.Background()
	} //^ not necessary as ctx is always passed in search.go, branch predictor will skip this if

	if err := limit.RateLimit.Wait(ctx); err != nil {
		return err
	}

	if options.UserAgent == "" {
		options.UserAgent = useragent.RandomUserAgent()
	}
	log.Trace().Msgf("%v: UserAgent: %v", seName, options.UserAgent)

	var col *colly.Collector
	if options.MaxPages == 1 {
		col = colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent)) // so there is no thread creation overhead
	} else {
		col = colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent), colly.Async(true))
	}
	pagesCol := colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent), colly.Async(true))

	var retError error

	pagesCol.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil { // dont fully understand this
			log.Error().Msgf("%v: Pages Collector; Error OnRequest %v", seName, r)
			r.Abort()
			retError = err
			return
		}
		r.Ctx.Put("originalURL", r.URL.String())
	})

	pagesCol.OnError(func(r *colly.Response, err error) {
		log.Error().Msgf("%v: Pages Collector; Error OnError:\nURL: %v\nError: %v", seName, r.Ctx.Get("originalURL"), err)
		log.Trace().Msgf("%v: HTML Response:\n%v", seName, string(r.Body))
		retError = err
	})

	pagesCol.OnResponse(func(r *colly.Response) {
		urll := r.Ctx.Get("originalURL")

		bucket.SetResultResponse(urll, r, relay, seName)
	})

	col.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil { // dont fully understand this
			log.Error().Msgf("%v: SE Collector; Error OnRequest %v", seName, r)
			r.Abort()
			retError = err
			return
		}
		var reqBody []byte
		r.Body.Read(reqBody)
		r.Ctx.Put("body", string(reqBody))
	})

	col.OnError(func(r *colly.Response, err error) {
		log.Error().Msgf("%v: SE Collector; Error OnError:\nURL: %v\nError: %v", seName, r.Request.URL.String(), err)
		log.Trace().Msgf("%v: HTML Response:\n%v", seName, string(r.Body))
		retError = err
	})

	col.OnHTML("div.filters > table > tbody", func(e *colly.HTMLElement) {
		var linkText string
		var titleText string
		var descText string
		var rrank int

		var reqBody string = e.Request.Ctx.Get("body")
		var page int
		fmt.Sscanf(reqBody, "q="+query+"&dc=%d", &page)
		page = page/20 + 1

		e.DOM.Children().Each(func(i int, row *goquery.Selection) {
			switch i % 4 {
			case 0:
				rankText := strings.TrimSpace(row.Children().First().Text())
				fmt.Sscanf(rankText, "%d", &rrank)

				linkElement := row.Find("td > a.result-link")
				linkHref, _ := linkElement.Attr("href")
				linkText = strings.TrimSpace(linkHref)
				titleText = strings.TrimSpace(linkElement.Text())
			case 1:
				descText = strings.TrimSpace(row.Find("td.result-snippet").Text())
			case 3:
				if linkText != "" && linkText != "#" && titleText != "" {
					res := structures.Result{
						URL:          linkText,
						Rank:         -1,
						SERank:       rrank,
						SEPage:       page,
						SEOnPageRank: (i/4 + 1),
						Title:        titleText,
						Description:  descText,
						SearchEngine: seName,
					}

					bucket.SetResult(&res, relay, options, pagesCol)
				}
			}
		})
	})

	col.PostRaw(seURL, []byte("q="+query+"&dc=1"))
	for i := 1; i < options.MaxPages; i++ {
		col.PostRaw(seURL, []byte("q="+query+"&dc="+strconv.Itoa(i*20)))
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
