# qrsync ![goReportCard](https://goreportcard.com/badge/github.com/michaelclapham/qrsync)
A way of syncing web clients using your phone

## Server Instructions
Download and setup the Go programming language then:
```console
cd qrsync/server
go run qrsync-server
```

Now go to localhost:8080/qr to get a QR code

To run on a different domain use
go run qrsync-server -domain example.org

This will still run on localhost, but your links will use the domain

## Client Instructions
The client is written in typescript and uses no dependencies so to make changes simply run
```console
cd client
tsc -w
```
The client web UI is served by the main go server

## Web Admin Intructions
Copy ZXing minified js from: https://github.com/aleris/zxing-typescript/blob/master/docs/examples/zxing.qrcodereader.min.js to ./web/admin/lib
Admin web project is written in Typescript (version included in package.json) so run
```console
npm install --save-dev
tsc -w
```
To start compiling web admin in watch mode.
The admin web UI is served by the main go server