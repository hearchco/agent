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
	"github.com/hearchco/hearchco/src/search/parse"
	"github.com/rs/zerolog/log"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, &options, &settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		linkEl := dom.Find(dompaths.Link)
		linkHref, hrefExists := linkEl.Attr("href")
		var linkText string

		if linkHref[0] == 'h' {
			//normal link
			linkText = parse.ParseURL(linkHref)
		} else {
			//telemetry link, e.g. //web.search.ch/r/redirect?event=website&origin=result!u377d618861533351/https://de.wikipedia.org/wiki/Charles_Paul_Wilp
			linkText = parse.ParseURL("http" + strings.Split(linkHref, "http")[1]) //works for https, dont worry
		}

		titleText := strings.TrimSpace(linkEl.Text())
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())

		if hrefExists && linkText != "" && linkText != "#" && titleText != "" {
			page := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
			bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol)
			pageRankCounter[page]++
		}
	})

	col.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "Sorry for the CAPTCHA") {
			log.Error().
				Str("engine", Info.Name.String()).
				Msg("Returned CAPTCHA")
		}
	})

	safeSearchParam := getSafeSearch(options)

	errChannel := make(chan error, options.MaxPages)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	_sedefaults.DoPostRequest(Info.URL, strings.NewReader("query="+query+"&country=web&language=all"+safeSearchParam), colCtx, col, Info.Name, errChannel)
	col.Wait() // wait so I can get the JSESSION cookie back

	for i := 1; i < options.MaxPages; i++ {
		pageStr := strconv.Itoa(i + 1)
		colCtx = colly.NewContext()
		colCtx.Put("page", pageStr)

		// query not needed as its saved in the session
		_sedefaults.DoGetRequest(pageURL+pageStr, pageURL+pageStr, colCtx, col, Info.Name, errChannel)
	}

	retErrors := make([]error, 0)
	for i := 0; i < options.MaxPages; i++ {
		err := <-errChannel
		if err != nil {
			retErrors = append(retErrors, err)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retErrors
}

func getSafeSearch(options engines.Options) string {
	if options.SafeSearch {
		return "&safeSearch=true"
	}
	return "&safeSearch=false"
}
