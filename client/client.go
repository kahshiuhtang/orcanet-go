package client

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	pb "peer-node/fileshare"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func requestFile(client pb.FileShareClient, fileDesc *pb.FileDesc) string {
	log.Printf("Requesting IP For File (%s)", fileDesc.FileName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := client.PlaceFileRequest(ctx, fileDesc)
	if err != nil {
		log.Fatalf("client.requestFileStorage failed: %v", err)
	}
	log.Println(feature)
	return ""
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

func RequestFileFromProducer(baseURL string, filename string) bool {
	encodedParams := url.Values{}
	encodedParams.Add("filename", filename)
	queryString := encodedParams.Encode()

	// Construct the URL with the query string
	urlWithQuery := fmt.Sprintf("%s?%s", baseURL, queryString)
	resp, err := http.Get(urlWithQuery)

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

func SendCurrency(baseURL string, amount int) bool {
	// Send the POST request
	jsonData, err := json.Marshal(amount)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	return false
}

var (
	marketAddr = flag.String("addr", "localhost:50051", "The market address in the format of host:port")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(*marketAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewFileShareClient(conn)

}
