package result

import (
	"github.com/rs/zerolog/log"
)

type GeneralScraped struct {
	url         string
	title       string
	description string
	rank        RankScraped
}

func (r GeneralScraped) URL() string {
	if r.url == "" {
		log.Panic().Msg("url is empty")
		// ^PANIC - Assert because the url should never be empty.
	}

	return r.url
}

func (r GeneralScraped) Title() string {
	if r.title == "" {
		log.Panic().Msg("title is empty")
		// ^PANIC - Assert because the title should never be empty.
	}

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
