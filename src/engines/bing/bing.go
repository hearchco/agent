package bing

import (
	"context"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/rank"
	"github.com/tminaorg/brzaguza/src/sedefaults"
	"github.com/tminaorg/brzaguza/src/structures"
	"github.com/tminaorg/brzaguza/src/utility"
)

const SEDomain string = "www.bing.com"

const seName string = "Bing"
const seURL string = "https://www.bing.com/search?q="
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

	col.OnHTML("ol#b_results > li.b_algo", func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, _ := dom.Find("h2 > a").Attr("href")
		linkText := utility.ParseURL(linkHref)
		titleText := strings.TrimSpace(dom.Find("h2 > a").Text())
		descText := strings.TrimSpace(dom.Find("div.b_caption").Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			if descText == "" {
				descText = strings.TrimSpace(dom.Find("p.b_algoSlug").Text())
			}
			if strings.Contains(descText, "Web") {
				descText = strings.Split(descText, "Web")[1]
			}

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
		} else {
			log.Trace().Msgf("%v: Matched Result, but couldn't retrieve data.\nURL:%v\nTitle:%v\nDescription:%v", seName, linkText, titleText, descText)
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))
	col.Request("GET", seURL+query, nil, colCtx, nil)
	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		col.Request("GET", seURL+query+"&first="+strconv.Itoa(i*10+1), nil, colCtx, nil)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
