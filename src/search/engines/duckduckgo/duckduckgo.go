package duckduckgo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

	localeCookie := getLocale(options)

	col.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", localeCookie)
	})

	col.OnHTML(dompaths.ResultsContainer, func(e *colly.HTMLElement) {
		var linkText, linkScheme, titleText, descText string
		var hrefExists bool
		var rrank int

		pageIndex := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)
		page := pageIndex + options.Pages.Start + 1

		e.DOM.Children().Each(func(i int, row *goquery.Selection) {
			switch i % 4 {
			case 0:
				rankText := strings.TrimSpace(row.Children().First().Text())
				fmt.Sscanf(rankText, "%d", &rrank)
				var linkHref string
				linkHref, hrefExists = row.Find(dompaths.Link).Attr("href")
				if strings.Contains(linkHref, "https") {
					linkScheme = "https://"
				} else {
					linkScheme = "http://"
				}
				titleText = _sedefaults.SanitizeTitle(row.Find(dompaths.Title).Text())
			case 1:
				descText = _sedefaults.SanitizeDescription(row.Find(dompaths.Description).Text())
			case 2:
				rawURL := linkScheme + row.Find("td > span.link-text").Text()
				linkText = _sedefaults.SanitizeURL(rawURL)
			case 3:
				if !hrefExists {
					log.Error().
						Str("engine", Info.Name.String()).
						Str("url", linkText).
						Str("title", titleText).
						Str("description", descText).
						Str("link selector", dompaths.Link).
						Msg("duckduckgo.Search(): href attribute doesn't exist on matched URL element")
					return
				}

				res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, (i/4 + 1))
				bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol, nEnabledEngines)
			}
		})
	})

	retErrors := make([]error, 0, options.Pages.Max)

	// starts from at least 0
	for i := options.Pages.Start; i < options.Pages.Start+options.Pages.Max; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i-options.Pages.Start))

		var err error
		// i == 0 is the first page
		if i == 0 {
			urll := Info.URL + "?q=" + query
			anonUrll := Info.URL + "?q=" + anonymize.String(query)
			err = _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
		} else {
			requestData := strings.NewReader("q=" + query + "&dc=" + strconv.Itoa(i*20))
			err = _sedefaults.DoPostRequest(Info.URL, requestData, colCtx, col, Info.Name)
		}

		if err != nil {
			retErrors = append(retErrors, err)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retErrors[:len(retErrors):len(retErrors)]
}

func getLocale(options engines.Options) string {
	spl := strings.SplitN(strings.ToLower(options.Locale), "_", 2)
	return "kl=" + spl[1] + "-" + spl[0]
}
