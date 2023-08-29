package bucket

import (
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/bucket/result"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/engines"
	"github.com/tminaorg/brzaguza/src/rank"
)

type Relay struct {
	ResultMap         map[string]*result.Result
	Mutex             sync.RWMutex
	EngineDoneChannel chan bool
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
			Rank:          -1,
			Title:         seResult.Title,
			Description:   seResult.Description,
			EngineRanks:   engineRanks,
			TimesReturned: 1,
			Response:      nil,
		}

		if config.InsertDefaultRank {
			result.Rank = rank.DefaultRank(seResult.Rank.Rank, seResult.Rank.Page, seResult.Rank.OnPageRank)
		}

		relay.Mutex.Lock()
		relay.ResultMap[result.URL] = &result
		relay.Mutex.Unlock()
	} else {
		relay.Mutex.Lock()
		mapRes.EngineRanks[mapRes.TimesReturned] = seResult.Rank //can go out of bounds if the same results is returned multiple times by the same engine (e.g. swisscows, presearch ("banana death" -> slate result))
		mapRes.TimesReturned++
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
		relay.Mutex.Unlock()
		log.Error().Msgf("URL not in map when adding response! Should not be possible. URL: %v", link)
		return
	}

	mapRes.Response = response

	resCopy := *mapRes
	rankAddr := &(mapRes.Rank)
	relay.Mutex.Unlock()
	rank.SetRank(&resCopy, rankAddr, &(relay.Mutex)) //copy contains pointer to response
}

func MakeSEResult(urll string, title string, description string, searchEngineName engines.Name, seRank int, sePage int, seOnPageRank int) *engines.RetrievedResult {
	ser := engines.RetrievedRank{
		SearchEngine: searchEngineName,
		Rank:         seRank,
		Page:         sePage,
		OnPageRank:   seOnPageRank,
	}
	res := engines.RetrievedResult{
		URL:         urll,
		Title:       title,
		Description: description,
		Rank:        ser,
	}
	return &res
}
