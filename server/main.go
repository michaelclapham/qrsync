package main

import (
	"flag"
	"log"
)

func main() {
	ssl := flag.Bool("ssl", true, "Use ssl")
	port := flag.Int("port", 443, "Port to run server on")
	flag.Parse()
	var app = App{}
	app.Initialize("../web/client", "../web/admin", 1)
	log.Fatal(app.ListenOnPort(*port, *ssl))
}
