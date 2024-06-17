package result

type GeneralScraped struct {
	url         string
	title       string
	description string
	rank        RankScraped
}

func (r GeneralScraped) URL() string {
	return r.url
}

func (r GeneralScraped) Title() string {
	return r.title
}

func (r GeneralScraped) Description() string {
	return r.description
}

func (r GeneralScraped) Rank() RankScraped {
	return r.rank
}

func (r GeneralScraped) Convert(erCap int) Result {
	engineRanks := make([]Rank, 0, erCap)
	engineRanks = append(engineRanks, r.Rank().Convert())
	return &General{
		generalJSON{
			URL:         r.URL(),
			Title:       r.Title(),
			Description: r.Description(),
			EngineRanks: engineRanks,
		},
	}
}
