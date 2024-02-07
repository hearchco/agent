package googleimages

import (
	"context"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/engines/_sedefaults"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) error {
	if err := _sedefaults.Prepare(Info.Name, &options, &settings, &Support, &Info, &ctx); err != nil {
		return err
	}

	var col *colly.Collector
	var pagesCol *colly.Collector
	var retError error

	_sedefaults.InitializeCollectors(&col, &pagesCol, &options, &timings)

	_sedefaults.PagesColRequest(Info.Name, pagesCol, ctx)
	_sedefaults.PagesColError(Info.Name, pagesCol)
	_sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	_sedefaults.ColRequest(Info.Name, col, ctx)
	_sedefaults.ColError(Info.Name, col)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	col.OnResponse(func(e *colly.Response) {
		body := string(e.Body)
		index := strings.Index(body, "{\"ischj\":")
		if index != -1 {
			body = body[index:]
			var jsonResponse JsonResponse
			if err := json.Unmarshal([]byte(body), &jsonResponse); err != nil {
				log.Error().
					Str("body", body).
					Msg("googleimages.Search() -> col.OnResponse: failed parsing response: failed unmarshalling JSON")
			} else {
				page := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

				metadata := jsonResponse.ISCHJ.Metadata[0]

				origImg := metadata.OriginalImage
				original := result.Image{
					URL:    origImg.Url,
					Height: origImg.Height,
					Width:  origImg.Width,
				}

				thmbImg := metadata.Thumbnail
				thumbnail := result.Image{
					URL:    thmbImg.Url,
					Height: thmbImg.Height,
					Width:  thmbImg.Width,
				}

				resultJson := metadata.Result
				textInGridJson := metadata.TextInGrid

				res := bucket.MakeSEImageResult(resultJson.ReferrerUrl, resultJson.PageTitle, textInGridJson.Snippet, resultJson.SiteTitle, original, thumbnail, Info.Name, page, pageRankCounter[page]+1)
				bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
				pageRankCounter[page]++
			}
		} else {
			log.Error().
				Str("body", body).
				Msg("googleimages.Search() -> col.OnResponse: failed parsing response: failed finding start of JSON")
		}
	})

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	_sedefaults.DoGetRequest(Info.URL+query+params+"1", colCtx, col, Info.Name, &retError)

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		_sedefaults.DoGetRequest(Info.URL+query+params+strconv.Itoa(i*10), colCtx, col, Info.Name, &retError)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}
