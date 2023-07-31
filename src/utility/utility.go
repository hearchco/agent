package utility

import (
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

func ParseURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Error().Err(err).Msgf("Couldn't parse URL: %v", rawURL)
		return rawURL
	}
	return parsedURL.String()
}
