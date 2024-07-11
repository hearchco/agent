package exchange

import (
	"github.com/hearchco/agent/src/exchange/engines"
	"github.com/hearchco/agent/src/exchange/engines/currencyapi"
	"github.com/hearchco/agent/src/exchange/engines/exchangerateapi"
	"github.com/hearchco/agent/src/exchange/engines/frankfurter"
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the enginer command to generate them again.
	var x [1]struct{}
	_ = x[engines.UNDEFINED-(0)]
	_ = x[engines.CURRENCYAPI-(1)]
	_ = x[engines.EXCHANGERATEAPI-(2)]
	_ = x[engines.FRANKFURTER-(3)]
}

const enginerLen = 4

func exchangerArray() [enginerLen]engines.Exchanger {
	var engineArray [enginerLen]engines.Exchanger
	engineArray[engines.CURRENCYAPI] = currencyapi.New()
	engineArray[engines.EXCHANGERATEAPI] = exchangerateapi.New()
	engineArray[engines.FRANKFURTER] = frankfurter.New()
	return engineArray
}
