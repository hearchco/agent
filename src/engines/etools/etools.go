package etools

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

const SEDomain string = "www.etools.ch"

const seName string = "Etools"
const seURL string = "https://www.etools.ch/searchSubmit.do"

// const seGETURL string = "https://www.etools.ch/searchSubmit.do?query="
// https://www.etools.ch/search.do?page=<page number>&query=<query>
const sePAGEURL string = "https://www.etools.ch/search.do?page="

const resultsPerPage int = 10

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

	var col *colly.Collector = colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(options.UserAgent)) //site breaks if this is Async
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
		urll := r.Ctx.Get("originalURL") //because i may have followed redirects

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

	col.OnHTML("table.result > tbody > tr", func(e *colly.HTMLElement) {
		dom := e.DOM

		linkEl := dom.Find("td.record > a")
		linkHref, _ := linkEl.Attr("href")
		var linkText string

		if linkHref[0] == 'h' {
			//normal link
			linkText = utility.ParseURL(linkHref)
		} else {
			//telemetry link, e.g. //web.search.ch/r/redirect?event=website&origin=result!u377d618861533351/https://de.wikipedia.org/wiki/Charles_Paul_Wilp
			linkText = utility.ParseURL("http" + strings.Split(linkHref, "http")[1]) //works for https, dont worry
		}

		titleText := strings.TrimSpace(linkEl.Text())
		descText := strings.TrimSpace(dom.Find("td.record > div.text").Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)
			seRankString := strings.TrimSpace(dom.Find("td[class=\"count help\"]").Text())
			seRank, convErr := strconv.Atoi(seRankString)
			if convErr != nil {
				log.Error().Err(convErr).Msgf("%v: SERank string to int conversion error. URL: %v, SERank string: %v", seName, linkText, seRankString)
			}

			//var onPageRank int = e.Index // this should also work, but is a bit more volatile
			var onPageRank int = (seRank-1)%resultsPerPage + 1

			res := structures.Result{
				URL:          linkText,
				Rank:         -1,
				SERank:       seRank,
				SEPage:       page,
				SEOnPageRank: onPageRank,
				Title:        titleText,
				Description:  descText,
				SearchEngine: seName,
			}
			if config.InsertDefaultRank {
				res.Rank = rank.DefaultRank(res.SERank, res.SEPage, res.SEOnPageRank)
			}

			bucket.SetResult(&res, relay, options, pagesCol)
		}
	})

	col.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "Sorry for the CAPTCHA") {
			log.Error().Msgf("%v: Returned captcha.", seName)
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	//also takes "&token=5d8d98d9a968388eeb4191afa00ca469"
	col.Request("POST", seURL, strings.NewReader("query="+query+"&country=web&language=all"), colCtx, nil)
	col.Wait() //wait so I can get the JSESSION cookie back

	for i := 1; i < options.MaxPages; i++ {
		pageStr := strconv.Itoa(i + 1)
		colCtx = colly.NewContext()
		colCtx.Put("page", pageStr)
		col.Request("GET", sePAGEURL+pageStr, nil, colCtx, nil)

		//col.Wait()

		//time.Sleep(200 * time.Millisecond)
		//a delay can help reduce response volatility for this site
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
