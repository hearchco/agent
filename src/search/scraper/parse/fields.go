package parse

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/scraper"
)

// Fetches from DOM via dompaths. Returns url, title and description.
func RawFieldsFromDOM(dom *goquery.Selection, dompaths scraper.DOMPaths, seName engines.Name) (string, string, string) {
	descText := dom.Find(dompaths.Description).Text()
	titleDom := dom.Find(dompaths.Title)
	titleText := titleDom.Text()

	// Title and URL selector are often the same.
	var linkDom *goquery.Selection
	if dompaths.URL == dompaths.Result {
		linkDom = titleDom
	} else {
		linkDom = dom.Find(dompaths.URL)
	}

	linkText, hrefExists := linkDom.Attr("href")
	if !hrefExists {
		log.Error().
			Caller().
			Str("engine", seName.String()).
			Str("url", linkText).
			Str("title", titleText).
			Str("description", descText).
			Msgf("Href attribute doesn't exist on matched URL element (%v)", dompaths.URL)

		return "", "", ""
	}

	return linkText, titleText, descText
}

// Fetches from DOM via dompaths and sanitizes. Returns url, title and description.
func FieldsFromDOM(dom *goquery.Selection, dompaths scraper.DOMPaths, seName engines.Name) (string, string, string) {
	return SanitizeFields(RawFieldsFromDOM(dom, dompaths, seName))
}

func SanitizeURL(urlText string) string {
	return ParseURL(urlText)
}

func SanitizeTitle(titleText string) string {
	return ParseTextWithHTML(strings.TrimSpace(titleText))
}

func SanitizeDescription(descText string) string {
	return ParseTextWithHTML(strings.TrimSpace(descText))
}

func SanitizeFields(linkText string, titleText string, descText string) (string, string, string) {
	return SanitizeURL(linkText), SanitizeTitle(titleText), SanitizeDescription(descText)
}
