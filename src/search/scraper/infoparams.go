package scraper

import "github.com/hearchco/agent/src/search/engines"

type Info struct {
	Name    engines.Name
	Domain  string
	URL     string
	Origins []engines.Name
}

type Params struct {
	Page       string
	Locale     string
	LocaleSec  string
	SafeSearch string
}
