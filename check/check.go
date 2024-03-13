package main

import (
	"fmt"

	"github.com/hearchco/hearchco/check/check_engines"
)

func main() {
	fmt.Println("Checking packages:")
	check_engines.Check()
}
