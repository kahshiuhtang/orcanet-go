package client

import (
	"context"
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

func RequestFileFromMarket(client pb.FileShareClient, fileDesc *pb.FileDesc) string {
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
	// client := pb.NewFileShareClient(conn)

}
