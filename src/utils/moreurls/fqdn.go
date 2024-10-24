package moreurls

import (
	"net/url"

	"github.com/rs/zerolog/log"
)

// Returns the fully qualified domain name of the URL.
func FQDN(urll string) string {
	// Parse the URL.
	u, err := url.Parse(urll)
	if err != nil {
		log.Panic().
			Err(err).
			Str("url", urll).
			Msg("Failed to parse the URL")
		// ^PANIC - Assert correct URL.
	}

	return u.Hostname()
}
