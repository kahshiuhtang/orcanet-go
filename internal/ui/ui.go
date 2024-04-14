package ui

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type GetFileJSONBody struct {
	Filename string `json:"filename"`
	CID      string `json:"cid"`
}

func getFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			// For JSON content type, decode the JSON into a struct
			var payload GetFileJSONBody
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&payload); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "500 - Internal Server Error: %v\n", err)
				return
			}
			if payload.Filename == "" && payload.CID == "" {
				fmt.Fprintf(w, "400 - Missing CID and Filename field in request\n")
				return
			}
			notFound := false

			if _, err := os.Stat("files/stored" + payload.Filename); !os.IsNotExist(err) {
				notFound = true

			}
			if _, err := os.Stat("files/requested" + payload.Filename); !os.IsNotExist(err) {
				notFound = true
			}
			if _, err := os.Stat("files" + payload.Filename); !os.IsNotExist(err) {
				notFound = true
			}
			if !notFound {
				http.Error(w, "File not found in local storage", http.StatusNotFound)
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "400 - Bad Request: Unsupported content type: %s\n", contentType)
			return
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "405 - Method Not Allowed: Only POST requests are supported\n")
		return
	}
}

type FileInfo struct {
	Filename     string `json:"filename"`
	Filesize     int    `json:"filesize"`
	Filehash     string `json:"filehash"`
	Lastmodified string `json:"lastmodified"`
	Filecontent  string `json:"filecontent"`
}

func getFileInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		queryParams := r.URL.Query()

		// Retrieve specific query parameters by key
		filename := queryParams.Get("filename")
		if st, err := os.Stat("files/" + filename); !os.IsNotExist(err) {
			fileData, err := os.ReadFile("files/" + filename)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			lenData := len(fileData)
			base64Encode := base64.StdEncoding.EncodeToString(fileData)
			fileInfoResp := FileInfo{
				Filename:     filename,
				Filesize:     lenData,
				Filehash:     "",
				Lastmodified: st.ModTime().String(),
				Filecontent:  base64Encode,
			}
			jsonData, err := json.Marshal(fileInfoResp)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "405 - Method Not Allowed: Only GET requests are supported\n")
		return
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {

}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			// For JSON content type, decode the JSON into a struct
			var payload GetFileJSONBody
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&payload); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "500 - Internal Server Error: %v\n", err)
				return
			}
			fmt.Println("filename:", payload.Filename)
			if payload.Filename == "" && payload.CID == "" {
				fmt.Fprintf(w, "400 - Missing CID and Filename field in request\n")
				return
			}
			fileDir := "./files"
			var filePath string

			// Check if the file exists in the "stored" directory
			storedFilePath := filepath.Join(fileDir, "stored", payload.Filename)
			if _, err := os.Stat(storedFilePath); err == nil {
				filePath = storedFilePath
			}
			// Check if the file exists in the "requested" directory
			requestedFilePath := filepath.Join(fileDir, "requested", payload.Filename)
			if _, err := os.Stat(requestedFilePath); err == nil {
				filePath = requestedFilePath
			}
			fmt.Println("filePath:", filePath)

			// Attempt to delete the file
			err := os.Remove(filePath)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Println("File deleted successfully.")
			return

		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "400 - Bad Request: Unsupported content type: %s\n", contentType)
			return
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "405 - Method Not Allowed: Only POST requests are supported\n")
		return
	}

}

func InitServer() {
	http.HandleFunc("/getFile", getFile)
	http.HandleFunc("/getFileInfo", getFileInfo)
	http.HandleFunc("/uploadFile", uploadFile)
	http.HandleFunc("/deleteFile", deleteFile)
}
