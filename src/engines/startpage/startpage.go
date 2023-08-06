package startpage

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

const SEDomain string = "www.startpage.com"

const seName string = "Startpage"
const seURL string = "https://www.startpage.com/sp/search?q="
const resPerPage int = 10

const useSafeSearch bool = false

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

	col.OnHTML("section > div.w-gl__result", func(e *colly.HTMLElement) {
		dom := e.DOM

		linkHref, _ := dom.Find("a.result-link").Attr("href")
		linkText := utility.ParseURL(linkHref)
		titleText := strings.TrimSpace(dom.Find("a.w-gl__result-title").Text())
		descText := strings.TrimSpace(dom.Find("p.w-gl__description").Text())

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
		} else {
			log.Trace().Msgf("%v: Matched Result, but couldn't retrieve data.\nURL:%v\nTitle:%v\nDescription:%v", seName, linkText, titleText, descText)
		}
	})

	col.OnResponse(func(r *colly.Response) {
		if strings.Contains(string(r.Body), "to prevent possible abuse of our service") {
			log.Error().Msgf("%v: Request blocked by engine due to scraping.", seName)
		} else if strings.Contains(string(r.Body), "This page cannot function without javascript") {
			log.Error().Msgf("%v: Engine couldn't load requests, needs javascript.", seName)
		}
	})

	var safeSearchParameter string = ""
	if !useSafeSearch {
		safeSearchParameter += "&qadf=none"
	}

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))
	col.Request("GET", seURL+query+safeSearchParameter, nil, colCtx, nil)
	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		col.Request("GET", seURL+query+"&page="+strconv.Itoa(i+1)+safeSearchParameter, nil, colCtx, nil)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
