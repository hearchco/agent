package bing

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

func removeTelemetry(urll string) (string, error) {
	if !strings.HasPrefix(urll, "https://www.bing.com/ck/a?") {
		return urll, nil
	}

	parsedUrl, err := url.Parse(urll)
	if err != nil {
		return "", fmt.Errorf("failed parsing URL: %w", err)
	}

	// Get the first value of "u" parameter and remove "a1" from the beginning.
	encodedUrl := parsedUrl.Query().Get("u")[2:]

	cleanUrl, err := base64.RawURLEncoding.DecodeString(encodedUrl)
	if err != nil {
		return "", fmt.Errorf("failed decoding base64: %w", err)
	}

	return string(cleanUrl), nil
}
