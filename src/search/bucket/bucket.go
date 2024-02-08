package bucket

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

type Relay struct {
	ResultMap map[string]*result.Result
	Mutex     sync.RWMutex
}

func AddSEResult(seResult *result.RetrievedResult, seName engines.Name, relay *Relay, options *engines.Options, pagesCol *colly.Collector) {
	log.Trace().
		Str("engine", seName.String()).
		Str("title", seResult.Title).
		Str("url", seResult.URL).
		Msg("Got result")

	relay.Mutex.RLock()
	mapRes, exists := relay.ResultMap[seResult.URL]
	relay.Mutex.RUnlock()

	if !exists {
		engineRanks := make([]result.RetrievedRank, len(config.EnabledEngines))
		engineRanks[0] = seResult.Rank
		result := result.Result{
			URL:           seResult.URL,
			Rank:          0,
			Title:         seResult.Title,
			Description:   seResult.Description,
			EngineRanks:   engineRanks,
			TimesReturned: 1,
			Response:      nil,
		}

		relay.Mutex.Lock()
		relay.ResultMap[result.URL] = &result
		relay.Mutex.Unlock()
	} else {
		alreadyIn := false
		relay.Mutex.RLock()
		for ind := range mapRes.EngineRanks { // this could also be done by changing EngineRanks to a map
			if seName == mapRes.EngineRanks[ind].SearchEngine {
				alreadyIn = true
				break
			}
		}
		relay.Mutex.RUnlock()

		relay.Mutex.Lock()
		if !alreadyIn {
			mapRes.EngineRanks[mapRes.TimesReturned] = seResult.Rank
			mapRes.TimesReturned++
		}
		if len(mapRes.Description) < len(seResult.Description) {
			mapRes.Description = seResult.Description
		}
		relay.Mutex.Unlock()
	}

	if !exists && options.VisitPages {
		if err := pagesCol.Visit(seResult.URL); err != nil {
			log.Error().
				Err(err).
				Str("url", seResult.URL).
				Msg("bucket.AddSEResult(): failed visiting")
		}
	}
}

func SetResultResponse(link string, response *colly.Response, relay *Relay, seName engines.Name) error {
	log.Trace().
		Str("engine", seName.String()).
		Str("link", link).
		Msg("Got response")

	relay.Mutex.Lock()
	mapRes, exists := relay.ResultMap[link]

	if !exists {
		relay.Mutex.Unlock()
		relay.Mutex.RLock()
		err := fmt.Errorf("bucket.SetResultResponse(): URL not in map when adding response, should not be possible. URL: %v.\nRelay: %v", link, relay)
		relay.Mutex.RUnlock()
		return err
	} else {
		mapRes.Response = response
		relay.Mutex.Unlock()
	}

	return nil
}

func MakeSEResult(urll string, title string, description string, searchEngineName engines.Name, sePage int, seOnPageRank int) *result.RetrievedResult {
	ser := result.RetrievedRank{
		SearchEngine: searchEngineName,
		Rank:         0,
		Page:         uint(sePage),
		OnPageRank:   uint(seOnPageRank),
	}
	res := result.RetrievedResult{
		URL:         urll,
		Title:       title,
		Description: description,
		Rank:        ser,
	}
	return &res
}
