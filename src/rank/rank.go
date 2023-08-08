package rank

import (
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/tminaorg/brzaguza/src/structures"
)

// TLDR: you must mutex.Lock when changing *rankAddr, you probably dont need to mutex.RLock() when reading result
// (in reality even *rankAddr shouldnt need a lock, but go would definately complain about simultanious read/write because of it)
func SetRank(result *structures.Result, rankAddr *int, mutex *sync.RWMutex) {

	//mutex.RLock()
	reqUrl := result.Response.Request.URL.String() //dummy code, if error here, uncomment lock
	//mutex.RUnlock()

	if reqUrl != result.URL { //dummy code
		log.Trace().Msgf("(This is ok) Request URL not same as result.URL \\/ %v | %v", reqUrl, result.URL)
	}

	rrank := result.SearchEngines[0].Page*100 + result.SearchEngines[0].OnPageRank

	mutex.Lock()
	*rankAddr = rrank
	mutex.Unlock()

	log.Trace().Msgf("Set rank to %v for %v: %v", rrank, result.Title, result.URL)
}

func DefaultRank(seRank int, sePage int, seOnPageRank int) int {
	return sePage*100 + seOnPageRank
}
