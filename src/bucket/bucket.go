package bucket

import (
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket/result"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
)

type Relay struct {
	ResultMap map[string]*result.Result
	Mutex     sync.RWMutex
}

func AddSEResult(seResult *engines.RetrievedResult, seName engines.Name, relay *Relay, options *engines.Options, pagesCol *colly.Collector) {
	log.Trace().Msgf("%v: Got Result -> %v: %v", seName, seResult.Title, seResult.URL)

	relay.Mutex.RLock()
	mapRes, exists := relay.ResultMap[seResult.URL]
	relay.Mutex.RUnlock()

	if !exists {
		engineRanks := make([]engines.RetrievedRank, len(config.EnabledEngines))
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
		pagesCol.Visit(seResult.URL)
	}
}

func SetResultResponse(link string, response *colly.Response, relay *Relay, seName engines.Name) {
	log.Trace().Msgf("%v: Got Response -> %v", seName, link)

	relay.Mutex.Lock()
	mapRes, exists := relay.ResultMap[link]

	if !exists {
		log.Error().Msgf("URL not in map when adding response! Should not be possible. URL: %v", link)
	} else {
		mapRes.Response = response
	}

	relay.Mutex.Unlock()
}

func MakeSEResult(urll string, title string, description string, searchEngineName engines.Name, sePage int, seOnPageRank int) *engines.RetrievedResult {
	ser := engines.RetrievedRank{
		SearchEngine: searchEngineName,
		Rank:         0,
		Page:         uint(sePage),
		OnPageRank:   uint(seOnPageRank),
	}
	res := engines.RetrievedResult{
		URL:         urll,
		Title:       title,
		Description: description,
		Rank:        ser,
	}
	return &res
}
