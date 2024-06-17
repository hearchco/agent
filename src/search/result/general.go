package result

type General struct {
	generalJSON
}

type generalJSON struct {
	URL         string  `json:"url"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Rank        int     `json:"rank"`
	Score       float64 `json:"score"`
	EngineRanks []Rank  `json:"engine_ranks"`
}

func (r General) URL() string {
	return r.generalJSON.URL
}

func (r General) Title() string {
	return r.generalJSON.Title
}

func (r General) Description() string {
	return r.generalJSON.Description
}

func (r *General) SetDescription(desc string) {
	r.generalJSON.Description = desc
}

func (r General) Rank() int {
	return r.generalJSON.Rank
}

func (r *General) SetRank(rank int) {
	r.generalJSON.Rank = rank
}

func (r General) Score() float64 {
	return r.generalJSON.Score
}

func (r *General) SetScore(score float64) {
	r.generalJSON.Score = score
}

func (r General) EngineRanks() []Rank {
	return r.generalJSON.EngineRanks
}

func (r *General) ShrinkEngineRanks() {
	ranksLen := len(r.generalJSON.EngineRanks)
	r.generalJSON.EngineRanks = r.generalJSON.EngineRanks[:ranksLen:ranksLen]
}

func (r *General) AppendEngineRanks(rank Rank) {
	if r.generalJSON.EngineRanks == nil {
		r.generalJSON.EngineRanks = make([]Rank, 0)
	}

	r.generalJSON.EngineRanks = append(r.generalJSON.EngineRanks, rank)
}

func (r General) ConvertToOutput(salt string) ResultOutput {
	return r
}
