# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

#CGO cross compiler for  pi
CC=arm-linux-gnueabihf-gcc

# Binary names
BINARY_NAME=door_control
BINARY_UNIX=$(BINARY_NAME)_rpi

all: build build-pi
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v 
		./$(BINARY_NAME)
deps:



# Cross compilation
build-pi:
		CC=$(CC) GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=1 $(GOBUILD) -v -x -o $(BINARY_UNIX)  -ldflags="-extld=$(CC) "

docker-build:
		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v