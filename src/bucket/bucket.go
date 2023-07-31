package bucket

import (
	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/rank"
	"github.com/tminaorg/brzaguza/src/structures"
)

func SetResult(result *structures.Result, relay *structures.Relay, options *structures.Options, pagesCol *colly.Collector) {
	log.Trace().Msgf("%v: Got Result -> %v: %v", result.SearchEngine, result.Title, result.URL)

	relay.Mutex.Lock()
	mapRes, exists := relay.ResultMap[result.URL]

	if !exists {
		relay.ResultMap[result.URL] = result
	} else if len(mapRes.Description) < len(result.Description) {
		mapRes.Description = result.Description
	}
	relay.Mutex.Unlock()

	if !exists && options.VisitPages {
		pagesCol.Visit(result.URL)
	}
}

func SetResultResponse(link string, response *colly.Response, relay *structures.Relay, seName string) {
	log.Trace().Msgf("%v: Got Response -> %v", seName, link)

	relay.Mutex.Lock()
	mapRes, exists := relay.ResultMap[link]

	if !exists {
		log.Error().Msgf("URL not in map when adding response! Should not be possible. URL: %v", link)
		relay.Mutex.Unlock()
		return
	}

	mapRes.Response = response

	resCopy := *mapRes
	rankAddr := &(mapRes.Rank)
	relay.Mutex.Unlock()
	rank.SetRank(&resCopy, rankAddr, &(relay.Mutex)) //copy contains pointer to response
}
