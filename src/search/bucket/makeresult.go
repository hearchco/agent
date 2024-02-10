package bucket

import (
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

// Returns nil on invalid data
func MakeSEResult(urll string, title string, description string, searchEngineName engines.Name, sePage int, seOnPageRank int) *result.RetrievedResult {
	if urll == "" || title == "" {
		log.Error().
			Str("engine", searchEngineName.String()).
			Str("url", urll).
			Str("title", title).
			Str("description", description).
			Msg("bucket.MakeSEResult(): invalid result, some fields are empty.")
		return nil
	}

	ser := result.RetrievedRank{
		SearchEngine: searchEngineName,
		Rank:         0,
		Page:         uint(sePage),
		OnPageRank:   uint(seOnPageRank),
	}
	res := result.RetrievedResult{
		URL:         urll,
		Title:       title,
		Description: description,
		Rank:        ser,
	}
	return &res
}
