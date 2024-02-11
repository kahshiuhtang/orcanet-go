package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Consumer struct {
	currentCoins float64
}

func (cons *Consumer) SetupConsumer() bool {
	return true
}

// Use requestFile if the function is internal, otherwise name it RequestFile
func (cons *Consumer) RequestFileFromMarket() bool {
	return false
}

func (cons *Consumer) RequestFileFromProducer(address string) bool {
	resp, err := http.Get(address)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
	fmt.Println(sb)
	return false
}

func (cons *Consumer) SendCurrency() bool {
	return false
}
