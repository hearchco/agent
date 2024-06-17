package rank

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/hearchco/agent/src/search/result"
)

type ByScore []result.Result

func (r ByScore) Len() int           { return len(r) }
func (r ByScore) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByScore) Less(i, j int) bool { return r[i].Score() > r[j].Score() }

type ByPageAndOnPageRank []*result.Rank

func (r ByPageAndOnPageRank) Len() int      { return len(r) }
func (r ByPageAndOnPageRank) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r ByPageAndOnPageRank) Less(i, j int) bool {
	if r[i].Page() != r[j].Page() {
		return r[i].Page() < r[j].Page()
	}

	if r[i].OnPageRank() != r[j].OnPageRank() {
		return r[i].OnPageRank() < r[j].OnPageRank()
	}

	log.Panic().
		Caller().
		Str("comparableA", fmt.Sprintf("%v", r[i])).
		Str("comparableB", fmt.Sprintf("%v", r[j])).
		Msg("Failed at ranking: same page and onpagerank")
	// ^PANIC

	panic("Failed at ranking: same page and onpagerank")
}
