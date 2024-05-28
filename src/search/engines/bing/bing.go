package bing

import (
	"context"
	"encoding/base64"
	"net/url"
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

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		linkText, titleText, descText := _sedefaults.FieldsFromDOM(e.DOM, dompaths, Info.Name) // the telemetry link is a valid link so it can be sanitized
		linkText = _sedefaults.SanitizeURL(removeTelemetry(linkText))

		if descText == "" {
			descText = e.DOM.Find("p.b_algoSlug").Text()
		}
		if strings.Contains(descText, "Web") {
			descText = strings.Split(descText, "Web")[1]
		}
		descText = _sedefaults.SanitizeDescription(descText)

		pageIndex := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)
		page := pageIndex + options.Pages.Start + 1

		res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[pageIndex]+1)
		valid := bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol, nEnabledEngines)
		if valid {
			pageRankCounter[pageIndex]++
		}
	})

	retErrors := make([]error, 0, options.Pages.Max)

	// static params
	localeParam := getLocale(options)

	// starts from at least 0
	for i := options.Pages.Start; i < options.Pages.Start+options.Pages.Max; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i-options.Pages.Start))

		// dynamic params
		pageParam := ""
		// i == 0 is the first page
		if i > 0 {
			pageParam = "&first=" + strconv.Itoa(i*10+1)
		}

		urll := Info.URL + query + pageParam + localeParam
		anonUrll := Info.URL + anonymize.String(query) + pageParam + localeParam

		err := _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
		if err != nil {
			retErrors = append(retErrors, err)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retErrors[:len(retErrors):len(retErrors)]
}

func removeTelemetry(link string) string {
	if strings.HasPrefix(link, "https://www.bing.com/ck/a?") {
		parsedUrl, err := url.Parse(link)
		if err != nil {
			log.Error().Err(err).Str("url", link).Msg("bing.removeTelemetry(): error parsing url")
			return ""
		}

		// get the first value of u parameter and remove "a1" in front
		encodedUrl := parsedUrl.Query().Get("u")[2:]

		cleanUrl, err := base64.RawURLEncoding.DecodeString(encodedUrl)
		if err != nil {
			log.Error().Err(err).Msg("bing.removeTelemetry(): failed decoding string from base64")
		}
		return string(cleanUrl)
	}
	return link
}

func getLocale(options engines.Options) string {
	spl := strings.SplitN(strings.ToLower(options.Locale), "_", 2)
	return "&setlang=" + spl[0] + "&cc=" + spl[1]
}
