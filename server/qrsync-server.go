package main

import (
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
	ID   string
	Name string
}

var clientMap map[string]Client

// NewClient Creates new Client
func newClient(id string) Client {
	c := new(Client)
	c.ID = id
	c.Name = fmt.Sprintf("unnamed client %d", len(clientMap)+1)
	return *c
}

func getOrCreateClient(id string) Client {
	client, ok := clientMap[id]
	if !ok {
		client = newClient(id)
		clientMap[id] = client
	}
	return client
}

func addClientHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	clientID := urlParts[len(urlParts)-1]
	client := getOrCreateClient(clientID)
	fmt.Fprintf(w, "Hello client id %s!\n", client.ID)
	fmt.Fprintf(w, "Your path is %s\n", r.URL.Path)
	for id, client := range clientMap {
		fmt.Fprintf(w, "Other client %s %s \n", id, client.Name)
	}
}

func clientListHandler(w http.ResponseWriter, r *http.Request) {
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
var port = flag.Int("port", 80, "Port to run server on")

func qrHandler(w http.ResponseWriter, r *http.Request) {
	var png []byte
	clientID := "c" + fmt.Sprint(randomNumber())
	link := fmt.Sprintf("http://%v:8080/client/c%v", *domain, clientID)
	fmt.Println(link)
	png, err := qrcode.Encode(link, qrcode.Medium, 256)
	if err != nil {
		log.Println("Error generating QR code")
		log.Println(err.Error())
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func main() {
	flag.Parse()
	fmt.Printf("Domain for links is %v \n", *domain)
	clientMap = make(map[string]Client)
	fileServer := http.FileServer(http.Dir("../client"))
	http.Handle("/", fileServer)
	http.HandleFunc("/client/", addClientHandler)
	http.HandleFunc("/qr", qrHandler)
	http.HandleFunc("/clients", clientListHandler)
	fmt.Println("Starting http server on port ", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
