package google

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/search/parse"
	"github.com/tminaorg/brzaguza/src/sedefaults"
	"github.com/tminaorg/brzaguza/src/structures"
)

// This should be in SESettings
var timings config.SETimings = config.SETimings{
	Timeout:     10 * time.Second, // the default in colly
	PageTimeout: 5 * time.Second,
	Delay:       100 * time.Millisecond,
	RandomDelay: 50 * time.Millisecond,
	Parallelism: 2, //two requests will be sent to the server, 100 + [0,50) milliseconds apart from the next two
}

func Search(ctx context.Context, query string, relay *structures.Relay, options *structures.Options, settings *config.SESettings) error {
	if err := sedefaults.FunctionPrepare(Info.Name, options, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	sedefaults.InitializeCollectors(&col, &pagesCol, options, &timings)

	sedefaults.PagesColRequest(Info.Name, pagesCol, &ctx, &retError)
	sedefaults.PagesColError(Info.Name, pagesCol)
	sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	sedefaults.ColRequest(Info.Name, col, &ctx, &retError)
	sedefaults.ColError(Info.Name, col, &retError)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResPerPage)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, _ := dom.Find(dompaths.Link).Attr("href")
		linkText := parse.ParseURL(linkHref)
		titleText := strings.TrimSpace(dom.Find(dompaths.Title).Text())
		descText := strings.TrimSpace(dom.Find(dompaths.Description).Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			var pageStr string = e.Request.Ctx.Get("page")
			page, _ := strconv.Atoi(pageStr)

			res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, -1, page, pageRankCounter[page]+1)
			bucket.AddSEResult(res, Info.Name, relay, options, pagesCol)
			pageRankCounter[page]++
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))
	col.Request("GET", Info.URL+query, nil, colCtx, nil)
	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		col.Request("GET", Info.URL+query+"&start="+strconv.Itoa(i*10), nil, colCtx, nil)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
