package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var app App

var expectedQR = "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEACAMAAABrrFhUAAAABlBMVEX///8AAABVwtN+AAABp0lEQVR42uzdYU7EIBCA0e79L70naArD0A7wvn+apus+TZwU0EuSJEmSJEmSJEmSJEmSJEmSfsHu7nN33+zXAwAAAAAAAADkAWRd3wuV/fUBAAAAAAAAAOIAT4PI08etg1D09QAAAAAAAAAA9QBaH2AAAAAAAAAAAPYdhLIBAQAAAAAAAAB1F0aig882K0MAAAAAAADABgDRjYtZg9KyO0UBAAAAAACAhQHeArx2DQAAAAAAANjg9/isNzT7/gAAAAAAAACA+Qslowcnej9f5icCAAAAAAAA2Big9QVbF0h6NzyW+Y4DAAAAAAAABwNEB5fR60YPYgIAAAAAAAAAvlsYiT5A6b0OAAAAAAAAAFCvrA0QowMYAAAAAAAAACD+BrIOMvYenBx9sAIAAAAAAAAAyBtkZg0+5QchAAAAAAAAAMBrfxi5zBMgAAAAAAAAAED6P1gAAAAAAAAAAJwHEL2u3MFJAAAAAAAA4ECA6PXR+5Q7MQIAAAAAAAAcCBDd4DBrEFpmpygAAAAAAACwMIAkSZIkSZIkSZIkSZIkSZKkZfoHAAD//2LuOQH73iDVAAAAAElFTkSuQmCC"

func TestMain(m *testing.M) {
	app = App{}
	os.Mkdir("./test_tmp", os.ModeDir)
	app.Initialize("./test_tmp", "./test_tmp", -1)
	code := m.Run()
	os.RemoveAll("./test_tmp")
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}

func TestCreateClient(t *testing.T) {
	t.Log("test create client")
	req, _ := http.NewRequest("GET", "/api/client", nil)
	response := executeRequest(req)
	bodyString := response.Body.String()
	t.Log(bodyString)
	var client Client
	err := json.Unmarshal(response.Body.Bytes(), &client)
	if err != nil {
		t.Fatal("Creating client does not return valid json")
	}
	if client.ID != "c1" {
		t.Fatal("Expected client to have id c1")
	}
	if client.QR != expectedQR {
		t.Fatal("Expected client to have QR code of string: c1")
	}
}

func TestGetClient(t *testing.T) {
	t.Log("test create client")
	req, _ := http.NewRequest("GET", "/api/client", nil)
	createResponse := executeRequest(req)
	req, _ = http.NewRequest("GET", "/api/client/c1", nil)
	getResponse := executeRequest(req)
	createBodyString := createResponse.Body.String()
	getBodyString := getResponse.Body.String()
	if createBodyString != getBodyString {
		t.Fatal("Expected get client call to return client just created")
	}
}

func TestGetClients(t *testing.T) {
	t.Log("test create client")
	req, _ := http.NewRequest("GET", "/api/client", nil)
	executeRequest(req)
	app.MockRandom = -2
	req, _ = http.NewRequest("GET", "/api/client", nil)
	executeRequest(req)
	req, _ = http.NewRequest("GET", "/api/clients", nil)
	getResponse := executeRequest(req)
	var actualClientMap map[string]Client
	err := json.Unmarshal(getResponse.Body.Bytes(), &actualClientMap)
	if err != nil {
		t.Fatal("Error parsing client map json")
	}
}

func TestGetNonexistantClient(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/client/bob", nil)
	getResponse := executeRequest(req)
	bodyString := getResponse.Body.String()
	if bodyString != "No client exists with id bob" {
		t.Fatal("Expected correct error msg for nonexistant client")
	}
}

func TestSetClient(t *testing.T) {
	t.Log("test set client")
	client := Client{
		ID:      "c1",
		Name:    "Bob",
		GoToURL: "http://scratch.mit.edu",
		QR:      "",
	}
	jsonBytes, err := json.Marshal(client)
	if err != nil {
		t.Fatal("Failed to marshal test json")
	}
	req, _ := http.NewRequest("POST", "/api/client/c1", bytes.NewReader(jsonBytes))
	setResponse := executeRequest(req)
	var returnedClient Client
	err = json.Unmarshal(setResponse.Body.Bytes(), &returnedClient)
	if err != nil {
		t.Fatal("Error parsing client map json")
	}
	if returnedClient != client {
		t.Fatal("Service should return same client posted to it")
	}
	req, _ = http.NewRequest("GET", "/api/client/c1", nil)
	getResponse := executeRequest(req)
	err = json.Unmarshal(getResponse.Body.Bytes(), &returnedClient)
	if err != nil || returnedClient != client {
		t.Fatal("Get client service should return client that was posted")
	}
}

func TestServeWebIndex(t *testing.T) {
	expectedBodyString := "<html><body>Hi</body></html>"
	ioutil.WriteFile(filepath.Join("./test_tmp", "index.html"), []byte(expectedBodyString), os.ModePerm)
	prefixes := []string{"/client", "/admin"}
	for _, prefix := range prefixes {
		req, _ := http.NewRequest("GET", prefix, nil)
		getResponse := executeRequest(req)
		redirectLocation := getResponse.HeaderMap["Location"][0]
		if getResponse.Code != 301 || redirectLocation != fmt.Sprint(prefix, "/") {
			t.Fatal("Failed to redirect to correct location for prefix ", prefix)
		}
		req, _ = http.NewRequest("GET", redirectLocation, nil)
		getResponse = executeRequest(req)
		bodyString := getResponse.Body.String()
		if bodyString != expectedBodyString {
			t.Fatal("Failed to serve index.html for redirect location ", redirectLocation)
		}
	}
}

func TestServeWebFile(t *testing.T) {
	expectedBodyString := "function() { console.log(5) }"
	ioutil.WriteFile(filepath.Join("./test_tmp", "main.js"), []byte(expectedBodyString), os.ModePerm)
	prefixes := []string{"/client", "/admin"}
	for _, prefix := range prefixes {
		req, _ := http.NewRequest("GET", fmt.Sprint(prefix, "/main.js"), nil)
		getResponse := executeRequest(req)
		bodyString := getResponse.Body.String()
		if bodyString != expectedBodyString {
			t.Fatal("Failed to return contents of file (main.js) on prefix ", prefix)
		}
	}
}
