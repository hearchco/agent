package currency

import (
	"sync"
)

type CurrencyMap struct {
	currs map[Currency][]float64
	lock  sync.RWMutex
}

func NewCurrencyMap() CurrencyMap {
	return CurrencyMap{
		currs: make(map[Currency][]float64),
	}
}

func (c *CurrencyMap) Append(currs map[Currency]float64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for curr, rate := range currs {
		c.currs[curr] = append(c.currs[curr], rate)
	}
}

func (c *CurrencyMap) Extract() map[Currency]float64 {
	c.lock.RLock()
	defer c.lock.RUnlock()

	avg := make(map[Currency]float64)
	for curr, rates := range c.currs {
		var sum float64
		for _, rate := range rates {
			sum += rate
		}
		avg[curr] = sum / float64(len(rates))
	}
	return avg
}
