package bucket

import (
	"sync"

	"github.com/hearchco/hearchco/src/search/result"
)

type Relay struct {
	ResultMap map[string]*result.Result
	Mutex     sync.RWMutex
}
