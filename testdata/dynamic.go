// +build dynamic

package main

import (
	"fmt"

	"github.com/petderek/dflag"
)

var flags = struct {
	Dynamic string
	Static  string `value:"staticdef"`
	Shadow  string `value:"xyz"`
}{}

// to test this:
// go run -tags=dynamic ./dynamic.go -dynamic foo -static bar
func main() {
	// setting values before calling parse will take precedence
	// order is:
	// 1. flag
	// 2. dynamic value assigned to struct
	// 3. static tag value

	flags.Dynamic = "dynamicdef"
	flags.Shadow = "shadowdef"
	dflag.Parse(&flags)

	fmt.Printf("%s\n%s\n%s\n", flags.Dynamic, flags.Static, flags.Shadow)
}
