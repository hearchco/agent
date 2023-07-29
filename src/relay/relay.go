package relay

import (
	"github.com/tminaorg/brzaguza/src/structures"
)

var ResultChannel chan structures.Result = make(chan structures.Result)
var RankChannel chan structures.ResultRank = make(chan structures.ResultRank)
var EngineDoneChannel chan bool = make(chan bool)
var ResultMap map[string]*structures.Result
