// +build urlparse

package main

import (
	"log"
	"net"
	"net/url"
	"os"

	"github.com/petderek/dflag"
)

var flags = struct {
	Url  url.URL `required:"false" usage:"the url to check"`
	Ping bool
}{}

// to test this:
// go run -tags=urlparse ./testdata -url www.google.com -port 80
func main() {
	dflag.Parse(&flags, dflag.Logger(log.Default()))
	actualUrl := flags.Url

	// assumptions!!
	addr := actualUrl.Host
	if actualUrl.Port() == "" {
		if actualUrl.Scheme == "https" {
			addr += ":443"
		} else if actualUrl.Scheme == "http" {
			addr += ":80"
		}
	}

	if !flags.Ping {
		log.Println("Address is ", addr)
		os.Exit(0)
	}

	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("success: ", c.RemoteAddr().String())
	c.Close()
}
