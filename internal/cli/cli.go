package cli

import (
	"bufio"
	"fmt"
	"net"
	orcaClient "orca-peer/internal/client"
	orcaServer "orca-peer/internal/server"
	orcaStatus "orca-peer/internal/status"
	"os"
	"strings"
)

func StartCLI() {
	fmt.Println("Welcome to Orcanet!")
	fmt.Println("Dive In and Explore! Type 'help' for available commands.")
	port := getPort()
	serverReady := make(chan bool)
	confirming := false
	confirmation := ""
	go orcaServer.StartServer(port, serverReady, &confirming, &confirmation)
	<-serverReady

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

		if *confirming {
			switch command {
			case "yes":
				*confirmation = "yes"
			default:
				*confirmation = "no"
			}
			continue
		}

		switch command {
		case "get":
			if len(args) == 3 {
				go orcaClient.GetFileOnce(args[0], args[1], args[2])
			} else {
				fmt.Println("Usage: get [ip] [port] [filename>]")
				fmt.Println()
			}
		case "store":
			if len(args) == 3 {
				go orcaClient.RequestStorage(args[0], args[1], args[2])
			} else {
				fmt.Println("Usage: store [ip] [port] [filename]")
				fmt.Println()
			}
		case "import":
			if len(args) == 1 {
				go orcaClient.ImportFile(args[0])
			} else {
				fmt.Println("Usage: import [filepath]")
				fmt.Println()
			}
		case "location":
			fmt.Println(orcaStatus.GetLocationData())
		case "network":
			fmt.Println("Testing Network Speeds...")
			networkData := orcaStatus.GetNetworkInfo()
			if networkData.Success {
				fmt.Printf("Latency: %fms, Download: %fMbps, Upload: %fMbps\n", networkData.LatencyMs, networkData.DownloadSpeedMbps, networkData.UploadSpeedMbps)
			} else {
				fmt.Println("Unable to test network speeds. Please try again")
			}

		case "list":
			// TO-DO
		case "exit":
			fmt.Println("Exiting...")
			return
		case "help":
			fmt.Println("COMMANDS:")
			fmt.Println(" get [ip] [port] [filename]     Request a file")
			fmt.Println(" store [ip] [port] [filename]   Request storage of a file")
			fmt.Println(" import [filepath]              Import a file")
			fmt.Println(" list                           List all files you are storing")
			fmt.Println(" location						 Print your location")
			fmt.Println(" network						 Test speed of network")
			fmt.Println(" exit                           Exit the program")
			fmt.Println()
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
			fmt.Println()
		}
	}
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
