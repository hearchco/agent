package presearch

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/search/parse"
	"github.com/tminaorg/brzaguza/src/sedefaults"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) error {
	if err := sedefaults.Prepare(Info.Name, &options, &settings, &Support, &Info, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	sedefaults.InitializeCollectors(&col, &pagesCol, &options, &timings)

	sedefaults.PagesColRequest(Info.Name, pagesCol, &ctx, &retError)
	sedefaults.PagesColError(Info.Name, pagesCol)
	sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	sedefaults.ColError(Info.Name, col, &retError)

	safeSearch := getSafeSearch(options.SafeSearch)

	col.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			log.Error().Msgf("%v: SE Collector; Error OnRequest %v", Info.Name, r)
			r.Abort()
			retError = err
			return
		}

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
				log.Error().Err(err).Msgf("%v: Failed body unmarshall to json:\n%v", Info.Name, string(r.Body))
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
			if err != nil && !engines.IsTimeoutError(err) {
				log.Error().Err(err).Msgf("%v: failed requesting with API", Info.Name)
			}
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))
	colCtx.Put("isAPI", "false")

	err := col.Request("GET", Info.URL+query, nil, colCtx, nil)
	if err != nil && !engines.IsTimeoutError(err) {
		log.Error().Err(err).Msgf("%v: failed requesting with GET method", Info.Name)
	}

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		colCtx.Put("isAPI", "false")

		err := col.Request("GET", Info.URL+query+"&page="+strconv.Itoa(i+1), nil, colCtx, nil)
		if err != nil && !engines.IsTimeoutError(err) {
			log.Error().Err(err).Msgf("%v: failed requesting with GET method on page", Info.Name)
		}
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
