package main

import (
	"flag"
	orcaCLI "orca-peer/internal/cli"
	orcaTest "orca-peer/test"
)

var test bool

func main() {
	flag.BoolVar(&test, "test", false, "Create test server with no CLI.")
	flag.Parse()
	if test {
		orcaTest.RunTestServer()
	} else {
		orcaCLI.StartCLI()
	}
}
