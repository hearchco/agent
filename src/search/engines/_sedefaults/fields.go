package _sedefaults

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/parse"
	"github.com/rs/zerolog/log"
)

// Fetches from DOM via dompaths. Returns (url, title, description)
func RawFieldsFromDOM(dom *goquery.Selection, dompaths *engines.DOMPaths, seName engines.Name) (string, string, string) {
	descText := dom.Find(dompaths.Description).Text()
	titleDom := dom.Find(dompaths.Title)
	titleText := titleDom.Text()

	// Title and Link selector are often the same, utilize this
	var linkDom *goquery.Selection
	if dompaths.Link == dompaths.Result {
		linkDom = titleDom
	} else {
		linkDom = dom.Find(dompaths.Link)
	}

	linkText, hrefExists := linkDom.Attr("href")

	if !hrefExists {
		log.Error().
			Str("engine", seName.String()).
			Str("url", linkText).
			Str("title", titleText).
			Str("description", descText).
			Msgf("_sedefaults.RawFieldsFromDOM(): href attribute doesn't exist on matched URL element (%v)", dompaths.Link)

		return "", "", ""
	}

	return linkText, titleText, descText
}

// Fetches from DOM via dompaths and sanitizes. Returns (url, title, description)
func FieldsFromDOM(dom *goquery.Selection, dompaths *engines.DOMPaths, seName engines.Name) (string, string, string) {
	return SanitizeFields(RawFieldsFromDOM(dom, dompaths, seName))
}

func SanitizeURL(urlText string) string {
	return parse.ParseURL(urlText)
}

func SanitizeTitle(titleText string) string {
	return parse.ParseTextWithHTML(strings.TrimSpace(titleText))
}

func SanitizeDescription(descText string) string {
	return parse.ParseTextWithHTML(strings.TrimSpace(descText))
}

func SanitizeFields(linkText string, titleText string, descText string) (string, string, string) {
	return SanitizeURL(linkText), SanitizeTitle(titleText), SanitizeDescription(descText)
}
