# Commands used to generate SSL private key and certificate locally.
# Will use different commands on actual Linux server
$env:OPENSSL_CONF = "C:\Program Files (x86)\GnuWin32\share\openssl.cnf"
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650