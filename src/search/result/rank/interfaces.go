package rank

import (
	"github.com/hearchco/agent/src/search/engines"
)

type scoreRanker interface {
	Score() float64
}

type scoreEngineRanker[T ranker] interface {
	scoreRanker

	EngineRanks() []T
}

type ranker interface {
	SearchEngine() engines.Name
	Rank() int
}
