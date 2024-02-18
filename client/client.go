package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	pb "peer-node/fileshare"
	"time"
)

type Consumer struct {
	CurrentCoins float64
}

func requestFile(client pb.FileShareClient, fileDesc *pb.FileDesc) {
	log.Printf("Requesting IP For File (%s)", fileDesc.FileName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := client.PlaceFileRequest(ctx, fileDesc)
	if err != nil {
		log.Fatalf("client.requestFileStorage failed: %v", err)
	}
	log.Println(feature)
}

func requestFileStorage(client pb.FileShareClient, fileDesc *pb.FileDesc) {
	log.Printf("Requesting IP For Storage For File (%s)", fileDesc.FileName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := client.PlaceFileStoreRequest(ctx, fileDesc)
	if err != nil {
		log.Fatalf("client.requestFileStorage failed: %v", err)
	}
	log.Println(feature)
}

func (cons *Consumer) SetupConsumer() bool {
	return true
}

// Use requestFile if the function is internal, otherwise name it RequestFile
// Should use RPC
func (cons *Consumer) RequestFileFromMarket(priceOffer float64) bool {
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
	fmt.Println(sb)
	return false
}

func (cons *Consumer) SendCurrency() bool {
	return false
}
