package bucket

import (
	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/config"
	"github.com/tminaorg/brzaguza/src/rank"
	"github.com/tminaorg/brzaguza/src/structures"
)

func AddSEResult(seResult *structures.SEResult, seName string, relay *structures.Relay, options *structures.Options, pagesCol *colly.Collector) {
	log.Trace().Msgf("%v: Got Result -> %v: %v", seName, seResult.Title, seResult.URL)

	relay.Mutex.RLock()
	mapRes, exists := relay.ResultMap[seResult.URL]
	relay.Mutex.RUnlock()

	if !exists {
		searchEngines := make([]structures.SERank, config.NumberOfEngines)
		searchEngines[0] = seResult.Rank
		result := structures.Result{
			URL:           seResult.URL,
			Rank:          -1,
			Title:         seResult.Title,
			Description:   seResult.Description,
			SearchEngines: searchEngines,
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
		mapRes.SearchEngines[mapRes.TimesReturned] = seResult.Rank
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

func AddSEResultResponse(link string, response *colly.Response, relay *structures.Relay, seName string) {
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

func MakeSEResult(urll string, title string, description string, searchEngineName string, seRank int, sePage int, seOnPageRank int) *structures.SEResult {
	ser := structures.SERank{
		SearchEngine: searchEngineName,
		Rank:         seRank,
		Page:         sePage,
		OnPageRank:   seOnPageRank,
	}
	res := structures.SEResult{
		URL:         urll,
		Title:       title,
		Description: description,
		Rank:        ser,
	}
	return &res
}
