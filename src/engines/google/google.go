package google

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/rank"
	"github.com/tminaorg/brzaguza/src/search"
	"github.com/tminaorg/brzaguza/src/search/useragent"
	"github.com/tminaorg/brzaguza/src/structures"
)

const url string = "https://www.google.com/search?q="

func Search(ctx context.Context, query string, relay *structures.Relay, options *structures.Options) error {
	if ctx == nil {
		ctx = context.Background()
	} //^ not necessary as ctx is always passed in search.go, branch predictor will skip this if

	if err := search.RateLimit.Wait(ctx); err != nil {
		return err
	}

	if options.UserAgent == "" {
		options.UserAgent = useragent.DefaultUserAgent()
	}
	log.Trace().Msgf("%v\n", options.UserAgent)

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
			res := structures.Result{
				Rank:        -1,
				URL:         linkText,
				Title:       titleText,
				Description: descText,
			}

			_, exists := relay.ResultMap[res.URL]

			if !exists || len(relay.ResultMap[res.URL].Description) < len(descText) {
				relay.ResultChannel <- res
			}

			if !exists && options.VisitPages {
				pagesCol.Visit(linkText)
			}
		}
	})

	col.Visit(url + query)
	for i := 1; i < options.MaxPages; i++ {
		col.Visit(url + query + "&start=" + strconv.Itoa(i*10))
	}

	col.Wait() //order should be irrelevant, right? \/
	pagesCol.Wait()

	relay.EngineDoneChannel <- true

	return retError
}
