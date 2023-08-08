package etools

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/sedefaults"
	"github.com/tminaorg/brzaguza/src/structures"
	"github.com/tminaorg/brzaguza/src/utility"
)

const SEDomain string = "www.etools.ch"

const seName string = "Etools"
const seURL string = "https://www.etools.ch/searchSubmit.do"

// const seGETURL string = "https://www.etools.ch/searchSubmit.do?query="
// https://www.etools.ch/search.do?page=<page number>&query=<query>
const sePAGEURL string = "https://www.etools.ch/search.do?page="

const resultsPerPage int = 10

func Search(ctx context.Context, query string, relay *structures.Relay, options *structures.Options) error {
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

	sedefaults.ColRequest(seName, col, &ctx, &retError)
	sedefaults.ColError(seName, col, &retError)

	col.OnHTML("table.result > tbody > tr", func(e *colly.HTMLElement) {
		dom := e.DOM

		linkEl := dom.Find("td.record > a")
		linkHref, _ := linkEl.Attr("href")
		var linkText string

		if linkHref[0] == 'h' {
			//normal link
			linkText = utility.ParseURL(linkHref)
		} else {
			//telemetry link, e.g. //web.search.ch/r/redirect?event=website&origin=result!u377d618861533351/https://de.wikipedia.org/wiki/Charles_Paul_Wilp
			linkText = utility.ParseURL("http" + strings.Split(linkHref, "http")[1]) //works for https, dont worry
		}

		titleText := strings.TrimSpace(linkEl.Text())
		descText := strings.TrimSpace(dom.Find("td.record > div.text").Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)
			seRankString := strings.TrimSpace(dom.Find("td[class=\"count help\"]").Text())
			seRank, convErr := strconv.Atoi(seRankString)
			if convErr != nil {
				log.Error().Err(convErr).Msgf("%v: SERank string to int conversion error. URL: %v, SERank string: %v", seName, linkText, seRankString)
			}

			//var onPageRank int = e.Index // this should also work, but is a bit more volatile
			var onPageRank int = (seRank-1)%resultsPerPage + 1

			res := bucket.MakeSEResult(linkText, titleText, descText, seName, seRank, page, onPageRank)
			bucket.AddSEResult(res, seName, relay, options, pagesCol)
		}
	})

	col.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "Sorry for the CAPTCHA") {
			log.Error().Msgf("%v: Returned captcha.", seName)
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	//also takes "&token=5d8d98d9a968388eeb4191afa00ca469"
	col.Request("POST", seURL, strings.NewReader("query="+query+"&country=web&language=all"), colCtx, nil)
	col.Wait() //wait so I can get the JSESSION cookie back

	for i := 1; i < options.MaxPages; i++ {
		pageStr := strconv.Itoa(i + 1)
		colCtx = colly.NewContext()
		colCtx.Put("page", pageStr)
		col.Request("GET", sePAGEURL+pageStr, nil, colCtx, nil)

		//col.Wait()

		//time.Sleep(200 * time.Millisecond)
		//a delay can help reduce response volatility for this site
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
