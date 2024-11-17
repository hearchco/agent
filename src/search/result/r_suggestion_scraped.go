package result

import (
	"github.com/hearchco/agent/src/search/engines"
)

type SuggestionScraped struct {
	value string
	rank  RankSimpleScraped
}

func NewSuggestionScraped(value string, seName engines.Name, rank int) SuggestionScraped {
	r := NewRankSimpleScraped(seName, rank)
	return SuggestionScraped{
		value,
		r,
	}
}

func (s SuggestionScraped) Key() string {
	return s.Value()
}

func (s SuggestionScraped) Value() string {
	return s.value
}

func (s SuggestionScraped) Rank() RankSimpleScraped {
	return s.rank
}

func (s SuggestionScraped) Convert(erCap int) Suggestion {
	engineRanks := make([]RankSimple, 0, erCap)
	engineRanks = append(engineRanks, s.Rank().Convert())
	return Suggestion{
		suggestionJSON{
			Value:       s.Value(),
			EngineRanks: engineRanks,
		},
	}
}
