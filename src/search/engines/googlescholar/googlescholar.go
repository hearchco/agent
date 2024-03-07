package googlescholar

import (
	"context"
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

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, &options, &settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, hrefExists := dom.Find(dompaths.Link).Attr("href")
		linkText := removeTelemetry(parse.ParseURL(linkHref))
		titleText := strings.TrimSpace(dom.Find(dompaths.Title).Text())
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())
		citeInfo := strings.TrimSpace(dom.Find("div.gs_a").Text())

		descText = citeInfo + " || " + descText

		if hrefExists && linkText != "" && linkText != "#" && titleText != "" {
			page := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, pageRankCounter[page]+1)
			bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol)
			pageRankCounter[page]++
		}
	})

	retErrors := make([]error, options.MaxPages)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	urll := Info.URL + query
	anonUrll := Info.URL + anonymize.String(query)
	err = _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
	retErrors[0] = err

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		urll := Info.URL + query + "&start=" + strconv.Itoa(i*10)
		anonUrll := Info.URL + anonymize.String(query) + "&start=" + strconv.Itoa(i*10)
		err = _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
		retErrors[i] = err
	}

	col.Wait()
	pagesCol.Wait()

	return _sedefaults.NonNilErrorsFromSlice(retErrors)
}

func removeTelemetry(link string) string {
	parsedURL, err := url.Parse(link)
	if err != nil {
		log.Error().
			Err(err).
			Str("link", link).
			Msg("error parsing link")
		return link
	}

	// remove seemingly unused params in query
	q := parsedURL.Query()
	for _, key := range []string{"dq", "lr", "oi", "ots", "sig"} {
		q.Del(key)
	}
	parsedURL.RawQuery = q.Encode()
	return parsedURL.String()
}
