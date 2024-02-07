package bucket

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/rs/zerolog/log"
)

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
