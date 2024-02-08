package swisscows

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/bucket"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/engines"
	"github.com/hearchco/hearchco/src/search/parse"
	"github.com/hearchco/hearchco/src/sedefaults"
	"github.com/rs/zerolog/log"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) error {
	if err := sedefaults.Prepare(Info.Name, &options, &settings, &Support, &Info, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	sedefaults.InitializeCollectors(&col, &pagesCol, &settings, &options, &timings)

	sedefaults.PagesColRequest(Info.Name, pagesCol, ctx)
	sedefaults.PagesColError(Info.Name, pagesCol)
	sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	sedefaults.ColRequest(Info.Name, col, ctx)
	sedefaults.ColError(Info.Name, col)

	col.OnRequest(func(r *colly.Request) {
		if r.Method == "OPTIONS" {
			return
		}

		var qry string = "?" + r.URL.RawQuery
		nonce, sig, err := generateAuth(qry)
		if err != nil {
			log.Error().Err(err).Msg("swisscows.Search() -> col.OnRequest: failed building request: failed generating auth")
			return
		}

		r.Headers.Set("X-Request-Nonce", nonce)
		r.Headers.Set("X-Request-Signature", sig)
		r.Headers.Set("Pragma", "no-cache")
	})

	col.OnResponse(func(r *colly.Response) {
		query := r.Request.URL.Query().Get("query")
		urll := r.Request.URL.String()
		anonUrll := anonymize.Substring(urll, query)
		log.Trace().
			Str("url", anonUrll).
			Str("nonce", r.Request.Headers.Get("X-Request-Nonce")).
			Str("signature", r.Request.Headers.Get("X-Request-Signature")).
			Msg("swisscows.Search() -> col.OnResponse()")

		var pageStr string = r.Ctx.Get("page")
		page, _ := strconv.Atoi(pageStr)

		var parsedResponse SCResponse
		err := json.Unmarshal(r.Body, &parsedResponse)
		if err != nil {
			log.Error().
				Err(err).
				Str("body", string(r.Body)).
				Msg("swisscows.Search() -> col.OnResponse(): failed body unmarshall to json")

			return
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

	localeParam := getLocale(&options)

	var colCtx *colly.Context

	for i := 0; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		//col.Request("OPTIONS", seAPIURL+"freshness=All&itemsCount="+strconv.Itoa(sResCount)+"&offset="+strconv.Itoa(i*10)+"&query="+query+localeURL, nil, colCtx, nil)
		//col.Wait()

		urll := Info.URL + "freshness=All&itemsCount=" + strconv.Itoa(settings.RequestedResultsPerPage) + "&offset=" + strconv.Itoa(i*10) + "&query=" + query + localeParam
		anonUrll := Info.URL + "freshness=All&itemsCount=" + strconv.Itoa(settings.RequestedResultsPerPage) + "&offset=" + strconv.Itoa(i*10) + "&query=" + anonymize.String(query) + localeParam
		sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, &retError)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}

func getLocale(options *engines.Options) string {
	return "&region=" + strings.Replace(options.Locale, "_", "-", 1)
}

/*
var pageRankCounter []int = make([]int, options.MaxPages*Info.ResPerPage)
col.OnHTML("div.web-results > article.item-web", func(e *colly.HTMLElement) {
	dom := e.DOM

	linkHref, hrefExists := dom.Find("a.site").Attr("href")
	linkText := parse.ParseURL(linkHref)
	titleText := strings.TrimSpace(dom.Find("h2.title").Text())
	descText := strings.TrimSpace(dom.Find("p.description").Text())

	if hrefExists && linkText != "" && linkText != "#" && titleText != "" {
		var pageStr string = e.Request.Ctx.Get("page")
		page, _ := strconv.Atoi(pageStr)

		res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, -1, page, pageRankCounter[page]+1)
		bucket.AddSEResult(res, Info.Name, relay, options, pagesCol)
		pageRankCounter[page]++
	} else {
		log.Trace().
			Str("engine", Info.Name.String()).
			Str("url", linkText).
			Str("title", titleText).
			Str("description", descText).
			Msg("Matched Result, but couldn't retrieve data")
	}
})
*/
