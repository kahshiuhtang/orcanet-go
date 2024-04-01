package main

import (
	"encoding/json"
	"flag"
	"fmt"
	orcaCLI "orca-peer/internal/cli"
	orcaHash "orca-peer/internal/hash"
	orcaTest "orca-peer/test"
	"os"
)

var test bool
var boostrapNodeAddress string

func main() {
	flag.BoolVar(&test, "test", false, "Create test server with no CLI.")
	flag.StringVar(&boostrapNodeAddress, "bootstrap", "", "Give address to boostrap.")
	flag.Parse()
	publicKey, privateKey := orcaHash.LoadInKeys()
	cost := orcaHash.GeneratePriceBytes(100)
	signed, err := orcaHash.SignFile(cost, privateKey)
	if err != nil {
		fmt.Println("Error signing file")
	}
	jsonData := map[string]interface{}{
		"key1": 10,
	}
	// Marshal JSON object to byte array
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	suc := orcaHash.VerifySignature(signed, jsonBytes, publicKey)
	if suc == nil {
		fmt.Println("Valid Transaction")
	} else {
		fmt.Println("Invalid Transaction")
	}
	os.MkdirAll("./files/stored/", 0755)
	if test {
		orcaTest.RunTestServer()
	} else {
		orcaCLI.StartCLI(&boostrapNodeAddress, publicKey, privateKey)
	}
}
