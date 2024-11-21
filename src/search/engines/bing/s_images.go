package bing

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines/options"
	"github.com/hearchco/agent/src/search/result"
	"github.com/hearchco/agent/src/search/scraper"
	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/hearchco/agent/src/utils/moreurls"
)

func (se Engine) ImageSearch(query string, opts options.Options, resChan chan result.ResultScraped) ([]error, bool) {
	foundResults := atomic.Bool{}
	retErrors := make([]error, 0, opts.Pages.Max)
	pageRankCounter := scraper.NewPageRankCounter(opts.Pages.Max)

	se.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Cookie", fmt.Sprintf("_EDGE_CD=%s", localeCookieString(opts.Locale)))
		r.Headers.Add("Cookie", fmt.Sprintf("_EDGE_S=%s", localeAltCookieString(opts.Locale)))
	})

	se.OnHTML(imgDompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		var jsonMetadata imgJsonMetadata
		metadataS, metadataExists := dom.Find(imgDompaths.Metadata.Path).Attr(imgDompaths.Metadata.Attr)
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

		titleText := strings.TrimSpace(dom.Find(imgDompaths.Title).Text())
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
		imgFormatS := strings.TrimSpace(dom.Find(imgDompaths.ImgFormatStr).Text())
		if imgFormatS == "" {
			log.Trace().
				Caller().
				Str("engine", se.Name.String()).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Msg("Couldn't find image format (probably a video)")
			return
		}

		// Extract only the resolution using regex (<digit><char><digit>).
		regex := regexp.MustCompile(`(\d+)[^\d]*(\d+)`)
		match := regex.FindStringSubmatch(imgFormatS)
		if len(match) != 3 {
			log.Error().
				Caller().
				Str("engine", se.Name.String()).
				Strs("match", match).
				Str("imgFormatS", imgFormatS).
				Str("jsonMetadata", metadataS).
				Msg("Failed to extract image format")
			return
		}
		origHS, origWS := match[1], match[2]
		log.Trace().
			Caller().
			Str("engine", se.Name.String()).
			Str("imgFormatS", imgFormatS).
			Str("height", origHS).
			Str("width", origWS).
			Msg("Extracted image format")

		// Convert the height to integer.
		origH, err := strconv.Atoi(origHS)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("engine", se.Name.String()).
				Str("height", origHS).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Str("imgFormatS", imgFormatS).
				Msg("Failed to convert original height to int")
			return
		}

		// Convert the width to integer.
		origW, err := strconv.Atoi(origWS)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Str("engine", se.Name.String()).
				Str("width", origWS).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Str("imgFormatS", imgFormatS).
				Msg("Failed to convert original width to int")
			return
		}

		found := false
		var thmbHS, thmbWS string
		for _, thmb := range imgDompaths.Thumbnail {
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

		source := strings.TrimSpace(dom.Find(imgDompaths.Source).Text())
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
				Str("url", jsonMetadata.ImageURL).
				Str("title", titleText).
				Str("desc", jsonMetadata.Desc).
				Int("page", page).
				Int("rank", pageRankCounter.GetPlusOne(pageIndex)).
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

	for i := range opts.Pages.Max {
		pageNum0 := i + opts.Pages.Start
		ctx := colly.NewContext()
		ctx.Put("page", strconv.Itoa(i))

		// Build the parameters.
		params := moreurls.NewParams(
			paramQueryK, query,
			imgParamAsyncK, imgParamAsyncV,
			paramPageK, strconv.Itoa(pageNum0*35+1),
			imgParamCountK, imgParamCountV,
		)

		// Build the url.
		urll := moreurls.Build(imageSearchURL, params)

		// Build anonymous url, by anonymizing the query.
		params.Set(paramQueryK, anonymize.String(query))
		anonUrll := moreurls.Build(imageSearchURL, params)

		// Send the request.
		if err := se.Get(ctx, urll, anonUrll); err != nil {
			retErrors = append(retErrors, err)
		}
	}

	se.Wait()
	close(resChan)
	return retErrors[:len(retErrors):len(retErrors)], foundResults.Load()
}
