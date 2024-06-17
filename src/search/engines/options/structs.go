package options

// User provided options for every search engine.
type Options struct {
	Pages      Pages
	Locale     Locale
	SafeSearch bool
}

// Start must be 0-based index.
// Max must be greater than 0.
type Pages struct {
	Start int
	Max   int
}
