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
	"github.com/hearchco/hearchco/src/search/parse"
	"github.com/rs/zerolog/log"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) []error {
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

		var pageStr string = e.Request.Ctx.Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			log.Error().
				Err(err).
				Str("engine", Info.Name.String()).
				Str("page", pageStr).
				Msg("Failed to convert page number")
			return
		}

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
				titleText = strings.TrimSpace(row.Find(dompaths.Title).Text())
			case 1:
				descText = strings.TrimSpace(row.Find(dompaths.Description).Text())
			case 2:
				rawURL := linkScheme + row.Find("td > span.link-text").Text()
				linkText = parse.ParseURL(rawURL)
			case 3:
				if hrefExists && linkText != "" && linkText != "#" && titleText != "" {
					res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, page, (i/4 + 1))
					bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol)
				}
			}
		})
	})

	retErrors := make([]error, options.Pages.Start+options.Pages.Max)

	// starts from at least 0
	for i := options.Pages.Start; i < options.Pages.Start+options.Pages.Max; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		// i == 0 is the first page
		if i <= 0 {
			urll := Info.URL + "?q=" + query
			anonUrll := Info.URL + "?q=" + anonymize.String(query)
			err = _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
		} else {
			requestData := strings.NewReader("q=" + query + "&dc=" + strconv.Itoa(i*20))
			err = _sedefaults.DoPostRequest(Info.URL, requestData, colCtx, col, Info.Name)
		}

		retErrors[i] = err
	}

	col.Wait()
	pagesCol.Wait()

	realRetErrors := make([]error, 0)
	for _, err := range retErrors {
		if err != nil {
			realRetErrors = append(realRetErrors, err)
		}
	}
	return realRetErrors
}

func getLocale(options engines.Options) string {
	spl := strings.SplitN(strings.ToLower(options.Locale), "_", 2)
	return "kl=" + spl[1] + "-" + spl[0]
}
