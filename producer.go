package main

// formatting and printing values to the console.
import (
	"fmt"
	"net/http"
)

// Used for build HTTP servers and clients.
type Producer struct {
	portNum      string
	currentCoins float64
}

// Send some response back to client to
func verify(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "VERIFY")
}
func sendFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "SEND FILE")
}
func handleCoin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "HANDLE COIN")
}
func (prod *Producer) SetupProducer(port string) bool {
	http.HandleFunc("/", verify)
	http.HandleFunc("/reqFile", sendFile)
	http.HandleFunc("/recvCoin", handleCoin)
	fmt.Println("[Client]: Listening On Port" + port)
	fmt.Println("[Client]: Press CTRL + C to quit.")
	http.ListenAndServe(port, nil)
	return true
}

func (prod *Producer) SendFile() bool {
	return false
}
