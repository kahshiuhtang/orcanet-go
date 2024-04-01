package server

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"orca-peer/internal/hash"
	"os"
	"path/filepath"

	"aead.dev/minisign"
)

const keyServerAddr = "serverAddr"

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
func handleTransaction(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Print the received byte string
	fmt.Println("Received byte string:", string(body))

	// Respond with a success message
	fmt.Fprintf(w, "Byte string received successfully\n")
}

// Start HTTP server
func StartServer(port string, serverReady chan bool, confirming *bool, confirmation *string) {
	server := Server{
		storage: hash.NewDataStore("files/stored/"),
	}
	http.HandleFunc("/requestFile/", func(w http.ResponseWriter, r *http.Request) {
		server.sendFile(w, r, confirming, confirmation)
	})
	http.HandleFunc("/storeFile/", func(w http.ResponseWriter, r *http.Request) {
		server.storeFile(w, r, confirming, confirmation)
	})

	fmt.Printf("Listening on port %s...\n\n", port)
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

	// Open the file
	file_data, err := server.storage.GetFile(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}
	fmt.Println(bs)

	// Generate a new minisign private / public key pair.
	publicKey, privateKey, err := minisign.GenerateKey(rand.Reader)
	if err != nil {
		panic(err) // TODO: error handling
	}
	fmt.Println("Public Key:", publicKey)

	// Sign bytes with the private key
	message := []byte(bs)
	signature := minisign.Sign(privateKey, message)
	fmt.Println("singature:", signature)

	if !minisign.Verify(publicKey, message, signature) {
		log.Fatalln("signature verification failed")
	} else {
		fmt.Println("signature verification succeeded")
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
	w.Header().Set("X-Signature", string(signature))
	w.Header().Set("X-PublicKey", publicKey.String())
	w.Header().Set("X-Message", string(message))

	// Copy file contents to response body
	_, err = io.Copy(w, bytes.NewBuffer(file_data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("\nFile %s sent!\n\n> ", filename)
}

type FileData struct {
	FileName string `json:"filename"`
	Content  []byte `json:"content"`
}

func (server *Server) storeFile(w http.ResponseWriter, r *http.Request, confirming *bool, confirmation *string) {
	// Parse JSON object from request body
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
	fmt.Printf("\nStored file %s hash %s!\n\n> ", fileData.FileName, file_hash)
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
