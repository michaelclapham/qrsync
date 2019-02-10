# qrsync ![goReportCard](https://goreportcard.com/badge/github.com/michaelclapham/qrsync)
A way of syncing web clients using your phone
Launch an admin device on a phone (or any device with a camera), then launch clients in web browsers you wish to control.
Enter the URL you'd like your clients to go to in the admin UI, then scan the QR code in each clients web browser for it to redirect to it.

## Server Instructions
Download and setup the Go programming language then:
```console
cd qrsync/server
go build
./server 
```
(or server.exe if you're on Windows)

Install openssl and use commands in gen_ssl.ps1 (syntax is the same as bash) to generate server.crt and server.key

Now go to https://localhost/client for each web browser you'd like to control (clients), and load up https://localhost/admin on the phone you'd like to control them from.

To run on a different port use
./server -port 7000

Server must be running on HTTPS for WebRTC (video capture for QR codes) to work.

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