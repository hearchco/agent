package google

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/search/limit"
	"github.com/tminaorg/brzaguza/src/search/useragent"
	"github.com/tminaorg/brzaguza/src/structures"
)

const seName string = "Google"
const seURL string = "https://www.google.com/search?q="
const resPerPage int = 10

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
	})

	col.OnError(func(r *colly.Response, err error) {
		log.Error().Msgf("%v: SE Collector; Error OnError:\nURL: %v\nError: %v", seName, r.Request.URL.String(), err)
		log.Trace().Msgf("%v: HTML Response:\n%v", seName, string(r.Body))
		retError = err
	})

	var pageRankCounter []int = make([]int, options.MaxPages*resPerPage)

	col.OnHTML("div.g", func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, _ := dom.Find("a").Attr("href")
		linkText := strings.TrimSpace(linkHref)
		titleText := strings.TrimSpace(dom.Find("div > div > div > a > h3").Text())
		descText := strings.TrimSpace(dom.Find("div > div > div > div:first-child > span:first-child").Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			pageNum := getPageNum(e.Request.URL.String())
			res := structures.Result{
				URL:          linkText,
				Rank:         -1,
				SERank:       -1,
				SEPage:       pageNum,
				SEOnPageRank: pageRankCounter[pageNum] + 1,
				Title:        titleText,
				Description:  descText,
				SearchEngine: seName,
			}
			pageRankCounter[pageNum]++

			bucket.SetResult(&res, relay, options, pagesCol)
		}
	})

	col.Visit(seURL + query + "&start=0")
	for i := 1; i < options.MaxPages; i++ {
		col.Visit(seURL + query + "&start=" + strconv.Itoa(i*10))
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}

func getPageNum(uri string) int {
	urll, err := url.Parse(uri)
	if err != nil {
		fmt.Println(err)
	}
	qry := urll.Query()
	startString := qry.Get("start")
	startInt, _ := strconv.Atoi(startString)
	return startInt/10 + 1
}
