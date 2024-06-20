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

func (s SuggestionScraped) Convert(erCap int) Suggestion {
	engineRanks := make([]RankSimple, 0, erCap)
	engineRanks = append(engineRanks, s.Rank().Convert())
	return Suggestion{
		suggestionJSON{
			s.Value(),
			engineRanks,
		},
	}
}
