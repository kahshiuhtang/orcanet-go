package client

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"orca-peer/internal/hash"
	orcaHash "orca-peer/internal/hash"
	"os"
	"path/filepath"
)

type Client struct {
	name_map hash.NameMap
}

func NewClient(path string) *Client {
	return &Client{
		name_map: *hash.NewNameStore(path),
	}
}

type FileData struct {
	FileName string `json:"filename"`
	Content  []byte `json:"content"`
}

func (client *Client) ImportFile(filePath string) {
	// Extract filename from the provided file path
	_, fileName := filepath.Split(filePath)
	if fileName == "" {
		fmt.Print("\nProvided path is a directory, not a file\n\n> ")
		return
	}

	// Open the source file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Print("\nFile does not exist\n\n> ")
		return
	}
	defer file.Close()

	// Create the directory if it doesn't exist
	err = os.MkdirAll("./files", 0755)
	if err != nil {
		return
	}

	// Save the file to the destination directory with the same filename
	destinationPath := filepath.Join("./files", fileName)
	destinationFile, err := os.OpenFile(destinationPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer destinationFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destinationFile, file)
	if err != nil {
		return
	}

	fmt.Printf("\nFile '%s' imported successfully!\n\n> ", fileName)
}

type Data struct {
	Bytes               []byte `json:"bytes"`
	UnlockedTransaction []byte `json:"transaction"`
	PublicKey           string `json:"public_key"`
}

func SendTransaction(price float64, ip string, port string, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) {
	cost := orcaHash.GeneratePriceBytes(price)
	byteBuffer := bytes.NewBuffer(cost)
	pubKeyString, err := orcaHash.ExportRsaPublicKeyAsPemStr(publicKey)
	if err != nil {
		fmt.Println("Error sending public key in header:", err)
		return
	}
	data := Data{
		Bytes:               byteBuffer.Bytes(),
		UnlockedTransaction: cost,
		PublicKey:           string(pubKeyString),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/sendTransaction", ip, port), bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	fmt.Println("Verifying Signature...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	} else {
		fmt.Println("Send Request")
	}
	defer resp.Body.Close()

}
func (client *Client) GetFileOnce(ip, port, filename string) {
	file_hash := client.name_map.GetFileHash(filename)
	if file_hash == "" {
		fmt.Println("Error: do not have hash for the file")
		return
	}
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/requestFile/%s", ip, port, file_hash))
	if err != nil {
		fmt.Printf("Error: %s\n\n", err)
		return
	}
	fmt.Println("Response:")
	fmt.Println(resp)
	fmt.Println("ResponseBody:")
	fmt.Println(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		fmt.Printf("\nError: %s\n> ", body)
		return
	}

	// Create the directory if it doesn't exist
	err = os.MkdirAll("./files/requested/", 0755)
	if err != nil {
		panic(err)
	}

	// Create file
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

	fmt.Printf("\nFile %s downloaded successfully!\n\n> ", filename)
}

func (client *Client) RequestStorage(ip, port, filename string) {
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

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		fmt.Printf("\nError: %s\n> ", body)
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	client.name_map.PutFileHash(filename, string(body))

	fmt.Println(string(body))
	fmt.Print("> ")
}
