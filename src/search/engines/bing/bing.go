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
	"github.com/hearchco/hearchco/src/search/parse"
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

	_sedefaults.PagesColRequest(ctx, Info.Name, pagesCol)
	_sedefaults.PagesColError(Info.Name, pagesCol)
	_sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	_sedefaults.ColRequest(ctx, Info.Name, col)
	_sedefaults.ColError(Info.Name, col)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, hrefExists := dom.Find(dompaths.Link).Attr("href")
		linkText := parse.ParseURL(linkHref)
		linkText = removeTelemetry(linkText)

		titleText := strings.TrimSpace(dom.Find(dompaths.Title).Text())
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())

		if hrefExists && linkText != "" && linkText != "#" && titleText != "" {
			if descText == "" {
				descText = strings.TrimSpace(dom.Find("p.b_algoSlug").Text())
			}
			if strings.Contains(descText, "Web") {
				descText = strings.Split(descText, "Web")[1]
			}

			page := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
			bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
			pageRankCounter[page]++
		} else {
			log.Trace().
				Str("engine", Info.Name.String()).
				Str("url", linkText).
				Str("title", titleText).
				Str("description", descText).
				Msg("Matched result, but couldn't retrieve data")
		}
	})

	localeParam := getLocale(&options)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	urll := Info.URL + query + localeParam
	anonUrll := Info.URL + anonymize.String(query) + localeParam
	_sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, &retError)

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		urll := Info.URL + query + "&first=" + strconv.Itoa(i*10+1) + localeParam
		anonUrll := Info.URL + anonymize.String(query) + "&first=" + strconv.Itoa(i*10+1) + localeParam
		_sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, &retError)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
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
		return parse.ParseURL(string(cleanUrl))
	}
	return link
}

func getLocale(options *engines.Options) string {
	spl := strings.SplitN(strings.ToLower(options.Locale), "_", 2)
	return "&setlang=" + spl[0] + "&cc=" + spl[1]
}
