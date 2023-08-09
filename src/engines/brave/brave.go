package brave

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/sedefaults"
	"github.com/tminaorg/brzaguza/src/structures"
	"github.com/tminaorg/brzaguza/src/utility"
)

const SEDomain string = "search.brave.com"

const seName string = "Brave"
const seURL string = "https://search.brave.com/search?q="
const resPerPage int = 20

func Search(ctx context.Context, query string, relay *structures.Relay, options *structures.Options) error {
	if err := sedefaults.FunctionPrepare(seName, options, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	sedefaults.InitializeCollectors(&col, &pagesCol, options, nil)

	sedefaults.PagesColRequest(seName, pagesCol, &ctx, &retError)
	sedefaults.PagesColError(seName, pagesCol)
	sedefaults.PagesColResponse(seName, pagesCol, relay)

	sedefaults.ColRequest(seName, col, &ctx, &retError)
	sedefaults.ColError(seName, col, &retError)

	var pageRankCounter []int = make([]int, options.MaxPages*resPerPage)

	col.OnHTML("div#results > div[class*=\"snippet fdb\"][data-type=\"web\"]", func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, _ := dom.Find("a.result-header").Attr("href")
		linkText := utility.ParseURL(linkHref)
		titleText := strings.TrimSpace(dom.Find("a.result-header > span.snippet-title").Text())
		descText := strings.TrimSpace(dom.Find("div.snippet-content > p.snippet-description").Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			if descText == "" {
				descText = strings.TrimSpace(dom.Find("div.product > div.flex-hcenter > div > div[class=\"text-sm text-gray\"]").Text())
			}
			if descText == "" {
				descText = strings.TrimSpace(dom.Find("p.snippet-description").Text())
			}

			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)

			res := bucket.MakeSEResult(linkText, titleText, descText, seName, -1, page, pageRankCounter[page]+1)
			bucket.AddSEResult(res, seName, relay, options, pagesCol)
			pageRankCounter[page]++
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))
	col.Request("GET", seURL+query+"&source=web", nil, colCtx, nil)
	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		col.Request("GET", seURL+query+"&spellcheck=0&offset="+strconv.Itoa(i), nil, colCtx, nil)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
