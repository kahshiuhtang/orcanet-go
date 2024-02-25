package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath") 

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


