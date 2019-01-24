# qrsync ![goReportCard](https://goreportcard.com/badge/github.com/michaelclapham/qrsync)
A way of syncing web clients using your phone

## Server Instructions
Download and setup the Go programming language then:
cd qrsync/server
go run qrsync-server

Now go to localhost:8080/qr to get a QR code

To run on a different domain use
go run qrsync-server -domain example.org

This will still run on localhost, but your links will use the domain

## Client Intructions
Coming soon