package google

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/rank"
	"github.com/tminaorg/brzaguza/src/sedefaults"
	"github.com/tminaorg/brzaguza/src/structures"
	"github.com/tminaorg/brzaguza/src/utility"
)

const SEDomain string = "www.google.com"

const seName string = "Google"
const seURL string = "https://www.google.com/search?q="
const resPerPage int = 10

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

	var pageRankCounter []int = make([]int, options.MaxPages*resPerPage)

	col.OnHTML("div.g", func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, _ := dom.Find("a").Attr("href")
		linkText := utility.ParseURL(linkHref)
		titleText := strings.TrimSpace(dom.Find("div > div > div > a > h3").Text())
		descText := strings.TrimSpace(dom.Find("div > div > div > div:first-child > span:first-child").Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)

			res := structures.Result{
				URL:          linkText,
				Rank:         -1,
				SERank:       -1,
				SEPage:       page,
				SEOnPageRank: pageRankCounter[page] + 1,
				Title:        titleText,
				Description:  descText,
				SearchEngine: seName,
			}
			if config.InsertDefaultRank {
				res.Rank = rank.DefaultRank(res.SERank, res.SEPage, res.SEOnPageRank)
			}
			pageRankCounter[page]++

			bucket.SetResult(&res, relay, options, pagesCol)
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))
	col.Request("GET", seURL+query, nil, colCtx, nil)
	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		col.Request("GET", seURL+query+"&start="+strconv.Itoa(i*10), nil, colCtx, nil)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
