build:
	go build main.go

server: build
	./main --type=server

client: build
	./main --type=client

hc: build
	./main --type=heartbeatclient

hs: build
	./main --type=heartbeatserver

cds: build 
	./main --type=commanddispatcherserver

cdc: build 
	./main --type=commanddispatcherclient
