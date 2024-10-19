package moreurls

import (
	"net/url"

	"github.com/rs/zerolog/log"
)

// Constructs a URL with the given parameters.
func Build(urll string, params Params) string {
	// Parse the URL.
	u, err := url.Parse(urll)
	if err != nil {
		log.Panic().
			Err(err).
			Str("url", urll).
			Msg("Failed to parse the URL")
		// ^PANIC - Assert correct URL
	}

	// Convert the parameters to encoded RawQuery keeping the order of keys.
	u.RawQuery = params.QueryEscape()

	return u.String()
}
