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

// Client representation
type Client struct {
	ID      string `json:"id"`
	QR      string `json:"qr"`
	Name    string `json:"name"`
	GoToURL string `json:"gotoUrl"`
}

// App Stores the state of our web server
type App struct {
	Router       *mux.Router
	ClientMap    map[string]Client
	MockRandom   int
	ClientUIPath string
	AdminUIPath  string
}

func (a *App) randomPositiveInt() int {
	if a.MockRandom < 0 {
		return -1 * a.MockRandom
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	return rand.New(s1).Int()
}

// Initialize sets up the app
func (a *App) Initialize(clientUIPath string, adminUIPath string, randomSeed int) {
	a.Router = mux.NewRouter()
	a.ClientMap = make(map[string]Client)
	a.MockRandom = randomSeed
	a.ClientUIPath = clientUIPath
	a.AdminUIPath = adminUIPath
	a.initializeRoutes()
}

func respondAndLogError(w http.ResponseWriter, errorType string, err error) {
	w.WriteHeader(500)
	errorString := ""
	if err != nil {
		errorString = err.Error()
	}
	fmt.Fprint(w, errorType, errorString)
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
		respondAndLogError(w, fmt.Sprint("No client exists with id ", clientID), nil)
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
	a.ClientMap[clientID] = inputClient
	jsonBytes, jsonErr := json.MarshalIndent(a.ClientMap[clientID], "", "    ")
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

func (a *App) setupFileHandler(prefix string, dir string) {
	fileServer := http.FileServer(http.Dir(dir))
	slashOnEnd := fmt.Sprint(prefix, "/")
	a.Router.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, slashOnEnd, 301)
	})
	a.Router.PathPrefix(slashOnEnd).Handler(http.StripPrefix(slashOnEnd, fileServer))
}

func (a *App) initializeRoutes() {
	a.setupFileHandler("/client", a.ClientUIPath)
	a.setupFileHandler("/admin", a.AdminUIPath)
	a.Router.HandleFunc("/api/client", a.createClient).Methods("GET")
	a.Router.HandleFunc("/api/client/{id}", a.getClient).Methods("GET")
	a.Router.HandleFunc("/api/client/{id}", a.setClient).Methods("POST")
	a.Router.HandleFunc("/api/clients", a.getClients).Methods("GET")
}

// ListenOnPort Starts the app listening on the provided port
func (a *App) ListenOnPort(port int) error {
	fmt.Println("Starting https server on port ", port)
	return http.ListenAndServeTLS(fmt.Sprint(":", port), "ssl/server.crt", "ssl/server.key", a.Router)
}
