package presearch

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

	safeSearch := getSafeSearch(options.SafeSearch)

	col.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", "use_local_search_results=false")
		r.Headers.Add("Cookie", "ai_results_disable=1")
		r.Headers.Add("Cookie", "use_safe_search="+safeSearch)
	})

	col.OnResponse(func(r *colly.Response) {
		var pageStr string = r.Request.Ctx.Get("page")
		page, _ := strconv.Atoi(pageStr)

		var apiStr string = r.Request.Ctx.Get("isAPI")
		isApi, _ := strconv.ParseBool(apiStr)

		if isApi {
			//json response
			var pr PresearchResponse
			err := json.Unmarshal(r.Body, &pr)
			if err != nil {
				log.Error().
					Err(err).
					Str("SEName", Info.Name.String()).
					Str("Body", string(r.Body)).
					Msg("Failed body unmarshall to json")
			}

			counter := 1
			for _, result := range pr.Results.StandardResults {
				goodURL := parse.ParseURL(result.Link)
				goodTitle := parse.ParseTextWithHTML(result.Title)
				goodDesc := parse.ParseTextWithHTML(result.Desc)

				res := bucket.MakeSEResult(goodURL, goodTitle, goodDesc, Info.Name, page, counter)
				bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
				counter += 1
			}
		} else {
			//html response, forward call to API
			suff := strings.SplitN(string(r.Body), "window.searchId = \"", 2)[1]
			searchId := strings.SplitN(suff, "\"", 2)[0]

			nextCtx := colly.NewContext()
			nextCtx.Put("page", strconv.Itoa(page))
			nextCtx.Put("isAPI", "true")
			err := col.Request("GET", "https://presearch.com/results?id="+searchId, nil, nextCtx, nil)
			if engines.IsTimeoutError(err) {
				log.Trace().
					Err(err).
					Str("SEName", Info.Name.String()).
					Msg("failed requesting with API")
			} else if err != nil {
				log.Error().
					Err(err).
					Str("SEName", Info.Name.String()).
					Msg("failed requesting with API")
			}
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))
	colCtx.Put("isAPI", "false")

	sedefaults.DoGetRequest(Info.URL+query, colCtx, col, Info.Name, &retError)

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		colCtx.Put("isAPI", "false")

		sedefaults.DoGetRequest(Info.URL+query+"&page="+strconv.Itoa(i+1), colCtx, col, Info.Name, &retError)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}

func getSafeSearch(ss bool) string {
	if ss {
		return "true"
	}
	return "false"
}
