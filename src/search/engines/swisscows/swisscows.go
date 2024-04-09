package swisscows

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/engines/_sedefaults"
	"github.com/rs/zerolog/log"
)

type Engine struct{}

func New() Engine {
	return Engine{}
}

func (e Engine) Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings, salt string, nEnabledEngines int) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, &options, &settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

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

		pageIndex := _sedefaults.PageFromContext(r.Request.Ctx, Info.Name)
		page := pageIndex + options.Pages.Start + 1

		var parsedResponse SCResponse
		err := json.Unmarshal(r.Body, &parsedResponse)
		if err != nil {
			log.Error().
				Err(err).
				Bytes("body", r.Body).
				Msg("swisscows.Search() -> col.OnResponse(): failed body unmarshall to json")

			return
		}

		counter := 1
		for _, result := range parsedResponse.Items {
			goodLink, goodTitle, goodDesc := _sedefaults.SanitizeFields(result.URL, result.Title, result.Desc)

			res := bucket.MakeSEResult(goodLink, goodTitle, goodDesc, Info.Name, page, counter)
			valid := bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol, nEnabledEngines)
			if valid {
				counter += 1
			}
		}
	})

	retErrors := make([]error, 0, options.Pages.Max)

	// static params
	localeParam := getLocale(options)
	itemsParam := "freshness=All&itemsCount=" + strconv.Itoa(settings.RequestedResultsPerPage)

	// starts from at least 0
	for i := options.Pages.Start; i < options.Pages.Start+options.Pages.Max; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i-options.Pages.Start))
		//col.Request("OPTIONS", seAPIURL+"freshness=All&itemsCount="+strconv.Itoa(sResCount)+"&offset="+strconv.Itoa(i*10)+"&query="+query+localeURL, nil, colCtx, nil)
		//col.Wait()

		// dynamic params
		offsetParam := "&offset=" + strconv.Itoa(i*10)

		urll := Info.URL + itemsParam + offsetParam + "&query=" + query + localeParam
		anonUrll := Info.URL + itemsParam + offsetParam + "&query=" + anonymize.String(query) + localeParam

		err := _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
		if err != nil {
			retErrors = append(retErrors, err)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retErrors[:len(retErrors):len(retErrors)]
}

func getLocale(options engines.Options) string {
	return "&region=" + strings.Replace(options.Locale, "_", "-", 1)
}

/*
var pageRankCounter []int = make([]int, options.Pages.Max*Info.ResPerPage)
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
		bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol, nEnabledEngines)
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
