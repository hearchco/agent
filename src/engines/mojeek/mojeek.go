package mojeek

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/rank"
	"github.com/tminaorg/brzaguza/src/search/limit"
	"github.com/tminaorg/brzaguza/src/search/useragent"
	"github.com/tminaorg/brzaguza/src/structures"
	"github.com/tminaorg/brzaguza/src/utility"
)

const SEDomain string = "www.mojeek.com"

const seName string = "Mojeek"
const seURL string = "https://www.mojeek.com/search?q="
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
		log.Debug().Msgf("%v: Pages Collector; Error OnError:\nURL: %v\nError: %v", seName, r.Ctx.Get("originalURL"), err)
		log.Trace().Msgf("%v: HTML Response:\n%v", seName, string(r.Body))
		//retError = err
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

	col.OnHTML("ul.results-standard > li", func(e *colly.HTMLElement) {
		dom := e.DOM

		titleEl := dom.Find("h2 > a.title")
		linkHref, _ := titleEl.Attr("href")
		linkText := utility.ParseURL(linkHref)
		titleText := strings.TrimSpace(titleEl.Text())
		descText := strings.TrimSpace(dom.Find("p.s").Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)

			res := structures.Result{
				URL:          linkText,
				Rank:         -1,
				SERank:       (page-1)*resPerPage + pageRankCounter[page] + 1,
				SEPage:       page,
				SEOnPageRank: pageRankCounter[page] + 1,
				Title:        titleText,
				Description:  descText,
				SearchEngine: seName,
			}
			if config.InsertDefaultRank {
				res.Rank = rank.DefaultRank(res.SERank, res.SEPage, res.SEOnPageRank)
			}
			pageRankCounter[page]++

			bucket.SetResult(&res, relay, options, pagesCol)
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))
	col.Request("GET", seURL+query, nil, colCtx, nil)
	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		col.Request("GET", seURL+query+"&s="+strconv.Itoa(i*10+1), nil, colCtx, nil)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
