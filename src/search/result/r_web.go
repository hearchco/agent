package result

import (
	"time"

	"github.com/hearchco/agent/src/utils/anonymize"
	"github.com/rs/zerolog/log"
)

type Web struct {
	webJSON
}

type webJSON struct {
	URL         string  `json:"url"`
	FQDN        string  `json:"fqdn"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Rank        int     `json:"rank"`
	Score       float64 `json:"score"`
	EngineRanks []Rank  `json:"engine_ranks"`
}

func (r Web) Key() string {
	return r.URL()
}

func (r Web) URL() string {
	if r.webJSON.URL == "" {
		log.Panic().Msg("URL is empty")
		// ^PANIC - Assert because the URL should never be empty.
	}

	return r.webJSON.URL
}

func (r Web) FQDN() string {
	if r.webJSON.FQDN == "" {
		log.Panic().Msg("FQDN is empty")
		// ^PANIC - Assert because the FQDN should never be empty.
	}

	return r.webJSON.FQDN
}

func (r Web) Title() string {
	if r.webJSON.Title == "" {
		log.Panic().Msg("Title is empty")
		// ^PANIC - Assert because the Title should never be empty.
	}

	return r.webJSON.Title
}

func (r Web) Description() string {
	return r.webJSON.Description
}

func (r *Web) SetDescription(desc string) {
	r.webJSON.Description = desc
}

func (r Web) Rank() int {
	return r.webJSON.Rank
}

func (r *Web) SetRank(rank int) {
	r.webJSON.Rank = rank
}

func (r Web) Score() float64 {
	return r.webJSON.Score
}

func (r *Web) SetScore(score float64) {
	r.webJSON.Score = score
}

func (r Web) EngineRanks() []Rank {
	if r.webJSON.EngineRanks == nil {
		log.Panic().Msg("EngineRanks is nil")
		// ^PANIC - Assert because the EngineRanks should never be nil.
	}

	return r.webJSON.EngineRanks
}

func (r *Web) InitEngineRanks() {
	r.webJSON.EngineRanks = make([]Rank, 0)
}

func (r *Web) ShrinkEngineRanks() {
	if r.webJSON.EngineRanks == nil {
		log.Panic().Msg("EngineRanks is nil")
		// ^PANIC - Assert because the EngineRanks should never be nil.
	}

	ranksLen := len(r.webJSON.EngineRanks)
	r.webJSON.EngineRanks = r.webJSON.EngineRanks[:ranksLen:ranksLen]
}

func (r *Web) AppendEngineRanks(rank Rank) {
	if r.webJSON.EngineRanks == nil {
		log.Panic().Msg("EngineRanks is nil")
		// ^PANIC - Assert because the EngineRanks should never be nil.
	}

	r.webJSON.EngineRanks = append(r.webJSON.EngineRanks, rank)
}

func (r Web) ConvertToOutput(secret string) ResultOutput {
	fqdnHash, fqdnTimestamp := anonymize.CalculateHMACBase64(r.FQDN(), secret, time.Now())

	return WebOutput{
		webOutputJSON{
			r,
			fqdnHash,
			fqdnTimestamp,
		},
	}
}
