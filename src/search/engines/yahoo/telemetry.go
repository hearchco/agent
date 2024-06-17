package yahoo

import (
	"net/url"
	"strings"
)

func removeTelemetry(urll string) (string, error) {
	if !strings.Contains(urll, "://r.search.yahoo.com/") {
		return urll, nil
	}

	suff := strings.SplitAfterN(urll, "/RU=http", 2)[1]
	urll = "http" + strings.SplitN(suff, "/RK=", 2)[0]

	newLink, err := url.QueryUnescape(urll)
	if err != nil {
		return "", err
	}

	return newLink, nil
}
