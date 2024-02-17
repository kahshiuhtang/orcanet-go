package main

import (
	"fmt"
	"os"
	pnc "peer-node/client"
	pns "peer-node/server"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("SeaTurtle Client", "Client that allows users to produce/consume files.")
	var producer *bool = parser.Flag("p", "producer", &argparse.Options{
		Required: false,
		Default:  true,
	})
	var consumer *bool = parser.Flag("c", "consumer", &argparse.Options{
		Required: false,
		Default:  true,
	})
	var port *string = parser.String("p", "port", &argparse.Options{
		Required: true,
		Default:  ":8080",
	})
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		fmt.Print(parser.Usage(err))
	}
	if *producer && *consumer {
		fmt.Println("[Error]: Can be producer or consumer, not both.")
		return
	}
	if *producer {
		var client = pns.Producer{CurrentCoins: 0}
		client.SetupProducer(*port)
	}
	if *consumer {
		var client = pnc.Consumer{CurrentCoins: 0.0}
		client.SetupConsumer()
		client.RequestFileFromProducer("http://127.0.0.1:9095/")
	}
}
