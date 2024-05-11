package presearch

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

func (e Engine) Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.CategoryTimings, salt string, nEnabledEngines int) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, &options, &settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

	safeSearch := getSafeSearch(options.SafeSearch)

	col.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", "use_local_search_results=false")
		r.Headers.Add("Cookie", "ai_results_disable=1")
		r.Headers.Add("Cookie", "use_safe_search="+safeSearch)
	})

	col.OnResponse(func(r *colly.Response) {
		pageIndex := _sedefaults.PageFromContext(r.Request.Ctx, Info.Name)
		page := pageIndex + options.Pages.Start + 1

		var apiStr string = r.Request.Ctx.Get("isAPI")
		isApi, _ := strconv.ParseBool(apiStr)

		if isApi {
			//json response
			var pr PresearchResponse
			err := json.Unmarshal(r.Body, &pr)
			if err != nil {
				log.Error().
					Err(err).
					Str("engine", Info.Name.String()).
					Bytes("body", r.Body).
					Msg("Failed body unmarshall to json")
			}

			counter := 1
			for _, result := range pr.Results.StandardResults {
				goodURL, goodTitle, goodDesc := _sedefaults.SanitizeFields(result.Link, result.Title, result.Desc)

				res := bucket.MakeSEResult(goodURL, goodTitle, goodDesc, Info.Name, page, counter)
				valid := bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol, nEnabledEngines)
				if valid {
					counter += 1
				}
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
					Str("engine", Info.Name.String()).
					Msg("failed requesting with API")
			} else if err != nil {
				log.Error().
					Err(err).
					Str("engine", Info.Name.String()).
					Msg("failed requesting with API")
			}
		}
	})

	retErrors := make([]error, 0, options.Pages.Max)

	// starts from at least 0
	for i := options.Pages.Start; i < options.Pages.Start+options.Pages.Max; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i-options.Pages.Start))
		colCtx.Put("isAPI", "false")

		// dynamic params
		pageParam := ""
		// i == 0 is the first page
		if i > 0 {
			pageParam = "&page=" + strconv.Itoa(i+1)
		}

		urll := Info.URL + query + pageParam
		anonUrll := Info.URL + anonymize.String(query) + pageParam

		err := _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
		if err != nil {
			retErrors = append(retErrors, err)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retErrors[:len(retErrors):len(retErrors)]
}

func getSafeSearch(ss bool) string {
	if ss {
		return "true"
	}
	return "false"
}
