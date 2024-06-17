package scraper

type DOMPaths struct {
	ResultsContainer string
	Result           string
	URL              string
	Title            string
	Description      string
}

type DOMPathsImages struct {
	DOMPaths

	OriginalSize struct {
		Height string
		Width  string
	}
	ThumbnailSize struct {
		Height string
		Width  string
	}
	ThumbnailURL string
	SourceName   string
	SourceURL    string
}
