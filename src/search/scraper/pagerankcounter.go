package scraper

import (
	"sync/atomic"
)

// A goroutine-safe counter for PageRank.
type PageRankCounter struct {
	counts []atomic.Int32
}

// Create a new PageRankCounter.
func NewPageRankCounter(pages int) PageRankCounter {
	return PageRankCounter{counts: make([]atomic.Int32, pages)}
}

// Increment the count for a page.
func (prc *PageRankCounter) Increment(page int) {
	prc.counts[page].Add(1)
}

// Get the count for a page + 1.
func (prc *PageRankCounter) GetPlusOne(page int) int {
	return int(prc.counts[page].Load() + 1)
}
