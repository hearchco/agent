package duckduckgo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

	sedefaults.ColRequest(Info.Name, col, &ctx, &retError)
	sedefaults.ColError(Info.Name, col, &retError)

	col.OnHTML(dompaths.ResultsContainer, func(e *colly.HTMLElement) {
		var linkText string
		var linkScheme string
		var titleText string
		var descText string
		var rrank int

		var pageStr string = e.Request.Ctx.Get("page")
		page, _ := strconv.Atoi(pageStr)

		e.DOM.Children().Each(func(i int, row *goquery.Selection) {
			switch i % 4 {
			case 0:
				rankText := strings.TrimSpace(row.Children().First().Text())
				fmt.Sscanf(rankText, "%d", &rrank)
				linkHref, _ := row.Find(dompaths.Link).Attr("href")
				if strings.Contains(linkHref, "https") {
					linkScheme = "https://"
				} else {
					linkScheme = "http://"
				}
				titleText = strings.TrimSpace(row.Find(dompaths.Title).Text())
			case 1:
				descText = strings.TrimSpace(row.Find(dompaths.Description).Text())
			case 2:
				rawURL := linkScheme + row.Find("td > span.link-text").Text()
				linkText = parse.ParseURL(rawURL)
			case 3:
				if linkText != "" && linkText != "#" && titleText != "" {
					res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, (i/4 + 1))
					bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
				}
			}
		})
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	err := col.Request("GET", Info.URL+"?q="+query, nil, colCtx, nil)
	if err != nil && !engines.IsTimeoutError(err) {
		log.Error().Err(err).Msgf("%v: failed requesting with GET method", Info.Name)
	}

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		err := col.Request("POST", Info.URL, strings.NewReader("q="+query+"&dc="+strconv.Itoa(i*20)), colCtx, nil)
		if err != nil && !engines.IsTimeoutError(err) {
			log.Error().Err(err).Msgf("%v: failed requesting with POST method on page", Info.Name)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
