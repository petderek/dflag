// +build example

package main

import (
	"fmt"

	"github.com/petderek/dflag"
)

var flags = struct {
	Count    int    `name:"c" value:"10" usage:"the number of times to print the word"`
	Word     string `value:"foo" usage:"the word to print"`
	NewLines bool   `value:"true" usage:"if we should print each word on its own line"`
}{}

// to test this:
// go run -tags=example ./example.go -c 5 -word bar -newlines=false
func main() {
	dflag.Parse(&flags)
	for i := 0; i < flags.Count; i++ {
		if flags.NewLines {
			fmt.Println(flags.Word)
		} else {
			fmt.Print(flags.Word)
		}
	}
	if !flags.NewLines {
		fmt.Print("\n")
	}
}
