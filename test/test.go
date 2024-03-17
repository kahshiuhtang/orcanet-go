package test

import (
	"log"
	orcaClient "orca-peer/internal/client"
	pb "orca-peer/internal/fileshare"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type testFilePeerServer struct {
	pb.UnimplementedFileShareServer
	savedAddress map[string][]*pb.StorageIP
	savedFiles   map[string][]*pb.FileDesc // read-only after initialized

	mu sync.Mutex // protects routeNotes
}

func main() {
	serverIP := "localhost:50051"
	go SetupTestMarket()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(serverIP, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewFileShareClient(conn)
	orcaClient.RequestFileFromMarket(client, &pb.FileDesc{})
}
