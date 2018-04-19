SHELL := /bin/bash

.PHONY: generate
generate:
	go generate ./...

.PHONY: generate-certificates
generate-certificates:
	openssl req -nodes -x509 -newkey -sha256 -newkey rsa:2048 -keyout internal/tls/server.key -days 3650 -subj "/C=US/ST=CA/O=Acme, Inc./CN=grpc-helloworld.192.168.99.100.nip.io" -reqexts SAN -config <(cat /etc/ssl/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:grpc-helloworld.192.168.99.100.nip.io,DNS:grpc-helloworld.192.168.99.100.nip.io:443,DNS:grpc-helloworld.192.168.99.100.nip.io:8443,DNS:localhost,DNS:localhost:443,DNS:localhost:8443,IP:127.0.0.1")) -out internal/tls/server.pem
