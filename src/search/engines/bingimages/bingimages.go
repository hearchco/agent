package bingimages

import (
	"context"
	"encoding/base64"
	"net/url"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/engines/_sedefaults"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/hearchco/hearchco/src/search/search/parse"
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

	col.OnHTML(dompaths.Result, func(e *colly.HTMLElement) {
		dom := e.DOM

		var jsonMetadata JsonMetadata
		metadataS, metadataExists := dom.Find(dompaths.Metadata).Attr("m")
		if err := json.Unmarshal([]byte(metadataS), &jsonMetadata); err != nil {
			log.Error().
				Err(err).
				Msg("bingimages.Search() -> onHTML: failed to unmarshal metadata")
		}

		titleText := strings.TrimSpace(dom.Find(dompaths.Title).Text())

		// this returns "2000 x 1500 · jpeg"
		imgFormatS := strings.TrimSpace(dom.Find(dompaths.ImgFormatStr).Text())
		var imgH, imgW int
		if imgFormatS == "" {
			log.Error().
				Str("engine", Info.Name.String()).
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
					Msg("bingimages.Search() -> onHTML: Failed to convert original height to int")
			}
			imgW, err = strconv.Atoi(imgFormat[1])
			if err != nil {
				imgW = 0
				log.Error().
					Err(err).
					Str("engine", Info.Name.String()).
					Str("width", imgFormat[1]).
					Msg("bingimages.Search() -> onHTML: Failed to convert original width to int")
			}
		}

		thmbHS, thmbHSExists := dom.Find(dompaths.ThumbnailHeight).Attr("height")
		var thmbH int
		if thmbHSExists {
			var err error
			if thmbH, err = strconv.Atoi(thmbHS); err != nil {
				thmbH = 0
				log.Error().
					Err(err).
					Str("engine", Info.Name.String()).
					Str("height", thmbHS).
					Msg("bingimages.Search() -> onHTML: Failed to convert thumbnail height to int")
			}
		} else {
			log.Error().
				Str("engine", Info.Name.String()).
				Msg("bingimages.Search() -> onHTML: Couldn't find thumbnail height")
		}

		thmbWS, thmbWSExists := dom.Find(dompaths.ThumbnailWidth).Attr("width")
		var thmbW int
		if thmbWSExists {
			var err error
			thmbW, err = strconv.Atoi(thmbHS)
			if err != nil {
				thmbW = 0
				log.Error().
					Err(err).
					Str("engine", Info.Name.String()).
					Str("width", thmbWS).
					Msg("bingimages.Search() -> onHTML: Failed to convert thumbnail width to int")
			}
		} else {
			log.Error().
				Str("engine", Info.Name.String()).
				Msg("bingimages.Search() -> onHTML: Couldn't find thumbnail width")
		}

		source := strings.TrimSpace(dom.Find(dompaths.Source).Text())

		if metadataExists && titleText != "" && source != "" {
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
			bucket.AddSEResult(res, Info.Name, relay, &options, pagesCol)
			pageRankCounter[page]++
		} else {
			log.Trace().
				Str("engine", Info.Name.String()).
				Str("url", jsonMetadata.Purl).
				Str("title", titleText).
				Str("description", jsonMetadata.Desc).
				Msg("Matched result, but couldn't retrieve data")
		}
	})

	localeParam := getLocale(&options)

	colCtx := colly.NewContext()
	colCtx.Put("page", strconv.Itoa(1))

	_sedefaults.DoGetRequest(Info.URL+query+params[0]+"&first=1"+params[1]+localeParam, colCtx, col, Info.Name, &retError)

	for i := 1; i < options.MaxPages; i++ {
		colCtx = colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i+1))

		_sedefaults.DoGetRequest(Info.URL+query+params[0]+"&first="+strconv.Itoa(i*10+1)+params[1]+localeParam, colCtx, col, Info.Name, &retError)
	}

	col.Wait()
	pagesCol.Wait()

	return retError
}

func removeTelemetry(link string) string {
	if strings.HasPrefix(link, "https://www.bingimages.com/ck/a?") {
		parsedUrl, err := url.Parse(link)
		if err != nil {
			log.Error().Err(err).Str("url", link).Msg("bingimages.removeTelemetry(): error parsing url")
			return ""
		}

		// get the first value of u parameter and remove "a1" in front
		encodedUrl := parsedUrl.Query().Get("u")[2:]

		cleanUrl, err := base64.RawURLEncoding.DecodeString(encodedUrl)
		if err != nil {
			log.Error().Err(err).Msg("bingimages.removeTelemetry(): failed decoding string from base64")
		}
		return parse.ParseURL(string(cleanUrl))
	}
	return link
}

func getLocale(options *engines.Options) string {
	spl := strings.SplitN(strings.ToLower(options.Locale), "_", 2)
	return "&setlang=" + spl[0] + "&cc=" + spl[1]
}
