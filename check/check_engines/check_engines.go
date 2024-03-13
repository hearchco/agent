package check_engines

import (
	"fmt"

	"github.com/hearchco/hearchco/src/search/engines"
)

func Check() {
	fmt.Println("Checking engines:")

	names := engines.Names()
	prettyNames := engines.PrettyNames()
	if len(names) != len(prettyNames) {
		panic("PrettyNames and _NameValues have different lengths")
	}

	for i := range names {
		fmt.Printf("\t%s: %s\n", names[i], prettyNames[i])
	}
}
