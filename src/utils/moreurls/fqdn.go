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

	// Check if the hostname is empty.
	if u.Hostname() == "" {
		log.Panic().
			Str("url", urll).
			Msg("Hostname is empty")
		// ^PANIC - Assert non-empty URL.
	}

	return u.Hostname()
}
