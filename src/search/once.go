package search

import (
	"sync"
	"sync/atomic"

	"github.com/hearchco/agent/src/search/engines"
)

type onceWrapper struct {
	once    *sync.Once
	errored atomic.Bool
	scraped atomic.Bool
}

func initOnceWrapper(engs []engines.Name) map[engines.Name]*onceWrapper {
	onceWrapMap := make(map[engines.Name]*onceWrapper, len(engs))
	for _, eng := range engs {
		onceWrapMap[eng] = &onceWrapper{
			once:    &sync.Once{},
			errored: atomic.Bool{},
			scraped: atomic.Bool{},
		}
	}
	return onceWrapMap
}

func (ow *onceWrapper) Do(f func()) {
	ow.once.Do(f)
}

func (ow *onceWrapper) Errored() {
	if !ow.errored.Load() {
		ow.errored.Store(true)
	}
}

func (ow *onceWrapper) Scraped() {
	if !ow.scraped.Load() {
		ow.scraped.Store(true)
	}
}

func (ow *onceWrapper) Success() bool {
	return !ow.errored.Load() && ow.scraped.Load()
}
