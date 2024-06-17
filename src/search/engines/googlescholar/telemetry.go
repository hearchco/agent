package googlescholar

import (
	"net/url"
)

// Remove seemingly unused params in query.
func removeTelemetry(link string) (string, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return link, err
	}

	q := parsedURL.Query()
	for _, key := range []string{"dq", "lr", "oi", "ots", "sig"} {
		q.Del(key)
	}
	parsedURL.RawQuery = q.Encode()

	return parsedURL.String(), nil
}
