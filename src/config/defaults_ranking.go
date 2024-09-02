package config

import (
	"github.com/hearchco/agent/src/search/engines"
)

func initCategoryRanking() CategoryRanking {
	return CategoryRanking{
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
}

func EmptyRanking(engs []engines.Name) CategoryRanking {
	rnk := initCategoryRanking()

	for _, eng := range engs {
		rnk.Engines[eng.String()] = CategoryEngineRanking{
			Mul:   1,
			Const: 0,
		}
	}

	return rnk
}

func ReqPrefOthRanking(req []engines.Name, pref []engines.Name, oth []engines.Name) CategoryRanking {
	rnk := initCategoryRanking()

	// First set the least important engines
	for _, eng := range oth {
		rnk.Engines[eng.String()] = CategoryEngineRanking{
			Mul:   1,
			Const: 0,
		}
	}

	// Afterwards overwrite with the preferred engines
	for _, eng := range pref {
		rnk.Engines[eng.String()] = CategoryEngineRanking{
			Mul:   1.25,
			Const: 0,
		}
	}

	// Finally overwrite with the required engines
	for _, eng := range req {
		rnk.Engines[eng.String()] = CategoryEngineRanking{
			Mul:   1.5,
			Const: 0,
		}
	}

	return rnk
}
