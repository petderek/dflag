// +build fizzbuzz

package main

import (
	"fmt"

	"github.com/petderek/dflag"
)

var flags = struct {
	FizzOn int `dflag:"required"`
	BuzzOn int `dflag:"required"`
}{}

// to test this:
// go run -tags=fizzbuzz ./fizzbuzz.go -fizzon 3 -buzzon 5
func main() {
	dflag.Parse(&flags)
	for i := 1; i <= 25; i++ {
		fizzy := i%flags.FizzOn == 0
		buzzy := i%flags.BuzzOn == 0

		fmt.Print(i, " ")
		if fizzy {
			fmt.Print("fizz")
		}
		if buzzy {
			fmt.Print("buzz")
		}
		fmt.Print("\n")
	}
}
