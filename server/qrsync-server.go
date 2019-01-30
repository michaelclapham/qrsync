package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

// Client representation
type Client struct {
	ID      string `json:"id"`
	QR      string `json:"qr"`
	Name    string `json:"name"`
	GoToURL string `json:"gotoUrl"`
}

var clientMap map[string]Client

func clientHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	clientID := urlParts[len(urlParts)-1]

	// If POST modify client
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var inputClient Client
		err := decoder.Decode(&inputClient)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, "{ \"error\": \"error in input\"}")
		}
		clientMap[clientID] = inputClient
	}

	// Always return client
	jsonBytes, err := json.Marshal(clientMap[clientID])
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err.Error())
	} else {
		w.Write(jsonBytes)
	}
}

func clientListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var updatedClients []Client
		err := decoder.Decode(&updatedClients)
		if err != nil {
			fmt.Fprint(w, "{ \"error\": \"error decoding json\"}")
		} else {
			for _, client := range updatedClients {
				clientMap[client.ID] = client
			}
		}

	}

	jsonBytes, err := json.Marshal(clientMap)
	if err != nil {
		fmt.Fprint(w, "{ \"error\": \"error marshalling json\"}")
	} else {
		w.Write(jsonBytes)
	}
}

func randomNumber() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	return rand.New(s1).Int()
}

var domain = flag.String("domain", "localhost", "Domain name for links")
var port = flag.Int("port", 443, "Port to run server on")

// Handles a request to create a new client with a corresponding QR code containing the clients ID
func CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	var png []byte
	clientID := "c" + fmt.Sprint(randomNumber())
	link := fmt.Sprint(clientID)
	fmt.Println(link)
	png, qrError := qrcode.Encode(link, qrcode.Medium, 256)
	if qrError != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error generating QR code")
		log.Println("Error generating QR code")
		log.Println(qrError.Error())
	}
	base64Png := base64.StdEncoding.EncodeToString(png)
	res := Client{
		clientID,
		base64Png,
		"",
		"",
	}
	clientMap[clientID] = res
	jsonBytes, _ := json.MarshalIndent(res, "", "    ")
	w.Write(jsonBytes)
}

func main() {
	flag.Parse()
	fmt.Printf("Domain for links is %v \n", *domain)
	clientMap = make(map[string]Client)
	fileServer := http.FileServer(http.Dir("../web"))
	http.Handle("/", fileServer)
	http.HandleFunc("/api/newclient", CreateClientHandler)
	http.HandleFunc("/api/client/", clientHandler)
	http.HandleFunc("/api/clients", clientListHandler)
	fmt.Println("Starting https server on port ", *port)
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), "ssl/server.crt", "ssl/server.key", nil))
}
