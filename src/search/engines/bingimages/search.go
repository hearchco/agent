package bingimages

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/hearchco/agent/src/utils/morestrings"
)

func (se Engine) Search(query string, opts options.Options, resChan chan result.ResultScraped) ([]error, bool) {
	foundResults := atomic.Bool{}
	retErrors := make([]error, 0, opts.Pages.Max)
	pageRankCounter := scraper.NewPageRankCounter(opts.Pages.Max)

	se.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		var jsonMetadata jsonMetadata
		metadataS, metadataExists := dom.Find(dompaths.Metadata.Path).Attr(dompaths.Metadata.Attr)
		if !metadataExists {
			log.Error().
				Str("engine", se.Name.String()).
				Msg("Matched result, but couldn't retrieve data")
			return
		}

		if err := json.Unmarshal([]byte(metadataS), &jsonMetadata); err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("jsonMetadata", metadataS).
				Msg("Failed to unmarshal metadata")
			return
		}

		if jsonMetadata.ImageURL == "" || jsonMetadata.PageURL == "" || jsonMetadata.ThumbnailURL == "" {
			log.Error().
				Caller().
				Str("engine", se.Name.String()).
				Str("jsonMetadata", metadataS).
				Str("url", jsonMetadata.PageURL).
				Str("original", jsonMetadata.ImageURL).
				Str("thumbnail", jsonMetadata.ThumbnailURL).
				Msg("Couldn't find image, thumbnail, or page URL")
			return
		}

		titleText := strings.TrimSpace(dom.Find(dompaths.Title).Text())
		if titleText == "" {
			// Could also use the json data ("t" field), it seems to include weird/erroneous characters though (particularly '\ue000' and '\ue001').
			log.Error().
				Caller().
				Str("engine", se.Name.String()).
				Str("jsonMetadata", metadataS).
				Msg("Couldn't find title")
			return
		}

		// This returns "2000 x 1500 · jpeg".
		imgFormatS := strings.TrimSpace(dom.Find(dompaths.ImgFormatStr).Text())
		if imgFormatS == "" {
			log.Trace().
				Caller().
				Str("engine", se.Name.String()).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Msg("Couldn't find image format (probably a video)")
			return
		}

		// Convert to "2000x1500·jpeg".
		imgFormatS = strings.ReplaceAll(imgFormatS, " ", "")
		// Remove everything after 2000x1500.
		imgFormatS = strings.Split(imgFormatS, "·")[0]
		// Create height and width.
		imgFormat := strings.Split(imgFormatS, "x")

		origH, err := strconv.Atoi(imgFormat[0])
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("engine", se.Name.String()).
				Str("height", imgFormat[0]).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Str("imgFormatS", imgFormatS).
				Msg("Failed to convert original height to int")
			return
		}

		origW, err := strconv.Atoi(imgFormat[1])
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("engine", se.Name.String()).
				Str("width", imgFormat[1]).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Str("imgFormatS", imgFormatS).
				Msg("Failed to convert original width to int")
			return
		}

		found := false
		var thmbHS, thmbWS string
		for _, thmb := range dompaths.Thumbnail {
			var thmbHExists, thmbWExists bool
			thmbHS, thmbHExists = dom.Find(thmb.Path).Attr(thmb.Height)
			thmbWS, thmbWExists = dom.Find(thmb.Path).Attr(thmb.Width)
			if thmbHExists && thmbWExists {
				found = true
				break
			}
		}

		if !found {
			log.Error().
				Caller().
				Str("engine", se.Name.String()).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Str("height", thmbHS).
				Str("width", thmbWS).
				Msg("Couldn't find thumbnail format")
			return
		}

		thmbH, err := strconv.Atoi(thmbHS)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("engine", se.Name.String()).
				Str("height", thmbHS).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Msg("Failed to convert thumbnail height to int")
			return
		}

		thmbW, err := strconv.Atoi(thmbWS)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("engine", se.Name.String()).
				Str("width", thmbWS).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Msg("Failed to convert thumbnail width to int")
			return
		}

		source := strings.TrimSpace(dom.Find(dompaths.Source).Text())
		if source == "" {
			log.Error().
				Caller().
				Str("engine", se.Name.String()).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Msg("Couldn't find source")
			return
		}

		pageIndex := se.PageFromContext(e.Request.Ctx)
		page := pageIndex + opts.Pages.Start + 1

		r, err := result.ConstructImagesResult(
			se.Name, jsonMetadata.ImageURL, titleText, jsonMetadata.Desc, page, pageRankCounter.GetPlusOne(pageIndex),
			origH, origW, thmbH, thmbW, jsonMetadata.ThumbnailURL, source, jsonMetadata.PageURL,
		)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("result", fmt.Sprintf("%v", r)).
				Msg("Failed to construct result")
		} else {
			log.Trace().
				Caller().
				Int("page", page).
				Int("rank", pageRankCounter.GetPlusOne(pageIndex)).
				Str("result", fmt.Sprintf("%v", r)).
				Msg("Sending result to channel")
			resChan <- r
			pageRankCounter.Increment(pageIndex)
			if !foundResults.Load() {
				foundResults.Store(true)
			}
		}
	})

	se.OnResponse(func(r *colly.Response) {
		if len(r.Body) == 0 {
			log.Error().
				Str("engine", se.Name.String()).
				Msg("Got empty response, probably too many requests")
		}
	})

	// Static params.
	paramLocale := localeParamString(opts.Locale)

	for i := range opts.Pages.Max {
		pageNum0 := i + opts.Pages.Start
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(i))

		// Dynamic params.
		paramPage := fmt.Sprintf("%v=%v", paramKeyPage, pageNum0*35+1)

		combinedParams := morestrings.JoinNonEmpty("&", "&", paramAsync, paramPage, paramCount, paramLocale)

		urll := fmt.Sprintf("%v?q=%v%v", searchURL, query, combinedParams)
		anonUrll := fmt.Sprintf("%v?q=%v%v", searchURL, anonymize.String(query), combinedParams)

		if err := se.Get(ctx, urll, anonUrll); err != nil {
			retErrors = append(retErrors, err)
		}
	}

	se.Wait()
	close(resChan)
	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
}
