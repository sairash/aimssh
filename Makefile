build:
	go build -o binary/pomossh

ssh: build
	./binary/pomossh -ssh true

pomo: build
	./binary/pomossh 

test:
	go test ./... -v