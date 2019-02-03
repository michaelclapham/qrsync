package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var app App

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
	var client Client
	err := json.Unmarshal(response.Body.Bytes(), &client)
	if err != nil {
		t.Fatal("Creating client does not return valid json")
	}
}
