package swisscows

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/sedefaults"
	"github.com/tminaorg/brzaguza/src/structures"
	"github.com/tminaorg/brzaguza/src/utility"
)

type SCItem struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Desc       string `json:"description"`
	URL        string `json:"url"`
	DisplayURL string `json:"displayUrl"`
}

type SCResponse struct {
	Items []SCItem `json:"items"`
}

const SEDomain string = "swisscows.com"

const seName string = "Swisscows"
const seAPIURL string = "https://api.swisscows.com/web/search?"
const sResCount int = 10
const locale string = "en-US"

const defaultResultsPerPage int = 10

// const seURL string = "https://swisscows.com/en/web?query="

func Search(ctx context.Context, query string, relay *structures.Relay, options *structures.SEOptions, settings *config.SESettings) error {
	if err := sedefaults.FunctionPrepare(seName, options, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	sedefaults.InitializeCollectors(&col, &pagesCol, options)

	sedefaults.PagesColRequest(seName, pagesCol, &ctx, &retError)
	sedefaults.PagesColError(seName, pagesCol)
	sedefaults.PagesColResponse(seName, pagesCol, relay)

	col.OnRequest(func(r *colly.Request) {
		if err := (ctx).Err(); err != nil {
			log.Error().Msgf("%v: SE Collector; Error OnRequest %v", seName, r)
			r.Abort()
			retError = err
			return
		}

		if r.Method == "OPTIONS" {
			return
		}

		var qry string = "?" + r.URL.RawQuery
		nonce, sig := generateAuth(qry)

		//log.Debug().Msgf("qry: %v\nnonce: %v\nsignature: %v", qry, nonce, sig)

		r.Headers.Set("X-Request-Nonce", nonce)
		r.Headers.Set("X-Request-Signature", sig)
		r.Headers.Set("Pragma", "no-cache")
	})

	col.OnError(func(r *colly.Response, err error) {
		log.Error().Msgf("%v: SE Collector - OnError.\nMethod: %v\nURL: %v\nError: %v", seName, r.Request.Method, r.Request.URL.String(), err)
		log.Error().Msgf("%v: HTML Response written to %v%v_col.log.html", seName, config.LogDumpLocation, seName)
		writeErr := os.WriteFile(config.LogDumpLocation+seName+"_col.log.html", r.Body, 0644)
		if writeErr != nil {
			log.Error().Err(writeErr)
		}
		retError = err
	})

	var pageRankCounter []int = make([]int, options.MaxPages*sResCount)

	// not used
	col.OnHTML("div.web-results > article.item-web", func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, _ := dom.Find("a.site").Attr("href")
		linkText := utility.ParseURL(linkHref)
		titleText := strings.TrimSpace(dom.Find("h2.title").Text())
		descText := strings.TrimSpace(dom.Find("p.description").Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)

			res := bucket.MakeSEResult(linkText, titleText, descText, seName, -1, page, pageRankCounter[page]+1)
			bucket.AddSEResult(res, seName, relay, options, pagesCol)
			pageRankCounter[page]++
		} else {
			log.Trace().Msgf("%v: Matched Result, but couldn't retrieve data.\nURL:%v\nTitle:%v\nDescription:%v", seName, linkText, titleText, descText)
		}
	})

	col.OnResponse(func(r *colly.Response) {
		log.Trace().Msgf("URL: %v\nNonce: %v\nSig: %v", r.Request.URL.String(), r.Request.Headers.Get("X-Request-Nonce"), r.Request.Headers.Get("X-Request-Signature"))

		var pageStr string = r.Ctx.Get("page")
		page, _ := strconv.Atoi(pageStr)

		var parsedResponse SCResponse
		err := json.Unmarshal(r.Body, &parsedResponse)
		if err != nil {
			log.Error().Err(err).Msgf("%v: Failed body unmarshall to json:\n%v", seName, string(r.Body))
		}

		counter := 0
		for _, result := range parsedResponse.Items {
			goodURL := utility.ParseURL(result.URL)
			title := utility.ParseTextWithHTML(result.Title)
			desc := utility.ParseTextWithHTML(result.Desc)

			res := bucket.MakeSEResult(goodURL, title, desc, seName, -1, page, counter%defaultResultsPerPage+1)
			bucket.AddSEResult(res, seName, relay, options, pagesCol)
			counter += 1
		}
	})

	var colCtx *colly.Context
	for i := 0; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		//col.Request("OPTIONS", seAPIURL+"freshness=All&itemsCount="+strconv.Itoa(sResCount)+"&offset="+strconv.Itoa(i*10)+"&query="+query+"&region="+locale, nil, colCtx, nil)
		//col.Wait()
		col.Request("GET", seAPIURL+"freshness=All&itemsCount="+strconv.Itoa(sResCount)+"&offset="+strconv.Itoa(i*10)+"&query="+query+"&region="+locale, nil, colCtx, nil)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
