package etools

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/search/parse"
	"github.com/tminaorg/brzaguza/src/sedefaults"
	"github.com/tminaorg/brzaguza/src/structures"
)

func Search(ctx context.Context, query string, relay *structures.Relay, options structures.Options, settings config.SESettings) error {
	if err := sedefaults.Prepare(Info.Name, &options, &settings, &Support, &Info, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	sedefaults.InitializeCollectors(&col, &pagesCol, &options, nil)

	sedefaults.PagesColRequest(Info.Name, pagesCol, &ctx, &retError)
	sedefaults.PagesColError(Info.Name, pagesCol)
	sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	sedefaults.ColRequest(Info.Name, col, &ctx, &retError)
	sedefaults.ColError(Info.Name, col, &retError)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		linkEl := dom.Find(dompaths.Link)
		linkHref, _ := linkEl.Attr("href")
		var linkText string

		if linkHref[0] == 'h' {
			//normal link
			linkText = parse.ParseURL(linkHref)
		} else {
			//telemetry link, e.g. //web.search.ch/r/redirect?event=website&origin=result!u377d618861533351/https://de.wikipedia.org/wiki/Charles_Paul_Wilp
			linkText = parse.ParseURL("http" + strings.Split(linkHref, "http")[1]) //works for https, dont worry
		}

		titleText := strings.TrimSpace(linkEl.Text())
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)
			seRankString := strings.TrimSpace(dom.Find("td[class=\"count help\"]").Text())
			seRank, convErr := strconv.Atoi(seRankString)
			if convErr != nil {
				log.Error().Err(convErr).Msgf("%v: SERank string to int conversion error. URL: %v, SERank string: %v", Info.Name, linkText, seRankString)
			}

			//var onPageRank int = e.Index // this should also work, but is a bit more volatile
			var onPageRank int = (seRank-1)%Info.ResultsPerPage + 1

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, seRank, page, onPageRank)
			bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
		}
	})

	col.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "Sorry for the CAPTCHA") {
			log.Error().Msgf("%v: Returned captcha.", Info.Name)
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	col.Request("POST", Info.URL, strings.NewReader("query="+query+"&country=web&language=all"), colCtx, nil)
	col.Wait() //wait so I can get the JSESSION cookie back

	for i := 1; i < options.MaxPages; i++ {
		pageStr := strconv.Itoa(i + 1)
		colCtx = colly.NewContext()
		colCtx.Put("page", pageStr)
		col.Request("GET", pageURL+pageStr, nil, colCtx, nil)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
