package result

import (
	"github.com/hearchco/agent/src/search/engines"
)

type Result interface {
	Key() string
	URL() string
	Title() string
	Description() string
	SetDescription(string)
	Rank() int
	SetRank(int)
	Score() float64
	SetScore(float64)
	EngineRanks() []Rank
	InitEngineRanks()
	ShrinkEngineRanks()
	AppendEngineRanks(Rank)
	ConvertToOutput(string) ResultOutput
	Shorten(int, int) Result
}

type ResultScraped interface {
	Key() string
	URL() string
	Title() string
	Description() string
	Rank() RankScraped
	Convert(int) Result
}

type ConcReceiver[T Ranker] interface {
	Rank() T
}

type Ranker interface {
	SearchEngine() engines.Name
}

type ConcMapper[T any, V any] interface {
	AddOrUpgrade(V)
	ExtractWithResponders() ([]T, []engines.Name)
}
