package main

import (
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
	id   string
	name string
}

// NewClient Creates new Client
func NewClient(id string) *Client {
	c := new(Client)
	c.id = id
	c.name = ""
	return c
}

var clientList map[string]*Client

func clientHandler(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	clientID := urlParts[len(urlParts)-1]
	client := NewClient(clientID)
	clientList = append(clientList, client)
	fmt.Fprintf(w, "Hello client id %s!\n", clientID)
	fmt.Fprintf(w, "Your path is %s", r.URL.Path)
	for i, val := range clientList {
		fmt.Fprintf(w, "Other client %d %s", i, &val)
	}
}

func randomNumber() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	return rand.New(s1).Int()
}

var domain = flag.String("domain", "localhost", "Domain name for links")

func qrHandler(w http.ResponseWriter, r *http.Request) {
	var png []byte
	clientId := "c" + fmt.Sprint(randomNumber())
	link := fmt.Sprintf("http://%v:8080/client/c%v", *domain, clientId)
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
	http.HandleFunc("/client/", clientHandler)
	http.HandleFunc("/qr", qrHandler)
	fmt.Println("Starting http server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
