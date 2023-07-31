package rank

import (
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/structures"
)

// only one call for setResultResponse for a page will ever be performed, so there will never be a read and a write to result.Response at the same time
// thus not locking relay.ResultMap[result.URL] in this function, because the result parameter is a copy is safe. It may, however, read memory that the map sees,
// at the same time that the map is being written to, while this should be fine in this use-case, go may throw an error.
// TLDR: you must mutex.Lock when changing *rankAddr, you probably dont need to mutex.RLock() when reading result
// (in reality even *rankAddr shouldnt need a lock, but go would definately complain about simultanious read/write because of it)
func SetRank(result *structures.Result, rankAddr *int, mutex *sync.RWMutex) {

	mutex.RLock()
	reqUrl := result.Response.Request.URL.String() //dummy code
	mutex.RUnlock()

	if reqUrl != result.URL { //dummy code
		log.Trace().Msg("Request URL not same as result.URL \\/")
	}

	rrank := result.SEPage*100 + result.SEPageRank

	mutex.Lock()
	*rankAddr = rrank
	mutex.Unlock()

	log.Trace().Msgf("Set rank to %v for %v: %v", rrank, result.Title, result.URL)
}
