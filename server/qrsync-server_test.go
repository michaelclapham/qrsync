package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
func TestCreateClientHandler(t *testing.T) {

	req := httptest.NewRequest("GET", "https://localhost", nil)
	w := httptest.NewRecorder()
	CreateClientHandler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	if resp.StatusCode != 200 {
		t.Fail()
		fmt.Printf("Status code should be 200")
	}

	var client Client
	err := json.Unmarshal(body, &client)
	if err != nil {
		fmt.Print(err)
		t.Fail()
		fmt.Printf("Response Body was not valid client JSON")
	}
}
