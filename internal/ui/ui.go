package ui

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type GetFileJSONBody struct {
	Filename string `json:"filename"`
	CID      string `json:"cid"`
}

type UploadFileJSONBody struct {
	Filepath string `json:"filepath"`
}

var backend *Backend

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
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "400 - Missing CID and Filename field in request\n")
				return
			}
			fileaddress := ""

			if _, err := os.Stat("files/stored/" + payload.Filename); !os.IsNotExist(err) {
				fileaddress = "files/stored/" + payload.Filename
			}
			if _, err := os.Stat("files/requested/" + payload.Filename); !os.IsNotExist(err) && fileaddress == "" {
				fileaddress = "files/requested/" + payload.Filename
			}
			if _, err := os.Stat("files/" + payload.Filename); !os.IsNotExist(err) && fileaddress == "" {
				fileaddress = "files/" + payload.Filename
			}
			if fileaddress != "" {

			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "500 - Internal Server Error: Cannot find file inside files directory")
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
				w.WriteHeader(http.StatusInternalServerError)
				writeStatusUpdate(w, "Failed to read in file from given path")
				return
			}
			lenData := len(fileData)
			base64Encode := base64.StdEncoding.EncodeToString(fileData)
			hash := sha256.Sum256(fileData)

			// Encode the hash as a hexadecimal string
			hexHash := hex.EncodeToString(hash[:])

			fileInfoResp := FileInfo{
				Filename:     filename,
				Filesize:     lenData,
				Filehash:     hexHash,
				Lastmodified: st.ModTime().String(),
				Filecontent:  base64Encode,
			}
			jsonData, err := json.Marshal(fileInfoResp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				writeStatusUpdate(w, "Failed to convert JSON Data into a string")
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "405 - Method Not Allowed: Only GET requests are supported\n")
		return
	}
}
func writeStatusUpdate(w http.ResponseWriter, message string) {
	responseMsg := map[string]interface{}{
		"status": message,
	}
	responseMsgJsonString, err := json.Marshal(responseMsg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseMsgJsonString)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			var payload UploadFileJSONBody
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&payload); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "500 - Internal Server Error: %v\n", err)
				return
			}
			fileData, err := os.ReadFile(payload.Filepath)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			if _, err := os.Stat(payload.Filepath); !os.IsNotExist(err) {
				sourceFile, err := os.Open(payload.Filepath)
				if err != nil {
					fmt.Println("Error opening source file:", err)
					return
				}
				defer sourceFile.Close()
				hash := sha256.Sum256(fileData)

				// Encode the hash as a hexadecimal string
				hexHash := hex.EncodeToString(hash[:])

				// Create the destination file in the destination folder
				destinationFilePath := "files/" + hexHash
				destinationFile, err := os.Create(destinationFilePath)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "500 - Internal Server Error: %v\n", err)
					return
				}
				defer destinationFile.Close()

				_, err = io.Copy(destinationFile, sourceFile)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "500 - Internal Server Error: %v\n", err)
					return
				}
				w.WriteHeader(http.StatusOK)
				writeStatusUpdate(w, "Successfully uploaded file from local computer into files directory")
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "500 - Internal Server Error: %v\n", err)
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
func joinStrings(strings []string, delimiter string) string {
	if len(strings) == 0 {
		return ""
	}
	result := strings[0]
	for _, s := range strings[1:] {
		result += delimiter + s
	}
	return result
}
func getActivities(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		allActivities, err := backend.GetActivities()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 - Internal Server Error: %v\n", err)
			return
		}
		var activityStrings []string
		for _, activity := range allActivities {
			activityString, err := json.Marshal(activity)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			activityStrings = append(activityStrings, string(activityString))
		}
		jsonArrayString := "[" + joinStrings(activityStrings, ",") + "]"
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonArrayString))

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "405 - Method Not Allowed: Only GET requests are supported\n")
		return
	}

}

func setActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			var payload Activity
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&payload); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "500 - Internal Server Error: %v\n", err)
				return
			}
			backend.SetActivity(payload)
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

type RemoveActivityJSONBody struct {
	Id int `json:"id"`
}

func removeActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			var payload RemoveActivityJSONBody
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&payload); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "500 - Internal Server Error: %v\n", err)
				return
			}
			backend.RemoveActivity(payload.Id)
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

type UpdateActivityJSONBody struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func updateActivityName(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			var payload UpdateActivityJSONBody
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&payload); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "500 - Internal Server Error: %v\n", err)
				return
			}
			backend.UpdateActivityName(payload.Id, payload.Name)
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

type WriteFileJSONBody struct {
	Base64File       string `json:"base64File"`
	Filesize         string `json:"fileSize"`
	OriginalFileName string `json:"originalFileName"`
}

func writeFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			var payload WriteFileJSONBody
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&payload); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "500 - Internal Server Error: %v\n", err)
				return
			}
			backend.UploadFile(payload.Base64File, payload.OriginalFileName, payload.Filesize)
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
	backend = NewBackend()
	http.HandleFunc("/getFile", getFile)
	http.HandleFunc("/getFileInfo", getFileInfo)
	http.HandleFunc("/uploadFile", uploadFile)
	http.HandleFunc("/deleteFile", deleteFile)
	http.HandleFunc("/updateActivityName", updateActivityName)
	http.HandleFunc("/removeActivity", removeActivity)
	http.HandleFunc("/setActivity", setActivity)
	http.HandleFunc("/getActivities", getActivities)
	http.HandleFunc("/writeFile", writeFile)
}
