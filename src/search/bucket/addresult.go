package bucket

import (
	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/config"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
	"github.com/rs/zerolog/log"
)

// passing pointers down is on stack
func AddSEResult(seResult *result.RetrievedResult, seName engines.Name, relay *Relay, options engines.Options, pagesCol *colly.Collector) {
	log.Trace().
		Str("engine", seName.String()).
		Str("title", seResult.Title).
		Str("url", seResult.URL).
		Msg("Got result")

	relay.Mutex.RLock()
	mapRes, exists := relay.ResultMap[seResult.URL]
	relay.Mutex.RUnlock()

	if !exists {
		// create engine ranks slice with capacity of enabled engines
		engineRanks := make([]result.RetrievedRank, 0, len(config.EnabledEngines))
		engineRanks = append(engineRanks, seResult.Rank)

		result := result.Result{
			URL:         seResult.URL,
			Rank:        0,
			Title:       seResult.Title,
			Description: seResult.Description,
			EngineRanks: engineRanks,
			ImageResult: seResult.ImageResult,
		}

		relay.Mutex.Lock()
		relay.ResultMap[result.URL] = &result
		relay.Mutex.Unlock()
	} else {
		alreadyIn := false
		index := 0

		relay.Mutex.RLock()
		for ind := range mapRes.EngineRanks { // this could also be done by changing EngineRanks to a map
			if seName == mapRes.EngineRanks[ind].SearchEngine {
				alreadyIn = true
				index = ind
				break
			}
		}
		relay.Mutex.RUnlock()

		relay.Mutex.Lock()
		if !alreadyIn {
			mapRes.EngineRanks = append(mapRes.EngineRanks, seResult.Rank)
		} else if mapRes.EngineRanks[index].Page > seResult.Rank.Page {
			mapRes.EngineRanks[index].Page = seResult.Rank.Page
			mapRes.EngineRanks[index].OnPageRank = seResult.Rank.OnPageRank
		} else if mapRes.EngineRanks[index].Page == seResult.Rank.Page && mapRes.EngineRanks[index].OnPageRank > seResult.Rank.OnPageRank {
			mapRes.EngineRanks[index].OnPageRank = seResult.Rank.OnPageRank
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
