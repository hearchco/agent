package mojeek

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/rank"
	sedefaults "github.com/tminaorg/brzaguza/src/se-defaults"
	"github.com/tminaorg/brzaguza/src/structures"
	"github.com/tminaorg/brzaguza/src/utility"
)

const SEDomain string = "www.mojeek.com"

const seName string = "Mojeek"
const seURL string = "https://www.mojeek.com/search?q="
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

	col.OnHTML("ul.results-standard > li", func(e *colly.HTMLElement) {
		dom := e.DOM

		titleEl := dom.Find("h2 > a.title")
		linkHref, _ := titleEl.Attr("href")
		linkText := utility.ParseURL(linkHref)
		titleText := strings.TrimSpace(titleEl.Text())
		descText := strings.TrimSpace(dom.Find("p.s").Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)

			res := structures.Result{
				URL:          linkText,
				Rank:         -1,
				SERank:       (page-1)*resPerPage + pageRankCounter[page] + 1,
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
		col.Request("GET", seURL+query+"&s="+strconv.Itoa(i*10+1), nil, colCtx, nil)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
