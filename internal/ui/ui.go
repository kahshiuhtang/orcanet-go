package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
			if _, err := os.Stat("files/" + payload.Filename); !os.IsNotExist(err) {
				notFound = true
			}
			if _, err := os.Stat("files/requested/" + payload.Filename); !os.IsNotExist(err) {
				notFound = true
			}
			if _, err := os.Stat("files/stored/" + payload.Filename); !os.IsNotExist(err) {
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

func getFileInfo(w http.ResponseWriter, r *http.Request) {

}

func uploadFile(w http.ResponseWriter, r *http.Request) {

}

func deleteFile(w http.ResponseWriter, r *http.Request) {

}

func InitServer() {
	http.HandleFunc("/getFile", getFile)
	http.HandleFunc("/getFileInfo", getFileInfo)
	http.HandleFunc("/uploadFile", uploadFile)
	http.HandleFunc("/deleteFile", deleteFile)
}
