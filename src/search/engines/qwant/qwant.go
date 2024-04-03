package qwant

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/bucket"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/engines/_sedefaults"
	"github.com/rs/zerolog/log"
)

type Engine struct{}

func New() Engine {
	return Engine{}
}

func (e Engine) Search(ctx context.Context, query string, relay *bucket.Relay, options engines.Options, settings config.Settings, timings config.Timings, salt string, enabledEngines int) []error {
	ctx, err := _sedefaults.Prepare(ctx, Info, Support, &options, &settings)
	if err != nil {
		return []error{err}
	}

	col, pagesCol := _sedefaults.InitializeCollectors(ctx, Info.Name, options, settings, timings, relay)

	col.OnResponse(func(r *colly.Response) {
		var pageStr string = r.Ctx.Get("page")
		if pageStr == "" {
			// If I'm using GET as the first page
			return
		}

		pageIndex := _sedefaults.PageFromContext(r.Request.Ctx, Info.Name)
		page := pageIndex + options.Pages.Start + 1

		var parsedResponse QwantResponse
		err := json.Unmarshal(r.Body, &parsedResponse)
		if err != nil {
			log.Error().
				Err(err).
				Str("engine", Info.Name.String()).
				Bytes("body", r.Body).
				Msg("Failed body unmarshall to json")
		}

		mainline := parsedResponse.Data.Res.Items.Mainline
		counter := 1
		for _, group := range mainline {
			if group.Type != "web" {
				continue
			}
			for _, result := range group.Items {
				goodLink, goodTitle, goodDesc := _sedefaults.SanitizeFields(result.URL, result.Title, result.Description)

				res := bucket.MakeSEResult(goodLink, goodTitle, goodDesc, Info.Name, page, counter)
				valid := bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol, enabledEngines)
				if valid {
					counter += 1
				}
			}
		}
	})

	retErrors := make([]error, 0, options.Pages.Max)

	// static params
	localeParam := getLocale(options)
	deviceParam := getDevice(options)
	safeSearchParam := getSafeSearch(options)
	countParam := "&count=" + strconv.Itoa(settings.RequestedResultsPerPage)

	// starts from at least 0
	for i := options.Pages.Start; i < options.Pages.Start+options.Pages.Max; i++ {
		colCtx := colly.NewContext()
		colCtx.Put("page", strconv.Itoa(i-options.Pages.Start))

		// dynamic params
		offsetParam := "&offset=" + strconv.Itoa(i*settings.RequestedResultsPerPage)

		urll := Info.URL + query + countParam + localeParam + offsetParam + safeSearchParam
		anonUrll := Info.URL + anonymize.String(query) + countParam + localeParam + offsetParam + deviceParam + safeSearchParam

		err := _sedefaults.DoGetRequest(urll, anonUrll, colCtx, col, Info.Name)
		if err != nil {
			retErrors = append(retErrors, err)
		}
	}

	col.Wait()
	pagesCol.Wait()

	return retErrors[:len(retErrors):len(retErrors)]
}

// qwant returns this array when an invalid locale is supplied
var validLocales = [...]string{"bg_bg", "br_fr", "ca_ad", "ca_es", "ca_fr", "co_fr", "cs_cz", "cy_gb", "da_dk", "de_at", "de_ch", "de_de", "ec_ca", "el_gr", "en_au", "en_ca", "en_gb", "en_ie", "en_my", "en_nz", "en_us", "es_ad", "es_ar", "es_cl", "es_co", "es_es", "es_mx", "es_pe", "et_ee", "eu_es", "eu_fr", "fc_ca", "fi_fi", "fr_ad", "fr_be", "fr_ca", "fr_ch", "fr_fr", "gd_gb", "he_il", "hu_hu", "it_ch", "it_it", "ko_kr", "nb_no", "nl_be", "nl_nl", "pl_pl", "pt_ad", "pt_pt", "ro_ro", "sv_se", "th_th", "zh_cn", "zh_hk"}

func getLocale(options engines.Options) string {
	locale := strings.ToLower(options.Locale)
	for _, vl := range validLocales {
		if locale == vl {
			return "&locale=" + locale
		}
	}
	log.Warn().
		Str("locale", options.Locale).
		Strs("validLocales", validLocales[:]).
		Msg("qwant.getLocale(): Invalid qwant locale supplied. Falling back to en_US. Qwant supports these (disregard specific formatting)")
	return "&locale=" + strings.ToLower(config.DefaultLocale)
}

func getDevice(options engines.Options) string {
	if options.Mobile {
		return "&device=mobile"
	}
	return "&device=desktop"
}

func getSafeSearch(options engines.Options) string {
	if options.SafeSearch {
		return "&safesearch=1"
	}
	return "&safesearch=0"
}

/*
col.OnHTML("div[data-testid=\"sectionWeb\"] > div > div", func(e *colly.HTMLElement) {
	//first page
	idx := e.Index

	dom := e.DOM
	baseDOM := dom.Find("div[data-testid=\"webResult\"] > div > div > div > div > div")
	hrefElement := baseDOM.Find("a[data-testid=\"serTitle\"]")
	linkHref, hrefExists := hrefElement.Attr("href")
	linkText := parse.ParseURL(linkHref)
	titleText := strings.TrimSpace(hrefElement.Text())
	descText := strings.TrimSpace(baseDOM.Find("div > span").Text())

	if hrefExists && linkText != "" && linkText != "#" && titleText != "" {
		var pageStr string = e.Request.Ctx.Get("page")
		page, _ := strconv.Atoi(pageStr)

		res := bucket.MakeSEResult(linkText, titleText, descText, Info.Name, -1, page, idx+1)
		bucket.AddSEResult(&res, Info.Name, relay, options, pagesCol, enabledEngines)
	} else {
		log.Info().
			Str("link", linkText).
			Str("title", titleText).
			Str("desc", descText).
			Msg("Not good!")
	}
})
*/
