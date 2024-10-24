package result

type Result interface {
	Key() string
	URL() string
	FQDN() string
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

type ConcMapper[T any] interface {
	AddOrUpgrade(T)
}
