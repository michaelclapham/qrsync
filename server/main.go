package main

import (
	"flag"
	"log"
)

func main() {
	ssl := flag.Bool("ssl", false, "If true use SSL")
	defaultPort := 4002
	if *ssl {
		defaultPort = 4001
	}
	port := flag.Int("port", defaultPort, "Port to run server on")
	flag.Parse()
	var app = App{}
	app.Initialize("../web/client", "../web/admin", 1)
	log.Fatal(app.ListenOnPort(*port, *ssl))
}
