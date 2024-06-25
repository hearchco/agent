package moreurls

import (
	"fmt"
	"regexp"
)

const (
	faviconVerifierPrefix = "favicon-verify://"
	uriPattern            = "^(http(s?))(://)([^/]+)"
)

// Prepends the URI with the favicon-verify prefix.
func GetURIToVerify(urll string) (string, error) {
	uri, err := GetURI(urll)
	if err != nil {
		return "", err
	} else {
		return faviconVerifierPrefix + uri, nil
	}
}

// Extracts the URI from the URL.
// https://www.example.com/some/path -> https://www.example.com
func GetURI(urll string) (string, error) {
	re := regexp.MustCompile(uriPattern)
	ss := re.FindString(urll)
	if ss == "" {
		return "", fmt.Errorf("failed to extract URI from URL")
	} else {
		return ss, nil
	}
}
