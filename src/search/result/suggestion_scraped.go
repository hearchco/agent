package result

type SuggestionScraped struct {
	value string
	rank  RankSimpleScraped
}

func (s SuggestionScraped) Value() string {
	return s.value
}

func (s SuggestionScraped) Rank() RankSimpleScraped {
	return s.rank
}

func (r SuggestionScraped) Convert(erCap int) Suggestion {
	engineRanks := make([]RankSimple, 0, erCap)
	engineRanks = append(engineRanks, r.Rank().Convert())
	return Suggestion{
		suggestionJSON{
			r.Value(),
			engineRanks,
		},
	}
}
