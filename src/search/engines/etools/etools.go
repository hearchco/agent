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

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) error {
	if err := _sedefaults.Prepare(Info.Name, &options, &settings, &Support, &Info, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	_sedefaults.InitializeCollectors(&col, &pagesCol, &settings, &options, &timings)

	_sedefaults.PagesColRequest(Info.Name, pagesCol, ctx)
	_sedefaults.PagesColError(Info.Name, pagesCol)
	_sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	_sedefaults.ColRequest(Info.Name, col, ctx)
	_sedefaults.ColError(Info.Name, col)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		linkText, titleText, descText := _sedefaults.RawFieldsFromDOM(e.DOM, &dompaths, Info.Name) // telemetry url isnt valid link so cant pass it to FieldsFromDOM (?)

		if linkText[0] != 'h' {
			//telemetry link, e.g. //web.search.ch/r/redirect?event=website&origin=result!u377d618861533351/https://de.wikipedia.org/wiki/Charles_Paul_Wilp
			linkText = "http" + strings.Split(linkText, "http")[1] //works for https, dont worry
		}

		linkText, titleText, descText = _sedefaults.SanitizeFields(linkText, titleText, descText)

		page := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

		res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
		bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
		pageRankCounter[page]++
	})

	col.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "Sorry for the CAPTCHA") {
			log.Error().
				Str("engine", Info.Name.String()).
				Msg("Returned CAPTCHA")
		}
	})

	safeSearchParam := getSafeSearch(&options)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	_sedefaults.DoPostRequest(Info.URL, strings.NewReader("query="+query+"&country=web&language=all"+safeSearchParam), colCtx, col, Info.Name, &retError)
	col.Wait() // wait so I can get the JSESSION cookie back

	for i := 1; i < options.MaxPages; i++ {
		pageStr := strconv.Itoa(i + 1)
		colCtx = colly.NewContext()
		colCtx.Put("page", pageStr)

		// query not needed as its saved in the session
		_sedefaults.DoGetRequest(pageURL+pageStr, pageURL+pageStr, colCtx, col, Info.Name, &retError)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}

func getSafeSearch(options *engines.Options) string {
	if options.SafeSearch {
		return "&safeSearch=true"
	}
	return "&safeSearch=false"
}
