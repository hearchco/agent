package result

import (
	"github.com/rs/zerolog/log"
)

type Suggestion struct {
	suggestionJSON
}

type suggestionJSON struct {
	Value       string       `json:"value"`
	Rank        int          `json:"rank"`
	Score       float64      `json:"score"`
	EngineRanks []RankSimple `json:"engine_ranks"`
}

func (s Suggestion) Value() string {
	return s.suggestionJSON.Value
}

func (s Suggestion) Rank() int {
	return s.suggestionJSON.Rank
}

func (s *Suggestion) SetRank(rank int) {
	s.suggestionJSON.Rank = rank
}

func (s Suggestion) Score() float64 {
	return s.suggestionJSON.Score
}

func (s *Suggestion) SetScore(score float64) {
	s.suggestionJSON.Score = score
}

func (s Suggestion) EngineRanks() []RankSimple {
	if s.suggestionJSON.EngineRanks == nil {
		log.Panic().Msg("EngineRanks is nil")
		// ^PANIC - Assert because the EngineRanks should never be nil.
	}

	return s.suggestionJSON.EngineRanks
}

func (s *Suggestion) ShrinkEngineRanks() {
	if s.suggestionJSON.EngineRanks == nil {
		log.Panic().Msg("EngineRanks is nil")
		// ^PANIC - Assert because the EngineRanks should never be nil.
	}

	ranksLen := len(s.suggestionJSON.EngineRanks)
	s.suggestionJSON.EngineRanks = s.suggestionJSON.EngineRanks[:ranksLen:ranksLen]
}

func (s *Suggestion) AppendEngineRanks(rank RankSimple) {
	if s.suggestionJSON.EngineRanks == nil {
		log.Panic().Msg("EngineRanks is nil")
		// ^PANIC - Assert because the EngineRanks should never be nil.
	}

	s.suggestionJSON.EngineRanks = append(s.suggestionJSON.EngineRanks, rank)
}
