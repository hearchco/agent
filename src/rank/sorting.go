package rank

import (
	"fmt"

	"github.com/hearchco/hearchco/src/bucket/result"
	"github.com/rs/zerolog/log"
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

	log.Error().
		Str("comparableA", fmt.Sprintf("%v", r[i])).
		Str("comparableB", fmt.Sprintf("%v", r[j])).
		Msg("rank.(r ByRetrievedRank)Less(): failed at ranking")
	return true
}
