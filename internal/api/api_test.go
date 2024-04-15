package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHashWrongRequestType(t *testing.T) {
	req, err := http.NewRequest("GET", "/hash", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hashFile)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	fmt.Println(rr.Body)
	//expected := "Hello, World!"
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}

func TestHash(t *testing.T) {
	requestBody := map[string]string{"filepath": "tester.txt"}

	// Marshal the JSON payload
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/hash", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hashFile)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var responseData map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseData)
	if err != nil {
		t.Errorf("Bad response from server, unable to parse JSON")
	}
	// fileData, err := os.ReadFile("files/tester.txt")
	// if err != nil {
	// 	t.Errorf("Unable to load in file for testing")
	// }
	// expected := sha256.Sum256(fileData)
	// expected = expected
}

func TestGetFile(t *testing.T) {
	requestBody := map[string]string{"filename": "tester.txt", "cid": ""}

	// Marshal the JSON payload
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/getFile", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getFile)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var responseData map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseData)
	if err != nil {
		t.Errorf("Bad response from server, unable to parse JSON")
	}
	// fileData, err := os.ReadFile("files/tester.txt")
	// if err != nil {
	// 	t.Errorf("Unable to load in file for testing")
	// }
	// expected := sha256.Sum256(fileData)
	// expected = expected
}
