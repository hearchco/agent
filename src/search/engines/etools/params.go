package etools

const (
	// Variable params.
	paramQueryK      = "query"
	paramPageK       = "page"
	paramSafeSearchK = "safeSearch" // Can be "true" or "false".

	// Constant params.
	paramCountryK, paramCountryV   = "country", "web"
	paramLanguageK, paramLanguageV = "language", "all"
)

func safeSearchValue(safesearch bool) string {
	if safesearch {
		return "true"
	} else {
		return "false"
	}
}
