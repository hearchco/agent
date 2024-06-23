package config

import (
	"github.com/hearchco/agent/src/search/engines"
)

func EmptyRanking(engs []engines.Name) CategoryRanking {
	rnk := CategoryRanking{
		REXP:    0.5,
		A:       1,
		B:       0,
		C:       1,
		D:       0,
		TRA:     1,
		TRB:     0,
		TRC:     1,
		TRD:     0,
		Engines: map[string]CategoryEngineRanking{},
	}

	for _, eng := range engs {
		rnk.Engines[eng.String()] = CategoryEngineRanking{
			Mul:   1,
			Const: 0,
		}
	}

	return rnk
}
