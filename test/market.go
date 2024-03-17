package test

import (
	"fmt"
	"log"
	"net"

	pb "orca-peer/internal/fileshare"

	"google.golang.org/grpc"
)

func (s *testFilePeerServer) PlaceFileRequest(file *pb.FileDesc, stream pb.FileShare_PlaceFileRequestServer) error {
	//addresses := s.savedAddress[file.FileNameHash]
	addresses := s.savedAddress["file1"]
	for _, feature := range addresses {
		if err := stream.Send(feature); err != nil {
			return err
		}
	}
	return nil
}
func newServer() *testFilePeerServer {
	s := &testFilePeerServer{savedFiles: make(map[string][]*pb.FileDesc)}
	s.savedAddress = make(map[string][]*pb.StorageIP)
	s.savedFiles = make(map[string][]*pb.FileDesc)
	s.savedAddress["file1"] = make([]*pb.StorageIP, 0)
	s.savedAddress["file1"] = append(s.savedAddress["file1"], &pb.StorageIP{Success: true, Address: "localhost:9876", UserId: "server-1", FileName: "Test1.txt", FileByteSize: 100, FileCost: 1.0, IsLastCandidate: false})
	s.savedAddress["file1"] = append(s.savedAddress["file1"], &pb.StorageIP{Success: true, Address: "localhost:9877", UserId: "server-1", FileName: "Test1.txt", FileByteSize: 100, FileCost: 2.0, IsLastCandidate: false})
	s.savedAddress["file1"] = append(s.savedAddress["file1"], &pb.StorageIP{Success: true, Address: "localhost:9876", UserId: "server-1", FileName: "Test1.txt", FileByteSize: 100, FileCost: 5.0, IsLastCandidate: true})
	return s
}

func SetupTestMarket() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterFileShareServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
