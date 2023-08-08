package structures

import (
	"sync"
)

type Relay struct {
	ResultMap         map[string]*Result
	Mutex             sync.RWMutex
	EngineDoneChannel chan bool
}
