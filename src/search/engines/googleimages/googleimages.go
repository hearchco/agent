package googleimages

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/engines/_sedefaults"
	"github.com/rs/zerolog/log"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, &options, &settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

	// disable User Agent since Google Images responds with fake data if UA is correct
	col.UserAgent = ""

	var pageRankCounter = make([]int, options.MaxPages*Info.ResultsPerPage)

	col.OnResponse(func(e *colly.Response) {
		body := string(e.Body)
		index := strings.Index(body, "{\"ischj\":")

		if index == -1 {
			log.Error().
				Str("body", body).
				Msg("googleimages.Search() -> col.OnResponse: failed parsing response: failed finding start of JSON")
			return
		}

		body = body[index:]
		var jsonResponse JsonResponse
		if err := json.Unmarshal([]byte(body), &jsonResponse); err != nil {
			log.Error().
				Str("body", body).
				Msg("googleimages.Search() -> col.OnResponse: failed parsing response: failed unmarshalling JSON")
			return
		}

		page := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

		for _, metadata := range jsonResponse.ISCHJ.Metadata {
			origImg := metadata.OriginalImage
			thmbImg := metadata.Thumbnail
			resultJson := metadata.Result
			textInGridJson := metadata.TextInGrid

			if resultJson.ReferrerUrl != "" && origImg.Url != "" && thmbImg.Url != "" {
				res := bucket.MakeSEImageResult(
					origImg.Url, resultJson.PageTitle, textInGridJson.Snippet,
					resultJson.SiteTitle, resultJson.ReferrerUrl, thmbImg.Url,
					origImg.Height, origImg.Width, thmbImg.Height, thmbImg.Width,
					Info.Name, page, pageRankCounter[page]+1,
				)
				bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol)
				pageRankCounter[page]++
			} else {
				log.Error().
					Str("engine", Info.Name.String()).
					Str("jsonMetadata", fmt.Sprintf("%v", metadata)).
					Str("url", resultJson.ReferrerUrl).
					Str("original", origImg.Url).
					Str("thumbnail", thmbImg.Url).
					Msg("googleimages.Search() -> onHTML: Couldn't find image URL")
			}
		}
	})

	errChannel := make(chan error, options.MaxPages)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	urll := Info.URL + query + params + "1"
	anonUrll := Info.URL + anonymize.String(query) + params + "1"
	_sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, errChannel)

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		urll := Info.URL + query + params + strconv.Itoa(i*10)
		anonUrll := Info.URL + anonymize.String(query) + params + strconv.Itoa(i*10)
		_sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, errChannel)
	}

	retErrors := make([]error, 0)
	for i := 0; i < options.MaxPages; i++ {
		err := <-errChannel
		if err != nil {
			retErrors = append(retErrors, err)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retErrors
}
