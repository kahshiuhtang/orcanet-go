package main

// formatting and printing values to the console.
import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Used for build HTTP servers and clients.
type Producer struct {
	portNum      string
	currentCoins float64
	fileMappings map[string]string
}

// Send some response back to client to
func verify(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "VERIFY")
}
func sendFile(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "SEND FILE")
	// Extract filename from URL path
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
func handleCoin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "HANDLE COIN")
}
func (prod *Producer) SetupProducer(port string) bool {
	prod.fileMappings = make(map[string]string)
	http.HandleFunc("/", verify)
	http.HandleFunc("/reqFile/", sendFile)
	http.HandleFunc("/recvCoin", handleCoin)
	fmt.Println("[Client]: Listening On Port" + port)
	fmt.Println("[Client]: Press CTRL + C to quit.")
	http.ListenAndServe(port, nil)
	return true
}

func (prod *Producer) SendFile() bool {
	return false
}
