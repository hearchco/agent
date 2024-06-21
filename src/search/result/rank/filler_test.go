package rank

import (
	"testing"

	"github.com/hearchco/agent/src/search/engines"
	"github.com/hearchco/agent/src/search/result"
)

type ranksPair struct {
	orig     []result.RankScraped
	expected []result.RankScraped
}

func TestFillEngineRankRank(t *testing.T) {
	// Each elements represents a pair of original and expected engine ranks.
	// The number of elements represents the number of results.
	ranksTests := [...]ranksPair{
		{
			[]result.RankScraped{
				result.NewRankScraped(engines.GOOGLE, 0, 1, 1),
				result.NewRankScraped(engines.BING, 0, 1, 1),
				result.NewRankScraped(engines.MOJEEK, 0, 1, 3),
			},
			[]result.RankScraped{
				result.NewRankScraped(engines.GOOGLE, 1, 1, 1),
				result.NewRankScraped(engines.BING, 1, 1, 1),
				result.NewRankScraped(engines.MOJEEK, 2, 1, 3),
			},
		},
		{
			[]result.RankScraped{
				result.NewRankScraped(engines.MOJEEK, 0, 1, 1),
			},
			[]result.RankScraped{
				result.NewRankScraped(engines.MOJEEK, 1, 1, 1),
			},
		},
		{
			[]result.RankScraped{
				result.NewRankScraped(engines.GOOGLE, 0, 2, 1),
				result.NewRankScraped(engines.BING, 0, 3, 5),
			},
			[]result.RankScraped{
				result.NewRankScraped(engines.GOOGLE, 2, 2, 1),
				result.NewRankScraped(engines.BING, 2, 3, 5),
			},
		},
	}

	// Adding the ranks to the results, afterwards adding the results into the slice of results.
	resultsOrig := make(Results, 0, len(ranksTests))
	resultsExpected := make(Results, 0, len(ranksTests))
	for _, rankPair := range ranksTests {
		var resOrig result.Result = &result.General{}
		var resExpected result.Result = &result.General{}
		resOrig.InitEngineRanks()
		resExpected.InitEngineRanks()

		for _, rank := range rankPair.orig {
			resOrig.AppendEngineRanks(rank.Convert())
		}
		for _, rank := range rankPair.expected {
			resExpected.AppendEngineRanks(rank.Convert())
		}

		resultsOrig = append(resultsOrig, resOrig)
		resultsExpected = append(resultsExpected, resExpected)
	}

	// Creating the tests.
	tests := [...]testPair{
		{
			resultsOrig,
			resultsExpected,
		},
	}

	// Making sure that the tests exist.
	if len(tests) == 0 {
		t.Errorf("Bad tests made: len(tests) == 0")
	}

	// Making sure that the tests are made correctly.
	for _, test := range tests {
		if len(test.orig) != len(test.expected) {
			t.Errorf("Bad tests made: len(tests.orig) != len(tests.expected)")
		}

		for i := range test.orig {
			if len(test.orig[i].EngineRanks()) != len(test.expected[i].EngineRanks()) {
				t.Errorf("Bad tests made: len(tests.orig[%v].EngineRanks) != len(tests.expected[%v].EngineRanks)", i, i)
			}

			for j := range test.orig[i].EngineRanks() {
				if test.orig[i].EngineRanks()[j].SearchEngine() != test.expected[i].EngineRanks()[j].SearchEngine() {
					t.Errorf("Bad tests made: test.orig[%v].EngineRanks[%v].SearchEngine != test.expected[%v].EngineRanks[%v].SearchEngine", i, j, i, j)
				}
			}
		}
	}

	// Running the tests.
	for _, test := range tests {
		test.orig.fillEngineRankRank()

		for i := range test.orig {
			for j := range test.orig[i].EngineRanks() {
				if test.orig[i].EngineRanks()[j].Rank() != test.expected[i].EngineRanks()[j].Rank() {
					t.Errorf("fillEngineRankRank() = %v, want %v", test.orig[i].EngineRanks()[j].Rank(), test.expected[i].EngineRanks()[j].Rank())
				}
			}
		}
	}
}
