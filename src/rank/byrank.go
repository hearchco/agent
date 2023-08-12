package rank

import "github.com/tminaorg/brzaguza/src/bucket/result"

type ByRank []result.Result

func (r ByRank) Len() int           { return len(r) }
func (r ByRank) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByRank) Less(i, j int) bool { return r[i].Rank < r[j].Rank }
