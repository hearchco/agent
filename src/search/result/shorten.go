package result

// Changes the title and description of the result to be at most N and M characters long respectively.
func (r General) Shorten(maxTitleLength int, maxDescriptionLength int) Result {
	return &General{
		generalJSON{
			URL:         r.URL(),
			Title:       shortString(r.Title(), maxTitleLength),
			Description: shortString(r.Description(), maxDescriptionLength),
			Rank:        r.Rank(),
			Score:       r.Score(),
			EngineRanks: r.EngineRanks(),
		},
	}
}

func (r Images) Shorten(maxTitleLength int, maxDescriptionLength int) Result {
	return &Images{
		imagesJSON{
			General{
				generalJSON{
					URL:         r.URL(),
					Title:       shortString(r.Title(), maxTitleLength),
					Description: shortString(r.Description(), maxDescriptionLength),
					Rank:        r.Rank(),
					Score:       r.Score(),
					EngineRanks: r.EngineRanks(),
				},
			},
			r.OriginalSize(),
			r.ThumbnailSize(),
			r.ThumbnailURL(),
			r.SourceName(),
			r.SourceURL(),
		},
	}
}

func shortString(s string, n int) string {
	if n < 0 {
		return s
	}

	suffix := "..."
	if n-len(suffix) <= 0 {
		suffix = "" // No room for suffix.
	}

	if len(s) > n {
		short := firstNchars(s, n-len(suffix))
		return short + suffix
	}

	return s
}

func firstNchars(str string, n int) string {
	v := []rune(str)
	if n < 0 || n >= len(v) {
		return str
	}
	return string(v[:n])
}
