package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	api "orca-peer/internal/api"
	"orca-peer/internal/hash"
	"os"
	"path/filepath"
	"time"
)

const keyServerAddr = "serverAddr"

var (
	eventChannel chan bool
)

type Server struct {
	storage *hash.DataStore
}

func Init() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/requestFile", getFile)
	mux.HandleFunc("/sendTransaction", handleTransaction)

	ctx := context.Background()
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}
	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Server closed\n")
		} else if err != nil {
			fmt.Printf("Error listening for server: %s\n", err)
		}
	}()
}

type TransactionFile struct {
	Bytes               []byte `json:"bytes"`
	UnlockedTransaction []byte `json:"transaction"`
	PublicKey           string `json:"public_key"`
}
type Transaction struct {
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
	Uuid      string  `json:"uuid"`
}
type StatusResponse struct {
	Status string `json:"status"`
}

func sendStatusResponse(w http.ResponseWriter, status string, statusCode int) {
	responseData := StatusResponse{
		Status: status,
	}

	// Marshal the ResponseData struct into JSON
	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response with a 200 status code
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonResponse)
	if err != nil {
		fmt.Println("Failed to write JSON response:", err)
	}
}

func handleTransaction(w http.ResponseWriter, r *http.Request) {
	// Read the Request Body
	fmt.Println("Handling a transaction...")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading Request Body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Print the received byte string
	var data TransactionFile
	fmt.Println("Received byte string:")
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	publicKey, err := hash.ParseRsaPublicKeyFromPemStr(data.PublicKey)
	if err != nil {
		fmt.Println("Unable to unmarshalling public key:", err)
		return
	}
	timestamp := time.Now()
	timestampStr := timestamp.Format(time.RFC3339)
	err = os.WriteFile("./files/transactions/"+timestampStr, body, 0644)
	if err != nil {
		fmt.Println("Error writing transaction to file:", err)
		sendStatusResponse(w, "", http.StatusInternalServerError)
		return
	}
	fmt.Println("Data in struct:", data)
	error := hash.VerifySignature(data.UnlockedTransaction, data.Bytes, publicKey)
	if error != nil {
		fmt.Println("Properly Hashed Transaction")
	} else {
		fmt.Println("Did not properly hash transaction")
		sendStatusResponse(w, "Public Key does not unlock the locked transaction", http.StatusBadRequest)
		return
	}
	var transaction Transaction
	err = json.Unmarshal(data.UnlockedTransaction, &transaction)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		sendStatusResponse(w, "Error Marshalling JSON", http.StatusBadRequest)
		return
	} else {
		fmt.Println("Transaction JSON:")
		fmt.Println("Sent Amount:")
		fmt.Println(transaction.Price)
		fmt.Println("UUID:")
		fmt.Println(transaction.Uuid)
		sendStatusResponse(w, "Successfully stored transaction", http.StatusOK)
	}
	eventChannel <- true
	fmt.Println("> ")
}

// Start HTTP server
func StartServer(port string, serverReady chan bool, confirming *bool, confirmation *string) {
	eventChannel = make(chan bool)
	server := Server{
		storage: hash.NewDataStore("files/stored/"),
	}
	api.InitServer()
	http.HandleFunc("/requestFile/", func(w http.ResponseWriter, r *http.Request) {
		server.sendFile(w, r, confirming, confirmation)
	})
	http.HandleFunc("/storeFile/", func(w http.ResponseWriter, r *http.Request) {
		server.storeFile(w, r, confirming, confirmation)
	})
	http.HandleFunc("/sendTransaction", handleTransaction)

	fmt.Printf("Listening on port %s...\n", port)
	serverReady <- true
	http.ListenAndServe(":"+port, nil)
}

func (server *Server) sendFile(w http.ResponseWriter, r *http.Request, confirming *bool, confirmation *string) {
	// Extract filename from URL path
	filename := r.URL.Path[len("/requestFile/"):]

	// Ask for confirmation
	*confirming = true
	fmt.Printf("\nYou have just received a request to send file '%s'. Do you want to send the file? (yes/no): ", filename)

	// Check if confirmation is received
	for *confirmation != "yes" {
		if *confirmation != "" {
			http.Error(w, fmt.Sprintf("Client declined to send file '%s'.", filename), http.StatusUnauthorized)
			*confirmation = ""
			*confirming = false
			return
		}
	}
	*confirmation = ""
	*confirming = false

	file, err := os.Open("./files/" + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type
	contentType := "application/octet-stream"
	switch {
	case filename[len(filename)-4:] == ".txt":
		contentType = "text/plain"
	case filename[len(filename)-5:] == ".json":
		contentType = "application/json"
	case filename[len(filename)-4:] == ".mp4":
		contentType = "video/mp4"
	}

	// Set content disposition header
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", contentType)

	const chunkSize = 1024
	fmt.Println("File size: ")
	fmt.Println(stat.Size())
	if stat.Size() > chunkSize {
		fmt.Println("Must serve in chunks")
		buffer := make([]byte, chunkSize)
		for {
			//	time.Sleep(1 * time.Second)
			// Read 10 bytes from the file
			n, err := file.Read(buffer)
			// fmt.Println("n:", n)
			// fmt.Println("buffer:", buffer)
			// fmt.Println()
			if err != nil {
				// Check if it's the end of the file
				if err.Error() == "EOF" {
					break
				}
				// Handle other errors
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			fmt.Println("Sending chunk...")
			// Write the 10-byte chunk to the response
			w.Write(buffer[:n])
			<-eventChannel
			//w.Write([]byte("\n@@@@\n"))
		}
	} else {
		fmt.Println("sending in one piece")
		// Copy file contents to Response Body
		_, err = io.Copy(w, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	fmt.Printf("\nFile %s sent!\n> ", filename)
}

type FileData struct {
	FileName string `json:"filename"`
	Content  []byte `json:"content"`
}

func (server *Server) storeFile(w http.ResponseWriter, r *http.Request, confirming *bool, confirmation *string) {
	// Parse JSON object from Request Body
	var fileData FileData
	err := json.NewDecoder(r.Body).Decode(&fileData)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// Ask for confirmation
	*confirming = true
	fmt.Printf("\nYou have just received a request to store file '%s'. Do you want to store the file? (yes/no): ", fileData.FileName)

	// Check if confirmation is received
	for *confirmation != "yes" {
		if *confirmation != "" {
			http.Error(w, fmt.Sprintf("Client declined to store file '%s'.", fileData.FileName), http.StatusUnauthorized)
			*confirmation = ""
			*confirming = false
			return
		}
	}
	*confirmation = ""
	*confirming = false

	// Create file
	file_hash, err := server.storage.PutFile(fileData.Content)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", file_hash)
	fmt.Printf("\nStored file %s hash %s!\n> ", fileData.FileName, file_hash)
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /root request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func getFile(w http.ResponseWriter, r *http.Request) {
	// Get the context from the request
	ctx := r.Context()

	// Check if the "filename" query parameter is present
	hasFilename := r.URL.Query().Has("filename")

	// Retrieve the value of the "filename" query parameter
	filename := r.URL.Query().Get("filename")

	// Print information about the request
	fmt.Printf("%s: got /file request. filename(%t)=%s\n",
		ctx.Value(keyServerAddr),
		hasFilename, filename,
	)

	// Check if the "filename" parameter is present
	if hasFilename {
		// Check if the file exists in the local directory
		filePath := filepath.Join(".", filename)
		if _, err := os.Stat(filePath); err == nil {
			// Serve the file using http.ServeFile
			http.ServeFile(w, r, filePath)
			fmt.Printf("Served %s to client\n", filename)
			return
		} else if os.IsNotExist(err) {
			// File not found
			http.Error(w, "File not found", http.StatusNotFound)
			return
		} else {
			// Other error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		// Write a response indicating that no filename was found
		io.WriteString(w, "No filename found\n")
	}
}
