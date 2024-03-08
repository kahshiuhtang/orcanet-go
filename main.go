package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const keyServerAddr = "serverAddr"

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

		// Open the file
		filePath := filepath.Join(".", filename)
		file, err := os.Open(filePath)
		if err != nil {
			// Handle error (file not found, permission issues, etc.)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		// Check if the "chunksize" query parameter is present
		hasChunksize := r.URL.Query().Has("chunksize")

		if hasChunksize {
			// Retrieve the value of the "filename" query parameter
			chunkSizeParam := r.URL.Query().Get("chunksize")
			// Convert the value to an integer
			chunkSize, err := strconv.Atoi(chunkSizeParam)
			if err != nil {
				// Handle the error (invalid value)
				fmt.Printf("Invalid chunksize parameter value: %v\n", chunkSizeParam)
				return
			}

			fmt.Printf("chunkSize: %d\n", chunkSize)

			buffer := make([]byte, chunkSize)

			for {
				// Read 10 bytes from the file
				n, err := file.Read(buffer)
				if err != nil {
					// Check if it's the end of the file
					if err.Error() == "EOF" {
						break
					}
					// Handle other errors
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				// Write the 10-byte chunk to the response
				w.Write(buffer[:n])
				w.Write([]byte("\n***\n"))
			}
		} else {
			fmt.Printf("No chunk size specified\n")
			// Serve the file using http.ServeFile
			http.ServeFile(w, r, filePath)
			fmt.Printf("Served %s to client\n", filename)
			return
		}

	} else {
		// Write a response indicating that no filename was found
		io.WriteString(w, "No filename found\n")
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/requestFile", getFile)

	ctx := context.Background()
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}
