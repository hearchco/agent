package bingimages

import (
	"context"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/engines/_sedefaults"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

func Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, &options, &settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(options, settings, timings)

	_sedefaults.PagesColRequest(ctx, Info.Name, pagesCol)
	_sedefaults.PagesColError(Info.Name, pagesCol)
	_sedefaults.PagesColResponse(Info.Name, pagesCol, relay)

	_sedefaults.ColRequest(ctx, Info.Name, col)
	_sedefaults.ColError(Info.Name, col)

	var pageRankCounter []int = make([]int, options.MaxPages*Info.ResultsPerPage)

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		var jsonMetadata JsonMetadata
		metadataS, metadataExists := dom.Find(dompaths.Metadata).Attr("m")
		if !metadataExists {
			log.Trace().
				Str("engine", Info.Name.String()).
				Msg("Matched result, but couldn't retrieve data")
			return // result useless without URL (metadata)
		}

		err := json.Unmarshal([]byte(metadataS), &jsonMetadata)
		if err != nil {
			log.Error().
				Err(err).
				Msg("bingimages.Search() -> onHTML: failed to unmarshal metadata")
			return // result useless without URL (metadata)
		} else if jsonMetadata.Purl == "" || jsonMetadata.Murl == "" || jsonMetadata.Turl == "" {
			log.Error().
				Str("engine", Info.Name.String()).
				Str("jsonMetadata", metadataS).
				Str("url", jsonMetadata.Purl).
				Str("original", jsonMetadata.Murl).
				Str("thumbnail", jsonMetadata.Turl).
				Msg("bingimages.Search() -> onHTML: Couldn't find image URL")
			return // result useless without URL (metadata)
		}

		titleText := strings.TrimSpace(dom.Find(dompaths.Title).Text())
		if titleText == "" {
			log.Error().
				Str("engine", Info.Name.String()).
				Str("jsonMetadata", metadataS).
				Msg("bingimages.Search() -> onHTML: Couldn't find title")
		}

		// this returns "2000 x 1500 · jpeg"
		imgFormatS := strings.TrimSpace(dom.Find(dompaths.ImgFormatStr).Text())
		var imgH, imgW int
		if imgFormatS == "" {
			// TODO: this happens when bingimages returns youtube link, should it be Trace or do we skip youtube links as images?
			log.Error().
				Str("engine", Info.Name.String()).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Msg("bingimages.Search() -> onHTML: Couldn't find image format")
		} else {
			// convert to "2000x1500·jpeg"
			imgFormatS = strings.ReplaceAll(imgFormatS, " ", "")
			// remove everything after 2000x1500
			imgFormatS = strings.Split(imgFormatS, "·")[0]
			// create height and width
			imgFormat := strings.Split(imgFormatS, "x")

			var err error
			imgH, err = strconv.Atoi(imgFormat[0])
			if err != nil {
				imgH = 0
				log.Error().
					Err(err).
					Str("engine", Info.Name.String()).
					Str("height", imgFormat[0]).
					Str("jsonMetadata", metadataS).
					Str("title", titleText).
					Str("imgFormatS", imgFormatS).
					Msg("bingimages.Search() -> onHTML: Failed to convert original height to int")
			}
			imgW, err = strconv.Atoi(imgFormat[1])
			if err != nil {
				imgW = 0
				log.Error().
					Err(err).
					Str("engine", Info.Name.String()).
					Str("width", imgFormat[1]).
					Str("jsonMetadata", metadataS).
					Str("title", titleText).
					Str("imgFormatS", imgFormatS).
					Msg("bingimages.Search() -> onHTML: Failed to convert original width to int")
			}
		}

		thmbHS, thmbHSExists := dom.Find(dompaths.ThumbnailHeight).Attr("height")
		var thmbH int
		if !thmbHSExists {
			log.Error().
				Str("engine", Info.Name.String()).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Msg("bingimages.Search() -> onHTML: Couldn't find thumbnail height")
		} else {
			var err error
			if thmbH, err = strconv.Atoi(thmbHS); err != nil {
				thmbH = 0
				log.Error().
					Err(err).
					Str("engine", Info.Name.String()).
					Str("height", thmbHS).
					Str("jsonMetadata", metadataS).
					Str("title", titleText).
					Msg("bingimages.Search() -> onHTML: Failed to convert thumbnail height to int")
			}
		}

		thmbWS, thmbWSExists := dom.Find(dompaths.ThumbnailWidth).Attr("width")
		var thmbW int
		if !thmbWSExists {
			log.Error().
				Str("engine", Info.Name.String()).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Msg("bingimages.Search() -> onHTML: Couldn't find thumbnail width")
		} else {
			var err error
			thmbW, err = strconv.Atoi(thmbHS)
			if err != nil {
				thmbW = 0
				log.Error().
					Err(err).
					Str("engine", Info.Name.String()).
					Str("width", thmbWS).
					Str("jsonMetadata", metadataS).
					Str("title", titleText).
					Msg("bingimages.Search() -> onHTML: Failed to convert thumbnail width to int")
			}
		}

		source := strings.TrimSpace(dom.Find(dompaths.Source).Text())
		if source == "" {
			log.Error().
				Str("engine", Info.Name.String()).
				Str("jsonMetadata", metadataS).
				Str("title", titleText).
				Msg("bingimages.Search() -> onHTML: Couldn't find source")
		}

		page := _sedefaults.PageFromContext(e.Request.Ctx, Info.Name)

		original := result.Image{
			URL:    jsonMetadata.Murl,
			Height: uint(imgH),
			Width:  uint(imgW),
		}
		thumbnail := result.Image{
			URL:    jsonMetadata.Turl,
			Height: uint(thmbH),
			Width:  uint(thmbW),
		}

		res := bucket.MakeSEImageResult(jsonMetadata.Purl, titleText, jsonMetadata.Desc, source, original, thumbnail, Info.Name, page, pageRankCounter[page]+1)
		bucket.AddSEResult(res, Info.Name, relay, options, pagesCol)
		pageRankCounter[page]++
	})

	errChannel := make(chan error, 1)

	localeParam := getLocale(options)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	urll := Info.URL + query + params[0] + "&first=1" + params[1] + localeParam
	anonUrll := Info.URL + anonymize.String(query) + params[0] + "&first=1" + params[1] + localeParam
	_sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name, errChannel)

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		urll := Info.URL + query + params[0] + "&first=" + strconv.Itoa(i*10+1) + params[1] + localeParam
		anonUrll := Info.URL + anonymize.String(query) + params[0] + "&first=" + strconv.Itoa(i*10+1) + params[1] + localeParam
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

func getLocale(options engines.Options) string {
	spl := strings.SplitN(strings.ToLower(options.Locale), "_", 2)
	return "&setlang=" + spl[0] + "&cc=" + spl[1]
}
