package rank

import (
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/structures"
)

func SetRank(result *structures.Result) {
	result.Rank = result.SEPage*100 + result.SEPageRank

	log.Trace().Msgf("Set rank to %v for %v: %v", result.Rank, result.Title, result.URL)
}
