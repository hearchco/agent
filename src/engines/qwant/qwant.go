package qwant

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
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

	sedefaults.InitializeCollectors(&col, &pagesCol, &options, &timings)

	sedefaults.PagesColRequest(Info.Name, pagesCol, ctx)
	sedefaults.PagesColError(Info.Name, pagesCol)
	sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	sedefaults.ColRequest(Info.Name, col, ctx)
	sedefaults.ColError(Info.Name, col)

	col.OnResponse(func(r *colly.Response) {
		var pageStr string = r.Ctx.Get("page")
		if pageStr == "" {
			// If I'm using GET as the first page
			return
		}

		page, _ := strconv.Atoi(pageStr)

		var parsedResponse QwantResponse
		err := json.Unmarshal(r.Body, &parsedResponse)
		if err != nil {
			log.Error().Err(err).Msgf("%v: Failed body unmarshall to json:\n%v", Info.Name, string(r.Body))
		}

		mainline := parsedResponse.Data.Res.Items.Mainline
		counter := 1
		for _, group := range mainline {
			if group.Type != "web" {
				continue
			}
			for _, result := range group.Items {
				goodURL := parse.ParseURL(result.URL)

				res := bucket.MakeSEResult(goodURL, result.Title, result.Description, Info.Name, page, counter)
				bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
				counter += 1
			}
		}
	})

	locale := getLocale(&options)
	nRequested := settings.RequestedResultsPerPage
	device := getDevice(&options)
	safeSearch := getSafeSearch(&options)

	for i := 0; i < options.MaxPages; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		reqString := Info.URL + query + "&count=" + strconv.Itoa(nRequested) + "&locale=" + locale + "&offset=" + strconv.Itoa(i*nRequested) + "&device=" + device + "&safesearch=" + safeSearch

		sedefaults.DoGetRequest(reqString, colCtx, col, Info.Name, &retError)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}

func getLocale(options *engines.Options) string {
	locale := options.Locale
	locale = strings.ToLower(locale)
	locale = strings.ReplaceAll(locale, "-", "_")
	return locale
}

func getDevice(options *engines.Options) string {
	if options.Mobile {
		return "mobile"
	}
	return "desktop"
}

func getSafeSearch(options *engines.Options) string {
	if options.SafeSearch {
		return "1"
	}
	return "0"
}

/*
col.OnHTML("div[data-testid=\"sectionWeb\"] > div > div", func(e *colly.HTMLElement) {
	//first page
	idx := e.Index

	dom := e.DOM
	baseDOM := dom.Find("div[data-testid=\"webResult\"] > div > div > div > div > div")
	hrefElement := baseDOM.Find("a[data-testid=\"serTitle\"]")
	linkHref, hrefExists := hrefElement.Attr("href")
	linkText := parse.ParseURL(linkHref)
	titleText := strings.TrimSpace(hrefElement.Text())
	descText := strings.TrimSpace(baseDOM.Find("div > span").Text())

	if hrefExists && linkText != "" && linkText != "#" && titleText != "" {
		var pageStr string = e.Request.Ctx.Get("page")
		page, _ := strconv.Atoi(pageStr)

		res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, -1, page, idx+1)
		bucket.AddSEResult(res, Info.Name, relay, options, pagesCol)
	} else {
		log.Info().Msgf("Not Good! %v\n%v\n%v", linkText, titleText, descText)
	}
})
*/
