package result

import (
	"github.com/rs/zerolog/log"
)

type Suggestion struct {
	suggestionJSON
}

type suggestionJSON struct {
	Value       string       `json:"value"`
	EngineRanks []RankSimple `json:"engine_ranks"`
}

func (s Suggestion) Value() string {
	return s.suggestionJSON.Value
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
