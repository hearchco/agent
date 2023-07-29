package google

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/sourcegraph/conc"
	"github.com/tminaorg/brzaguza/src/rank"
	"github.com/tminaorg/brzaguza/src/relay"
	"github.com/tminaorg/brzaguza/src/search"
	"github.com/tminaorg/brzaguza/src/search/useragent"
	"github.com/tminaorg/brzaguza/src/structures"
)

const url string = "https://www.google.com/search?q="

const avgResultsPerRequest int = 10 // cant be zero - divide by zero error

func Search(ctx context.Context, query string, options *structures.Options, worker *conc.WaitGroup) error {
	if ctx == nil {
		ctx = context.Background()
	} //^ not necessary as ctx is always passed in search.go, branch predictor will skip this if

	if err := search.RateLimit.Wait(ctx); err != nil {
		return err
	}

	if options.UserAgent == "" {
		options.UserAgent = useragent.DefaultUserAgent()
	}

	var numPages int
	if options.JustFirstPage {
		numPages = 1
	} else {
		numPages = options.Limit / avgResultsPerRequest
	}

	var col *colly.Collector

	if numPages == 1 {
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
	})

	pagesCol.OnError(func(r *colly.Response, err error) {
		retError = err
	})

	pagesCol.OnResponse(func(r *colly.Response) {
		rnk := rank.RankPage(r)
		relay.RankChannel <- structures.ResultRank{
			URL:  r.Request.URL.String(),
			Rank: rnk,
		}
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

	col.OnHTML("div.g", func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, _ := dom.Find("a").Attr("href")
		linkText := strings.TrimSpace(linkHref)
		titleText := strings.TrimSpace(dom.Find("div > div > div > a > h3").Text())
		descText := strings.TrimSpace(dom.Find("div > div > div > div:first-child > span:first-child").Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			var res structures.Result = structures.Result{
				Rank:        -1,
				URL:         linkText,
				Title:       titleText,
				Description: descText,
			}
			relay.ResultChannel <- res
			pagesCol.Visit(linkText)
		}
	})

	col.Visit(url + query)
	for i := 2; i <= numPages; i++ {
		col.Visit(url + query + "&start=" + strconv.Itoa(i*10))
	}

	if numPages != 1 { //order should be irrelevant, right? \/
		col.Wait()
	}
	pagesCol.Wait()

	relay.EngineDoneChannel <- true

	return retError
}
