package swisscows

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/search/parse"
	"github.com/tminaorg/brzaguza/src/sedefaults"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings) error {
	if err := sedefaults.Prepare(Info.Name, &options, &settings, &Support, &Info, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	sedefaults.InitializeCollectors(ctx, &col, &pagesCol, &options, nil)

	sedefaults.PagesColRequest(Info.Name, pagesCol, &ctx, &retError)
	sedefaults.PagesColError(Info.Name, pagesCol)
	sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	col.OnRequest(func(r *colly.Request) {
		if err := (ctx).Err(); err != nil {
			log.Error().Msgf("%v: SE Collector; Error OnRequest %v", Info.Name, r)
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
		log.Error().Err(err).Msgf("%v: SE Collector - OnError.\nMethod: %v\nURL: %v", Info.Name, r.Request.Method, r.Request.URL.String())
		log.Error().Msgf("%v: HTML Response written to %v%v_col.log.html", Info.Name, config.LogDumpLocation, Info.Name)
		writeErr := os.WriteFile(config.LogDumpLocation+string(Info.Name)+"_col.log.html", r.Body, 0644)
		if writeErr != nil {
			log.Error().Err(writeErr)
		}
		retError = err
	})

	col.OnResponse(func(r *colly.Response) {
		log.Trace().Msgf("URL: %v\nNonce: %v\nSig: %v", r.Request.URL.String(), r.Request.Headers.Get("X-Request-Nonce"), r.Request.Headers.Get("X-Request-Signature"))

		var pageStr string = r.Ctx.Get("page")
		page, _ := strconv.Atoi(pageStr)

		var parsedResponse SCResponse
		err := json.Unmarshal(r.Body, &parsedResponse)
		if err != nil {
			log.Error().Err(err).Msgf("%v: Failed body unmarshall to json:\n%v", Info.Name, string(r.Body))
		}

		counter := 1
		for _, result := range parsedResponse.Items {
			goodURL := parse.ParseURL(result.URL)
			title := parse.ParseTextWithHTML(result.Title)
			desc := parse.ParseTextWithHTML(result.Desc)

			res := bucket.MakeSEResult(goodURL, title, desc, Info.Name, page, counter)
			bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
			counter += 1
		}
	})

	var locale string = getLocale(&options)

	var colCtx *colly.Context
	for i := 0; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		//col.Request("OPTIONS", seAPIURL+"freshness=All&itemsCount="+strconv.Itoa(sResCount)+"&offset="+strconv.Itoa(i*10)+"&query="+query+"&region="+locale, nil, colCtx, nil)
		//col.Wait()
		if err := col.Request("GET", Info.URL+"freshness=All&itemsCount="+strconv.Itoa(settings.RequestedResultsPerPage)+"&offset="+strconv.Itoa(i*10)+"&query="+query+"&region="+locale, nil, colCtx, nil); err != nil {
			log.Error().Err(err).Msg("swisscows: failed requesting with GET method")
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}

func getLocale(options *engines.Options) string {
	return options.Locale
}

/*
var pageRankCounter []int = make([]int, options.MaxPages*Info.ResPerPage)
col.OnHTML("div.web-results > article.item-web", func(e *colly.HTMLElement) {
	dom := e.DOM

	linkHref, _ := dom.Find("a.site").Attr("href")
	linkText := parse.ParseURL(linkHref)
	titleText := strings.TrimSpace(dom.Find("h2.title").Text())
	descText := strings.TrimSpace(dom.Find("p.description").Text())

	if linkText != "" && linkText != "#" && titleText != "" {
		var pageStr string = e.Request.Ctx.Get("page")
		page, _ := strconv.Atoi(pageStr)

		res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, -1, page, pageRankCounter[page]+1)
		bucket.AddSEResult(res, Info.Name, relay, options, pagesCol)
		pageRankCounter[page]++
	} else {
		log.Trace().Msgf("%v: Matched Result, but couldn't retrieve data.\nURL:%v\nTitle:%v\nDescription:%v", Info.Name, linkText, titleText, descText)
	}
})
*/
