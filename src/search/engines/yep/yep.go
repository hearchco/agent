package yep

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
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, options, settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

	pageRankCounter := make([]int, options.Pages.Max)

	col.OnRequest(func(r *colly.Request) {
		r.Headers.Del("Accept")
	})

	col.OnResponse(func(r *colly.Response) {
		body := string(r.Body)
		start := "[\"Ok\","
		end := ']'
		index := strings.Index(body, start)

		if index != 0 || body[len(body)-1] != byte(end) {
			log.Error().
				Str("body", body).
				Str("start", start).
				Str("end", string(end)).
				Str("engine", Info.Name.String()).
				Msg("failed parsing response: failed finding start and/or end of JSON")
			return
		}

		// starts after start and ends before end
		// so after "[\"Ok\"," and before "]"
		resultsJson := body[len(start) : len(body)-1]
		var content JsonResponse
		if err := json.Unmarshal([]byte(resultsJson), &content); err != nil {
			log.Error().
				Err(err).
				Str("engine", Info.Name.String()).
				Str("body", body).
				Str("content", resultsJson).
				Msg("Failed unmarshalling content")
			return
		}

		pageIndex := _sedefaults.PageFromContext(r.Request.Ctx, Info.Name)
		page := pageIndex + options.Pages.Start + 1

		for _, result := range content.Results {
			if result.TType != "Organic" {
				continue
			}

			goodLink, goodTitle, goodDescription := _sedefaults.SanitizeFields(result.URL, result.Title, result.Snippet)

			res := bucket.MakeSEResult(goodLink, goodTitle, goodDescription, Info.Name, page, pageRankCounter[pageIndex]+1)
			valid := bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol, nEnabledEngines)
			if valid {
				pageRankCounter[pageIndex]++
			}
		}
	})

	retErrors := make([]error, 0, options.Pages.Max)

	// static params
	localeParam := getLocale(options)
	safeSearchParam := getSafeSearch(options)

	// starts from at least 0
	for i := options.Pages.Start; i < options.Pages.Start+options.Pages.Max; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i-options.Pages.Start))

		// dynamic params
		pageParam := ""
		// i == 0 is the first page
		if i > 0 {
			pageParam = "&limit=" + strconv.Itoa((i+2)*10+1)
		}

		urll := Info.URL + "client=web" + localeParam + pageParam + "&no_correct=false&q=" + query + safeSearchParam + "&type=web"
		anonUrll := Info.URL + "client=web" + localeParam + pageParam + "&no_correct=false&q=" + anonymize.String(query) + safeSearchParam + "&type=web"

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
	locale := strings.Split(options.Locale, "_")[1]
	return "&gl=" + locale
}

func getSafeSearch(options engines.Options) string {
	if options.SafeSearch {
		return "&safeSearch=strict"
	}
	return "&safeSearch=off"
}
