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
	} else {
		w.Write(jsonBytes)
	}
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/client", a.createClient).Methods("GET")
}

// ListenOnPort Starts the app listening on the provided port
func (a *App) ListenOnPort(port int) {
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", port), a.Router))
}
