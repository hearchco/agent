package rank

import (
	"github.com/rs/zerolog/log"
	"github.com/hearchco/hearchco/src/bucket/result"
)

type ByScore []result.Result

func (r ByScore) Len() int           { return len(r) }
func (r ByScore) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByScore) Less(i, j int) bool { return r[i].Score > r[j].Score }

type ByRetrievedRank []RankFiller

func (r ByRetrievedRank) Len() int      { return len(r) }
func (r ByRetrievedRank) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r ByRetrievedRank) Less(i, j int) bool {
	if r[i].RetRank.Page != r[j].RetRank.Page {
		return r[i].RetRank.Page < r[j].RetRank.Page
	}
	if r[i].RetRank.OnPageRank != r[j].RetRank.OnPageRank {
		return r[i].RetRank.OnPageRank < r[j].RetRank.OnPageRank
	}

	log.Error().Msgf("failed at ranking: %v, %v", r[i], r[j])
	return true
}
