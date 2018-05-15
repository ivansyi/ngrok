.PHONY: default server client deps fmt clean all release-all assets client-assets server-assets contributors
export GOPATH:=$(shell pwd)

BUILDTAGS=release
default: all

deps: assets
	go get -tags '$(BUILDTAGS)' -d -v ngrok/...

server: deps
	go install -tags '$(BUILDTAGS)' ngrok/main/ngrokd

fmt:
	go fmt ngrok/...

client: deps
	GOOS=linux GOARCH=arm go install -tags '$(BUILDTAGS)' ngrok/main/ngrok

assets: client-assets server-assets

bin/go-bindata:
	GOOS="" GOARCH="" go get github.com/jteeuwen/go-bindata/go-bindata

client-assets: bin/go-bindata
	bin/go-bindata -nomemcopy -pkg=assets -tags=$(BUILDTAGS) \
		-debug=$(if $(findstring debug,$(BUILDTAGS)),true,false) \
		-o=src/ngrok/client/assets/assets_$(BUILDTAGS).go \
		assets/client/...

server-assets: bin/go-bindata
	bin/go-bindata -nomemcopy -pkg=assets -tags=$(BUILDTAGS) \
		-debug=$(if $(findstring debug,$(BUILDTAGS)),true,false) \
		-o=src/ngrok/server/assets/assets_$(BUILDTAGS).go \
		assets/server/...

release-client: BUILDTAGS=release
release-client: client

release-server: BUILDTAGS=release
release-server: server

release-all: fmt release-client release-server

all: fmt client server

clean:
	go clean -i -r ngrok/...
	rm -rf src/ngrok/client/assets/ src/ngrok/server/assets/

contributors:
	echo "Contributors to ngrok, both large and small:\n" > CONTRIBUTORS
	git log --raw | grep "^Author: " | sort | uniq | cut -d ' ' -f2- | sed 's/^/- /' | cut -d '<' -f1 >> CONTRIBUTORS

rootca:
	openssl genrsa -out rootCA.key 2048
	openssl req -x509 -new -nodes -key rootCA.key -subj "/CN=ktvdaren.com" -days 5000 -out rootCA.pem
	cp rootCA.pem assets/client/tls/ngrokroot.crt

tls:
	openssl genrsa -out server.key 2048
	openssl req -new -key server.key -subj "/CN=ngrok.ktvdaren.com" -out server.csr
	openssl x509 -req -in server.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out server.crt -days 5000
	cp server.crt assets/server/tls/snakeoil.crt
	cp device.key assets/server/tls/snakeoil.key

cert:
	if [ "Y${DOMNAME}" = "Y" ];then echo "please specify env viariable: DOMNAME at first." && exit 1; else echo -n ""; fi
	openssl genrsa -out ${DOMNAME}-server.key 2048
	openssl req -new -key ${DOMNAME}-server.key -subj "/CN=*.${DOMNAME}" -out ${DOMNAME}-server.csr
	openssl x509 -req -in ${DOMNAME}-server.csr -CA assets/client/tls/ngrokroot.crt -CAkey rootCA.key -CAcreateserial -out ${DOMNAME}-server.crt -days 5000
