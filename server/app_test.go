package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var app App

var expectedQR = "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEACAMAAABrrFhUAAAABlBMVEX///8AAABVwtN+AAABp0lEQVR42uzdYU7EIBCA0e79L70naArD0A7wvn+apus+TZwU0EuSJEmSJEmSJEmSJEmSJEmSfsHu7nN33+zXAwAAAAAAAADkAWRd3wuV/fUBAAAAAAAAAOIAT4PI08etg1D09QAAAAAAAAAA9QBaH2AAAAAAAAAAAPYdhLIBAQAAAAAAAAB1F0aig882K0MAAAAAAADABgDRjYtZg9KyO0UBAAAAAACAhQHeArx2DQAAAAAAANjg9/isNzT7/gAAAAAAAACA+Qslowcnej9f5icCAAAAAAAA2Big9QVbF0h6NzyW+Y4DAAAAAAAABwNEB5fR60YPYgIAAAAAAAAAvlsYiT5A6b0OAAAAAAAAAFCvrA0QowMYAAAAAAAAACD+BrIOMvYenBx9sAIAAAAAAAAAyBtkZg0+5QchAAAAAAAAAMBrfxi5zBMgAAAAAAAAAED6P1gAAAAAAAAAAJwHEL2u3MFJAAAAAAAA4ECA6PXR+5Q7MQIAAAAAAAAcCBDd4DBrEFpmpygAAAAAAACwMIAkSZIkSZIkSZIkSZIkSZKkZfoHAAD//2LuOQH73iDVAAAAAElFTkSuQmCC"

func TestMain(m *testing.M) {
	app = App{}
	app.Initialize()
	app.MockRandom = -1
	code := m.Run()
	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}

func TestCreateClient(t *testing.T) {
	t.Log("test create client")
	req, _ := http.NewRequest("GET", "/client", nil)
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
	req, _ := http.NewRequest("GET", "/client", nil)
	createResponse := executeRequest(req)
	req, _ = http.NewRequest("GET", "/client/c1", nil)
	getResponse := executeRequest(req)
	createBodyString := createResponse.Body.String()
	getBodyString := getResponse.Body.String()
	if createBodyString != getBodyString {
		t.Fatal("Expected get client call to return client just created")
	}
}
