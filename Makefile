EXTENSION=out
DELETE_COMMAND=rm

ifeq ($(OS),Windows_NT)
EXTENSION=exe
DELETE_COMMAND=del
endif

all: client.$(EXTENSION) server.$(EXTENSION)

client.$(EXTENSION): api.pb.go
	go build -o client.$(EXTENSION) ./client

server.$(EXTENSION): api.pb.go
	go build -o server.$(EXTENSION) ./server

api.pb.go: 
	protoc -I ./api --go_out=plugins=grpc:./api ./api/api.proto

clean: 
	$(DELETE_COMMAND) *.exe