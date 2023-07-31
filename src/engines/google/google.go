package google

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/rank"
	"github.com/tminaorg/brzaguza/src/search/limit"
	"github.com/tminaorg/brzaguza/src/search/useragent"
	"github.com/tminaorg/brzaguza/src/structures"
)

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
	log.Trace().Msgf("%v", options.UserAgent)

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
			r.Abort()
			retError = err
			return
		}
		r.Ctx.Put("originalURL", r.URL.String())
	})

	pagesCol.OnError(func(r *colly.Response, err error) {
		retError = err
	})

	pagesCol.OnResponse(func(r *colly.Response) {
		urll := r.Ctx.Get("originalURL")

		setResultResponse(urll, r, relay)
	})

	col.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil { // dont fully understand this
			r.Abort()
			retError = err
			return
		}
	})

	col.OnError(func(r *colly.Response, err error) {
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
				Rank:        -1,
				SEPageRank:  pageRankCounter[pageNum],
				SEPage:      pageNum,
				URL:         linkText,
				Title:       titleText,
				Description: descText,
			}
			pageRankCounter[pageNum]++

			setResult(&res, relay, options, pagesCol)
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

func setResult(result *structures.Result, relay *structures.Relay, options *structures.Options, pagesCol *colly.Collector) {
	log.Trace().Msgf("Got Result %v: %v", result.Title, result.URL)

	relay.Mutex.Lock()
	mapRes, exists := relay.ResultMap[result.URL]

	if !exists {
		relay.ResultMap[result.URL] = result
	} else if len(mapRes.Description) < len(result.Description) {
		mapRes.Description = result.Description
	}
	relay.Mutex.Unlock()

	if !exists && options.VisitPages {
		pagesCol.Visit(result.URL)
	}
}

func setResultResponse(link string, response *colly.Response, relay *structures.Relay) {
	log.Trace().Msgf("Got Response %v", link)

	relay.Mutex.Lock()
	mapRes, exists := relay.ResultMap[link]

	if !exists {
		log.Error().Msgf("URL not in map when adding response! Should not be possible. URL: %v", link)
		relay.Mutex.Unlock()
		return
	}

	mapRes.Response = response

	resCopy := *mapRes
	rankAddr := &(mapRes.Rank)
	relay.Mutex.Unlock()
	rank.SetRank(&resCopy, rankAddr, &(relay.Mutex)) //copy contains pointer to response
}

func getPageNum(uri string) int {
	urll, err := url.Parse(uri)
	if err != nil {
		fmt.Println(err)
	}
	qry := urll.Query()
	startString := qry.Get("start")
	startInt, _ := strconv.Atoi(startString)
	return startInt / 10
}
