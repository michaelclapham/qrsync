package main

import "flag"

func main() {
	var port = flag.Int("port", 443, "Port to run server on")
	flag.Parse()
	var app = App{}
	app.Initialize()
	app.ListenOnPort(*port)
}
