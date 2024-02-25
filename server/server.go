package server

// formatting and printing values to the console.
import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"

	pb "peer-node/fileshare"

	"google.golang.org/grpc"
)

// Used for build HTTP servers and clients.
type fileShareServerNode struct {
	pb.UnimplementedFileShareServer
	savedFiles   map[string][]*pb.FileDesc // read-only after initialized
	mu           sync.Mutex                // protects routeNotes
	currentCoins float64
}

func sendFileToConsumer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		for k, v := range r.URL.Query() {
			fmt.Printf("%s: %s\n", k, v)
		}
		// file = r.URL.Query().Get("filename")
		w.Write([]byte("Received a GET request\n"))
	case "POST":
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", reqBody)
		w.Write([]byte("Received a POST request\n"))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
	w.Write([]byte("Received a GET request\n"))
	filename := r.URL.Path[len("/reqFile/"):]

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set content type
	contentType := "application/octet-stream"
	switch {
	case filename[len(filename)-4:] == ".txt":
		contentType = "text/plain"
	case filename[len(filename)-5:] == ".json":
		contentType = "application/json"
		// Add more cases for other file types if needed
	}

	// Set content disposition header
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", contentType)

	// Copy file contents to response body
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type sentCoinMessage struct {
	Amount string `json:"amount"`
}

func setupProducer(gRPCPort int, httpPort int) *fileShareServerNode {
	s := &fileShareServerNode{savedFiles: make(map[string][]*pb.FileDesc)}
	// s.loadMappings(*jsonDBFile) // Have a load and save mappings
	http.HandleFunc("/file", sendFileToConsumer)
	fmt.Println("[Server]: Listening On Port" + strconv.Itoa(httpPort))
	fmt.Println("[Server]: Press CTRL + C to quit.")
	go func() {
		for {
			http.ListenAndServe(":"+strconv.Itoa(httpPort), nil)
		}
	}()
	return s
}

// Can add back in TLS later
var (
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	gRPCPort   = flag.Int("gport", 50051, "The gRPC port for send/receive gRPC")
	httpPort   = flag.Int("hport", 50052, "The server port for listening for HTTP requests")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *gRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterFileShareServer(grpcServer, setupProducer(*gRPCPort, *httpPort))
	grpcServer.Serve(lis)
}
