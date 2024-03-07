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
	"github.com/hearchco/hearchco/src/search/parse"
	"github.com/rs/zerolog/log"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, &options, &settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

	var pageRankCounter []int = make([]int, options.Pages.Max*Info.ResultsPerPage)

	col.OnRequest(func(r *colly.Request) {
		r.Headers.Del("Accept")
	})

	col.OnResponse(func(r *colly.Response) {
		body := string(r.Body)
		index := strings.Index(body, "{\"results\":")

		if index == -1 || body[len(body)-1] != ']' {
			log.Error().
				Str("body", body).
				Str("engine", Info.Name.String()).
				Msg("failed parsing response: failed finding start and/or end of JSON")
			return
		}

		body = body[index : len(body)-1]
		var content JsonResponse
		if err := json.Unmarshal([]byte(body), &content); err != nil {
			log.Error().
				Err(err).
				Str("engine", Info.Name.String()).
				Str("content", body).
				Msg("Failed unmarshalling content")
			return
		}

		page := _sedefaults.PageFromContext(r.Request.Ctx, Info.Name)
		for _, result := range content.Results {
			if result.TType != "Organic" {
				continue
			}

			goodURL := parse.ParseURL(result.URL)
			goodTitle := parse.ParseTextWithHTML(result.Title)
			goodDescription := parse.ParseTextWithHTML(result.Snippet)

			res := bucket.MakeSEResult(goodURL, goodTitle, goodDescription, Info.Name, page, pageRankCounter[page]+1)
			bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol)
			pageRankCounter[page]++
		}
	})

	retErrors := make([]error, options.Pages.Start+options.Pages.Max)

	// static params
	localeParam := getLocale(options)
	safeSearchParam := getSafeSearch(options)

	// starts from at least 0
	for i := options.Pages.Start; i < options.Pages.Start+options.Pages.Max; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		// dynamic params
		pageParam := ""
		// i == 0 is the first page
		if i > 0 {
			pageParam = "&limit=" + strconv.Itoa((i+2)*10+1)
		}

		urll := Info.URL + "client=web" + localeParam + pageParam + "&no_correct=false&q=" + query + safeSearchParam + "&type=web"
		anonUrll := Info.URL + "client=web" + localeParam + pageParam + "&no_correct=false&q=" + anonymize.String(query) + safeSearchParam + "&type=web"

		err = _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
		retErrors[i] = err
	}

	col.Wait()
	pagesCol.Wait()

	return _sedefaults.NonNilErrorsFromSlice(retErrors)
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
