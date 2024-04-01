package main

import (
	"flag"
	orcaCLI "orca-peer/internal/cli"
	orcaTest "orca-peer/test"
	"os"
)

var test bool
var boostrapNodeAddress string

func main() {
	flag.BoolVar(&test, "test", false, "Create test server with no CLI.")
	flag.StringVar(&boostrapNodeAddress, "bootstrap", "", "Give address to boostrap.")
	flag.Parse()

	os.MkdirAll("./files/stored/", 0755)
	if test {
		orcaTest.RunTestServer()
	} else {
		orcaCLI.StartCLI(&boostrapNodeAddress)
	}
}
