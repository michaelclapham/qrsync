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

func clientHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	clientID := urlParts[len(urlParts)-1]
	client := getOrCreateClient(clientID)

	// If POST modify client
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var nClient Client
		err := decoder.Decode(&nClient)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, "{ \"error\": \"error marshalling json\"}")
		} else {
			w.WriteHeader(200)
			fmt.Fprint(w, "{ \"success\": \"added new client\"}")
		}
	} else {
		jsonBytes, err := json.Marshal(client)
		if err != nil {
			fmt.Fprint(w, "{ \"error\": \"error marshalling json\"}")
		} else {
			w.Write(jsonBytes)
		}
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

func createBlankClientHandler(w http.ResponseWriter, r *http.Request) {
	var png []byte
	clientID := "c" + fmt.Sprint(randomNumber())
	link := fmt.Sprintf("http://%v:8080/client/c%v", *domain, clientID)
	fmt.Println(link)
	png, qrError := qrcode.Encode(link, qrcode.Medium, 256)
	if qrError != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "Error generating QR code")
		log.Println("Error generating QR code")
		log.Println(qrError.Error())
	}
	base64Png := base64.StdEncoding.EncodeToString(png)
	res := struct {
		ID string `json:"id"`
		Qr string `json:"qr"`
	}{
		clientID,
		base64Png,
	}
	jsonBytes, _ := json.MarshalIndent(res, "", "    ")
	w.Write(jsonBytes)
}

func main() {
	flag.Parse()
	fmt.Printf("Domain for links is %v \n", *domain)
	clientMap = make(map[string]Client)
	fileServer := http.FileServer(http.Dir("../client"))
	http.Handle("/", fileServer)
	http.HandleFunc("/newclient", createBlankClientHandler)
	http.HandleFunc("/client/", clientHandler)
	http.HandleFunc("/clients", clientListHandler)
	fmt.Println("Starting http server on port ", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
