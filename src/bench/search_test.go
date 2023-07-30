package bench

import (
	"testing"

	"github.com/tminaorg/brzaguza/src/search"
)

const queryShort = "banana death"
const queryLong = "how long should i wait before marriage"

func BenchmarkSearchShort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		search.PerformSearch(queryShort, 1, false)
	}
}

func BenchmarkSearchLong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		search.PerformSearch(queryLong, 5, false)
	}
}
