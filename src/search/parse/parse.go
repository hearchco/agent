package parse

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

func ParseURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	rawURL, unescErr := url.QueryUnescape(rawURL) // if the url was part of a telemetry link, this will help.
	if unescErr != nil {
		log.Error().Err(unescErr).Msgf("Couldn't unescape URL: %v", rawURL)
		return rawURL
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Error().Err(err).Msgf("Couldn't parse URL: %v", rawURL)
		return rawURL
	}
	return parsedURL.String()
}

func ParseTextWithHTML(rawHTML string) string {
	var result string = ""
	htmlNode, perr := html.ParseFragment(strings.NewReader(rawHTML), nil)
	if perr != nil {
		log.Error().Err(perr).Msgf("Couldn't utility.ParseTextWithHTML: %v", rawHTML)
		return ""
	}
	for _, el := range htmlNode {
		sel := goquery.NewDocumentFromNode(el)
		result += sel.Text()
	}
	return result
}
