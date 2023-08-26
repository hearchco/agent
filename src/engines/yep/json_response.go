package yep

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (m *AMain) UnmarshalJSON(buf []byte) error {
	var status string
	tmp := []interface{}{&status, m}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		log.Error().Msgf("Wrong number of fields in Yep JSON response: %v != %v", g, e)
		return fmt.Errorf("wrong number of fields in Yep JSON response: %v != %v", g, e) //is this the best way to do this?
	}
	return nil
}

type YepResponse []interface{}

type AResult struct {
	URL       string `json:"url"`
	Title     string `json:"title"`
	TType     string `json:"type"`
	Snippet   string `json:"snippet"`
	VisualURL string `json:"visual_url"`
	FirstSeen string `json:"first_seen"`
}

type AMain struct {
	Total   int       `json:"total"`
	Results []AResult `json:"results"`
}
