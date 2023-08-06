package qwant

import (
	"context"
	"encoding/json"
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

const SEDomain string = "www.qwant.com"

const seName string = "Qwant"

// const seURL string = "https://www.qwant.com/?q="
const seAPIURL string = "https://api.qwant.com/v3/search/web?q="

const defaultResultsPerPage int = 10

const qSafeSearch string = "0" //sets safeSearch for Qwant
const qDevice string = "desktop"
const qLocale string = "en_us"
const qResCount int = 10 //ask this many results back

type QwantResults struct {
	Title       string `json:"title"`
	URL         string `json:"url"` //there is also a source field, what is it?
	Description string `json:"desc"`
}

type QwantMainlineItems struct {
	Type  string         `json:"type"`
	Items []QwantResults `json:"items"`
}

type QwantResponse struct {
	Status string `json:"status"`
	Data   struct {
		Res struct {
			Items struct {
				Mainline []QwantMainlineItems `json:"mainline"`
			} `json:"items"`
		} `json:"result"`
	} `json:"data"`
}

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

	col.OnResponse(func(r *colly.Response) {
		var pageStr string = r.Ctx.Get("page")
		if pageStr == "" {
			//the first page
			return
		}

		page, _ := strconv.Atoi(pageStr)

		var parsedResponse QwantResponse
		err := json.Unmarshal(r.Body, &parsedResponse)
		if err != nil {
			log.Error().Err(err).Msgf("%v: Failed body unmarshall to json:\n%v", seName, string(r.Body))
		}

		mainline := parsedResponse.Data.Res.Items.Mainline
		counter := 0
		for _, group := range mainline {
			if group.Type != "web" {
				continue
			}
			for _, result := range group.Items {

				goodURL := utility.ParseURL(result.URL)
				res := structures.Result{
					URL:          goodURL,
					Rank:         -1,
					SERank:       (page-1)*qResCount + counter,
					SEPage:       page,
					SEOnPageRank: counter%defaultResultsPerPage + 1,
					Title:        result.Title,
					Description:  result.Description,
					SearchEngine: seName,
				}
				if config.InsertDefaultRank {
					res.Rank = rank.DefaultRank(res.SERank, res.SEPage, res.SEOnPageRank)
				}

				bucket.SetResult(&res, relay, options, pagesCol)
				counter += 1
			}
		}
	})

	col.OnHTML("div[data-testid=\"sectionWeb\"] > div > div", func(e *colly.HTMLElement) {
		//first page
		idx := e.Index

		dom := e.DOM
		baseDOM := dom.Find("div[data-testid=\"webResult\"] > div > div > div > div > div")
		hrefElement := baseDOM.Find("a[data-testid=\"serTitle\"]")
		linkHref, _ := hrefElement.Attr("href")
		linkText := utility.ParseURL(linkHref)
		titleText := strings.TrimSpace(hrefElement.Text())
		descText := strings.TrimSpace(baseDOM.Find("div > span").Text())

		if linkText != "" && linkText != "#" && titleText != "" {
			res := structures.Result{
				URL:          linkText,
				Rank:         -1,
				SERank:       -1,
				SEPage:       1,
				SEOnPageRank: idx + 1,
				Title:        titleText,
				Description:  descText,
				SearchEngine: seName,
			}

			bucket.SetResult(&res, relay, options, pagesCol)
		} else {
			log.Info().Msgf("Not Good! %v\n%v\n%v", linkText, titleText, descText)
		}
	})

	for i := 0; i < options.MaxPages; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))
		reqString := seAPIURL + query + "&count=" + strconv.Itoa(qResCount) + "&locale=" + qLocale + "&offset=" + strconv.Itoa(i*qResCount) + "&device=" + qDevice + "&safesearch=" + qSafeSearch
		col.Request("GET", reqString, nil, colCtx, nil)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
