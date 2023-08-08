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
	"github.com/tminaorg/brzaguza/src/sedefaults"
	"github.com/tminaorg/brzaguza/src/structures"
	"github.com/tminaorg/brzaguza/src/utility"
)

const SEDomain string = "lite.duckduckgo.com"

const seName string = "DuckDuckGo"
const seURL string = "https://lite.duckduckgo.com/lite/"

func Search(ctx context.Context, query string, relay *structures.Relay, options *structures.SEOptions, settings *config.SESettings) error {
	if err := sedefaults.FunctionPrepare(seName, options, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	sedefaults.InitializeCollectors(&col, &pagesCol, options)

	sedefaults.PagesColRequest(seName, pagesCol, &ctx, &retError)
	sedefaults.PagesColError(seName, pagesCol)
	sedefaults.PagesColResponse(seName, pagesCol, relay)

	col.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil { // dont fully understand this
			log.Error().Msgf("%v: SE Collector; Error OnRequest %v", seName, r)
			r.Abort()
			retError = err
			return
		}
		if r.Body == nil {
			//This is the first page, so this isnt a POST request
			r.Ctx.Put("body", "q="+query+"&dc=1")
		} else {
			var reqBody []byte
			r.Body.Read(reqBody)
			r.Ctx.Put("body", string(reqBody))
		}
	})
	sedefaults.ColError(seName, col, &retError)

	col.OnHTML("div.filters > table > tbody", func(e *colly.HTMLElement) {
		var linkText string
		var linkScheme string
		var titleText string
		var descText string
		var rrank int

		var reqBody string = e.Request.Ctx.Get("body")
		var page int
		fmt.Sscanf(reqBody, "q="+query+"&dc=%d", &page)
		page = page/20 + 1

		e.DOM.Children().Each(func(i int, row *goquery.Selection) {
			switch i % 4 {
			case 0:
				rankText := strings.TrimSpace(row.Children().First().Text())
				fmt.Sscanf(rankText, "%d", &rrank)
				linkHref, _ := row.Find("a.result-link").Attr("href")
				if strings.Contains(linkHref, "https") {
					linkScheme = "https://"
				} else {
					linkScheme = "http://"
				}
				titleText = strings.TrimSpace(row.Find("td > a.result-link").Text())
			case 1:
				descText = strings.TrimSpace(row.Find("td.result-snippet").Text())
			case 2:
				rawURL := linkScheme + row.Find("td > span.link-text").Text()
				linkText = utility.ParseURL(rawURL)
			case 3:
				if linkText != "" && linkText != "#" && titleText != "" {
					res := bucket.MakeSEResult(linkText, titleText, descText, seName, rrank, page, (i/4 + 1))
					bucket.AddSEResult(res, seName, relay, options, pagesCol)
				}
			}
		})
	})

	col.Visit(seURL + "?q=" + query)
	//col.PostRaw(seURL, []byte("q="+query+"&dc=1"))
	for i := 1; i < options.MaxPages; i++ {
		col.PostRaw(seURL, []byte("q="+query+"&dc="+strconv.Itoa(i*20)))
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
