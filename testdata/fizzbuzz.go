// +build fizzbuzz

package main

import (
	"fmt"
	"os"

	"github.com/petderek/dflag"
)

var flags = struct {
	FizzOn int
	BuzzOn int
}{}

// to test this:
// go run -tags=fizzbuzz ./fizzbuzz.go -fizzon 3 -buzzon 5
func main() {
	dflag.Parse(&flags)
	if flags.FizzOn <= 0 || flags.BuzzOn <= 0 {
		fmt.Println("fizzon and buzzon must be set to positive numbers")
		os.Exit(1)
	}

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
