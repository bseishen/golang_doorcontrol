# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

#CGO cross compiler for arm
CC=arm-linux-gnueabihf-gcc
#CGO cross compiler for arm64
CC64=aarch64-linux-gnu-gcc

# Binary names
BINARY_NAME=door_control
BINARY_ARMV5=$(BINARY_NAME)_ARMV5
BINARY_ARMV8=$(BINARY_NAME)_ARMV8

all: build build-ARMv5 build-ARMv8
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_ARMV5)
		rm -f $(BINARY_ARMV8)
		rm -f $(BINARY_NAME)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v
		./$(BINARY_NAME)
deps:



# Cross compilation
build-ARMv5:
		CC=$(CC) GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=1 $(GOBUILD) -v -x -o $(BINARY_ARMV5)  -ldflags="-extld=$(CC) "

build-ARMv8:
		CC=$(CC64) GOOS=linux GOARCH=arm64 CGO_ENABLED=1 $(GOBUILD) -v -x -o $(BINARY_ARMV8)  -ldflags="-extld=$(CC64) "

docker-build:
		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
