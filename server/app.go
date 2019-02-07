package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	qrcode "github.com/skip2/go-qrcode"
)

// App Stores the state of our web server
type App struct {
	Router     *mux.Router
	ClientMap  map[string]Client
	MockRandom int
}

func (a *App) randomPositiveInt() int {
	if a.MockRandom < 0 {
		return -1 * a.MockRandom
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	return rand.New(s1).Int()
}

// Initialize sets up the app
func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.ClientMap = make(map[string]Client)
	a.MockRandom = -1
	a.initializeRoutes()
}

func respondAndLogError(w http.ResponseWriter, errorType string, err error) {
	w.WriteHeader(500)
	fmt.Fprint(w, errorType, err)
	log.Println(errorType, ":", err)
}

func (a *App) createClient(w http.ResponseWriter, r *http.Request) {
	var png []byte
	clientID := "c" + fmt.Sprint(a.randomPositiveInt())
	png, qrError := qrcode.Encode(clientID, qrcode.Medium, 256)
	if qrError != nil {
		respondAndLogError(w, "QR Code generation error", qrError)
		return
	}
	base64Png := base64.StdEncoding.EncodeToString(png)
	newClient := Client{
		clientID,
		base64Png,
		"",
		"",
	}
	a.ClientMap[clientID] = newClient
	jsonBytes, jsonErr := json.MarshalIndent(newClient, "", "    ")
	if jsonErr != nil {
		respondAndLogError(w, "JSON formatting error", jsonErr)
		return
	}
	w.Write(jsonBytes)
}

func (a *App) getClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID, ok := vars["id"]
	if !ok {
		respondAndLogError(w, fmt.Sprint("No id provided"), nil)
		return
	}
	client, ok := a.ClientMap[clientID]
	if !ok {
		respondAndLogError(w, fmt.Sprint("Not client with id ", clientID), nil)
		return
	}
	jsonBytes, jsonErr := json.MarshalIndent(client, "", "    ")
	if jsonErr != nil {
		respondAndLogError(w, "JSON formatting error", jsonErr)
		return
	}
	w.Write(jsonBytes)
}

func (a *App) setClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID, ok := vars["id"]
	if !ok {
		respondAndLogError(w, fmt.Sprint("No id provided"), nil)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var inputClient Client
	err := decoder.Decode(&inputClient)
	if err != nil {
		respondAndLogError(w, "JSON parsing error", err)
		return
	}
	clientMap[clientID] = inputClient
	jsonBytes, jsonErr := json.MarshalIndent(clientMap[clientID], "", "    ")
	if jsonErr != nil {
		respondAndLogError(w, "JSON formatting error", jsonErr)
		return
	}
	w.Write(jsonBytes)
}

func (a *App) getClients(w http.ResponseWriter, r *http.Request) {
	jsonBytes, jsonErr := json.MarshalIndent(a.ClientMap, "", "    ")
	if jsonErr != nil {
		respondAndLogError(w, "JSON formatting error", jsonErr)
		return
	}
	w.Write(jsonBytes)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/client", a.createClient).Methods("GET")
	a.Router.HandleFunc("/api/client/{id}", a.getClient).Methods("GET")
	a.Router.HandleFunc("/api/client/{id}", a.setClient).Methods("POST")
	a.Router.HandleFunc("/api/clients", a.getClients).Methods("GET")
}

// ListenOnPort Starts the app listening on the provided port
func (a *App) ListenOnPort(port int) {
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", port), a.Router))
}
