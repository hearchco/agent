package etools

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
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

func (e Engine) Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.CategoryTimings, salt string, nEnabledEngines int) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, options, settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

	pageRankCounter := make([]int, options.Pages.Max)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		linkText, titleText, descText := _sedefaults.RawFieldsFromDOM(e.DOM, dompaths, Info.Name) // telemetry url isnt valid link so cant pass it to FieldsFromDOM (?)

		// Need to perform this check here so the check below (linkText[0] != 'h') doesn't panic
		if linkText == "" {
			log.Error().
				Str("title", titleText).
				Str("description", descText).
				Msg("etools.Search(): invalid result, url is empty.")
			return
		}

		if linkText[0] != 'h' {
			//telemetry link, e.g. //web.search.ch/r/redirect?event=website&origin=result!u377d618861533351/https://de.wikipedia.org/wiki/Charles_Paul_Wilp
			linkText = "http" + strings.Split(linkText, "http")[1] //works for https, dont worry
		}

		linkText, titleText, descText = _sedefaults.SanitizeFields(linkText, titleText, descText)

		pageIndex := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)
		page := pageIndex + options.Pages.Start + 1

		res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[pageIndex]+1)
		valid := bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol, nEnabledEngines)
		if valid {
			pageRankCounter[pageIndex]++
		}
	})

	col.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "Sorry for the CAPTCHA") {
			log.Error().
				Str("engine", Info.Name.String()).
				Msg("Returned CAPTCHA")
		}
	})

	retErrors := make([]error, 0, options.Pages.Max)

	// static params
	safeSearchParam := getSafeSearch(options)

	// starts from at least 0
	for i := options.Pages.Start; i < options.Pages.Start+options.Pages.Max; i++ {
		pageStr := strconv.Itoa(i - options.Pages.Start)
		colCtx := colly.NewContext()
		colCtx.Put("page", pageStr)

		var err error
		// i == 0 is the first page
		if i == 0 {
			requestData := strings.NewReader("query=" + query + "&country=web&language=all" + safeSearchParam)
			err = _sedefaults.DoPostRequest(Info.URL, requestData, colCtx, col, Info.Name)
			col.Wait() // col.Wait() is needed to save the JSESSION cookie
		} else {
			// query is not needed as it's saved in the JSESSION cookie
			err = _sedefaults.DoGetRequest(pageURL+pageStr, pageURL+pageStr, colCtx, col, Info.Name)
		}

		if err != nil {
			retErrors = append(retErrors, err)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retErrors[:len(retErrors):len(retErrors)]
}

func getSafeSearch(options engines.Options) string {
	if options.SafeSearch {
		return "&safeSearch=true"
	}
	return "&safeSearch=false"
}
