package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

type FileData struct {
	FileName string `json:"filename"`
	Content  []byte `json:"content"`
}

func sendFile(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL path
	filename := r.URL.Path[len("/requestFile/"):]

	// Open the file
	file, err := os.Open("./files/stored/" + filename)
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

	fmt.Printf("\nFile %s sent!\n", filename)
}

func storeFile(w http.ResponseWriter, r *http.Request) {
	// Parse JSON object from request body
	var fileData FileData
	err := json.NewDecoder(r.Body).Decode(&fileData)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// Create file
	file, err := os.Create("./files/stored/" + fileData.FileName)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Write content to file
	_, err = file.Write(fileData.Content)
	if err != nil {
		http.Error(w, "Failed to write to file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Requested client stored file %s successfully!\n", fileData.FileName)
	fmt.Printf("\nStored file %s!\n", fileData.FileName)
}

func getFile(ip, port, filename string) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/requestFile/%s", ip, port, filename))
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n\n", resp.Status)
		return
	}

	out, err := os.Create("./files/requested/" + filename)
	if err != nil {
		return
	}
	defer out.Close()

	// Write response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}

	fmt.Printf("File %s downloaded successfully!\n\n", filename)
}

func requestStorage(ip, port, filename string) {
	// Read file content
	content, err := os.ReadFile("./files/documents/" + filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Create FileData struct
	fileData := FileData{
		FileName: filename,
		Content:  content,
	}

	// Marshal FileData to JSON
	jsonData, err := json.Marshal(fileData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Send POST request to store file
	resp, err := http.Post(fmt.Sprintf("http://%s:%s/storeFile/", ip, port), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println(string(body))
	fmt.Print()
}

// Ask user to enter a port and returns it
func getPort() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter a port number to start listening to requests: ")
	for {
		port, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}
		port = strings.TrimSpace(port)

		// Validate port
		listener, err := net.Listen("tcp", ":"+port)
		if err == nil {
			defer listener.Close()
			return port
		}

		fmt.Print("Invalid port. Please enter a different port: ")
	}
}

// Start HTTP server
func startServer(port string, serverReady chan<- bool) {
	http.HandleFunc("/requestFile/", sendFile)
	http.HandleFunc("/storeFile/", storeFile)

	fmt.Printf("Listening on port %s...\n\n", port)
	serverReady <- true
	http.ListenAndServe(":"+port, nil)
}

// Start CLI
func startCLI() {
	fmt.Println("Dive In and Explore! Type 'help' for available commands.")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
			continue
		}

		text = strings.TrimSpace(text)
		parts := strings.Fields(text)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		switch command {
		case "get":
			if len(args) == 3 {
				getFile(args[0], args[1], args[2])
			} else {
				fmt.Println("Usage: get <ip> <port> <filename>")
				fmt.Println()
			}
		case "store":
			if len(args) == 3 {
				requestStorage(args[0], args[1], args[2])
			} else {
				fmt.Println("Usage: store <ip> <port> <filename>")
				fmt.Println()
			}
		case "list":
			// TO-DO
		case "exit":
			fmt.Println("Exiting...")
			return
		case "help":
			fmt.Println("COMMANDS:")
			fmt.Println(" get <ip> <port> <filename>     Request a file")
			fmt.Println(" store <ip> <port> <filename>   Request storage of a file")
			fmt.Println(" list                           List all files you are storing")
			fmt.Println(" exit                           Exit the program")
			fmt.Println()
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
			fmt.Println()
		}
	}
}

func main() {
	fmt.Println("Welcome to Orcanet!")
	port := getPort()

	serverReady := make(chan bool)
	go startServer(port, serverReady)
	<-serverReady

	startCLI()
}
