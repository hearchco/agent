package yandex

//4-5 digit
//https://yandex.com/dev/xml/doc/en/concepts/get-request <- api
//https://yandex.com/search/?text=a+hard+man&p=2&rnd=4498

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

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

const seName string = "Yandex"
const seURL string = "https://yandex.com/search/?text="
const resPerPage int = 10

// const locale string = "en-US,en"
const sendRND bool = false

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

		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Referer", "https://yandex.com/")
		r.Headers.Set("DNT", "1")
		r.Headers.Set("Connection", "keep-alive")
		if len(col.Cookies("https://yandex.com")) == 0 {
			r.Headers.Set("Cookie", r.Ctx.Get("cook"))
		}
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "same-origin")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Sec-GPC", "1")
		r.Headers.Set("sec-ch-ua-platform", "Linux") // UserAgent specific
		r.Headers.Set("sec-ch-ua", "\"Google Chrome\";v=\"112\", \"Chromium\";v=\"112\", \"Not=A?Brand\";v=\"24\"")
		r.Headers.Set("sec-ch-ua-mobile", "?0")
		r.Headers.Set("TE", "trailers")

		log.Trace().Msgf("Request Headers:\n%v", r.Headers)
		log.Trace().Msgf("The Cookies:")
		ccookies := col.Cookies("https://yandex.com")
		for cookie := range ccookies {
			log.Trace().Msgf("%v", cookie)
		}
	})

	col.OnError(func(r *colly.Response, err error) {
		log.Error().Msgf("%v: SE Collector; Error OnError:\nURL: %v\nError: %v", seName, r.Request.URL.String(), err)
		log.Trace().Msgf("%v: HTML Response:\n%v", seName, string(r.Body))
		retError = err
	})

	var pageRankCounter []int = make([]int, options.MaxPages*resPerPage)

	col.OnHTML("ul#search-result > li", func(e *colly.HTMLElement) {
		dom := e.DOM

		log.Info().Msg("YANDEX GRACED US")

		linkHref, _ := dom.Find("div:first-child > a").Attr("href")
		linkText := utility.ParseURL(linkHref)
		titleText := strings.TrimSpace(dom.Find("div:first-child > a > h2 > span").Text())
		descText := strings.TrimSpace(dom.Find("div:nth-child(3) > div > span").Text())

		fmt.Printf("%v %v %v %v", linkHref, linkText, titleText, descText)

		if linkText != "" && linkText != "#" && titleText != "" {
			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)

			res := structures.Result{
				URL:          linkText,
				Rank:         -1,
				SERank:       -1,
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

	col.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "Please confirm that you and not a robot are sending requests") {
			log.Error().Msgf("%v: Returned captcha.", seName)
		} else {
			log.Info().Msgf("NOT IN CAPTHA WOOOOOOOO!")
		}

		log.Trace().Msgf("Response headers:\n%v", r.Headers)
	})

	rndText := ""
	if sendRND {
		randS := rand.New(rand.NewSource(time.Now().UnixNano()))
		rndText = "&rnd=" + strconv.Itoa(randS.Intn(99999-1000+1)+1000) //random number from [1000,99999]
	}

	var cookies string
	readCookies("./src/engines/yandex/cookie.db", &cookies, col)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))
	colCtx.Put("cook", cookies)

	col.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
	col.Request("GET", seURL+query, nil, colCtx, nil)
	//col.Request("GET", seURL+query+"&search_source=yacom_desktop_common&msid=1691021728532406-396109311027094500-balancer-l7leveler-kubr-yp-vla-50-BAL-8238", nil, colCtx, nil)

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		col.Request("GET", seURL+query+"&p="+strconv.Itoa(i*10)+rndText, nil, colCtx, nil)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}

func readCookies(filename string, cookies *string, col *colly.Collector) {
	dat, err := os.ReadFile("./" + filename)
	if err != nil {
		panic(err)
	}

	*cookies = string(dat)
}
